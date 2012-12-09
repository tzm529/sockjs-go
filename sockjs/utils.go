package sockjs

import (
	"encoding/json"
)

func aframe(messages ...string) []byte {
	s, _ := json.Marshal(&messages)
	return append([]byte{'a'}, s...)
}
