package main

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/urfave/cli/v3"
	_ "modernc.org/sqlite"
)

var db *sql.DB

func main() {
	var err error

	db, err = sql.Open("sqlite", "dev.db")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS snippets (
			id TEXT PRIMARY KEY,
			command TEXT NOT NULL,
			description TEXT
		)`)

	if err != nil {
		log.Fatal(err)
	}

	cmd := &cli.Command{
		Name:    "zei",
		Version: "v0.0.1",
		Usage:   "A command snippet cli",
		Action: func(ctx context.Context, c *cli.Command) error {
			log.Println("todo")
			return nil
		},
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
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func listSnippets(_ context.Context, _ *cli.Command) error {
	rows, err := db.Query(`SELECT * FROM snippets`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		var cmdText string
		var description string

		if err = rows.Scan(&id, &cmdText, &description); err != nil {
			return err
		}

		fmt.Printf("%v %v %v\n", id, cmdText, description)
	}

	return nil
}

func addSnippet(_ context.Context, _ *cli.Command) error {
	var id string
	var cmdText string
	var description string

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("id (ex. some-id): ")
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

	if _, err := db.Exec(`INSERT INTO snippets (id, command, description) VALUES (?, ?, ?);`, id, cmdText, description); err != nil {
		return err
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
