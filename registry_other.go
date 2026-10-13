//go:build !windows

package gowin

// GetRegRaw requires Windows.
func GetRegRaw(hkey, path, name string) ([]uint16, error) { return nil, ErrUnsupportedPlatform }

// GetDwordReg requires Windows.
func GetDwordReg(hkey, path, name string) (uint32, error) { return 0, ErrUnsupportedPlatform }

// WriteStringReg requires Windows.
func WriteStringReg(hkey, path, name, val string) error { return ErrUnsupportedPlatform }

// WriteDwordReg requires Windows.
func WriteDwordReg(hkey, path, name string, val uint32) error { return ErrUnsupportedPlatform }

// DeleteKey requires Windows.
func DeleteKey(hkey, path, name string) error { return ErrUnsupportedPlatform }
