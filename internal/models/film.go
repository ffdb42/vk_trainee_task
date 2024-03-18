package models

type Film struct {
	ID          int         `json:"id"`
	Name        *string     `json:"name"`
	Description *string     `json:"description"`
	ReleaseDate *CustomDate `json:"release_date"`
	Rating      *int        `json:"rating"`
}

func (a Film) CopyWith(from *Film) *Film {
	if from.Name != nil {
		a.Name = from.Name
	}
	if from.Description != nil {
		a.Description = from.Description
	}
	if from.ReleaseDate != nil {
		a.ReleaseDate = from.ReleaseDate
	}
	if from.Rating != nil {
		a.Rating = from.Rating
	}
	return &a
}
