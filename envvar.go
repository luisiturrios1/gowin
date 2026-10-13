// Modified 2026-10-12: modernized documentation and implementation.

package gowin

import "os"

// GetEnvVar returns a process environment variable, or an empty string if unset.
// It works on every platform supported by Go.
func GetEnvVar(name string) string { return os.Getenv(name) }

// WriteEnvVar changes a process environment variable. It does not persist the
// variable to the Windows registry or update other running processes.
func WriteEnvVar(name, val string) error { return os.Setenv(name, val) }
