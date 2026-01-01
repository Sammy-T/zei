package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/sammy-t/zei"
	"github.com/urfave/cli/v3"
)

func main() {
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

	return zei.ExecSnippet(id)
}

func listSnippets(_ context.Context, _ *cli.Command) error {
	snippets, err := zei.GetSnippets()
	if err != nil {
		return err
	}

	for _, snippet := range snippets {
		fmt.Printf("%v\n\n", snippet.DisplayText())
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

	return zei.AddSnippet(id, cmdText, description)
}

func removeSnippet(_ context.Context, c *cli.Command) error {
	if c.Args().Len() == 0 {
		return fmt.Errorf("invalid snippet id args")
	}

	ids := c.Args().Slice()

	return zei.RemoveSnippet(ids)
}
