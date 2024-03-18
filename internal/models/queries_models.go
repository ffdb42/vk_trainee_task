package models

type FilmPost struct {
	Film       Film  `json:"film"`
	ActorsList []int `json:"actors_ids"`
}

type FilmPut struct {
	FilmPost
	RemoveActors []int `json:"remove_actors_ids"`
}

type ActorRespond struct {
	Actor *Actor  `json:"actor"`
	Films []*Film `json:"films"`
}

type FilmRespond struct {
	Film   *Film    `json:"film"`
	Actors []*Actor `json:"actors"`
}
