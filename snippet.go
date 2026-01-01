package zei

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"time"

	cmdstr "github.com/sammy-t/zei/internal/cmdStr"
)

type Snippet struct {
	ID          string `gorm:"primaryKey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Command     string `gorm:"not null"`
	Description string
}

// DisplayText returns the main fields of the snippet as a "friendly" string.
func (s *Snippet) DisplayText() string {
	return fmt.Sprintf("[%v] %v\n%v", s.ID, s.Command, s.Description)
}

// Exec attempts to execute the snippet's command.
func (s *Snippet) Exec() error {
	cmdArgs := cmdstr.Split(s.Command, false)

	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)

	outPipe, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	go readPipe(outPipe)

	if err := cmd.Run(); err != nil {
		return err
	}

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
