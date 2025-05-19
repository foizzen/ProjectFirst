package internal

import (
	"fmt"
	"net/http"
	"os"
)

type Middleware func(http.Handler) http.Handler

type Middles struct {
    Middlewares []Middleware
}

func (midls *Middles) Add(mid Middleware) {
	midls.Middlewares = append(midls.Middlewares, mid)
}

func (midls *Middles) Handle(handler http.Handler) http.Handler {
	for _, mdwr := range midls.Middlewares {
		handler = mdwr(handler)
	}
	return handler
}

type Group struct {
	Middlewares []Middleware
}

func (midls *Middles) Group(a... Middleware) *Group {
	xd := Group{Middlewares: make([]Middleware, 0, len(a)+len(midls.Middlewares))}
	xd.Middlewares = append(xd.Middlewares, midls.Middlewares...)
	xd.Middlewares = append(xd.Middlewares, a...)
	return &xd
}

func (gr *Group) Handle(handler http.Handler) http.Handler {
	for _, mdwr := range gr.Middlewares {
		handler = mdwr(handler)
	}
	return handler
}

func (gr *Group) Add(mid Middleware) {
	gr.Middlewares = append(gr.Middlewares, mid)
}

func (gr *Group) Group(a... Middleware) *Group {
	xd := Group{Middlewares: make([]Middleware, 0, len(a)+len(gr.Middlewares))}
	xd.Middlewares = append(xd.Middlewares, gr.Middlewares...)
	xd.Middlewares = append(xd.Middlewares, a...)
	return &xd
}

func MiddleAuth(hand http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tkn, err := r.Cookie("tokenproj")
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		valid, err := CheckJWT(tkn.Value)
		if !valid {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "JWT token not valid or not exist: %s", err)
			return
		}
		hand.ServeHTTP(w, r)
	})
}

func MiddleRecoverLog(hand http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Fprintf(os.Stderr, "recovered: %s \n\n", err)
				w.WriteHeader(http.StatusInternalServerError)
			}
		} ()
		hand.ServeHTTP(w, r)
	})
}