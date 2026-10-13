package gowin

import (
	"fmt"
	"os"
	"sync"
	"testing"
)

func TestEnvironment(t *testing.T) {
	const name = "GOWIN_TEST_ENVIRONMENT"
	t.Setenv(name, "initial")
	for _, value := range []string{"", "hello", "México 😀", "%PATH%", "a=b"} {
		t.Run(value, func(t *testing.T) {
			if err := WriteEnvVar(name, value); err != nil {
				t.Fatal(err)
			}
			if got := GetEnvVar(name); got != value {
				t.Fatalf("got %q, want %q", got, value)
			}
		})
	}
	if err := os.Unsetenv(name); err != nil {
		t.Fatal(err)
	}
	if GetEnvVar(name) != "" {
		t.Fatal("unset variable must return empty string")
	}
}

func TestInvalidEnvironment(t *testing.T) {
	for _, tc := range []struct{ name, value string }{{"", "x"}, {"bad=name", "x"}, {"bad\x00name", "x"}, {"GOWIN_TEST_INVALID", "bad\x00value"}} {
		t.Run(fmt.Sprintf("%q/%q", tc.name, tc.value), func(t *testing.T) {
			if err := WriteEnvVar(tc.name, tc.value); err == nil {
				t.Fatal("expected error")
			}
		})
	}
}

func TestConcurrentEnvironment(t *testing.T) {
	const name = "GOWIN_TEST_CONCURRENT"
	t.Setenv(name, "")
	var wg sync.WaitGroup
	for i := range 8 {
		wg.Go(func() {
			for range 100 {
				if err := WriteEnvVar(name, fmt.Sprint(i)); err != nil {
					t.Error(err)
					return
				}
				_ = GetEnvVar(name)
			}
		})
	}
	wg.Wait()
}

var benchmarkEnv string

func BenchmarkGetEnvVar(b *testing.B) {
	b.Setenv("GOWIN_BENCHMARK", "value")
	b.ReportAllocs()
	for b.Loop() {
		benchmarkEnv = GetEnvVar("GOWIN_BENCHMARK")
	}
}
