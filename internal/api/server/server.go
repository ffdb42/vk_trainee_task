package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/ffdb42/vk_trainee_task/internal/constants"
	"github.com/ffdb42/vk_trainee_task/internal/db"
	"github.com/ffdb42/vk_trainee_task/internal/models"
	"github.com/ffdb42/vk_trainee_task/internal/utils"
	_ "github.com/swaggo/files"
	_ "github.com/swaggo/http-swagger"
	"golang.org/x/crypto/bcrypt"
)

type Server struct{}

// @Summary Sign up
// @Tags auth
// @Description Регистрация пользователя
// @ID sing-up
// @Accept json
// @Produce json
// @Param requestBody body models.SignUpRequest true "Пароль + юзернейм"
// @Success 200 {string} string "user signed up"
// @Failure 400 {string} string "error string"
// @Failure 500 {string} string "internal server error"
// @Router /sign-up/ [post]
func (s *Server) SignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("unexpected method"))
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("ERROR %v %v: cannot read request body: %v", r.Method, r.RequestURI, err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("cannot get request body"))
		return
	}
	var userMap map[string]interface{}
	err = json.Unmarshal(body, &userMap)
	if err != nil {
		log.Printf("ERROR %v %v: cannot unmarshal request body: %v", r.Method, r.RequestURI, err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("cannot get request body"))
		return
	}
	name, ok := userMap["name"].(string)
	if !ok {
		log.Printf("ERROR %v %v: name was not provided: %v", r.Method, r.RequestURI, userMap)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("name was not provided"))
		return
	}
	if len(name) == 0 || len(name) > 100 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("name length should be at least 1 and no more than 100 characters"))
		return
	}
	pass, ok := userMap["password"].(string)
	if !ok {
		log.Printf("ERROR %v %v: password was not provided: %v", r.Method, r.RequestURI, userMap)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("password was not provided"))
		return
	}
	if len(pass) == 0 || len(pass) > 100 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("pass length should be at least 1 and no more than 100 characters"))
		return
	}
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.MinCost)
	if err != nil {
		log.Printf("ERROR %v %v: cannot generate hash for pass: %v", r.Method, r.RequestURI, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal server error"))
		return
	}
	user := models.User{Name: name, Password: string(hashedPass), Role: constants.UserRole}
	err = db.Instance().AddUser(&user)
	if err != nil {
		log.Printf("ERROR %v %v: cannot add user to db: %v", r.Method, r.RequestURI, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal server error"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("user signed up"))
}

func (s *Server) ActorHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		id, err := utils.ParseID(r.RequestURI)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		if id > 0 {
			s.getActor(id, w, r)
		} else {
			s.getActors(w, r)
		}

	case http.MethodPost:
		s.postActor(w, r)
	case http.MethodPut:
		s.putActor(w, r)
	case http.MethodDelete:
		s.deleteActor(w, r)
	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Unexpected HTTP method"))
	}
}

func (s *Server) FilmHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		id, _ := utils.ParseID(r.RequestURI)

		if id > 0 {
			s.getFilm(id, w, r)
		} else {
			s.getFilms(w, r)
		}

	case http.MethodPost:
		s.postFilm(w, r)
	case http.MethodPut:
		s.putFilm(w, r)
	case http.MethodDelete:
		s.deleteFilm(w, r)
	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Unexpected HTTP method"))
	}
}

// @Summary Search
// @Tags search
// @Description Поиск фильмов по фрагменту из названия или фрагменту имени актера, который указан в титрах
// @ID search
// @Security BasicAuth
// @Produce json
// @Param search_by query string true "искомый фрагмент"
// @Success 200 {object} models.FilmsSearch
// @Failure 400 {string} string "error string"
// @Failure 401 {string} string "unauthtorized"
// @Failure 500 {string} string "internal server error"
// @Router /search/ [get]
func (s *Server) SearchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("not found"))
		return
	}
	query := r.URL.Query()
	search, ok := query["search_by"]
	if len(search) != 1 || !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid query param"))
		return
	}
	films, err := db.Instance().SearchForFilmByStringFragment(search[0])
	if err != nil {
		log.Printf("ERROR %v %v: cannot search for film: %v", r.Method, r.RequestURI, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal server error"))
		return
	}
	res, err := json.Marshal(map[string]any{"films": films})
	if err != nil {
		log.Printf("ERROR %v %v: cannot marshal to json: %v", r.Method, r.RequestURI, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal server error"))
	}
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

// @Summary Get actor
// @Tags actor
// @Description Поиск актера по id
// @ID get-actor
// @Security BasicAuth
// @Produce json
// @Param id path int true "id"
// @Success 200 {object} models.ActorRespond
// @Failure 400 {string} string "error string"
// @Failure 401 {string} string "unauthtorized"
// @Failure 404 {string} string "not found"
// @Failure 500 {string} string "internal server error"
// @Router /actor/{id} [get]
func (*Server) getActor(id int, w http.ResponseWriter, r *http.Request) {
	actor, err := db.Instance().GetActor(id)
	if err != nil {
		log.Printf("ERROR %v %v: cannot get value from db: %v", r.Method, r.RequestURI, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("cannot get actor"))
		return
	}
	respond := models.ActorRespond{Actor: actor, Films: []*models.Film{}}
	films, err := db.Instance().GetActorFilms(id)
	if err != nil {
		log.Printf("ERROR %v %v: cannot get films list from db: %v", r.Method, r.RequestURI, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("cannot get actor"))
		return
	}
	respond.Films = films
	res, err := json.Marshal(respond)
	if err != nil {
		log.Printf("ERROR %v %v: cannot marshal json: %v", r.Method, r.RequestURI, err)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("actor not found"))
	}
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

// @Summary Get actors
// @Tags actor
// @Description Получения списка актеров
// @ID get-actors
// @Security BasicAuth
// @Produce json
// @Success 200 {object} models.GetActors
// @Failure 400 {string} string "error string"
// @Failure 401 {string} string "unauthtorized"
// @Failure 404 {string} string "not found"
// @Failure 500 {string} string "internal server error"
// @Router /actor/ [get]
func (*Server) getActors(w http.ResponseWriter, r *http.Request) {
	actors, err := db.Instance().GetActors()
	if err != nil {
		log.Printf("ERROR %v %v: cannot get value from db: %v", r.Method, r.RequestURI, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("cannot get actors"))
		return
	}
	respond := []models.ActorRespond{}
	for _, v := range *actors {
		films, err := db.Instance().GetActorFilms(v.ID)
		if err != nil {
			log.Printf("ERROR %v %v: cannot get films list from db: %v", r.Method, r.RequestURI, err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("cannot get actor"))
			return
		}
		respond = append(respond, models.ActorRespond{Actor: &v, Films: films})
	}
	res, err := json.Marshal(map[string]any{"actors": respond})
	if err != nil {
		log.Printf("ERROR %v %v: cannot marshal json: %v", r.Method, r.RequestURI, err)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("actor not found"))
	}
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

// @Summary Add actor
// @Tags actor
// @Description Создание записи об актере
// @ID post-actor
// @Security BasicAuth
// @Accept json
// @Produce json
// @Param requestBody body models.ActorPost true "Информация об актере"
// @Success 200 {string} string "actor added"
// @Failure 400 {string} string "error string"
// @Failure 401 {string} string "unauthtorized"
// @Failure 403 {string} string "forbidden"
// @Failure 500 {string} string "internal server error"
// @Router /actor/ [post]
func (*Server) postActor(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("ERROR %v %v: cannot read request body: %v", r.Method, r.RequestURI, err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("cannot get request body"))
		return
	}
	var actor models.Actor
	err = json.Unmarshal(body, &actor)
	if err != nil {
		log.Printf("ERROR %v %v: cannot unmarshal request body: %v", r.Method, r.RequestURI, err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("cannot get request body"))
		return
	}
	if actor.FirstName == nil || len(*actor.FirstName) == 0 || len(*actor.FirstName) > 20 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("first name length should be at least 1 and no more than 20 characters"))
		return
	}
	if actor.LastName == nil || len(*actor.LastName) == 0 || len(*actor.LastName) > 20 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("last name length should be at least 1 and no more than 20 characters"))
		return
	}
	if actor.Sex == nil || (*actor.Sex != "m" && *actor.Sex != "f") {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("sex should be 'm' or 'f'"))
		return
	}
	err = db.Instance().AddActor(&actor)
	if err != nil {
		log.Printf("ERROR %v %v: cannot add value to db: %v", r.Method, r.RequestURI, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal server error"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("actor added"))
}

// @Summary Update actor
// @Tags actor
// @Description Изменение записи об актере
// @ID put-actor
// @Security BasicAuth
// @Accept json
// @Param id path int true "id"
// @Param requestBody body models.ActorPost true "Информация об актере"
// @Success 200 {string} string "actor updated"
// @Failure 400 {string} string "error string"
// @Failure 401 {string} string "unauthtorized"
// @Failure 403 {string} string "forbidden"
// @Failure 500 {string} string "internal server error"
// @Router /actor/{id} [put]
func (*Server) putActor(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseID(r.RequestURI)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	if id < 1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid id"))
		return
	}
	oldActor, err := db.Instance().GetActor(id)
	if err != nil {
		log.Printf("ERROR %v %v: cannot get value from db: %v", r.Method, r.RequestURI, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal server error"))
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("ERROR %v %v: cannot read request body: %v", r.Method, r.RequestURI, err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("cannot get request body"))
		return
	}
	var actor models.Actor
	err = json.Unmarshal(body, &actor)
	if err != nil {
		log.Printf("ERROR %v %v: cannot unmarshal request body: %v", r.Method, r.RequestURI, err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("cannot get request body"))
		return
	}
	if actor.FirstName == nil || len(*actor.FirstName) == 0 || len(*actor.FirstName) > 20 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("first name length should be at least 1 and no more than 20 characters"))
		return
	}
	if actor.LastName == nil || len(*actor.LastName) == 0 || len(*actor.LastName) > 20 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("last name length should be at least 1 and no more than 20 characters"))
		return
	}
	if actor.Sex == nil || (*actor.Sex != "m" && *actor.Sex != "f") {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("sex should be 'm' or 'f'"))
		return
	}
	updatedActor := oldActor.CopyWith(&actor)
	err = db.Instance().UpdateActor(updatedActor)
	if err != nil {
		log.Printf("ERROR %v %v: cannot update actor: %v", r.Method, r.RequestURI, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal server error"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("actor updated"))
}

// @Summary Delete actor
// @Tags actor
// @Description Удаление актера по id
// @ID delete-actor
// @Security BasicAuth
// @Param id path int true "id"
// @Success 200 {string} string "actor deleted"
// @Failure 400 {string} string "error string"
// @Failure 401 {string} string "unauthtorized"
// @Failure 404 {string} string "not found"
// @Failure 500 {string} string "internal server error"
// @Router /actor/{id} [delete]
func (*Server) deleteActor(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseID(r.RequestURI)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	if id < 1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid id"))
		return
	}
	n, err := db.Instance().DeleteActor(id)
	if err != nil {
		log.Printf("ERROR %v %v: cannot delete value from db: %v", r.Method, r.RequestURI, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal server error"))
		return
	}
	if n == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("not found"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("actor deleted"))
}

// @Summary Get film
// @Tags film
// @Description Поиск фильма по id
// @ID get-film
// @Security BasicAuth
// @Produce json
// @Param id path int true "id"
// @Success 200 {object} models.FilmRespond
// @Failure 400 {string} string "error string"
// @Failure 401 {string} string "unauthtorized"
// @Failure 404 {string} string "not found"
// @Failure 500 {string} string "internal server error"
// @Router /film/{id} [get]
func (*Server) getFilm(id int, w http.ResponseWriter, r *http.Request) {
	film, err := db.Instance().GetFilm(id)
	if err != nil {
		log.Printf("ERROR %v %v: cannot get value from db: %v", r.Method, r.RequestURI, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("cannot get actor"))
		return
	}
	respond := models.FilmRespond{Film: film, Actors: []*models.Actor{}}
	actors, err := db.Instance().GetFilmActors(id)
	if err != nil {
		log.Printf("ERROR %v %v: cannot get actors list from db: %v", r.Method, r.RequestURI, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("cannot get film"))
		return
	}
	respond.Actors = actors
	res, err := json.Marshal(respond)
	if err != nil {
		log.Printf("ERROR %v %v: cannot marshal json: %v", r.Method, r.RequestURI, err)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("actor not found"))
	}
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

// @Summary Get films
// @Tags film
// @Description Получения списка фильмов
// @ID get-films
// @Security BasicAuth
// @Produce json
// @Success 200 {object} models.GetFilms
// @Failure 400 {string} string "error string"
// @Failure 401 {string} string "unauthtorized"
// @Failure 404 {string} string "not found"
// @Failure 500 {string} string "internal server error"
// @Router /film/ [get]
func (*Server) getFilms(w http.ResponseWriter, r *http.Request) {
	var sortBy models.SortBy
	var sortOrder models.SortOrder
	query := r.URL.Query()
	querySortBy, ok := query["sort_by"]
	if !ok || len(querySortBy) != 1 {
		sortBy = constants.SortByRating
	} else {
		switch querySortBy[0] {
		case "name":
			sortBy = constants.SortByName
		case "release_date":
			sortBy = constants.SortByReleaseDate
		case "rating":
			fallthrough
		default:
			sortBy = constants.SortByRating
		}
	}
	querySortOrder, ok := query["sort_order"]
	if !ok || len(querySortOrder) != 1 {
		sortOrder = constants.SortDesc
	} else {
		switch querySortOrder[0] {
		case "ASC":
			sortOrder = constants.SortAsc
		case "DESC":
			fallthrough
		default:
			log.Printf("24 %v", querySortOrder)
			sortOrder = constants.SortDesc
		}
	}
	films, err := db.Instance().GetFilms(sortBy, sortOrder)
	if err != nil {
		log.Printf("ERROR %v %v: cannot get value from db: %v", r.Method, r.RequestURI, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("cannot get films"))
		return
	}
	respond := []models.FilmRespond{}
	for _, v := range *films {
		actors, err := db.Instance().GetFilmActors(v.ID)
		if err != nil {
			log.Printf("ERROR %v %v: cannot get actors list from db: %v", r.Method, r.RequestURI, err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("cannot get film"))
			return
		}
		respond = append(respond, models.FilmRespond{Film: &v, Actors: actors})
	}
	res, err := json.Marshal(respond)
	if err != nil {
		log.Printf("ERROR %v %v: cannot marshal json: %v", r.Method, r.RequestURI, err)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("actor not found"))
	}
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

// @Summary Add film
// @Tags film
// @Description Создание записи об фильме
// @ID post-film
// @Security BasicAuth
// @Accept json
// @Produce json
// @Param requestBody body models.FilmPostDoc true "Информация о фильме`"
// @Success 200 {string} string "film added"
// @Failure 400 {string} string "error string"
// @Failure 401 {string} string "unauthtorized"
// @Failure 403 {string} string "forbidden"
// @Failure 500 {string} string "internal server error"
// @Router /film/ [post]
func (*Server) postFilm(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("ERROR %v %v: cannot read request body: %v", r.Method, r.RequestURI, err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("cannot get request body"))
		return
	}
	var filmPost models.FilmPost
	err = json.Unmarshal(body, &filmPost)
	if err != nil {
		log.Printf("ERROR %v %v: cannot unmarshal request body: %v", r.Method, r.RequestURI, err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("cannot get request body"))
		return
	}
	if filmPost.Film.Name != nil && (len(*filmPost.Film.Name) < 1 || len(*filmPost.Film.Name) > 150) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("the length of the film name must be at least 1 and no more than 150 characters"))
		return
	}
	if filmPost.Film.Description != nil && len(*filmPost.Film.Description) > 1500 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("film's description len should not exceed 1500 symbols"))
		return
	}
	if filmPost.Film.Rating != nil && (*filmPost.Film.Rating < 0 || *filmPost.Film.Rating > 10) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("film's rating should be from 0 to 10"))
		return
	}
	filmID, err := db.Instance().AddFilm(&filmPost.Film)
	if err != nil {
		log.Printf("ERROR %v %v: cannot add value to db: %v", r.Method, r.RequestURI, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal server error"))
		return
	}
	for _, actorID := range filmPost.ActorsList {
		err := db.Instance().AddFilmsActors(actorID, filmID)
		if err != nil {
			log.Printf("ERROR %v %v: cannot add FilmActor: %v", r.Method, r.RequestURI, err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("cannot add actor with id %v", actorID)))
			return
		}
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("film added"))
}

// @Summary Update film
// @Tags film
// @Description Изменение записи о фильме
// @ID put-film
// @Security BasicAuth
// @Accept json
// @Param id path int true "id"
// @Param requestBody body models.FilmPostDoc true "Информация о фильме"
// @Success 200 {string} string "film updated"
// @Failure 400 {string} string "error string"
// @Failure 401 {string} string "unauthtorized"
// @Failure 403 {string} string "forbidden"
// @Failure 500 {string} string "internal server error"
// @Router /film/{id} [put]
func (*Server) putFilm(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseID(r.RequestURI)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	if id < 1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid id"))
		return
	}
	oldFilm, err := db.Instance().GetFilm(id)
	if err != nil {
		log.Printf("ERROR %v %v: cannot get value from db: %v", r.Method, r.RequestURI, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal server error"))
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("ERROR %v %v: cannot read request body: %v", r.Method, r.RequestURI, err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("cannot get request body"))
		return
	}
	var filmPost models.FilmPut
	err = json.Unmarshal(body, &filmPost)
	if err != nil {
		log.Printf("ERROR %v %v: cannot unmarshal request body: %v", r.Method, r.RequestURI, err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("cannot get request body"))
		return
	}
	if filmPost.Film.Name != nil && (len(*filmPost.Film.Name) < 1 || len(*filmPost.Film.Name) > 150) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("the length of the film name must be at least 1 and no more than 150 characters"))
		return
	}
	if filmPost.Film.Description != nil && len(*filmPost.Film.Description) > 1500 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("film's description len should not exceed 1500 symbols"))
		return
	}
	if filmPost.Film.Rating != nil && (*filmPost.Film.Rating < 0 || *filmPost.Film.Rating > 10) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("film's rating should be from 0 to 10"))
		return
	}
	updatedFilm := oldFilm.CopyWith(&filmPost.Film)
	err = db.Instance().UpdateFilm(updatedFilm)
	if err != nil {
		log.Printf("ERROR %v %v: cannot update actor: %v", r.Method, r.RequestURI, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal server error"))
		return
	}
	for _, actorID := range filmPost.ActorsList {
		err := db.Instance().AddFilmsActors(actorID, updatedFilm.ID)
		if err != nil {
			log.Printf("ERROR %v %v: cannot add FilmActor: %v", r.Method, r.RequestURI, err)
		}
	}
	for _, actorID := range filmPost.RemoveActors {
		err := db.Instance().DeleteFilmsActors(id, actorID)
		if err != nil {
			log.Printf("ERROR %v %v: cannot delete FilmActor: %v", r.Method, r.RequestURI, err)
		}
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("film updated"))
}

// @Summary Delete film
// @Tags film
// @Description Удаление фильма по id
// @ID delete-film
// @Security BasicAuth
// @Param id path int true "id"
// @Success 200 {string} string "film deleted"
// @Failure 400 {string} string "error string"
// @Failure 401 {string} string "unauthtorized"
// @Failure 404 {string} string "not found"
// @Failure 500 {string} string "internal server error"
// @Router /film/{id} [delete]
func (*Server) deleteFilm(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseID(r.RequestURI)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	if id < 1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid id"))
		return
	}
	n, err := db.Instance().DeleteFilm(id)
	if err != nil {
		log.Printf("ERROR %v %v: cannot delete value from db: %v", r.Method, r.RequestURI, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal server error"))
		return
	}
	if n == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("not found"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("film deleted"))
}
