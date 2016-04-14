package go_dropbox

import (
	"testing"
	"net/http"
	"net/http/httptest"
	"fmt"
	"reflect"
	"github.com/stretchr/testify/assert"
)

func TestNewDropbox(t *testing.T) {
	db := NewDropbox()

	a := reflect.TypeOf(db)
	b := reflect.TypeOf(&Dropbox{})

	assert.Equal(t, a, b,
		fmt.Sprintf("Expected same type, got %v, want %v", a, b),
	)

	c := db.mediaURL
	d := "https://api.dropboxapi.com/1/media/auto/"

	assert.Equal(t, c, d,
		fmt.Sprintf("Expected same value, got %v, want %v", c, d),
	)

	e := db.Locale
	f := "en"

	assert.Equal(t, e, f,
		fmt.Sprintf("Expected same value, got %v, want %v", e, f),
	)

	g := reflect.TypeOf(db.oAuth2Handler)
	h := reflect.TypeOf(&OAuth2Handler{})

	assert.Equal(t, g, h,
		fmt.Sprintf("Expected same type, got %v, want %v", g, h),
	)

	i := db.oAuth2Handler.authURL
	j := "https://api.dropbox.com/1/oauth2/token"

	assert.Equal(t, i, j,
		fmt.Sprintf("Expected same value, got %v, want %v", i, j),
	)
}

func TestExchangeToken(t *testing.T) {
	var sampleResponse = "{\"access_token\":\"my_unique_token\",\"token_type\":\"bearer\",\"uid\":\"12345\"}"

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, sampleResponse)
	}))

	defer testServer.Close()

	db := NewDropbox()
	db.SetAppInfo("my_app_key", "my_app_secret", "my_url")
	db.oAuth2Handler.authURL = testServer.URL
	//db.Debug = true

	token, err := db.ExchangeToken("this_code")

	if assert.Nil(t, err) {
		a := reflect.TypeOf(token)
		b := reflect.TypeOf(&Token{})

		assert.Equal(t, a, b,
			fmt.Sprintf("Expected same type, got %v, want %v", a, b),
		)

		c := token.Token
		d := "my_unique_token"

		assert.Equal(t, c, d,
			fmt.Sprintf("Expected same value, got %v, want %v", c, d),
		)
	}
}

func TestGetMediaURLReturn(t *testing.T) {
	var sampleResponse = "{\"url\":\"my_fantastic_url_string\",\"expires\":\"expiration_date_string\"}"

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, sampleResponse)
	}))

	defer testServer.Close()

	db := NewDropbox()
	db.SetAccessToken("my_access_token")
	db.mediaURL = testServer.URL

	sharedURL, dbErr := db.GetMediaURL("fantastic_song")

	if assert.Nil(t, dbErr) {
		a := reflect.TypeOf(sharedURL)
		b := reflect.TypeOf(&SharedURL{})

		assert.Equal(t, a, b,
			fmt.Sprintf("Expected same type, got %v, want %v", a, b),
		)

		c := sharedURL.URL
		d := "my_fantastic_url_string"

		assert.Equal(t, c, d,
			fmt.Sprintf("Expected same value, got %v, want %v", c, d),
		)

		e := sharedURL.Expires
		f := "expiration_date_string"

		assert.Equal(t, e, f,
			fmt.Sprintf("Expected same value, got %v, want %v", e, f),
		)
	}
}

func TestListFolderReturn(t *testing.T) {
	var sampleResponse = `
	{
		"entries": [
			{
				".tag": "file",
				"name": "Prime_Numbers.txt",
				"path_lower": "/homework/math/prime_numbers.txt",
				"path_display": "/Homework/math/Prime_Numbers.txt",
				"id": "id:a4ayc_80_OEAAAAAAAAAXw",
				"client_modified": "2015-05-12T15:50:38Z",
				"server_modified": "2015-05-12T15:50:38Z",
				"rev": "a1c10ce0dd78",
				"size": 7212,
				"sharing_info": {
					"read_only": true,
					"parent_shared_folder_id": "84528192421",
					"modified_by": "dbid:AAH4f99T0taONIb-OurWxbNQ6ywGRopQngc"
				},
				"property_groups": [
					{
						"template_id": "ptid:1a5n2i6d3OYEAAAAAAAAAYa",
						"fields": [
							{
								"name": "Security Policy",
								"value": "Confidential"
							}
						]
					}
				]
			},
			{
				".tag": "folder",
				"name": "math",
				"path_lower": "/homework/math",
				"path_display": "/Homework/math",
				"id": "id:a4ayc_80_OEAAAAAAAAAXz",
				"sharing_info": {
					"read_only": false,
					"parent_shared_folder_id": "84528192421"
				},
				"property_groups": [
					{
						"template_id": "ptid:1a5n2i6d3OYEAAAAAAAAAYa",
						"fields": [
							{
								"name": "Security Policy",
								"value": "Confidential"
							}
						]
					}
				]
			}
		],
		"cursor": "ZtkX9_EHj3x7PMkVuFIhwKYXEpwpLwyxp9vMKomUhllil9q7eWiAu",
		"has_more": false
	}`

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, sampleResponse)
	}))

	defer testServer.Close()

	db := NewDropbox()
	db.SetAccessToken("my_access_token")
	db.listFolderURL = testServer.URL

	folderMeta, dbErr := db.ListFolder()

	if assert.Nil(t, dbErr) {
		a := reflect.TypeOf(folderMeta)
		b := reflect.TypeOf(&Folder{})

		assert.Equal(t, a, b,
			fmt.Sprintf("Expected same type, got %v, want %v", a, b),
		)

		c := folderMeta.Cursor
		d := "ZtkX9_EHj3x7PMkVuFIhwKYXEpwpLwyxp9vMKomUhllil9q7eWiAu"

		assert.Equal(t, c, d,
			fmt.Sprintf("Expected same value, got %v, want %v", c, d),
		)

		e := folderMeta.HasMore
		f := false

		assert.Equal(t, e, f,
			fmt.Sprintf("Expected same value, got %v, want %v", e, f),
		)

		g := reflect.TypeOf(folderMeta.Entries)
		h := reflect.TypeOf([]*Entry{})

		assert.Equal(t, g, h,
			fmt.Sprintf("Expected same type, got %v, want %v", g, h),
		)

		i := len(folderMeta.Entries)
		j := 2

		assert.Equal(t, i, j,
			fmt.Sprintf("Expected same value, got %v, want %v", i, j),
		)

		k := folderMeta.Entries[0].Name
		l := "Prime_Numbers.txt"

		assert.Equal(t, k, l,
			fmt.Sprintf("Expected same value, got %v, want %v", k, l),
		)
	}
}

func TestAccessToken(t *testing.T) {
	var myToken = "my_access_token"

	db := NewDropbox()
	db.SetAccessToken(myToken)

	a := db.GetAccessToken()
	b := myToken

	assert.Equal(t, a, b,
		fmt.Sprintf("Expected same value, got %v, want %v", a, b),
	)
}

func TestGetAuthURL(t *testing.T) {
	db := NewDropbox()
	db.SetAppInfo("my_app_key", "my_app_secret", "my_url")

	a := db.GetAuthURL()
	b := "https://www.dropbox.com/1/oauth2/authorize?client_id=my_app_key&response_type=code&redirect_uri=my_url"

	assert.Equal(t, a, b,
		fmt.Sprintf("Expected same value, got %v, want %v", a, b),
	)
}
