//go:build windows

package gowin

import (
	"fmt"
	"strings"
	"syscall"

	"golang.org/x/sys/windows/registry"
)

func registryRoot(hkey, path, name string) (registry.Key, error) {
	if err := validateRegistryInput(hkey, path, name); err != nil {
		return 0, err
	}
	switch hkey {
	case "HKLM":
		return registry.LOCAL_MACHINE, nil
	case "HKCU":
		return registry.CURRENT_USER, nil
	case "HKCC":
		return registry.CURRENT_CONFIG, nil
	case "HKCR":
		return registry.CLASSES_ROOT, nil
	default:
		return registry.USERS, nil // HKU; validated above.
	}
}

func registryError(operation, hkey, path, name string, err error) error {
	return fmt.Errorf("%s %s\\%s (%q): %w", operation, hkey, path, name, err)
}

// GetRegRaw reads a registry value as little-endian UTF-16 units, preserving the
// legacy representation for all registry types and stripping trailing zero units.
// hkey must be HKLM, HKCU, HKCC, HKCR or HKU. An empty name reads the default value.
func GetRegRaw(hkey, path, name string) ([]uint16, error) {
	root, err := registryRoot(hkey, path, name)
	if err != nil {
		return nil, err
	}
	key, err := registry.OpenKey(root, path, registry.QUERY_VALUE)
	if err != nil {
		return nil, registryError("open", hkey, path, name, err)
	}
	defer key.Close()
	data, err := readRegistryValue(func(buffer []byte) (int, uint32, error) {
		return key.GetValue(name, buffer)
	}, registry.ErrShortBuffer)
	if err != nil {
		return nil, registryError("read", hkey, path, name, err)
	}
	return legacyUTF16(data), nil
}

// GetDwordReg reads a REG_DWORD as uint32. Other registry types return
// registry.ErrUnexpectedType, which can be inspected using errors.Is.
func GetDwordReg(hkey, path, name string) (uint32, error) {
	root, err := registryRoot(hkey, path, name)
	if err != nil {
		return 0, err
	}
	key, err := registry.OpenKey(root, path, registry.QUERY_VALUE)
	if err != nil {
		return 0, registryError("open", hkey, path, name, err)
	}
	defer key.Close()
	value, typ, err := key.GetIntegerValue(name)
	if err == nil && typ != registry.DWORD {
		err = registry.ErrUnexpectedType
	}
	if err != nil {
		return 0, registryError("read DWORD", hkey, path, name, err)
	}
	return uint32(value), nil
}

// WriteStringReg creates path if necessary and writes a REG_SZ value. An empty
// name selects the default value. Permissions depend on the key's Windows ACL.
func WriteStringReg(hkey, path, name, val string) error {
	root, err := registryRoot(hkey, path, name)
	if err != nil {
		return err
	}
	if strings.ContainsRune(val, 0) {
		return fmt.Errorf("registry string contains NUL: %w", syscall.EINVAL)
	}
	key, _, err := registry.CreateKey(root, path, registry.SET_VALUE)
	if err != nil {
		return registryError("create", hkey, path, name, err)
	}
	defer key.Close()
	if err := key.SetStringValue(name, val); err != nil {
		return registryError("write string", hkey, path, name, err)
	}
	return nil
}

// WriteDwordReg creates path if necessary and writes a REG_DWORD value.
func WriteDwordReg(hkey, path, name string, val uint32) error {
	root, err := registryRoot(hkey, path, name)
	if err != nil {
		return err
	}
	key, _, err := registry.CreateKey(root, path, registry.SET_VALUE)
	if err != nil {
		return registryError("create", hkey, path, name, err)
	}
	defer key.Close()
	if err := key.SetDWordValue(name, val); err != nil {
		return registryError("write DWORD", hkey, path, name, err)
	}
	return nil
}

// DeleteKey deletes the subkey name beneath path, including its values. It does
// not delete a value or recursively delete child keys. Registry view follows the
// process architecture, as for reads and writes.
func DeleteKey(hkey, path, name string) error {
	root, err := registryRoot(hkey, path, name)
	if err != nil {
		return err
	}
	key, err := registry.OpenKey(root, path, registry.ENUMERATE_SUB_KEYS)
	if err != nil {
		return registryError("open", hkey, path, name, err)
	}
	defer key.Close()
	if err := registry.DeleteKey(key, name); err != nil {
		return registryError("delete key", hkey, path, name, err)
	}
	return nil
}
