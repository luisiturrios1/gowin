package gowin

import (
	"bytes"
	"errors"
	"fmt"
	"slices"
	"strings"
	"syscall"
	"testing"
)

func TestLegacyUTF16(t *testing.T) {
	for _, tc := range []struct {
		name string
		data []byte
		want []uint16
	}{
		{"empty", nil, []uint16{}},
		{"terminator", []byte{0, 0}, []uint16{}},
		{"text", []byte{'a', 0, 'b', 0, 0, 0}, []uint16{'a', 'b'}},
		{"embedded NUL", []byte{'a', 0, 0, 0, 'b', 0, 0, 0}, []uint16{'a', 0, 'b'}},
		{"odd binary", []byte{1, 2, 3}, []uint16{0x201, 3}},
		{"DWORD legacy", []byte{255, 255, 255, 255}, []uint16{65535, 65535}},
		{"surrogate pair", []byte{0x3d, 0xd8, 0x00, 0xde, 0, 0}, []uint16{0xd83d, 0xde00}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := legacyUTF16(tc.data); !slices.Equal(got, tc.want) {
				t.Fatalf("got %x, want %x", got, tc.want)
			}
		})
	}
}

func TestValidateRegistryInput(t *testing.T) {
	for _, hive := range []string{"HKLM", "HKCU", "HKCC", "HKCR", "HKU"} {
		t.Run(hive, func(t *testing.T) {
			if err := validateRegistryInput(hive, "", ""); err != nil {
				t.Fatal(err)
			}
		})
	}
	for _, tc := range []struct {
		hive, path, name string
		invalidNUL       bool
	}{
		{"hkcu", "", "", false}, {"UNKNOWN", "", "", false},
		{"HKCU", "bad\x00path", "", true}, {"HKCU", "", "bad\x00name", true},
	} {
		t.Run(fmt.Sprintf("%q/%q/%q", tc.hive, tc.path, tc.name), func(t *testing.T) {
			err := validateRegistryInput(tc.hive, tc.path, tc.name)
			if err == nil {
				t.Fatal("expected error")
			}
			if tc.invalidNUL && !errors.Is(err, syscall.EINVAL) {
				t.Fatalf("missing EINVAL: %v", err)
			}
		})
	}
}

func TestReadRegistryValue(t *testing.T) {
	more := errors.New("more data")
	denied := errors.New("access denied")
	for _, size := range []int{0, 3, 64, 65, 100000} {
		t.Run(fmt.Sprint(size), func(t *testing.T) {
			want := bytes.Repeat([]byte{42}, size)
			got, err := readRegistryValue(func(buf []byte) (int, uint32, error) {
				if len(buf) < len(want) {
					return len(want), 0, fmt.Errorf("wrapped: %w", more)
				}
				copy(buf, want)
				return len(want), 0, nil
			}, more)
			if err != nil || !bytes.Equal(got, want) {
				t.Fatalf("size=%d error=%v", len(got), err)
			}
		})
	}
	for _, tc := range []struct {
		name string
		n    int
		err  error
		want error
	}{
		{"denied", 0, denied, denied}, {"no growth", 64, more, more},
		{"invalid negative", -1, nil, nil}, {"invalid success size", 65, nil, nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := readRegistryValue(func([]byte) (int, uint32, error) { return tc.n, 0, tc.err }, more)
			if err == nil || (tc.want != nil && !errors.Is(err, tc.want)) {
				t.Fatalf("got %v", err)
			}
		})
	}
	t.Run("continuous growth", func(t *testing.T) {
		calls := 0
		_, err := readRegistryValue(func(buf []byte) (int, uint32, error) { calls++; return len(buf) + 1, 0, more }, more)
		if calls != 8 || !errors.Is(err, more) {
			t.Fatalf("calls=%d err=%v", calls, err)
		}
	})
}

func FuzzLegacyUTF16(f *testing.F) {
	for _, seed := range [][]byte{nil, {0}, {'a', 0, 0, 0}, {0x3d, 0xd8, 0, 0xde}, {255, 255, 255}} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, data []byte) {
		original := bytes.Clone(data)
		got := legacyUTF16(data)
		// Independent reference: accumulate low and high bytes, then find the
		// last nonzero unit, as the legacy implementation did.
		want := make([]uint16, (len(data)+1)/2)
		last := 0
		for i, b := range data {
			want[i/2] |= uint16(b) << uint(8*(i%2))
		}
		for i, u := range want {
			if u != 0 {
				last = i + 1
			}
		}
		if !slices.Equal(got, want[:last]) || !bytes.Equal(data, original) {
			t.Fatal("decoding mismatch or modified input")
		}
	})
}

var benchmarkRaw []uint16

func BenchmarkLegacyUTF16(b *testing.B) {
	for _, n := range []int{32, 4096, 100000} {
		b.Run(fmt.Sprint(n), func(b *testing.B) {
			data := []byte(strings.Repeat("a\x00", n))
			b.ReportAllocs()
			b.SetBytes(int64(len(data)))
			for b.Loop() {
				benchmarkRaw = legacyUTF16(data)
			}
		})
	}
}
