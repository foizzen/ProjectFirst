package main

import (
	"fmt"
	"net/http"
)

func main() {
	serv := http.NewServeMux()

	// connecting handlers + middlewares:


	statichandler := http.StripPrefix("/static/", http.FileServer(http.Dir("./static")))
	serv.Handle("/static/", statichandler)

	fmt.Println("Server launched")
	http.ListenAndServe(":8080", serv)
}