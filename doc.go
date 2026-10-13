// Package gowin provides simple access to the Windows registry, shell folder
// paths and process environment variables. Registry operations require Windows;
// on other platforms they return ErrUnsupportedPlatform. Environment operations
// are portable. Windows service management is not implemented.
package gowin
