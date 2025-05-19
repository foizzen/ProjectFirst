package main

import (
	"ProjectFirst/internal"
	"fmt"
	"net/http"
	"os"
)

func main() {
	os.Stderr, _ = os.OpenFile("logs.log", os.O_RDWR | os.O_CREATE, 0644) 
	serv := http.NewServeMux()
	db, err := internal.GetDb()
	if err != nil {
		panic(err)
	}
	app := &internal.App{DB: db}
	mdls := internal.Middles{Middlewares: []internal.Middleware{internal.MiddleRecoverLog}}	
	mdls_protected := mdls.Group(internal.MiddleAuth)

	// connecting handlers + middlewares:
	serv.Handle("/games/create",    mdls.Handle(http.HandlerFunc(app.CreateGame)))
	serv.Handle("/games/{id}/edit", mdls.Handle(http.HandlerFunc(app.EditGame)))
	
	serv.Handle("/games/{id}/add",  mdls_protected.Handle(http.HandlerFunc(app.AddWaiter)))
	serv.Handle("/games/{id}/",     mdls.Handle(http.HandlerFunc(app.SomeGame)))

	serv.Handle("/login",  mdls.Handle(http.HandlerFunc(app.Log)))
	serv.Handle("/reg",    mdls.Handle(http.HandlerFunc(app.Reg)))  
	serv.Handle("/logout", mdls.Handle(http.HandlerFunc(app.Logout)))

	serv.Handle("/profile", mdls.Handle(http.HandlerFunc(app.Games)))
	serv.Handle("/",        mdls.Handle(http.HandlerFunc(app.Mainpage)))

	statichandler := http.StripPrefix("/static/", http.FileServer(http.Dir("./static")))
	serv.Handle("/static/", statichandler)

	fmt.Println("Server launched")
	err = http.ListenAndServe(":8080", serv)
	if err != nil {
		panic(err)
	}
}
