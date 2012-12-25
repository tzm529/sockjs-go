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

func infoHandler(h *Handler, w http.ResponseWriter, r *http.Request) {
	header := w.Header()
	header.Add("Content-Type", "application/json; charset=UTF-8")
	disableCache(header)
	preflight(header, r)
	w.WriteHeader(http.StatusOK)
	json, _ := json.Marshal(newInfoData(h.config.Websocket))
	w.Write(json)
}

func infoOptionsHandler(w http.ResponseWriter, r *http.Request) {
	h := w.Header()
	h.Add("Access-Control-Allow-Methods", "OPTIONS, GET")
	preflight(h, r)
	enableCache(h)
	w.WriteHeader(http.StatusNoContent)
}
