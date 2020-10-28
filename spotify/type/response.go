package spotify_type

import "time"

// Auto generator: https://mholt.github.io/json-to-go/

type RecentlyPlayed struct {
	Items   []Items `json:"items"`
	Next    string  `json:"next"`
	Cursors Cursors `json:"cursors"`
	Limit   int     `json:"limit"`
	Href    string  `json:"href"`
}

type ExternalUrls struct {
	Spotify string `json:"spotify"`
}

type Track struct {
	Artists          []Artists    `json:"artists"`
	AvailableMarkets []string     `json:"available_markets"`
	DiscNumber       int          `json:"disc_number"`
	DurationMs       int          `json:"duration_ms"`
	Explicit         bool         `json:"explicit"`
	ExternalUrls     ExternalUrls `json:"external_urls"`
	Href             string       `json:"href"`
	ID               string       `json:"id"`
	Name             string       `json:"name"`
	PreviewURL       string       `json:"preview_url"`
	TrackNumber      int          `json:"track_number"`
	Type             string       `json:"type"`
	URI              string       `json:"uri"`
}

type Context struct {
	URI          string       `json:"uri"`
	ExternalUrls ExternalUrls `json:"external_urls"`
	Href         string       `json:"href"`
	Type         string       `json:"type"`
}

type Items struct {
	Track    Track     `json:"track"`
	PlayedAt time.Time `json:"played_at"`
	Context  Context   `json:"context"`
}

type Cursors struct {
	After  string `json:"after"`
	Before string `json:"before"`
}

type CurrentlyPlaying struct {
	Context              Context `json:"context"`
	Timestamp            int64   `json:"timestamp"`
	ProgressMs           int     `json:"progress_ms"`
	IsPlaying            bool    `json:"is_playing"`
	CurrentlyPlayingType string  `json:"currently_playing_type"`
	Actions              Actions `json:"actions"`
	Item                 Item    `json:"item"`
}

type Disallows struct {
	Resuming bool `json:"resuming"`
}

type Actions struct {
	Disallows Disallows `json:"disallows"`
}

type Images struct {
	Height int    `json:"height"`
	URL    string `json:"url"`
	Width  int    `json:"width"`
}

type Album struct {
	AlbumType    string       `json:"album_type"`
	ExternalUrls ExternalUrls `json:"external_urls"`
	Href         string       `json:"href"`
	ID           string       `json:"id"`
	Images       []Images     `json:"images"`
	Name         string       `json:"name"`
	Type         string       `json:"type"`
	URI          string       `json:"uri"`
}

type Artists struct {
	ExternalUrls ExternalUrls `json:"external_urls"`
	Href         string       `json:"href"`
	ID           string       `json:"id"`
	Name         string       `json:"name"`
	Type         string       `json:"type"`
	URI          string       `json:"uri"`
}

type ExternalIds struct {
	Isrc string `json:"isrc"`
}

type Item struct {
	Album            Album        `json:"album"`
	Artists          []Artists    `json:"artists"`
	AvailableMarkets []string     `json:"available_markets"`
	DiscNumber       int          `json:"disc_number"`
	DurationMs       int          `json:"duration_ms"`
	Explicit         bool         `json:"explicit"`
	ExternalIds      ExternalIds  `json:"external_ids"`
	ExternalUrls     ExternalUrls `json:"external_urls"`
	Href             string       `json:"href"`
	ID               string       `json:"id"`
	Name             string       `json:"name"`
	Popularity       int          `json:"popularity"`
	PreviewURL       string       `json:"preview_url"`
	TrackNumber      int          `json:"track_number"`
	Type             string       `json:"type"`
	URI              string       `json:"uri"`
}
