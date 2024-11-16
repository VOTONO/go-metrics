// Package main is a custom static analysis tool using multichecker and analysis passes.
//
// This package loads a configuration file (config.json) to enable specific static checks
// provided by staticcheck and adds custom analyzers to enforce rules across Go code.
package main

import (
	"encoding/json"
	"os"
	"path/filepath"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/staticcheck"

	"github.com/VOTONO/go-metrics/cmd/staticlint/custom_analyzer"
)

// Config defines the name of the configuration file that lists the static analyzers to enable.
const Config = `config.json`

// ConfigData represents the structure of the configuration file.
//
// The config file should be a JSON file with a list of staticcheck analyzers
// to enable. Example:
//
//	{
//	    "Staticcheck": ["SA5000", "SA9004"]
//	}
type ConfigData struct {
	Staticcheck []string `json:"Staticcheck"`
}

// main is the entry point of the program. It reads the configuration file,
// selects specific static analyzers, and runs them as part of the multichecker.
//
// The function performs the following steps:
//  1. Determines the path of the configuration file based on the executable's location.
//  2. Reads and parses the configuration file into ConfigData.
//  3. Adds predefined analyzers and custom analyzers into the multichecker.
//  4. Initializes the multichecker with the selected analyzers.
//
// This setup allows fine-grained control over which staticcheck analyzers to run,
// while also adding custom checks.
func main() {
	appfile, err := os.Executable()
	if err != nil {
		panic(err) // Panic if executable path cannot be determined.
	}

	// Read configuration file.
	data, err := os.ReadFile(filepath.Join(filepath.Dir(appfile), Config))
	if err != nil {
		panic(err) // Panic if config file cannot be read.
	}

	// Parse configuration data into ConfigData struct.
	var cfg ConfigData
	if err = json.Unmarshal(data, &cfg); err != nil {
		panic(err) // Panic if JSON parsing fails.
	}

	// Initialize list of analyzers with built-in analyzers.
	mychecks := []*analysis.Analyzer{
		custom_analyzer.ErrCheckAnalyzer, // Custom analyzer to check for os.Exit in main.
		printf.Analyzer,                  // Checks for correct printf-style function calls.
		shadow.Analyzer,                  // Detects variable shadowing.
		structtag.Analyzer,               // Verifies struct field tags.
	}

	// Create a map to hold enabled staticcheck analyzers from config.
	checks := make(map[string]bool)
	for _, v := range cfg.Staticcheck {
		checks[v] = true
	}

	// Add staticcheck analyzers listed in config.
	for _, v := range staticcheck.Analyzers {
		if checks[v.Analyzer.Name] {
			mychecks = append(mychecks, v.Analyzer)
		}
	}

	// Run multichecker with the selected analyzers.
	multichecker.Main(
		mychecks...,
	)
}
