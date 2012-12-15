package sockjs

import (
	"encoding/json"
	"net/http"
)

func frame(prefix string, suffix string, strings ...[]byte) (f []byte) {
	var jsonin []string
	f = append(f, []byte(prefix)...)
	for _, v := range strings {
		jsonin = append(jsonin, string(v))
	}
	s, _ := json.Marshal(&jsonin)
	f = append(f, s...)
	f = append(f, []byte(suffix)...)
	return
}

func writeHttpClose(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`c[3000,"Go away!"]\n`))
}