package go_dropbox

import (
	"testing"
)

func createDropbox() *Dropbox {
	db := NewDropbox()
	db.SetAppInfo("my_app_key", "my_app_secret", "my_url")
	return db
}

func TestNewDropbox(t *testing.T) {
	var currentApiUrl = "https://api.dropbox.com/1";

	db := createDropbox()

	if db.APIURL != currentApiUrl {
		t.Fatalf("APIURL didn't match")
	}
}

func TestAccessToken(t *testing.T) {
	var myToken = "my_access_token"

	db := createDropbox()
	db.SetAccessToken(myToken)

	if db.GetAccessToken() != myToken {
		t.Fatalf("APIURL didn't match")
	}
}

func TestGetAuthURL(t *testing.T) {
	var urlToGetToken = "https://www.dropbox.com/1/oauth2/authorize?client_id=my_app_key&response_type=code&redirect_uri=my_url"

	db := createDropbox()

	if db.GetAuthURL() != urlToGetToken {
		t.Fatalf("AuthURL didn't match")
	}
}
