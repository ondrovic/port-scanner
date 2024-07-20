package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
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