package types

import "time"

type Participant struct {
	ID             string    `json:"id"`
	Active         bool      `json:"active"`
	AddedBy        string    `json:"added_by"`
	JoinedAt       time.Time `json:"joinedAt"`
	DisconnectedAt time.Time `json:"disconnectedAt"`
}

func NewParticipant(id string, addedBy string) *Participant {
	return &Participant{
		ID:       id,
		Active:   true,
		AddedBy:  addedBy,
		JoinedAt: time.Now(),
	}
}
