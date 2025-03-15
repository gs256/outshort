package links

import "time"

type ShortenRequest struct {
	Url string `json:"url"`
}

type UpsertLinkRequest struct {
	Url      string `json:"url"`
	Name     string `json:"name"`
	Alias    string `json:"alias"`
	Lifetime int    `json:"lifetime"`
}

type LinkResponse struct {
	Id           string    `json:"id"`
	Alias        string    `json:"alias"`
	OriginalUrl  string    `json:"originalUrl"`
	Name         string    `json:"name"`
	Lifetime     int       `json:"lifetime"`
	CreationDate time.Time `json:"creationDate"`
}
