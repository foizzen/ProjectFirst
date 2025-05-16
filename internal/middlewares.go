package internal

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type Token struct {
	headers   []byte
	payload   []byte
	signature []byte
}

func (t *Token) String() string {
	return string(t.headers)+"."+string(t.payload)+"."+string(t.signature)
}

type PFT struct {
	Username string `json:"username"`
}

type HFT struct {
	Method string `json:"alg"`
	Type   string `json:"typ"`
}

func CreateJWT(usr *User) (*Token, error) {
	tkn := &Token{}
	payload, err := json.Marshal(PFT{usr.Username})
	if err != nil {
		return nil, fmt.Errorf("token payload error: %s", err)
	}
	base64.RawStdEncoding.Encode(tkn.payload, payload)

	headers, err := json.Marshal(HFT{Method: "JWT", Type: "SHA256"})
	if err != nil {
		return nil, fmt.Errorf("token headers error: %s", err)
	}
	base64.RawStdEncoding.Encode(tkn.headers, headers)
	hash := sha256.New()
	hash.Write(tkn.headers)
	hash.Write(tkn.payload)
	if os.Getenv("SECRET_KEY") == "" {
		return nil, fmt.Errorf("can't get secret key")
	}
	hash.Write([]byte(os.Getenv("SECRET_KEY")))
	tkn.signature = hash.Sum(nil)
	return tkn, nil
}

func CheckJWT(stkn string) (bool, error) {
	comp := strings.Split(stkn, ".")
	tkn := &Token{}
	var err error
	tkn.headers, err = base64.RawStdEncoding.DecodeString(comp[0])
	if err != nil {
		return false, fmt.Errorf("headers decode err: %s", err)
	}
	head := &HFT{}
	err = json.Unmarshal(tkn.headers, head)
	if err != nil {
		return false, fmt.Errorf("error, not valid headers1: %s", err)
	}
	if (head.Method != "SHA256") || (head.Type != "JWT") {
		return false, nil
	}
	tkn.payload, err = base64.RawStdEncoding.DecodeString(comp[1])
	if err != nil {
		return false, fmt.Errorf("payload decode err: %s", err)
	}
	sig := sha256.New()
	sig.Write(tkn.headers)
	sig.Write(tkn.payload)
	sig.Write([]byte(os.Getenv("SECRET_KEY")))
	tkn.signature = sig.Sum(nil)
	if strings.Compare(string(tkn.signature), comp[2]) == 0 {
		return true, nil
	} else {
		return false, nil
	}
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

func MiddleLog(hand http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Print("AAAAAAAAAAa")
		hand.ServeHTTP(w, r)
	})
}