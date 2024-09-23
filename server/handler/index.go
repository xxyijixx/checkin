package handler

import (
	"fmt"
	"net/http"
)

func IndexHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
}
