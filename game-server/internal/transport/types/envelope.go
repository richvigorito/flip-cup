package types

import "encoding/json"

type Envelope struct {
  Type    string          `json:"type"`
  GameID  string          `json:"gameId,omitempty"`
  ID      string          `json:"id,omitempty"`
  Name    string          `json:"name,omitempty"`
  Payload json.RawMessage `json:"payload"`
}
