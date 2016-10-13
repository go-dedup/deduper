package command

import (
	//"encoding/json"
	"strings"

	"gopkg.in/suntong/deduper.v1/minhash"
)

// WriteCommand represents a command to persist a
// document ID and it's generated minhash value.
type WriteCommand struct {
	// ID is the document id
	ID string `json:"id"`

	// Value is the value to be written
	Value string `json:"value"`
}

// NewWriteCommand creates a new write command.
func NewWriteCommand(id string, value string) *WriteCommand {
	return &WriteCommand{
		ID:    id,
		Value: value,
	}
}

// CommandName returns the name of the command.
func (c *WriteCommand) CommandName() string {
	return "write"
}

// Apply writes a value to a key.
func (c *WriteCommand) Apply(server raft.Server) (interface{}, error) {
	mh := server.Context().(*minhash.MinHasher)
	mh.Add(c.ID, strings.NewReader(c.Value))
	return nil, nil
}

/*
func s() {
	matches := s.minhasher.FindSimilar(req.Body, threshold)
	_ = json.NewEncoder(w).Encode(matches)
}
*/
