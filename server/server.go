package server

import (
	"checkin/server/handler"
	"net/http"
)

func Run() {

	http.HandleFunc("/", handler.IndexHandler)

	http.ListenAndServe(":7788", nil)
}
