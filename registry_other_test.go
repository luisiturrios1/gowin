//go:build !windows

package gowin

import (
	"errors"
	"testing"
)

func TestUnsupportedPlatform(t *testing.T) {
	_, rawErr := GetRegRaw("HKCU", "", "test")
	_, stringErr := GetReg("HKCU", "", "test")
	_, dwordErr := GetDwordReg("HKCU", "", "test")
	_, folderErr := (&ShellFolders{}).Lookup("AppData")
	for _, err := range []error{rawErr, stringErr, dwordErr, folderErr, WriteStringReg("HKCU", "", "test", "x"), WriteDwordReg("HKCU", "", "test", 1), DeleteKey("HKCU", "", "test")} {
		if !errors.Is(err, ErrUnsupportedPlatform) {
			t.Fatalf("got %v", err)
		}
	}
	folders := &ShellFolders{}
	for _, lookup := range []func() string{folders.ProgramFiles, folders.AppData, folders.Desktop, folders.Documents, folders.StartMenu, folders.StartMenuPrograms} {
		if got := lookup(); got != "" {
			t.Fatalf("got %q", got)
		}
	}
	var nilFolders *ShellFolders
	if nilFolders.ProgramFiles() != "" {
		t.Fatal("ProgramFiles ignores receiver Context")
	}
}
