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

//go:embed templates/di.tmpl
var diTemplate string

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
	diDir           = "di"
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

	// Show log message
	utils.LogInfo(fmt.Sprintf("Generating structure for entity: %s", entityName))

	// Convert entity name to title case
	entityName = utils.StringToEntityName(entityName)

	// Convert entity name to lowercase
	instanceEntityName := utils.StringToInstanceName(entityName)

	// Convert entity name to directory name
	fileName := utils.StringToFileName(entityName)

	// Convert entity name to directory name (kebab-case)
	dirName := utils.StringToDirName(entityName)

	// log.Printf("entityName : %s\n", entityName)
	// log.Printf("instanceEntityName : %s\n", instanceEntityName)
	// log.Printf("fileName : %s\n", fileName)

	directories := []string{
		repositoriesDir,
		diDir,
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
		{fmt.Sprintf("%s/%s_repository.go", repositoriesDir, fileName), repositoryTemplate},
		{fmt.Sprintf("%s/%s_di.go", diDir, fileName), diTemplate},
		{fmt.Sprintf("%s/%s.go", entitiesDir, fileName), entityTemplate},
		{fmt.Sprintf("%s/%s_usecase.go", usecasesDir, fileName), usecaseTemplate},
		{fmt.Sprintf("%s/%s_handler.go", handlersDir, fileName), handlerTemplate},
		{fmt.Sprintf("%s/%s_routes.go", routesDir, fileName), routesTemplate},
	}

	// Generate files from templates
	for _, tmplInfo := range templates {
		if err := generateFile(
			tmplInfo.OutputPath,
			tmplInfo.Content,
			instanceEntityName,
			entityName,
			moduleName,
			dirName,
		); err != nil {
			return err
		}
	}

	utils.LogSuccess(fmt.Sprintf("Structure for '%s' generated successfully.", instanceEntityName))
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
	instanceEntityName string,
	entityName string,
	moduleName string,
	dirname string,
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
	if err := fileNameTmpl.Execute(&fileNameBuf, map[string]string{"Entity": instanceEntityName}); err != nil {
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
		"LowerEntity": instanceEntityName,
		"Entity":      entityName,
		"Module":      moduleName,
		"DirName": dirname,
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
