package zei

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os/exec"
	"time"

	"github.com/glebarez/sqlite"
	cmdstr "github.com/sammy-t/zei/internal/cmdStr"
	"gorm.io/gorm"
)

type Snippet struct {
	ID          string `gorm:"primaryKey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	Command     string         `gorm:"not null"`
	Description string
}

var db *gorm.DB

func init() {
	var err error

	db, err = gorm.Open(sqlite.Open("dev.db"), &gorm.Config{}) //// TODO: proper path
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&Snippet{})
}

// ExecSnippet attempts to execute the matching snippet.
func ExecSnippet(id string) error {
	var snippet Snippet

	if result := db.First(&snippet, "id = ?", id); result.Error != nil {
		return result.Error
	}

	fmt.Println(snippet) //// TODO: Format

	cmdArgs := cmdstr.Split(snippet.Command)

	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)

	outPipe, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	go func(pipe io.ReadCloser) {
		reader := bufio.NewReader(pipe)

		var line string
		var err error

		for err == nil {
			fmt.Print(line)
			line, err = reader.ReadString('\n')
		}
	}(outPipe)

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
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
