package links

import "time"

func ToLinks(linkModels []LinkModel) []LinkResponse {
	links := make([]LinkResponse, 0, len(linkModels))
	for _, model := range linkModels {
		links = append(links, ToLink(model))
	}
	return links
}

func ToLink(model LinkModel) LinkResponse {
	return NewLink(model.Id, model.Uid, model.Name, model.OriginalUrl, model.Alias, model.LifetimeSec, model.CreatedAt)
}

func NewLink(id int64, uid string, name string, originalUrl string, alias string, lifetime int, creationDate time.Time) LinkResponse {
	return LinkResponse{
		Id:           uid,
		Alias:        alias,
		OriginalUrl:  originalUrl,
		Name:         name,
		Lifetime:     lifetime,
		CreationDate: creationDate,
	}
}
