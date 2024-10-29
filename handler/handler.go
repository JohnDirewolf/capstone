package handler

import (
	"net/http"
)

func Page(w http.ResponseWriter, r *http.Request) {
	const page = `
	<html>
	<head></head>
	<body>
		<p> Yay Capstone Project!</p> 
		<p>I'm your DM, a servemux running on a Go server. </p>
	</body>
	</html>
	`
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(200)
	w.Write([]byte(page))
}
