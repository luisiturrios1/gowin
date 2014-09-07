// gowin - The windows OS interface 
//
// Copyright 2014-2020 Luis Iturrios. All rights reserved.
// Use of this source code is governed by a GPL-LICENCE
// license that can be found in the LICENSE file.
package gowin 

import (
	"os"
)

// Use to read value from windows environment variables by name
func GetEnvVar(name string)(val string){
	val = os.Getenv(name)
	return
}

// Use to write value on windows environment variable by name
func WriteEnvVar(name, val string)(err error){
	err = os.Setenv(name, val)
	return
}
