package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"port-scanner/internal/models"
)

// PrintResultAsJSON exports the function for use in other packages
func PrintResultAsJSON(result models.ScanResult) {
    jsonData, err := json.MarshalIndent(result, "", "  ")
    if err != nil {
        log.Fatalf("Error marshaling JSON: %v", err)
    }
    fmt.Println(string(jsonData))

    if result.Error != "" {
        os.Exit(1) // Exit with an error code if there was a DNS resolution error
    }
}

// ClearConsole exports the functions for use in other packages
func ClearConsole() {
	var clearCmd *exec.Cmd

	switch runtime.GOOS {
	case "linux", "darwin":
		clearCmd = exec.Command("clear")
	case "windows":
		clearCmd = exec.Command("cmd", "/c", "cls")
	default:
		fmt.Println("Unsupported platform")
		return
	}

	clearCmd.Stdout = os.Stdout
	clearCmd.Run()
}

// GetProjectRoot returns the root directory of the project containing a .env file.
func GetProjectRoot() (string, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", errors.New("could not get current file info")
	}

	// Get the directory of the current file
	dir := filepath.Dir(filename)

	// Traverse up the directory tree to find the project root
	for {
		if _, err := os.Stat(filepath.Join(dir, ".env")); err == nil {
			return dir, nil
		}

		// Move up one directory
		parent := filepath.Dir(dir)
		if parent == dir {
			break // reached the root of the filesystem
		}
		dir = parent
	}

	return "", errors.New("could not find project root containing .env")
}
