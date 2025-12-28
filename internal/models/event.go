package models

import (
	"github.com/google/uuid"
)

type EventStatus string

const (
	StatusScheduled EventStatus = "SCHEDULED"
	StatusFired     EventStatus = "FIRED"
	StatusCancelled EventStatus = "CANCELLED"
)

type Event struct {
	ID          string      `json:"id"`
	ReferenceID string      `json:"reference_id"`
	Topic       string      `json:"topic"`
	ExecuteAt   int64       `json:"execute_at"`
	Payload     []byte      `json:"payload"`
	Status      EventStatus `json:"status"`
	CreatedAt   int64       `json:"created_at"`
	UpdatedAt   int64       `json:"updated_at"`
}

func (e *Event) GenerateId() {
	e.ID = uuid.NewString()
}
