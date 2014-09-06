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
)

// Use to read value from windows registry indique the HKEY in the next definition
// Example: val, err := gowin.GetRegKey("HKLM", `Software\Microsoft\Windows\CurrentVersion\Explorer\Shell Folders`, "Common AppData")
// HKLM = HKEY_LOCAL_MACHINE 
// HKCC = HKEY_CURRENT_CONFIG 
// HKCR = HKEY_CLASSES_ROOT 
// HKCU = HKEY_CURRENT_USER 
// HKU = HKEY_USERS
func GetRegKey(hkey, path, name string)(val string, err error){
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
	var buffer [256]uint16
	n := uint32(len(buffer))
	err = syscall.RegQueryValueEx(handle, syscall.StringToUTF16Ptr(name), nil, &typ, (*byte)(unsafe.Pointer(&buffer[0])), &n)
	if err != nil {
		return
	}
	val = syscall.UTF16ToString(buffer[:])
	return
}

//TODO : make code for it 
// HKLM = HKEY_LOCAL_MACHINE 
// HKCC = HKEY_CURRENT_CONFIG 
// HKCR = HKEY_CLASSES_ROOT 
// HKCU = HKEY_CURRENT_USER 
// HKCU = HKEY_USERS HKU
func WriteRegKey(hkey, path, name string)(err error){
	return
}