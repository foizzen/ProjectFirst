package internal

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
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
		panic(fmt.Errorf("index.html not exist -- %s", err))
	}
	defer html.Close()

	ans, err := io.ReadAll(html)
	if err != nil {
		panic(fmt.Errorf("read index html err hahaha -- %s", err))
	}
	w.Write([]byte(ans))
}

func (a *App) Games(w http.ResponseWriter, r *http.Request) {
	games, err := GetGameWMW(a.DB, 5)
	if err != nil {
		panic(fmt.Errorf("getgameWMW not working -- %s", err))
	}
	gamesAns := make([]GamesAns, 0)
	for _, game := range games {
		gam := GamesAns{
			Href: fmt.Sprintf("/games/%d", game.Id), 
			Src:  fmt.Sprintf("/static/images/Games/%s", game.Images_url[0]), 
			Alt:  game.Images_dir,
		}
		gamesAns = append(gamesAns, gam)
	}
	resp, err := json.Marshal(gamesAns)
	if err != nil {
		panic(fmt.Errorf("marshal problem -- %s", err))
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
		panic(fmt.Errorf("server error in adding game -- %s", err))
	}
}

func (a *App) EditGame(w http.ResponseWriter, r *http.Request) {
	game := Game{}
	path := strings.Split(r.URL.Path, "/")
	que := r.URL.Query()
	var err error
	game.Id, err = strconv.Atoi(path[2])
	if err != nil {
		panic(fmt.Errorf("id in URL not valid -- %s", err))
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
		panic(fmt.Errorf("server error in edit game -- %s", err))
	}
}

func (a *App) SomeGame(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/")
	game_id, err := strconv.Atoi(path[2])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "can't parse URL: %s", err)
		return
	}
	game, err := GetSomeGame(a.DB, game_id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "game does't exist: %s", err)
	}
	temp, err := template.ParseFiles("index-link_temp.html")
	if err != nil {
		panic(fmt.Errorf("index-link err in parse -- %s", err))
	}
	err = temp.Execute(w, &struct {
		Game_name         string
		Images_url        []string
		Game_post_content string
		Waitings          int
		// Date_realease     string
		// Game_id           string
		// User_id           string
	}{
		Game_name:         game.Game_name,
		Images_url:        game.Images_url,
		Game_post_content: game.Game_post_content,
		Waitings:          game.Waitings,
	})
	if err != nil {
		panic(fmt.Errorf("index-link params err -- %s", err))
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
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "can't unpack json in request: %s", err)
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
		fmt.Fprintf(w, "user don't exist: %s", err)
		return
	}
	JWTtoken, err := CreateJWT(usr)
	if err != nil {
		panic(fmt.Errorf("error in create jwt token -- %s", err))
	}
	cook := &http.Cookie{
		Name:     "tokenproj",
		HttpOnly: true,
		Expires: time.Now().Add(5 * time.Minute),
		Value: JWTtoken.String(),
	}
	http.SetCookie(w, cook)
	fmt.Fprint(w, "login successful")
}

func (a *App) Reg(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(w, "only post method")
		return
	}
	usr := &User{}
	comp, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "can't read body of request: %s", err)
		return
	}
	err = json.Unmarshal(comp, usr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "can't decode json in body of request: %s", err)
		return
	}
	err = AddUser(a.DB, usr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "busy username: %s", err)
	}
	JWTtoken, err := CreateJWT(usr)
	if err != nil {
		panic(fmt.Errorf("error in create jwt token -- %s", err))
	}
	cook := &http.Cookie{
		Name:     "tokenproj",
		HttpOnly: true,
		Expires: time.Now().Add(5 * time.Minute),
		Value: JWTtoken.String(),
	}
	http.SetCookie(w, cook)
	fmt.Fprint(w, "reg + login successful")
}

func (a *App) Logout(w http.ResponseWriter, r *http.Request) {
	cook, err := r.Cookie("tokenproj")
	if err != nil {
		return
	}
	cook.Expires = time.Now().AddDate(0, 0, -1)
}

type Waiter struct {
	Game_id int `json:"game_id"`
	User_id int `json:"user_id"`
}

func (a *App) AddWaiter(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	wtr := &Waiter{}

	c, _ := r.Cookie("tokenproj")
	comp := strings.Split(c.Value, ".")
	d, err := base64.RawStdEncoding.DecodeString(comp[1])
	if err != nil {
		panic(fmt.Errorf("can't decode paylodad of token: %s", err))
	}
	err = json.Unmarshal(d, wtr)
	if err != nil {
		panic(fmt.Errorf("can't unmarshal paylodad of token: %s", err))
	}

	bdy, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "can't read the request body: %s", err)
		return
	}
	err = json.Unmarshal(bdy, wtr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "can't unmarshall json in request: %s", err)
		return
	}

	err = AddWaiter(a.DB, wtr.Game_id, wtr.User_id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "can't add waiter: %s", err)
		return
	}
}