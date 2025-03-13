package main

import "time"

func ToLinks(linkModels []LinkModel) []Link {
	links := make([]Link, 0, len(linkModels))
	for _, model := range linkModels {
		links = append(links, ToLink(model))
	}
	return links
}

func ToLink(model LinkModel) Link {
	return NewLink(model.Id, model.Uid, model.Name, model.OriginalUrl, model.Alias, model.LifetimeSec, model.CreatedAt)
}

func NewLink(id int64, uid string, name string, originalUrl string, alias string, lifetime int, creationDate time.Time) Link {
	return Link{
		Id:           uid,
		Alias:        alias,
		OriginalUrl:  originalUrl,
		Name:         name,
		Lifetime:     lifetime,
		CreationDate: creationDate,
	}
}
