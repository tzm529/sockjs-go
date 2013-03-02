package sockjs

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"unicode/utf8"

//	"time"
)

var dummyCookie = &http.Cookie{
	Name:  "JSESSIONID",
	Value: "dummy",
}

// callback format for htmlfile and jsonp protocols
var reCallback = regexp.MustCompile("[^a-zA-Z0-9-_.]")

var escapable = regexp.MustCompile("[\x00-\x1f\u200c-\u200f\u2028-\u202f\u2060-\u206f\ufff0-\uffff]")

func escaper(m []byte) []byte {
	r, _ := utf8.DecodeRune(m)
	return []byte(fmt.Sprintf(`\u%04x`, r))
}

func aframe(m ...[]byte) (f []byte) {
	strings := make([]string, len(m))
	for i := range m {
		strings[i] = string(m[i])
	}
	s, _ := json.Marshal(strings)
	s = escapable.ReplaceAllFunc(s, escaper)

	f = append(f, 'a')
	f = append(f, s...)
	return
}

func cframe(code int, m string) []byte {
	return []byte(fmt.Sprintf(`c[%d,"%s"]`, code, m))
}

func writeHttpClose(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "c[3000,\"Go away!\"]\n")
}

// addresses match, if the ip is the same, disregard port.
func verifyAddr(a, b string) bool {
	ai := strings.LastIndex(a, ":")
	bi := strings.LastIndex(b, ":")
	return a[:ai+1] == b[:bi+1]
}

// implicitly close the session when the handler function finishes
// in case it doesn't get explicitly closed.
func hfuncCloseWrapper(hfunc func(Session)) func(Session) {
	return func(s Session) {
		hfunc(s)
		s.End()
	}
}
