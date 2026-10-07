package core

import (
	"encoding/json"
)

type State struct {
	Goal    string
	History []Step

	Transcript []json.RawMessage
}
