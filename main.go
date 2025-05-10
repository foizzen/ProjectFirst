package main

import (
	"fmt"
	"net/http"
)

func main() {
	serv := http.NewServeMux()

	// connecting handlers + middlewares:
	// serv.HandleFunc("/post/{id}", )
	// serv.HandleFunc("/post/create", )
	// serv.HandleFunc("/post/{id}/edit", )
	// serv.HandleFunc("/login", )
	// serv.HandleFunc("/reg", )  // ??
	// serv.HandleFunc("/logout", )
	// serv.HandleFunc("/", )

	statichandler := http.StripPrefix("/static/", http.FileServer(http.Dir("./static")))
	serv.Handle("/static/", statichandler)

	fmt.Println("Server launched")
	http.ListenAndServe(":8080", serv)
}