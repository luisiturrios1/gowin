// Modified 2026-10-12: safe registry decoding and shared validation.

package gowin

import (
	"encoding/binary"
	"errors"
	"fmt"
	"strings"
	"syscall"
	"unicode/utf16"
)

// ErrUnsupportedPlatform indicates a Windows-only operation on another OS.
var ErrUnsupportedPlatform = errors.New("gowin: operation requires Windows")

func validateRegistryInput(hkey, path, name string) error {
	switch hkey {
	case "HKLM", "HKCU", "HKCC", "HKCR", "HKU":
	default:
		return fmt.Errorf("unknown HKEY: %q", hkey)
	}
	if strings.ContainsRune(path, 0) || strings.ContainsRune(name, 0) {
		return fmt.Errorf("registry path or name contains NUL: %w", syscall.EINVAL)
	}
	return nil
}

// GetReg reads a registry value using the legacy UTF-16 interpretation. It does
// not expand environment variables in REG_EXPAND_SZ. Use GetDwordReg for DWORDs.
func GetReg(hkey, path, name string) (string, error) {
	value, err := GetRegRaw(hkey, path, name)
	if err != nil {
		return "", err
	}
	for i, unit := range value {
		if unit == 0 {
			value = value[:i]
			break
		}
	}
	return string(utf16.Decode(value)), nil
}

// legacyUTF16 retains embedded NULs and removes trailing zero units, as the
// original GetRegRaw did. An odd final byte is padded with zero.
func legacyUTF16(data []byte) []uint16 {
	value := make([]uint16, len(data)/2+len(data)%2)
	for i := 0; i < len(data)/2; i++ {
		value[i] = binary.LittleEndian.Uint16(data[2*i:])
	}
	if len(data)%2 != 0 {
		value[len(value)-1] = uint16(data[len(data)-1])
	}
	end := len(value)
	for end > 0 && value[end-1] == 0 {
		end--
	}
	return value[:end]
}

// readRegistryValue tolerates growth between queries. Bound retries so another
// process continuously modifying the value cannot keep a call alive.
func readRegistryValue(query func([]byte) (int, uint32, error), moreData error) ([]byte, error) {
	buffer := make([]byte, 64)
	for range 8 {
		n, _, err := query(buffer)
		if err == nil {
			if n < 0 || n > len(buffer) {
				return nil, fmt.Errorf("invalid registry value size: %d", n)
			}
			return buffer[:n], nil
		}
		if !errors.Is(err, moreData) || n <= len(buffer) {
			return nil, err
		}
		buffer = make([]byte, n)
	}
	return nil, fmt.Errorf("registry value kept growing: %w", moreData)
}
