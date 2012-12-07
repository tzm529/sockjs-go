package sockjs

import (
	"encoding/json"
)

func aframe(messages ...string) []byte {
	s, err := json.Marshal(&messages)
	if err != nil { panic("sockjs: " + err.Error()) }
	return append([]byte{'a'}, s...)
}
