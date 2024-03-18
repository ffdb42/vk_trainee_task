package models

type SignUpRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type GetActors struct {
	Actors []ActorRespond `json:"actors"`
}

type GetFilms struct {
	Actors []FilmRespond `json:"actors"`
}

type FilmsSearch struct {
	Films []*Film `json:"films"`
}

type ActorPost struct {
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	Sex       *string `json:"sex"`
	Birthdate *string `json:"birthdate"`
}

type FilmPostDoc struct {
	Film       FilmDoc `json:"film"`
	ActorsList []int   `json:"actors_ids"`
}

type FilmDoc struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	ReleaseDate *string `json:"release_date"`
	Rating      *int    `json:"rating"`
}
