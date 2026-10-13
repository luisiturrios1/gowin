// Modified 2026-10-12: shared folder lookup and additive error-returning API.

package gowin

import "fmt"

const (
	// ALL selects shell folders shared by all users.
	ALL = true
	// USER selects the current user's shell folders.
	USER = false
)

// ShellFolders reads the same registry locations as legacy Gowin. Its zero value
// selects the current user. Do not modify Context concurrently with a lookup.
type ShellFolders struct {
	Context bool
}

// Lookup returns a folder's registry path and any read error. Supported names
// are ProgramFiles, AppData, Desktop, Documents, StartMenu and StartMenuPrograms.
// REG_EXPAND_SZ values retain their environment-variable placeholders.
func (s *ShellFolders) Lookup(name string) (string, error) {
	context := USER
	if name != "ProgramFiles" {
		context = s.Context
	}
	hkey, path, value, err := shellFolderRegistry(context, name)
	if err != nil {
		return "", err
	}
	return GetReg(hkey, path, value)
}

func shellFolderRegistry(context bool, name string) (hkey, path, value string, err error) {
	if name == "ProgramFiles" {
		return "HKLM", `SOFTWARE\Microsoft\Windows\CurrentVersion`, "ProgramFilesDir", nil
	}
	var common string
	switch name {
	case "AppData":
		value, common = "AppData", "Common AppData"
	case "Desktop":
		value, common = "Desktop", "Common Desktop"
	case "Documents":
		value, common = "Personal", "Common Documents"
	case "StartMenu":
		value, common = "Start Menu", "Common Start Menu"
	case "StartMenuPrograms":
		value, common = "Programs", "Common Programs"
	default:
		return "", "", "", fmt.Errorf("unknown shell folder: %q", name)
	}
	if context == ALL {
		return "HKLM", `Software\Microsoft\Windows\CurrentVersion\Explorer\Shell Folders`, common, nil
	}
	return "HKCU", `Software\Microsoft\Windows\CurrentVersion\Explorer\User Shell Folders`, value, nil
}

// ProgramFiles returns the ProgramFiles path, ignoring Context and read errors.
func (s *ShellFolders) ProgramFiles() string { val, _ := s.Lookup("ProgramFiles"); return val }

// AppData returns the selected AppData path, ignoring read errors.
func (s *ShellFolders) AppData() string { val, _ := s.Lookup("AppData"); return val }

// Desktop returns the selected Desktop path, ignoring read errors.
func (s *ShellFolders) Desktop() string { val, _ := s.Lookup("Desktop"); return val }

// Documents returns the selected Documents path, ignoring read errors.
func (s *ShellFolders) Documents() string { val, _ := s.Lookup("Documents"); return val }

// StartMenu returns the selected Start Menu path, ignoring read errors.
func (s *ShellFolders) StartMenu() string { val, _ := s.Lookup("StartMenu"); return val }

// StartMenuPrograms returns the selected Programs path, ignoring read errors.
func (s *ShellFolders) StartMenuPrograms() string {
	val, _ := s.Lookup("StartMenuPrograms")
	return val
}
