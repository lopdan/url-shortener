package shortener

type Redirect struct {
	Code      string // Key
	URL       string // Url
	CreatedAt int64  // Shortcut creation time
}