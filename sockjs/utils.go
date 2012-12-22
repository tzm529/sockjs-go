package sockjs

import (
	"fmt"
	"encoding/json"
	"net/http"
)

func aframe(suffix string, m ...[]byte) (f []byte) {
	strings := make([]string, len(m))
	for i := range m {
		strings[i] = string(m[i])
	}
	s, _ := json.Marshal(&strings)

	f = append(f, 'a')
	f = append(f, s...)
	f = append(f, []byte(suffix)...)
	return
}

func cframe(suffix string, code int, m string) (f []byte) {
	f = []byte(fmt.Sprintf(`c[%d,"%s"]%s`, code, m, suffix))
	return
}

func writeHttpClose(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`c[3000,"Go away!"]\n`))
}
