// gowin - The windows OS interface 
//
// Copyright 2014-2020 Luis Iturrios. All rights reserved.
// Use of this source code is governed by a GPL-LICENCE
// license that can be found in the LICENSE file.
package gowin


import (
	"errors"
	"syscall"
	"unsafe"
	"code.google.com/p/winsvc/winapi"
)

// TODO: Solve error in query DWORD registry
// Use to read value from windows registry the HKEY in the next definition
// HKLM = HKEY_LOCAL_MACHINE
// HKCC = HKEY_CURRENT_CONFIG 
// HKCR = HKEY_CLASSES_ROOT 
// HKCU = HKEY_CURRENT_USER 
// HKU = HKEY_USERS
func GetReg(hkey, path, name string)(val string, err error){
	var handle syscall.Handle
	switch hkey{
		case "HKLM":			
			err = syscall.RegOpenKeyEx(syscall.HKEY_LOCAL_MACHINE, syscall.StringToUTF16Ptr(path), 0, syscall.KEY_READ, &handle)
		case "HKCC":
			err = syscall.RegOpenKeyEx(syscall.HKEY_CURRENT_CONFIG, syscall.StringToUTF16Ptr(path), 0, syscall.KEY_READ, &handle)
		case "HKCR":
			err = syscall.RegOpenKeyEx(syscall.HKEY_CLASSES_ROOT, syscall.StringToUTF16Ptr(path), 0, syscall.KEY_READ, &handle)
		case "HKCU":
			err = syscall.RegOpenKeyEx(syscall.HKEY_CURRENT_USER, syscall.StringToUTF16Ptr(path), 0, syscall.KEY_READ, &handle)
		case "HKU":
			err = syscall.RegOpenKeyEx(syscall.HKEY_USERS, syscall.StringToUTF16Ptr(path), 0, syscall.KEY_READ, &handle)
		default:
			err = errors.New("Unknown HKEY: " + hkey)
			return
	}
	if err != nil {
		return
	}
	defer syscall.RegCloseKey(handle)
	var typ uint32
	var buffer [syscall.MAX_LONG_PATH]uint16
	n := uint32(len(buffer))
	err = syscall.RegQueryValueEx(handle, syscall.StringToUTF16Ptr(name), nil, &typ, (*byte)(unsafe.Pointer(&buffer[0])), &n)
	if err != nil {
		return
	}
	val = syscall.UTF16ToString(buffer[:])
	return
}

// Use to write string value to windows registry the HKEY in the next definition
// HKLM = HKEY_LOCAL_MACHINE
// HKCC = HKEY_CURRENT_CONFIG 
// HKCR = HKEY_CLASSES_ROOT 
// HKCU = HKEY_CURRENT_USER 
// HKCU = HKEY_USERS HKU
func WriteStringReg(hkey, path, name, val string)(err error){
	var handle syscall.Handle
	switch hkey{
	case "HKLM":
		err = syscall.RegOpenKeyEx(syscall.HKEY_LOCAL_MACHINE, syscall.StringToUTF16Ptr(""), 0, syscall.KEY_CREATE_SUB_KEY, &handle)
	case "HKCC":
		err = syscall.RegOpenKeyEx(syscall.HKEY_CURRENT_CONFIG, syscall.StringToUTF16Ptr(""), 0, syscall.KEY_CREATE_SUB_KEY, &handle)
	case "HKCR":
		err = syscall.RegOpenKeyEx(syscall.HKEY_CLASSES_ROOT, syscall.StringToUTF16Ptr(""), 0, syscall.KEY_CREATE_SUB_KEY, &handle)
	case "HKCU":
		err = syscall.RegOpenKeyEx(syscall.HKEY_CURRENT_USER, syscall.StringToUTF16Ptr(""), 0, syscall.KEY_CREATE_SUB_KEY, &handle)
	case "HKU":
		err = syscall.RegOpenKeyEx(syscall.HKEY_USERS, syscall.StringToUTF16Ptr(""), 0, syscall.KEY_CREATE_SUB_KEY, &handle)
	default:
		err = errors.New("Unknown HKEY: " + hkey)
		return
	}
	if err != nil {
		return
	}
	defer syscall.RegCloseKey(handle)
	var d uint32
	err = winapi.RegCreateKeyEx(handle, syscall.StringToUTF16Ptr(path), 0, nil, winapi.REG_OPTION_NON_VOLATILE, syscall.KEY_ALL_ACCESS, nil, &handle, &d)
	if err != nil {
		return
	}
	buf := syscall.StringToUTF16(val)
	err = winapi.RegSetValueEx(handle, syscall.StringToUTF16Ptr(name), 0, syscall.REG_SZ, (*byte)(unsafe.Pointer(&buf[0])), uint32(len(buf)*2))
	return
}

// Use to write uint32 value to windows registry the HKEY in the next definition
// HKLM = HKEY_LOCAL_MACHINE
// HKCC = HKEY_CURRENT_CONFIG
// HKCR = HKEY_CLASSES_ROOT
// HKCU = HKEY_CURRENT_USER
// HKCU = HKEY_USERS HKU
func WriteDwordReg(hkey, path, name string, val uint32)(err error){
	var handle syscall.Handle
	switch hkey{
	case "HKLM":
		err = syscall.RegOpenKeyEx(syscall.HKEY_LOCAL_MACHINE, syscall.StringToUTF16Ptr(""), 0, syscall.KEY_CREATE_SUB_KEY, &handle)
	case "HKCC":
		err = syscall.RegOpenKeyEx(syscall.HKEY_CURRENT_CONFIG, syscall.StringToUTF16Ptr(""), 0, syscall.KEY_CREATE_SUB_KEY, &handle)
	case "HKCR":
		err = syscall.RegOpenKeyEx(syscall.HKEY_CLASSES_ROOT, syscall.StringToUTF16Ptr(""), 0, syscall.KEY_CREATE_SUB_KEY, &handle)
	case "HKCU":
		err = syscall.RegOpenKeyEx(syscall.HKEY_CURRENT_USER, syscall.StringToUTF16Ptr(""), 0, syscall.KEY_CREATE_SUB_KEY, &handle)
	case "HKU":
		err = syscall.RegOpenKeyEx(syscall.HKEY_USERS, syscall.StringToUTF16Ptr(""), 0, syscall.KEY_CREATE_SUB_KEY, &handle)
	default:
		err = errors.New("Unknown HKEY: " + hkey)
		return
	}
	if err != nil {
		return
	}
	defer syscall.RegCloseKey(handle)
	var d uint32
	err = winapi.RegCreateKeyEx(handle, syscall.StringToUTF16Ptr(path), 0, nil, winapi.REG_OPTION_NON_VOLATILE, syscall.KEY_ALL_ACCESS, nil, &handle, &d)
	if err != nil {
		return
	}
	err = winapi.RegSetValueEx(handle, syscall.StringToUTF16Ptr(name), 0, syscall.REG_DWORD, (*byte)(unsafe.Pointer(&val)), uint32(unsafe.Sizeof(val)))
	return
}

// Use to remove key from windows registry the HKEY in the next definition
// HKLM = HKEY_LOCAL_MACHINE
// HKCC = HKEY_CURRENT_CONFIG
// HKCR = HKEY_CLASSES_ROOT
// HKCU = HKEY_CURRENT_USER
// HKCU = HKEY_USERS HKU
func DeleteKey(hkey, path, name string)(err error){
	var handle syscall.Handle
	switch hkey{
		case "HKLM":
			err = syscall.RegOpenKeyEx(syscall.HKEY_LOCAL_MACHINE, syscall.StringToUTF16Ptr(path), 0, syscall.KEY_READ, &handle)
		case "HKCC":
			err = syscall.RegOpenKeyEx(syscall.HKEY_CURRENT_CONFIG, syscall.StringToUTF16Ptr(path), 0, syscall.KEY_READ, &handle)
		case "HKCR":
			err = syscall.RegOpenKeyEx(syscall.HKEY_CLASSES_ROOT, syscall.StringToUTF16Ptr(path), 0, syscall.KEY_READ, &handle)
		case "HKCU":
			err = syscall.RegOpenKeyEx(syscall.HKEY_CURRENT_USER, syscall.StringToUTF16Ptr(path), 0, syscall.KEY_READ, &handle)
		case "HKU":
			err = syscall.RegOpenKeyEx(syscall.HKEY_USERS, syscall.StringToUTF16Ptr(path), 0, syscall.KEY_READ, &handle)
		default:
			err = errors.New("Unknown HKEY: " + hkey)
			return
	}
	if err != nil {
		return
	}
	defer syscall.RegCloseKey(handle)
	err = winapi.RegDeleteKey(handle, syscall.StringToUTF16Ptr(name))
	return
}
