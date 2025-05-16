package main

import (
	"ProjectFirst/internal"
	"fmt"
	"net/http"
)

func main() {
	serv := http.NewServeMux()
	db, err := internal.GetDb()
	if err != nil {
		panic(err)
	}
	app := &internal.App{DB: db}

	// connecting handlers + middlewares:
	// serv.HandleFunc("/post/{id}", )
	serv.Handle("/post/create", http.HandlerFunc(app.CreateGame))
	serv.Handle("/post/{id}/edit", http.HandlerFunc(app.EditGame))
	serv.Handle("/post/{id}/add", http.HandlerFunc(app.EditGame))
	serv.Handle("/login", http.HandlerFunc(app.Log))
	serv.Handle("/reg", http.HandlerFunc(app.Log))  // ??
	// serv.HandleFunc("/logout", )
	serv.Handle("/profile", http.HandlerFunc(app.Games))
	serv.Handle("/", http.HandlerFunc(app.Mainpage))

	statichandler := http.StripPrefix("/static/", http.FileServer(http.Dir("./static")))
	serv.Handle("/static/", statichandler)

	fmt.Println("Server launched")
	err = http.ListenAndServe(":8080", serv)
	if err != nil {
		panic(err)
	}
}
