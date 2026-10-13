// Run on Windows with: go run ./examples/basic
package main

import (
	"fmt"
	"log"

	"github.com/luisiturrios/gowin"
)

func main() {
	// Read-only: this example does not modify the registry.
	folders := gowin.ShellFolders{Context: gowin.USER}
	path, err := folders.Lookup("Documents")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Documents (registry value):", path)
	fmt.Println("USERPROFILE:", gowin.GetEnvVar("USERPROFILE"))
}
