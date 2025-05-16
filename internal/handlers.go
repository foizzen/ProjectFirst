package internal

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type GamesAns struct {
	Href string `json:"href"`
	Src  string `json:"src"`
	Alt  string `json:"alt"`
}

type App struct {
	DB *sql.DB
}

func (a *App) Mainpage(w http.ResponseWriter, r *http.Request) {
	html, err := os.OpenFile("index.html", os.O_RDWR, 0644)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "index.html not exist")
		return
	}
	defer html.Close()

	ans, err := io.ReadAll(html)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "read index html err hahaha: %s", err)
		return
	}
	w.Write([]byte(ans))
}

func (a *App) Games(w http.ResponseWriter, r *http.Request) {
	games, err := GetGameWMW(a.DB, 5)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "getgameWMW not working: %s", err)
	}
	gamesAns := make([]GamesAns, 0)
	for _, game := range games {
		gam := GamesAns{Href: fmt.Sprintf("/games/%d", game.Id), Src: fmt.Sprintf("/static/images/Games/%s/%s", game.Images_dir, game.Images_url[0]), Alt: game.Images_dir}
		gamesAns = append(gamesAns, gam)
	}
	resp, err := json.Marshal(gamesAns)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "marshal problem: %s", err)
		return
	}
	w.Write(resp)
}

func (a *App) CreateGame(w http.ResponseWriter, r *http.Request) {
	game := Game{}
	que := r.URL.Query()
	if que.Has("gn") {
		game.Game_name = que.Get("gn")
	}
	if que.Has("ac") {
		game.Author_corp = que.Get("ac")
	}
	if que.Has("imd") {
		game.Images_dir = que.Get("imd")
	}
	if que.Has("gpc") {
		game.Game_post_content = que.Get("gpc")
	}
	if que.Has("dr") {
		game.Date_realease = que.Get("dr")
	}
	err := AddGame(a.DB, &game)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "server error in adding game: %s", err)
		return
	}
}

func (a *App) EditGame(w http.ResponseWriter, r *http.Request) {
	game := Game{}
	path := strings.Split(r.URL.Path, "/")
	que := r.URL.Query()
	var err error
	game.Id, err = strconv.Atoi(path[2])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "id in URL not valid: %s", err)
		return
	}
	if que.Has("gn") {
		game.Game_name = que.Get("gn")
	}
	if que.Has("ac") {
		game.Author_corp = que.Get("ac")
	}
	if que.Has("imd") {
		game.Images_dir = que.Get("imd")
	}
	if que.Has("gpc") {
		game.Game_post_content = que.Get("gpc")
	}
	if que.Has("dr") {
		game.Date_realease = que.Get("dr")
	}
	err = EditGame(a.DB, &game)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "server error in edit game: %s", err)
		return
	}
}

func (a *App) Log(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(w, "only post method")
		return
	}
	jsn, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "bad data in request body or nothing: %s", err)
		return
	}
	usr := &User{}
	err = json.Unmarshal(jsn, usr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "can't unpack json: %s", err)
		return
	}
	res := SearchUser(a.DB, usr)
	if res == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "password incorrect: %s", err)
		return
	}
	if res == -1 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "user don't exist or server error: %s", err)
		return
	}
	JWTtoken, err := CreateJWT(usr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error in create jwt token: %s", err)
		return
	}
	cook := &http.Cookie{
		Name: "tokenproj",
		HttpOnly: true,
		Expires: time.Now().Add(5 * time.Minute),
		Value: JWTtoken.String(),
	}
	http.SetCookie(w, cook)
	fmt.Fprint(w, "login successful")
}
