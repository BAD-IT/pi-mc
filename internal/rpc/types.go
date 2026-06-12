package rpc

import "encoding/json"

// Command is sent from PI-mc to pi.
type Command struct {
	Type    string `json:"type"`
	ID      string `json:"id,omitempty"`
	Payload any    `json:"payload,omitempty"`
}

// Event is received from pi.
type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

// Specific Event Payloads

type MessageUpdate struct {
	Delta string `json:"delta"`
}

type ToolExecStart struct {
	ToolName string `json:"tool_name"`
}

type ToolExecEnd struct {
	ToolName string `json:"tool_name"`
	Result   string `json:"result"`
}

type QueueUpdate struct {
	Steering []string `json:"steering"`
	FollowUp []string `json:"follow_up"`
}

type Step struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
}

type DeclareSteps struct {
	Steps []Step `json:"steps"`
}
