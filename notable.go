package notable

import (
	"fmt"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Request received: %s %s", r.Method, r.URL.Path)
}
