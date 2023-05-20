package types

// Track  object response for track
type Track struct {
	ID     string `json:"id"`
	URL    string `json:"href"`
	Images []struct {
		URL string `json:"url"`
	} `json:"images"`
	Name    string   `json:"name"`
	Artists []Artist `json:"artists"`
}

type Artist struct {
	Name string `json:"name"`
}

// MeTrackResponse  /me/track response
type MeTrackResponse struct {
	Href string `json:"href"`

	Items []struct {
		Track Track `json:"track"`
	} `json:"items"`

	Previous string `json:"previous"`
	Next     string `json:"next"`
	Offset   string `json:"offset"`
	Total    string `json:"total"`
}
