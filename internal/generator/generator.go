package generator

import (
	"bufio"
	_ "embed"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/Masriadi/go-clean/internal/utils"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Embed template files
//
//go:embed templates/repository.tmpl
var repositoryTemplate string

//go:embed templates/entity.tmpl
var entityTemplate string

//go:embed templates/usecase.tmpl
var usecaseTemplate string

//go:embed templates/handler.tmpl
var handlerTemplate string

//go:embed templates/routes.tmpl
var routesTemplate string

// TemplateInfo holds the template and its output path
type TemplateInfo struct {
	OutputPath string
	Content    string
}

// Define directory constants
const (
	repositoriesDir = "data/repositories"
	entitiesDir     = "domain/entities"
	usecasesDir     = "domain/usecases"
	handlersDir     = "presentation/handlers/http/api/v1"
	routesDir       = "presentation/routes"
)

// GenerateStructure generates the directory structure and files for the given entity name
func GenerateStructure(entityName string, moduleName string) error {
	if entityName == "" || moduleName == "" {
		return fmt.Errorf("entityName and moduleName cannot be empty")
	}

	// Convert entity name to title case
	entityName = cases.Title(language.English).String(entityName)

	// Show log message
	utils.LogInfo(fmt.Sprintf("Generating structure for entity: %s", entityName))

	// Convert entity name to lowercase
	lowerEntityName := strings.ToLower(entityName)

	directories := []string{
		repositoriesDir,
		entitiesDir,
		usecasesDir,
		handlersDir,
		routesDir,
	}

	// Create directories
	if err := createDirectories(directories); err != nil {
		return err
	}

	// Define templates
	templates := []TemplateInfo{
		{fmt.Sprintf("%s/%s_repository.go", repositoriesDir, lowerEntityName), repositoryTemplate},
		{fmt.Sprintf("%s/%s.go", entitiesDir, lowerEntityName), entityTemplate},
		{fmt.Sprintf("%s/%s_usecase.go", usecasesDir, lowerEntityName), usecaseTemplate},
		{fmt.Sprintf("%s/%s_handler.go", handlersDir, lowerEntityName), handlerTemplate},
		{fmt.Sprintf("%s/%s_routes.go", routesDir, lowerEntityName), routesTemplate},
	}

	// Generate files from templates
	for _, tmplInfo := range templates {
		if err := generateFile(
			tmplInfo.OutputPath,
			tmplInfo.Content,
			lowerEntityName,
			entityName,
			moduleName,
		); err != nil {
			return err
		}
	}

	utils.LogSuccess(fmt.Sprintf("Structure for '%s' generated successfully.", lowerEntityName))
	return nil
}

// createDirectories creates the specified directories
func createDirectories(directories []string) error {
	for _, dir := range directories {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return fmt.Errorf("error creating directory %s: %w", dir, err)
		}
	}
	return nil
}

// generateFile generates a file from the given template content and entity name
func generateFile(
	outputPathTemplate,
	templateContent,
	lowerEntityName string,
	entityName string,
	moduleName string,
) error {
	funcMap := template.FuncMap{
		"lower": strings.ToLower,
		"upper": strings.ToUpper,
		"title": cases.Title(language.English).String,
	}

	// Process output path template
	fileNameTmpl, err := template.New("filename").Funcs(funcMap).Parse(outputPathTemplate)
	if err != nil {
		return fmt.Errorf("error parsing filename template: %w", err)
	}

	var fileNameBuf strings.Builder
	if err := fileNameTmpl.Execute(&fileNameBuf, map[string]string{"Entity": lowerEntityName}); err != nil {
		return fmt.Errorf("error generating filename: %w", err)
	}
	outPath := fileNameBuf.String()

	// Check if the file already exists
	if _, err := os.Stat(outPath); err == nil {
		// File exists, ask user for action
		if !askOverwrite(outPath) {
			fmt.Printf("Skipping file: %s\n", outPath)
			return nil
		}
	}

	// Parse template content
	tmpl, err := template.New("template").Funcs(funcMap).Parse(templateContent)
	if err != nil {
		return fmt.Errorf("error parsing template: %w", err)
	}

	// Create file
	outFile, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("error creating file %s: %w", outPath, err)
	}
	defer outFile.Close()

	// Execute template
	data := map[string]string{
		"LowerEntity": lowerEntityName,
		"Entity":      entityName,
		"Module":      moduleName,
	}
	if err := tmpl.Execute(outFile, data); err != nil {
		return fmt.Errorf("error writing to file %s: %w", outPath, err)
	}

	return nil
}

// askOverwrite prompts the user to decide whether to overwrite an existing file
func askOverwrite(filePath string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("File %s already exists. Do you want to overwrite it? (y/n): ", filePath)
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(response)
	return strings.ToLower(response) == "y"
}
