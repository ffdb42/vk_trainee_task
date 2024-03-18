package models

type Actor struct {
	ID        int         `json:"id"`
	FirstName *string     `json:"first_name"`
	LastName  *string     `json:"last_name"`
	Sex       *string     `json:"sex"`
	Birthdate *CustomDate `json:"birthdate"`
}

func (a Actor) CopyWith(from *Actor) *Actor {
	if from.FirstName != nil {
		a.FirstName = from.FirstName
	}
	if from.LastName != nil {
		a.LastName = from.LastName
	}
	if from.Sex != nil {
		a.Sex = from.Sex
	}
	if from.Birthdate != nil {
		a.Birthdate = from.Birthdate
	}
	return &a
}
