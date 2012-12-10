package sockjs

import (
	"encoding/json"
	"math/rand"
	"net/http"
)

type infoData struct {
	Websocket    bool     `json:"websocket"`
	CookieNeeded bool     `json:"cookie_needed"`
	Origins      []string `json:"origins"`
	Entropy      uint32   `json:"entropy"`
}

func newInfoData(ws bool) infoData {
	return infoData{
		Websocket:    ws,
		CookieNeeded: true,
		Origins:      []string{"*:*"},
		Entropy:      rand.Uint32(),
	}
}

func handleInfo(w http.ResponseWriter, r *http.Request, s *Handler) {
	h := w.Header()
	addCors(h, r)
	addContentTypeWithoutCache(h, "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json, err := json.Marshal(newInfoData(s.config.Websocket))
	if err != nil {
		panic(err)
	}
	w.Write(json)
}

func handleInfoOptions(w http.ResponseWriter, r *http.Request) {
	h := w.Header()
	addCorsAllowMethods(h, r, "OPTIONS, GET")
	addExpires(h)
	w.WriteHeader(http.StatusNoContent)
}
