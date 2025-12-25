package common

import "github.com/google/uuid"

type WakeCommand struct {
	EventID     uuid.UUID `json:"event_id"`
	BroadcastIP string    `json:"broadcast_ip"`
	MacAddress  string    `json:"mac"`
}

type PingCommand struct {
	EventID   uuid.UUID `json:"event_id"`
	Subnet    string    `json:"subnet"`
	IpAddress string    `json:"ip"`
}

type RunnerResponse struct {
	EventID uuid.UUID `json:"event_id"`
	OK      bool      `json:"ok"`
	Message string    `json:"msg"`
}
