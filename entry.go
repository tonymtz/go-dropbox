package go_dropbox

// no exported

type listFolderParameters struct {
	Path             string `json:"path"`
	Recursive        bool   `json:"recursive"`
	IncludeMediaInfo bool   `json:"include_media_info"`
	IncludeDeleted   bool   `json:"include_deleted"`
}

type sharesParameters struct {
	Locale   string `json:"locale"`
	ShortURL bool   `json:"short_url"`
}

// exported

type Entry struct {
	UID      string `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Path     string `json:"path_display,omitempty"`
	Revision string `json:"rev,omitempty"`
	Size     int    `json:"size,omitempty"`
}

type Folder struct {
	Entries []*Entry `json:"entries"`
	Cursor  string   `json:"cursor"`
	HasMore bool     `json:"has_more"`
}

type SharedURL struct {
	URL     string `json:"url"`
	Expires string `json:"expires"`
}

type DropboxError struct {
	ErrorSummary string `json:"error_summary"`
	StatusCode   int
}
