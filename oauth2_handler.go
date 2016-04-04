package go_dropbox

import (
	"fmt"
	"net/url"
	"net/http"
	"encoding/json"
)

const (
	AUTHORIZATION_URL = "https://www.dropbox.com/1/oauth2/authorize?client_id=%v&response_type=code&redirect_uri=%v"
	TOKEN_EXCHANGE_URL = "https://api.dropbox.com/1/oauth2/token"
)

type Token struct {
	UID   string        `json:"uid"`
	Token string        `json:"access_token"`
	Error *string       `json:"error"`
}

type OAuth2Handler struct {
	Key,
	Secret,
	RedirectURL     string
	Token           *Token
	SuccessCallback func()
	ErrorCallback   func()
}

func (h *OAuth2Handler) AuthCodeURL() string {
	return fmt.Sprintf(AUTHORIZATION_URL, h.Key, h.RedirectURL)
}

func (h *OAuth2Handler) TokenExchange(code string) (*Token, error) {
	data := url.Values{}

	data.Add("code", code)
	data.Add("grant_type", "authorization_code")
	data.Add("client_id", h.Key)
	data.Add("client_secret", h.Secret)
	data.Add("redirect_uri", h.RedirectURL)

	resp, err := http.PostForm(TOKEN_EXCHANGE_URL, data)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	token := &Token{}

	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&token)

	return token, nil
}
