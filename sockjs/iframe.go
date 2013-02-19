package sockjs

import (
	"net/http"
)

const (
	iframePageFormat string = `<!DOCTYPE html>
<html>
<head>
  <meta http-equiv="X-UA-Compatible" content="IE=edge" />
  <meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
  <script>
    document.domain = document.domain;
    _sockjs_onload = function(){SockJS.bootstrap_iframe();};
  </script>
  <script src="%s"></script>
</head>
<body>
  <h2>Don't panic!</h2>
  <p>This is a SockJS hidden iframe. It's used for cross domain magic.</p>
</body>
</html>`
)

func iframeHandler(h *handler, w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("If-None-Match") == h.config.iframeHash {
		w.WriteHeader(http.StatusNotModified)
		return
	}

	header := w.Header()
	header.Add("Content-Type", "text/html; charset=UTF-8")
	enableCache(header)
	header.Add("ETag", h.config.iframeHash)
	w.Write(h.config.iframePage)
}
