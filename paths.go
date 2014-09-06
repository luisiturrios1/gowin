// gowin - The windows OS interface 
//
// Copyright 2014-2020 Luis Iturrios. All rights reserved.
// Use of this source code is governed by a GPL-LICENCE
// license that can be found in the LICENSE file.
package gowin

import (

)

const (
	ALL=true//Use for make ShellFolders.Context
	USER=false//Use for make ShellFolders.Context
)

type ShellFolders struct {
	Context bool 
}

//TODO : make code for it
// Its don't use the Context value
func(s *ShellFolders)programFiles()(val string){
	return
}

//TODO : make code for it
func(s *ShellFolders)appData()(val string){
	return
}

//TODO : make code for it
func(s *ShellFolders)desktop()(val string){
	return
}

//TODO : make code for it
func(s *ShellFolders)documents()(val string){
	return
}

//TODO : make code for it
func(s *ShellFolders)startMenu()(val string){
	return
}

//TODO : make code for it
func(s *ShellFolders)startMenuPrograms()(val string){
	return
}