package internal

import (
	"net/http"
	"strings"
)

func middleAuth(hand http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Compare(r.Header.Get("Auth"), "token") != 0 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		hand.ServeHTTP(w, r)
	})
}