//go:build windows

package gowin

import (
	"crypto/rand"
	"errors"
	"fmt"
	"slices"
	"strings"
	"sync"
	"syscall"
	"testing"
	"unsafe"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
)

// Each test owns a random HKCU leaf. No admin rights or existing user keys are used.
func testRegistryKey(t testing.TB) string {
	t.Helper()
	path := `Software\GowinTest-` + rand.Text()
	key, _, err := registry.CreateKey(registry.CURRENT_USER, path, registry.SET_VALUE)
	if err != nil {
		t.Fatal(err)
	}
	if err := key.Close(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := registry.DeleteKey(registry.CURRENT_USER, path); err != nil {
			t.Errorf("cleanup %s: %v", path, err)
		}
	})
	return path
}

func TestRegistryStrings(t *testing.T) {
	path := testRegistryKey(t)
	for _, tc := range []struct{ name, value string }{
		{"", "default"}, {"empty", ""}, {"unicode", "México 日本 😀"},
		{"long", strings.Repeat("abc😀", 10000)}, {"literal", "%USERPROFILE%\\Documents"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if err := WriteStringReg("HKCU", path, tc.name, tc.value); err != nil {
				t.Fatal(err)
			}
			got, err := GetReg("HKCU", path, tc.name)
			if err != nil || got != tc.value {
				t.Fatalf("length=%d, want %d, err=%v", len(got), len(tc.value), err)
			}
		})
	}
	if err := WriteStringReg("HKCU", path, "unicode", "replacement"); err != nil {
		t.Fatal(err)
	}
	if got, err := GetReg("HKCU", path, "unicode"); err != nil || got != "replacement" {
		t.Fatalf("overwrite: %q %v", got, err)
	}
}

func TestRegistryDWORD(t *testing.T) {
	path := testRegistryKey(t)
	for _, want := range []uint32{0, 1, 65536, 0xffffffff} {
		t.Run(fmt.Sprint(want), func(t *testing.T) {
			if err := WriteDwordReg("HKCU", path, "number", want); err != nil {
				t.Fatal(err)
			}
			if got, err := GetDwordReg("HKCU", path, "number"); err != nil || got != want {
				t.Fatalf("got %d err %v", got, err)
			}
		})
	}
	raw, err := GetRegRaw("HKCU", path, "number")
	if err != nil || !slices.Equal(raw, []uint16{65535, 65535}) {
		t.Fatalf("legacy DWORD raw=%x err=%v", raw, err)
	}
	if err := WriteStringReg("HKCU", path, "text", "abc"); err != nil {
		t.Fatal(err)
	}
	if _, err := GetDwordReg("HKCU", path, "text"); !errors.Is(err, registry.ErrUnexpectedType) {
		t.Fatalf("type error: %v", err)
	}
	key, err := registry.OpenKey(registry.CURRENT_USER, path, registry.SET_VALUE)
	if err != nil {
		t.Fatal(err)
	}
	defer key.Close()
	if err := key.SetQWordValue("qword", 1); err != nil {
		t.Fatal(err)
	}
	if _, err := GetDwordReg("HKCU", path, "qword"); !errors.Is(err, registry.ErrUnexpectedType) {
		t.Fatalf("QWORD must not truncate: %v", err)
	}
}

func TestRegistryLegacyTypes(t *testing.T) {
	path := testRegistryKey(t)
	key, err := registry.OpenKey(registry.CURRENT_USER, path, registry.SET_VALUE)
	if err != nil {
		t.Fatal(err)
	}
	defer key.Close()
	if err := key.SetExpandStringValue("expand", `%USERPROFILE%\Documents`); err != nil {
		t.Fatal(err)
	}
	if got, err := GetReg("HKCU", path, "expand"); err != nil || got != `%USERPROFILE%\Documents` {
		t.Fatalf("expanded legacy value: %q %v", got, err)
	}
	if err := key.SetStringsValue("multi", []string{"one", "two"}); err != nil {
		t.Fatal(err)
	}
	if got, err := GetReg("HKCU", path, "multi"); err != nil || got != "one" {
		t.Fatalf("MULTI_SZ: %q %v", got, err)
	}
	if got, err := GetRegRaw("HKCU", path, "multi"); err != nil || !slices.Equal(got, []uint16{'o', 'n', 'e', 0, 't', 'w', 'o'}) {
		t.Fatalf("raw MULTI_SZ: %x %v", got, err)
	}
	if err := key.SetBinaryValue("odd", []byte{1, 2, 3}); err != nil {
		t.Fatal(err)
	}
	if got, err := GetRegRaw("HKCU", path, "odd"); err != nil || !slices.Equal(got, []uint16{0x201, 3}) {
		t.Fatalf("odd binary: %x %v", got, err)
	}
}

func TestRegistryErrorsAndValidation(t *testing.T) {
	path := testRegistryKey(t)
	_, err := GetReg("HKCU", path, "missing")
	if !errors.Is(err, registry.ErrNotExist) {
		t.Fatalf("missing value: %v", err)
	}
	var errno syscall.Errno
	if !errors.As(err, &errno) {
		t.Fatalf("missing wrapped Windows errno: %v", err)
	}
	if _, err := GetReg("HKCU", path+`\missing`, ""); !errors.Is(err, registry.ErrNotExist) {
		t.Fatalf("missing key: %v", err)
	}
	for _, tc := range []struct{ hive, path, name string }{
		{"unknown", path, "value"}, {"HKCU", path + "\x00bad", "value"}, {"HKCU", path, "bad\x00name"},
	} {
		t.Run(fmt.Sprintf("%q/%q/%q", tc.hive, tc.path, tc.name), func(t *testing.T) {
			_, rawErr := GetRegRaw(tc.hive, tc.path, tc.name)
			_, dwordErr := GetDwordReg(tc.hive, tc.path, tc.name)
			for _, err := range []error{rawErr, dwordErr, WriteStringReg(tc.hive, tc.path, tc.name, "x"), WriteDwordReg(tc.hive, tc.path, tc.name, 1), DeleteKey(tc.hive, tc.path, tc.name)} {
				if err == nil {
					t.Fatal("invalid input succeeded")
				}
			}
		})
	}
	if err := WriteStringReg("HKCU", path+`\invalid`, "value", "bad\x00value"); !errors.Is(err, syscall.EINVAL) {
		t.Fatalf("NUL value: %v", err)
	}
	if _, err := GetReg("HKCU", path+`\invalid`, "value"); !errors.Is(err, registry.ErrNotExist) {
		t.Fatalf("invalid input created key: %v", err)
	}
}

func TestRegistryRoots(t *testing.T) {
	for name, want := range map[string]registry.Key{"HKLM": registry.LOCAL_MACHINE, "HKCU": registry.CURRENT_USER, "HKCC": registry.CURRENT_CONFIG, "HKCR": registry.CLASSES_ROOT, "HKU": registry.USERS} {
		if got, err := registryRoot(name, "", ""); err != nil || got != want {
			t.Fatalf("%s: %v %v", name, got, err)
		}
	}
}

func TestDeleteKey(t *testing.T) {
	path := testRegistryKey(t)
	child := path + `\child`
	if err := WriteStringReg("HKCU", child, "value", "data"); err != nil {
		t.Fatal(err)
	}
	if err := DeleteKey("HKCU", path, "child"); err != nil {
		t.Fatal(err)
	}
	if _, err := GetReg("HKCU", child, "value"); !errors.Is(err, registry.ErrNotExist) {
		t.Fatalf("key still present: %v", err)
	}
	if err := DeleteKey("HKCU", path, "child"); !errors.Is(err, registry.ErrNotExist) {
		t.Fatalf("missing deletion: %v", err)
	}
}

func TestConcurrentRegistry(t *testing.T) {
	path := testRegistryKey(t)
	var wg sync.WaitGroup
	for i := range 8 {
		wg.Go(func() {
			name := fmt.Sprint(i)
			for value := range uint32(50) {
				if err := WriteDwordReg("HKCU", path, name, value); err != nil {
					t.Error(err)
					return
				}
				if got, err := GetDwordReg("HKCU", path, name); err != nil || got != value {
					t.Errorf("got %d err %v", got, err)
					return
				}
			}
		})
	}
	wg.Wait()
}

func TestRegistryWriteHandles(t *testing.T) {
	path := testRegistryKey(t)
	proc := windows.NewLazySystemDLL("kernel32.dll").NewProc("GetProcessHandleCount")
	count := func() uint32 {
		t.Helper()
		var n uint32
		// The Windows API writes one DWORD to this pointer for this call only.
		ok, _, err := proc.Call(uintptr(windows.CurrentProcess()), uintptr(unsafe.Pointer(&n)))
		if ok == 0 {
			t.Fatal(err)
		}
		return n
	}
	// Warm up lazy DLL loading and registry calls before measuring.
	if err := WriteStringReg("HKCU", path, "text", "warmup"); err != nil {
		t.Fatal(err)
	}
	before := count()
	for range 256 {
		if err := WriteStringReg("HKCU", path, "text", "value"); err != nil {
			t.Fatal(err)
		}
		if err := WriteDwordReg("HKCU", path, "number", 1); err != nil {
			t.Fatal(err)
		}
	}
	after := count()
	// Allow unrelated runtime handles, but catch the former one-handle-per-write leak.
	if after > before+16 {
		t.Fatalf("process handles grew from %d to %d", before, after)
	}
}

var benchmarkRegistryRaw []uint16

func BenchmarkRegistryRead(b *testing.B) {
	path := testRegistryKey(b)
	if err := WriteStringReg("HKCU", path, "value", strings.Repeat("x", 4096)); err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()
	for b.Loop() {
		value, err := GetRegRaw("HKCU", path, "value")
		if err != nil {
			b.Fatal(err)
		}
		benchmarkRegistryRaw = value
	}
}
