package sockjs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
)

// callback format for htmlfile and jsonp protocols
var reCallback = regexp.MustCompile("[^a-zA-Z0-9-_.]")

var escapable = regexp.MustCompile("[\x00-\x1f\u200c-\u200f\u2028-\u202f\u2060-\u206f\ufff0-\uffff]")
func escaper(m []byte) []byte {
	return []byte(fmt.Sprintf(`\u%04x`, bytes.Runes(m)[0]))
}

func aframe(prefix, suffix string, m ...[]byte) (f []byte) {
 	strings := make([]string, len(m))
	for i := range m {
		strings[i] = string(m[i])
	}
	s, _ := json.Marshal(strings)
	s = escapable.ReplaceAllFunc(s, escaper)

	f = append(f, prefix...)
	f = append(f, 'a')
	f = append(f, s...)
	f = append(f, suffix...)
	return
}

func cframe(prefix string, code int, m string, suffix string) []byte {
	return []byte(fmt.Sprintf(`%sc[%d,"%s"]%s`, prefix, code, m, suffix))
}

func writeHttpClose(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	w.Write(cframe("", 3000, "Go away!", "\n"))
}
