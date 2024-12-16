package remover

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Masriadi/go-clean/internal/utils"
)

// RemoveStructure removes the files for the given entity name
func RemoveStructure(entityName string) error {
	// Normalisasi nama entity (lowercase)
	entityName = strings.ToLower(entityName)

	// Daftar file yang akan dihapus
	files := []string{
		filepath.Join("data/repositories", fmt.Sprintf("%s_repository.go", entityName)),
		filepath.Join("domain/entities", fmt.Sprintf("%s.go", entityName)),
		filepath.Join("domain/usecases", fmt.Sprintf("%s_usecase.go", entityName)),
		filepath.Join("presentation/handlers/http/api/v1", fmt.Sprintf("%s_handler.go", entityName)),
		filepath.Join("presentation/routes", fmt.Sprintf("%s_routes.go", entityName)),
	}

	// Tanyakan konfirmasi sebelum menghapus
	if !confirmRemoval(files) {
		utils.LogInfo("Removal canceled.")
		return nil
	}

	// Hapus file
	for _, file := range files {
		if err := removeFile(file); err != nil {
			return err
		}
	}

	utils.LogSuccess(fmt.Sprintf("Files for '%s' removed successfully.", entityName))
	return nil
}

// confirmRemoval asks the user for confirmation before removing files
func confirmRemoval(files []string) bool {
	utils.LogInfo("The following files will be removed:")
	for _, file := range files {
		utils.LogInfo(file)
	}
	fmt.Print("Are you sure you want to proceed? (y/n): ")

	reader := bufio.NewReader(os.Stdin)
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(response)

	return strings.ToLower(response) == "y"
}

// removeFile removes the specified file
func removeFile(file string) error {
	// Cek apakah file ada
	if _, err := os.Stat(file); os.IsNotExist(err) {
		utils.LogInfo(fmt.Sprintf("File %s does not exist, skipping.", file))
		return nil
	}

	// Hapus file
	if err := os.Remove(file); err != nil {
		return fmt.Errorf("error removing file %s: %w", file, err)
	}

	return nil
}
