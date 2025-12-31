package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/urfave/cli/v3"
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

func main() {
	var err error

	db, err = gorm.Open(sqlite.Open("dev.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&Snippet{})

	cmd := &cli.Command{
		Name:        "zei",
		Version:     "v0.0.1",
		Description: "A command snippet cli",
		Usage:       "Execute snippet with ID",
		Action:      execSnippet,
		Commands: []*cli.Command{
			{
				Name: "snippet",
				Commands: []*cli.Command{
					{
						Name:    "list",
						Aliases: []string{"ls"},
						Usage:   "list snippets",
						Action:  listSnippets,
					},
					{
						Name:   "add",
						Usage:  "add a new snippet",
						Action: addSnippet,
					},
					{
						Name:      "remove",
						Aliases:   []string{"rm", "del"},
						Usage:     "remove snippet with ID",
						UsageText: "zei snippet remove ID",
						Action:    removeSnippet,
					},
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func execSnippet(_ context.Context, c *cli.Command) error {
	if c.Args().Len() != 1 {
		return fmt.Errorf("invalid snippet id args")
	}

	id := c.Args().First()

	var snippet Snippet

	if result := db.First(&snippet, "id = ?", id); result.Error != nil {
		return result.Error
	}

	fmt.Println(snippet) //// TODO: Format

	cmdArgs := SplitCmdStr(snippet.Command)

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

func listSnippets(_ context.Context, _ *cli.Command) error {
	var snippets []Snippet

	if result := db.Find(&snippets); result.Error != nil {
		return result.Error
	}

	for _, snippet := range snippets {
		fmt.Printf("%v\n", snippet) //// TODO: Format
	}

	return nil
}

func addSnippet(_ context.Context, _ *cli.Command) error {
	var id string
	var cmdText string
	var description string

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("id (ex. some-id): ") //// TODO: Validate
	scanner.Scan()
	id = scanner.Text()

	fmt.Print("command: ")
	scanner.Scan()
	cmdText = scanner.Text()

	fmt.Print("description: ")
	scanner.Scan()
	description = scanner.Text()

	confirmYesRe := regexp.MustCompile(`(?i)^y(es)?$|^$`)

	fmt.Printf("\nNew snippet\nid: %v\ncommand: %v\ndescription: %v\nSave? (Y/n): ", id, cmdText, description)
	scanner.Scan()

	if confirm := scanner.Text(); !confirmYesRe.MatchString(confirm) {
		return nil
	}

	if result := db.Create(&Snippet{ID: id, Command: cmdText, Description: description}); result.Error != nil {
		return result.Error
	}

	return nil
}

func removeSnippet(_ context.Context, c *cli.Command) error {
	if c.Args().Len() == 0 {
		return fmt.Errorf("invalid snippet id args")
	}

	ids := c.Args().Slice()

	if result := db.Delete(&Snippet{}, ids); result.Error != nil {
		return result.Error
	}

	return nil
}

func SplitCmdStr(cmdText string) []string {
	var parts []string

	var builder strings.Builder
	var quoting string

	var cmdLen = len(cmdText)

	quotesRe := regexp.MustCompile("\"|'|`")
	spaceRe := regexp.MustCompile(`\s`)

	for i, r := range cmdText {
		char := string(r)

		if quotesRe.MatchString(char) {
			switch quoting {
			case "": // Not already quoting, start.
				quoting = char
			case char: // Matches current quoting, end.
				quoting = ""
			}
		} else if quoting == "" && spaceRe.MatchString(char) {
			parts = append(parts, builder.String())
			builder.Reset()
			continue
		}

		builder.WriteRune(r)

		if i == cmdLen-1 {
			parts = append(parts, builder.String())
		}
	}

	return parts
}
