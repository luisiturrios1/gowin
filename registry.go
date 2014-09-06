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