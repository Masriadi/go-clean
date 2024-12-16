package main

import (
	"fmt"
	"os"

	"github.com/Masriadi/go-clean/internal/generator"
	"github.com/Masriadi/go-clean/internal/remover"
	"github.com/Masriadi/go-clean/internal/utils"
	"golang.org/x/mod/modfile"
)

// main function to execute commands
func main() {
	if len(os.Args) < 3 {
		utils.LogInfo("Usage: go-clean <command> <name>")
		return
	}

	command := os.Args[1]
	name := os.Args[2]

	// Get module name from go.mod
	moduleName, err := getModuleName()
	if err != nil {
		utils.LogError(fmt.Sprintf("Error reading module name: %v", err))
		return
	}
	utils.LogInfo(fmt.Sprintf("Module Name: %s", moduleName))

	switch command {
	case "gen":
		if err := generator.GenerateStructure(name, moduleName); err != nil {
			utils.LogError(fmt.Sprintf("Error generating structure: %v", err))
		} else {
			utils.LogSuccess(fmt.Sprintf("Structure for '%s' generated successfully.", name))
		}
	case "remove":
		if err := remover.RemoveStructure(name); err != nil {
			utils.LogError(fmt.Sprintf("Error removing structure: %v", err))
		} else {
			utils.LogSuccess(fmt.Sprintf("Structure for '%s' removed successfully.", name))
		}
	default:
		utils.LogError("Unknown command. Use 'gen' or 'remove'.")
	}
}

// getModuleName reads the go.mod file and returns the module name
func getModuleName() (string, error) {
	data, err := os.ReadFile("go.mod") // Use os.ReadFile instead of ioutil.ReadFile
	if err != nil {
		return "", err
	}

	modFile, err := modfile.Parse("go.mod", data, nil)
	if err != nil {
		return "", err
	}

	return modFile.Module.Mod.Path, nil
}
