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
	MEDIA_URL = "https://api.dropboxapi.com/1/media/auto/"
	LIST_FOLDER_URL = "https://api.dropboxapi.com/2/files/list_folder"
	SEARCH_URL = "https://api.dropboxapi.com/2/files/search"
)

// Dropbox client
type Dropbox struct {
	Debug         bool   // bool to dump request/response data
	Locale        string // Locale sent to the API to translate/format messages
	Token         *Token
	mediaURL      string
	listFolderURL string
	searchURL     string
	oAuth2Handler *OAuth2Handler
}

// NewDropbox returns a new Dropbox instance
func NewDropbox() *Dropbox {
	return &Dropbox{
		Locale: "en",
		mediaURL: MEDIA_URL,
		listFolderURL: LIST_FOLDER_URL,
		searchURL: SEARCH_URL,
		oAuth2Handler: newOAuth2Handler(),
	}
}

// SetAppInfo sets app_key & app_secret from your Dropbox app
func (db *Dropbox) SetAppInfo(appKey string, appSecret string, redirectURL string) {
	db.oAuth2Handler.setAppKeys(appKey, appSecret, redirectURL)
}

// SetAccessToken sets access token
func (db *Dropbox) SetAccessToken(accessToken string) {
	db.Token = &Token{
		Token: accessToken,
	}
}

func (db *Dropbox) GetAccessToken() string {
	return db.Token.Token
}

func (db *Dropbox) GetAuthURL() string {
	return db.oAuth2Handler.authCodeURL()
}

func (db *Dropbox) ExchangeToken(code string) (*Token, error) {
	return db.oAuth2Handler.tokenExchange(code)
}

// Shares a file for streaming (direct access)
func (db *Dropbox) GetMediaURL(file string) (*SharedURL, *DropboxError) {
	client := &http.Client{}

	data := mediaParameters{
		Locale: db.Locale,
	}

	encoded, _ := json.Marshal(data)

	request, _ := http.NewRequest(POST, db.mediaURL + file, bytes.NewBuffer(encoded))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", db.Token.Token))

	if db.Debug {
		dump, _ := httputil.DumpRequest(request, true)
		fmt.Println(string(dump))
	}

	response, err := client.Do(request)

	if err != nil {
		return nil, &DropboxError{
			StatusCode: http.StatusBadRequest,
			ErrorSummary: err.Error(),
		}
	}

	defer response.Body.Close()

	switch response.StatusCode {
	case 401:
		return nil, &DropboxError{
			StatusCode:http.StatusUnauthorized,
			ErrorSummary: INVALID_ACCESS_TOKEN,
		}
	default:
		dumpData, _ := ioutil.ReadAll(response.Body)

		if db.Debug {
			fmt.Printf("%s\n", string(dumpData))
		}

		var mediaURL SharedURL

		err := json.Unmarshal(dumpData, &mediaURL)

		if err != nil {
			return nil, &DropboxError{
				StatusCode: http.StatusServiceUnavailable,
				ErrorSummary: err.Error(),
			}
		}

		return &mediaURL, nil
	}
}

func (db *Dropbox) ListFolder() (*Folder, *DropboxError) {
	client := &http.Client{}

	data := listFolderParameters{
		Path: "",
		Recursive: true,
		IncludeMediaInfo: false,
		IncludeDeleted: false,
	}

	encoded, _ := json.Marshal(data)

	req, _ := http.NewRequest(POST, db.listFolderURL, bytes.NewBuffer(encoded))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", db.Token.Token))

	if db.Debug {
		dump, _ := httputil.DumpRequest(req, true)
		fmt.Println(string(dump))
	}

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

		if db.Debug {
			fmt.Printf("%s\n", string(dumpData))
		}

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

func (db *Dropbox) SearchMusic() (*Search, *DropboxError) {
	client := &http.Client{}

	data := searchParameters{
		Path: "",
		Query: ".mp3",
		Mode: &searchMode{Tag: "filename" },
		MaxResults: 1000,
	}

	encoded, _ := json.Marshal(data)

	req, _ := http.NewRequest(POST, db.searchURL, bytes.NewBuffer(encoded))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", db.Token.Token))

	if db.Debug {
		dump, _ := httputil.DumpRequest(req, true)
		fmt.Println(string(dump))
	}

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

		if db.Debug {
			fmt.Printf("%s\n", string(dumpData))
		}

		var metadata Search

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
