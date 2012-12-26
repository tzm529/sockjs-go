package sockjs

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
)

// callback format for htmlfile and jsonp protocols
var reCallback = regexp.MustCompile(`[^a-zA-Z0-9-_.]`)

func frame(prefix, suffix string, m ...[]byte) (f []byte) {
	strings := make([]string, len(m))
	for i := range m {
		strings[i] = string(m[i])
	}
	s, _ := json.Marshal(&strings)

	f = append(f, []byte(prefix)...)
	f = append(f, 'a')
	f = append(f, s...)
	f = append(f, []byte(suffix)...)
	return
}

func cframe(prefix string, code int, m string, suffix string) []byte {
	return []byte(fmt.Sprintf(`%sc[%d,"%s"]%s`, prefix, code, m, suffix))
}

func writeHttpClose(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	w.Write(cframe("", 3000, "Go away!", "\n"))
}
