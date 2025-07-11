package main

import (
	"encoding/json"
	"os"
	"regexp"
	"time"

	"github.com/google/uuid"
)

// SessionJSON represents a session in JSON format, excluding unexported fields
type SessionJSON struct {
	Commands      []CommandJSON `json:"commands"`
	SlackThreadTS string        `json:"slack_thread_ts,omitempty"`
}

// CommandJSON represents a command in JSON format
type CommandJSON struct {
	Timestamp time.Time `json:"timestamp"`
	Input     string    `json:"input"`
	Output    string    `json:"output"`
	Comment   string    `json:"comment,omitempty"`
	Redacted  bool      `json:"redacted"`
}

var r = regexp.MustCompile("\u001b]0;.*?\u0007(\u001b)")

func ohsh2nb(s *SessionJSON) *Nb {
	var nb Nb
	for _, comm := range s.Commands {
		var c Cell
		c.CellType = "code"
		c.ExecutionCount = 1
		c.ID = uuid.NewString()
		c.Source = []string{comm.Input}

		if len(comm.Output) > 0 && comm.Output[len(comm.Output)-1] == '\r' {
			out := []byte(comm.Output)
			out[len(out)-1] = '\n'
			comm.Output = string(out)
		}
		comm.Output = r.ReplaceAllString(comm.Output, "$1")
		c.Outputs = []Output{
			{
				Name:       "stdout",
				OutputType: "stream",
				Text:       []string{comm.Output},
			},
		}
		nb.Cells = append(nb.Cells, c)
	}
	nb.Metadata.Kernelspec.DisplayName = "Bash"
	nb.Metadata.Kernelspec.Language = "bash"
	nb.Metadata.Kernelspec.Name = "bash"
	nb.Metadata.LanguageInfo.CodemirrorMode.Name = "bash"
	nb.Metadata.LanguageInfo.FileExtension = ".sh"
	nb.Metadata.LanguageInfo.Mimetype = "text/x-sh"
	nb.Metadata.LanguageInfo.Name = "shell"
	nb.Nbformat = 4
	nb.NbformatMinor = 5
	return &nb
}
func main() {
	dec := json.NewDecoder(os.Stdin)
	var session SessionJSON
	if err := dec.Decode(&session); err != nil {
		panic(err)
	}
	nb := ohsh2nb(&session)
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "\t")
	if err := enc.Encode(nb); err != nil {
		panic(err)
	}

}
