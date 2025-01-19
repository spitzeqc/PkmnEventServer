package handlers

import (
	"net/http"
)

const connTestPage = "<html><head>HTML Page</head><body bgcolor=\"#FFFFFF\">tmp</body></html>"

/*
 * Handle requests to "/runConnTest"
 */
func HandleConnTest(w http.ResponseWriter, r *http.Request) {
	LogInfo("/runConnTest requested")
	if r.URL.String() == "/runConnTest" {
		w.Write( []byte(connTestPage) )
	}
}