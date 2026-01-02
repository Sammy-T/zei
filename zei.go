package zei

import (
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	appDirPath := filepath.Join(home, ".zei")
	dbPath := filepath.Join(appDirPath, "zei.db")

	if err = os.Mkdir(appDirPath, 0750); err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}

	db, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&Snippet{})
}

// IsValidId is a helper to determine if the snippet id is valid.
func IsValidId(id string) bool {
	idRe := regexp.MustCompile(`^[\w\d\-]+$`)
	return idRe.MatchString(id)
}

// GetSnippet retrieves the snippet with the provided id.
func GetSnippet(id string) (Snippet, error) {
	var snippet Snippet

	result := db.First(&snippet, "id = ?", id)
	return snippet, result.Error
}

// GetSnippets returns all stored snippets.
func GetSnippets() ([]Snippet, error) {
	var snippets []Snippet

	result := db.Order("id asc").Find(&snippets)
	return snippets, result.Error
}

// AddSnippet stores a new snippet in the database.
func AddSnippet(id string, cmdText string, description string) error {
	result := db.Create(&Snippet{ID: id, Command: cmdText, Description: description})
	return result.Error
}

// Update stores the updated snippet in the database.
func UpdateSnippet(id string, updated Snippet) error {
	result := db.Model(&Snippet{ID: id}).Updates(updated)
	return result.Error
}

// RemoveSnippet removes the matching snippet from the database.
func RemoveSnippet(ids []string) error {
	result := db.Delete(&Snippet{}, ids)
	return result.Error
}
