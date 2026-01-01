package zei

import (
	"log"
	"regexp"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	var err error

	db, err = gorm.Open(sqlite.Open("dev.db"), &gorm.Config{}) //// TODO: proper path
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

	if result := db.First(&snippet, "id = ?", id); result.Error != nil {
		return snippet, result.Error
	}

	return snippet, nil
}

// GetSnippets returns all stored snippets.
func GetSnippets() ([]Snippet, error) {
	var snippets []Snippet

	if result := db.Find(&snippets); result.Error != nil {
		return nil, result.Error
	}

	return snippets, nil
}

// AddSnippet stores a new snippet in the database.
func AddSnippet(id string, cmdText string, description string) error {
	if result := db.Create(&Snippet{ID: id, Command: cmdText, Description: description}); result.Error != nil {
		return result.Error
	}

	return nil
}

// RemoveSnippet removes the matching snippet from the database.
func RemoveSnippet(ids []string) error {
	if result := db.Delete(&Snippet{}, ids); result.Error != nil {
		return result.Error
	}

	return nil
}
