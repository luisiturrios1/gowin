package gowin

import "testing"

func TestShellFolderRegistry(t *testing.T) {
	for _, tc := range []struct{ name, user, common string }{
		{"AppData", "AppData", "Common AppData"}, {"Desktop", "Desktop", "Common Desktop"},
		{"Documents", "Personal", "Common Documents"}, {"StartMenu", "Start Menu", "Common Start Menu"},
		{"StartMenuPrograms", "Programs", "Common Programs"},
	} {
		for _, context := range []bool{USER, ALL} {
			label := "USER"
			if context == ALL {
				label = "ALL"
			}
			t.Run(tc.name+"/"+label, func(t *testing.T) {
				hkey, path, value, err := shellFolderRegistry(context, tc.name)
				wantKey, wantPath, wantValue := "HKCU", `Software\Microsoft\Windows\CurrentVersion\Explorer\User Shell Folders`, tc.user
				if context == ALL {
					wantKey, wantPath, wantValue = "HKLM", `Software\Microsoft\Windows\CurrentVersion\Explorer\Shell Folders`, tc.common
				}
				if err != nil || hkey != wantKey || path != wantPath || value != wantValue {
					t.Fatalf("got %q %q %q %v", hkey, path, value, err)
				}
			})
		}
	}
	for _, context := range []bool{USER, ALL} {
		hkey, path, value, err := shellFolderRegistry(context, "ProgramFiles")
		if err != nil || hkey != "HKLM" || path != `SOFTWARE\Microsoft\Windows\CurrentVersion` || value != "ProgramFilesDir" {
			t.Fatalf("ProgramFiles: %q %q %q %v", hkey, path, value, err)
		}
	}
	for _, name := range []string{"", "appdata", "unknown"} {
		if _, _, _, err := shellFolderRegistry(USER, name); err == nil {
			t.Fatalf("expected error for %q", name)
		}
		if _, err := (&ShellFolders{}).Lookup(name); err == nil {
			t.Fatalf("Lookup must reject %q", name)
		}
	}
}
