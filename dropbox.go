package go_dropbox

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"net/http/httputil"
	"encoding/json"
	"bytes"
)

const (
	POST = "POST"
	INVALID_ACCESS_TOKEN = "invalid_access_token"
	API_URL = "https://api.dropbox.com/1"
	LIST_FOLDER_URL = "https://api.dropboxapi.com/2/files/list_folder"
	SHARES_URL = "https://api.dropboxapi.com/1/shares/auto/%v"
)

// Dropbox client
type Dropbox struct {
	Locale        string // Locale sent to the API to translate/format messages.
	Token         *Token
	OAuth2Handler *OAuth2Handler
}

// NewDropbox returns a new Dropbox instance
func NewDropbox() *Dropbox {
	return &Dropbox{Locale: "en"}
}

// SetAppInfo sets app_key & app_secret from your Dropbox app
func (db *Dropbox)SetAppInfo(appKey string, appSecret string, redirectURL string) {
	oAuth2Handler := OAuth2Handler{
		Key:           appKey,
		Secret:        appSecret,
		RedirectURL:   redirectURL,
	}

	db.OAuth2Handler = &oAuth2Handler
}

// SetAccessToken sets access token
func (db *Dropbox) SetAccessToken(accessToken string) {
	db.Token = &Token{
		Token: accessToken,
	}
}

func (db *Dropbox) GetAuthURL() string {
	return db.OAuth2Handler.AuthCodeURL()
}

func (db *Dropbox) ExchangeToken(code string) (*Token, error) {
	return db.OAuth2Handler.TokenExchange(code)
}

func (db *Dropbox) GetURL(file string) (*SharedURL, *DropboxError) {
	client := &http.Client{}

	data := sharesParameters{
		ShortURL: false,
		Locale: db.Locale,
	}

	encoded, _ := json.Marshal(data)

	req, _ := http.NewRequest(POST, fmt.Sprintf(SHARES_URL, file), bytes.NewBuffer(encoded))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", db.Token.Token))

	dump, _ := httputil.DumpRequest(req, true)
	fmt.Println(string(dump))

	res, err := client.Do(req)

	if err != nil {
		return nil, &DropboxError{
			StatusCode: http.StatusBadRequest,
			ErrorSummary: err.Error(),
		}
	}

	defer res.Body.Close()

	switch res.StatusCode {
	case 401:
		return nil, &DropboxError{
			StatusCode:http.StatusUnauthorized,
			ErrorSummary: INVALID_ACCESS_TOKEN,
		}
	default:
		dumpData, _ := ioutil.ReadAll(res.Body)
		fmt.Printf("%s\n", string(dumpData))

		var sharedURL SharedURL

		err := json.Unmarshal(dumpData, &sharedURL)

		if err != nil {
			return nil, &DropboxError{
				StatusCode:http.StatusServiceUnavailable,
				ErrorSummary: err.Error(),
			}
		}

		return &sharedURL, nil
	}
}

func (db *Dropbox) ListFolder() (*Folder, *DropboxError) {
	client := &http.Client{}

	data := listFolderParameters{
		Path: "",
		Recursive: false,
		IncludeMediaInfo: false,
		IncludeDeleted: false,
	}

	encoded, _ := json.Marshal(data)

	req, _ := http.NewRequest(POST, LIST_FOLDER_URL, bytes.NewBuffer(encoded))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", db.Token.Token))

	dump, _ := httputil.DumpRequest(req, true)
	fmt.Println(string(dump))

	res, err := client.Do(req)

	if err != nil {
		return nil, &DropboxError{
			StatusCode: http.StatusBadRequest,
			ErrorSummary: err.Error(),
		}
	}

	defer res.Body.Close()

	switch res.StatusCode {
	case 401:
		return nil, &DropboxError{
			StatusCode:http.StatusUnauthorized,
			ErrorSummary: INVALID_ACCESS_TOKEN,
		}
	default:
		dumpData, _ := ioutil.ReadAll(res.Body)
		fmt.Printf("%s\n", string(dumpData))

		var metadata Folder

		err := json.Unmarshal(dumpData, &metadata)

		if err != nil {
			return nil, &DropboxError{
				StatusCode:http.StatusServiceUnavailable,
				ErrorSummary: err.Error(),
			}
		}

		return &metadata, nil
	}
}
