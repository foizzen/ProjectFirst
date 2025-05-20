package internal

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
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
	User_id  int    `json:"user_id"`
}

type HFT struct {
	Method string `json:"alg"`
	Type   string `json:"typ"`
}

// Create a JWT token for the user
func CreateJWT(usr *User) (*Token, error) {
	tkn := &Token{}
	payload, err := json.Marshal(PFT{usr.Username, usr.Id})
	if err != nil {
		return nil, fmt.Errorf("token payload error: %s", err)
	}
	tkn.payload = []byte(base64.RawURLEncoding.EncodeToString(payload))

	headers, err := json.Marshal(HFT{Method: "SHA256", Type: "JWT"})
	if err != nil {
		return nil, fmt.Errorf("token headers error: %s", err)
	}
	tkn.headers = []byte(base64.RawURLEncoding.EncodeToString(headers))
	hash := sha256.New()
	hash.Write(tkn.headers)
	hash.Write(tkn.payload)
	if os.Getenv("SECRET_KEY") == "" {
		return nil, fmt.Errorf("can't get secret key")
	}
	hash.Write([]byte(os.Getenv("SECRET_KEY")))
	tkn.signature = []byte(base64.RawURLEncoding.EncodeToString(hash.Sum(nil)))
	return tkn, nil
}

// Check validity JWT token 
func CheckJWT(stkn string) (bool, error) {
	comp := strings.Split(stkn, ".")
	tkn := &Token{}
	var err error
	tkn.headers = []byte(comp[0])
	headers, err := base64.RawURLEncoding.DecodeString(comp[0])
	if err != nil {
		return false, fmt.Errorf("headers decode err: %s", err)
	}
	head := &HFT{}
	err = json.Unmarshal(headers, head)
	if err != nil {
		return false, fmt.Errorf("error, not valid headers1: %s", err)
	}
	if (head.Method != "SHA256") || (head.Type != "JWT") {
		return false, fmt.Errorf("token govno))")
	}
	tkn.payload = []byte(comp[1])
	sig := sha256.New()
	sig.Write(tkn.headers)
	sig.Write(tkn.payload)
	if os.Getenv("SECRET_KEY") == "" {
		return false, fmt.Errorf("can't get secret key")
	}
	sig.Write([]byte(os.Getenv("SECRET_KEY")))
	tkn.signature = []byte(base64.RawURLEncoding.EncodeToString(sig.Sum(nil)))
	if strings.Compare(string(tkn.signature), comp[2]) == 0 {
		return true, nil
	} else {
		return false, nil
	} 
}

// Gets user_id from the token. It return 0 if something wrong
func GetUserid(tkn string) int {
	comp := strings.Split(tkn, ".")
	data, err := base64.RawURLEncoding.DecodeString(comp[2])
	if err != nil {
		return 0
	}
	usr := &User{}
	err = json.Unmarshal(data, usr)
	if err != nil {
		return 0 
	}
	return usr.Id
}