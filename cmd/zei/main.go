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
	"sync"
	"text/template"

	"github.com/pterm/pterm"
	"github.com/sammy-t/zei"
	cmdstr "github.com/sammy-t/zei/internal/cmdStr"
	tmpl "github.com/sammy-t/zei/internal/template"
	"github.com/urfave/cli/v3"
)

// Matches 'y(es)' or empty strings.
var confirmDefYesRe *regexp.Regexp = regexp.MustCompile(`(?i)^y(es)?$|^$`)

// Matches `|` with 0-1 surrounding spaces.
var pipeRe *regexp.Regexp = regexp.MustCompile(`\s?\|\s?`)

func main() {
	cmd := &cli.Command{
		Name:        "zei",
		Version:     "v1.3.1",
		Description: "A command snippet cli",
		Usage:       "Execute snippet with ID",
		Action:      execSnippet,
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
				Name:      "update",
				Usage:     "update snippet with ID",
				UsageText: "zei update ID",
				Action:    updateSnippet,
			},
			{
				Name:      "remove",
				Aliases:   []string{"rm", "del"},
				Usage:     "remove snippet with ID",
				UsageText: "zei remove ID",
				Action:    removeSnippet,
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

	snippet, err := zei.GetSnippet(id)
	if err != nil {
		return err
	}

	fmt.Printf("%v\n\n"+pterm.LightYellow("Execute '%v'? [Y/n]: "), colorSnippet(snippet), snippet.ID)

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	if confirm := scanner.Text(); !confirmDefYesRe.MatchString(confirm) {
		return nil
	}

	command := snippet.Command

	tmplFieldNames := tmpl.ParseFields(snippet.Command)

	// If the snippet's command is a template,
	// prompt the user for each field's value,
	// and build the command string.
	if len(tmplFieldNames) > 0 {
		fmt.Println(snippet.Command)

		var builder strings.Builder
		tmplVals := make(map[string]string)

		for _, name := range tmplFieldNames {
			fmt.Printf("%v: ", name)
			scanner.Scan()

			tmplVals[name] = scanner.Text()
		}

		tmplName := fmt.Sprintf("tmpl-%v", snippet.ID)

		cmdTmpl, err := template.New(tmplName).Parse(snippet.Command)
		if err != nil {
			return err
		}

		if err = cmdTmpl.Execute(&builder, tmplVals); err != nil {
			return err
		}

		command = builder.String()
	}

	// Check for piped commands
	piped := pipeRe.Split(command, -1)

	if len(piped) > 1 {
		return runPipedCmd(piped)
	}

	return runCmd(command)
}

func listSnippets(_ context.Context, _ *cli.Command) error {
	snippets, err := zei.GetSnippets()
	if err != nil {
		return err
	}

	for _, snippet := range snippets {
		fmt.Printf("%v\n\n", colorSnippet(snippet))
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

	if valid := zei.IsValidId(id); !valid {
		fmt.Println("Invalid snippet id. Valid characters include alphanumeric, '_', or '-'.")
		return addSnippet(context.TODO(), nil)
	}

	fmt.Print("command: ")
	scanner.Scan()
	cmdText = scanner.Text()

	fmt.Print("description: ")
	scanner.Scan()
	description = scanner.Text()

	fmt.Printf("\nNew snippet\nid: %v\ncommand: %v\ndescription: %v\nSave? [Y/n]: ", id, cmdText, description)
	scanner.Scan()

	if confirm := scanner.Text(); !confirmDefYesRe.MatchString(confirm) {
		return nil
	}

	return zei.AddSnippet(id, cmdText, description)
}

func updateSnippet(_ context.Context, c *cli.Command) error {
	if c.Args().Len() != 1 {
		return fmt.Errorf("invalid snippet id args")
	}

	id := c.Args().First()

	snippet, err := zei.GetSnippet(id)
	if err != nil {
		return err
	}

	var updated zei.Snippet
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Printf("Update '%v' snippet.\nLeave the field blank to keep its current value.\n\n", snippet.ID)

	fmt.Printf("id: %v\nNew id: ", snippet.ID)
	scanner.Scan()
	updated.ID = scanner.Text()

	fmt.Printf("command: %v\nNew command: ", snippet.Command)
	scanner.Scan()
	updated.Command = scanner.Text()

	fmt.Printf("description: %v\nNew description: ", snippet.Description)
	scanner.Scan()
	updated.Description = scanner.Text()

	return zei.UpdateSnippet(id, updated)
}

func removeSnippet(_ context.Context, c *cli.Command) error {
	if c.Args().Len() == 0 {
		return fmt.Errorf("invalid snippet id args")
	}

	ids := c.Args().Slice()

	return zei.RemoveSnippet(ids)
}

// runCommand executes the provided command string
func runCmd(command string) error {
	cmdArgs := cmdstr.Split(command, false)

	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)

	outPipe, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	go readPipe(outPipe)

	return cmd.Run()
}

// runPipedCmd connects and executes the commands
// of the provided piped command string
func runPipedCmd(piped []string) error {
	var cmds []*exec.Cmd

	for _, command := range piped {
		cmdArgs := cmdstr.Split(command, false)

		cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
		cmds = append(cmds, cmd)
	}

	cmdsLen := len(cmds)

	// Connect each command's output pipe
	// to the following one's input
	// and read the last command's output.
	for i, cmd := range cmds {
		outPipe, err := cmd.StdoutPipe()
		if err != nil {
			return err
		}

		if i+1 < cmdsLen {
			cmds[i+1].Stdin = outPipe
		} else {
			go readPipe(outPipe)
		}
	}

	var wg sync.WaitGroup

	// Start the commands in reverse order
	for i := 1; i <= cmdsLen; i++ {
		cmd := cmds[cmdsLen-i]

		if err := cmd.Start(); err != nil {
			return err
		}

		wg.Go(func() {
			if err := cmd.Wait(); err != nil {
				log.Print(err)
			}
		})
	}

	wg.Wait()
	return nil
}

// readPipe writes each line read from the provided pipe
// to standard output.
func readPipe(outPipe io.ReadCloser) {
	reader := bufio.NewReader(outPipe)

	var line string
	var err error

	for err == nil {
		fmt.Print(line)
		line, err = reader.ReadString('\n')
	}
}

// colorSnippet returns the main fields of the snippet
// as a color formatted string.
func colorSnippet(snippet zei.Snippet) string {
	return pterm.Sprintf("[%v] "+pterm.LightGreen("%v\n")+pterm.LightBlue("%v"), snippet.ID, snippet.Command, snippet.Description)
}
