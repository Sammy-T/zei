package zei

import (
	"fmt"
	"time"
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
