package sockjs

import (
	"encoding/json"
	
)

func aframe(messages ...[]byte) []byte {
	var jsonin []string
	for _, v := range messages {
		jsonin = append(jsonin, string(v))
	}
	s, _ := json.Marshal(&jsonin)
	return append([]byte{'a'}, s...)
}
