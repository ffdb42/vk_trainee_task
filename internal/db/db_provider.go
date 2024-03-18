package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/ffdb42/vk_trainee_task/internal/models"
	_ "github.com/lib/pq"
)

type DBProvider struct {
	db *sql.DB
}

var instance *DBProvider

func Init() {
	user := os.Getenv("POSTGRES_USER")
	pass := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	port := os.Getenv("DB_INT_PORT")
	connStr := fmt.Sprintf("host=db port= %v user=%v password=%v dbname=%v sslmode=disable", port, user, pass, dbName)

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		log.Fatalf("cannot open db: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("cannot ping db: %v", err)
	}

	instance = &DBProvider{db: db}

	log.Printf("established connection to db")
}

func Instance() *DBProvider {
	return instance
}

func (db *DBProvider) AddActor(actor *models.Actor) error {
	_, err := db.db.Exec(
		"INSERT INTO actors (first_name, last_name, sex, birthdate) values ($1, $2, $3, $4);",
		actor.FirstName,
		actor.LastName,
		actor.Sex,
		actor.Birthdate.Time,
	)
	return err
}

func (db *DBProvider) UpdateActor(actor *models.Actor) error {
	_, err := db.db.Exec(
		"UPDATE actors SET first_name = $1, last_name = $2, sex = $3, birthdate = $4 WHERE id = $5;",
		*actor.FirstName,
		*actor.LastName,
		*actor.Sex,
		actor.Birthdate.Time,
		actor.ID,
	)
	return err
}

func (db *DBProvider) GetActor(id int) (*models.Actor, error) {
	rows, err := db.db.Query("SELECT * FROM actors WHERE id = $1;", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := models.Actor{Birthdate: &models.CustomDate{}}
	if rows.Next() {
		err := rows.Scan(&res.ID, &res.FirstName, &res.LastName, &res.Sex, &res.Birthdate.Time)
		if err != nil {
			return nil, err
		}
	}
	return &res, nil
}

func (db *DBProvider) GetActors() (*[]models.Actor, error) {
	rows, err := db.db.Query("SELECT * FROM actors;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []models.Actor{}
	for rows.Next() {
		actor := models.Actor{Birthdate: &models.CustomDate{}}
		err := rows.Scan(&actor.ID, &actor.FirstName, &actor.LastName, &actor.Sex, &actor.Birthdate.Time)
		if err != nil {
			return nil, err
		}
		res = append(res, actor)
	}
	return &res, nil
}

func (db *DBProvider) DeleteActor(id int) (int64, error) {
	res, err := db.db.Exec("DELETE FROM actors WHERE id = $1;", id)
	if err != nil {
		return -1, err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return -1, err
	}
	return count, nil
}

func (db *DBProvider) AddFilm(film *models.Film) (int, error) {
	id := 0
	err := db.db.QueryRow(
		"INSERT INTO films (name, description, release_date, rating) values ($1, $2, $3, $4) RETURNING id;",
		film.Name,
		film.Description,
		film.ReleaseDate.Time,
		film.Rating,
	).Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func (db *DBProvider) UpdateFilm(film *models.Film) error {
	_, err := db.db.Exec(
		"UPDATE films SET name = $1, description = $2, release_date = $3, rating = $4 WHERE id = $5;",
		film.Name,
		*film.Description,
		film.ReleaseDate.Time,
		*film.Rating,
		film.ID,
	)
	return err
}

func (db *DBProvider) GetFilm(id int) (*models.Film, error) {
	rows, err := db.db.Query("SELECT * FROM films WHERE id = $1;", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := models.Film{ReleaseDate: &models.CustomDate{}}
	if rows.Next() {
		err := rows.Scan(&res.ID, &res.Name, &res.Description, &res.ReleaseDate.Time, &res.Rating)
		if err != nil {
			return nil, err
		}
	}
	return &res, nil
}

func (db *DBProvider) GetFilms(sortBy models.SortBy, sortOrder models.SortOrder) (*[]models.Film, error) {
	rows, err := db.db.Query(fmt.Sprintf("SELECT * FROM films ORDER BY %s %s;", sortBy, sortOrder))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []models.Film{}
	for rows.Next() {
		film := models.Film{ReleaseDate: &models.CustomDate{}}
		err := rows.Scan(&film.ID, &film.Name, &film.Description, &film.ReleaseDate.Time, &film.Rating)
		if err != nil {
			return nil, err
		}
		res = append(res, film)
	}
	return &res, nil
}

func (db *DBProvider) DeleteFilm(id int) (int64, error) {
	res, err := db.db.Exec("DELETE FROM films WHERE id = $1;", id)
	if err != nil {
		return -1, err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return -1, err
	}
	return count, nil
}

func (db *DBProvider) AddFilmsActors(actorID int, filmID int) error {
	_, err := db.db.Exec(
		"INSERT INTO films_actors (film_id, actor_id) values ($1, $2);",
		filmID,
		actorID,
	)
	return err
}

func (db *DBProvider) DeleteFilmsActors(filmID int, actorID int) error {
	_, err := db.db.Exec("DELETE FROM films_actors WHERE film_id = $1 AND actor_id = $2;", filmID, actorID)
	return err
}

func (db *DBProvider) AddUser(user *models.User) error {
	_, err := db.db.Exec(
		"INSERT INTO users (name, password, role) values ($1, $2, $3);",
		user.Name,
		user.Password,
		user.Role,
	)
	return err
}

func (db *DBProvider) GetUser(name string) (*models.User, error) {
	rows, err := db.db.Query("SELECT * FROM users WHERE name = $1;", name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := models.User{}
	if rows.Next() {
		err := rows.Scan(&res.ID, &res.Name, &res.Role, &res.Password)
		if err != nil {
			return nil, err
		}
		return &res, nil
	}
	return nil, nil
}

func (db *DBProvider) GetActorFilms(actorID int) ([]*models.Film, error) {
	rows, err := db.db.Query("SELECT films.* FROM films JOIN films_actors ON films.id = films_actors.film_id JOIN actors ON films_actors.actor_id = actors.id WHERE actors.id = $1", actorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []*models.Film{}
	for rows.Next() {
		film := models.Film{ReleaseDate: &models.CustomDate{}}
		err := rows.Scan(&film.ID, &film.Name, &film.Description, &film.ReleaseDate.Time, &film.Rating)
		if err != nil {
			return nil, err
		}
		res = append(res, &film)
	}
	return res, nil
}

func (db *DBProvider) GetFilmActors(filmID int) ([]*models.Actor, error) {
	rows, err := db.db.Query("SELECT actors.* FROM actors JOIN films_actors ON actors.id = films_actors.actor_id JOIN films ON films_actors.film_id = films.id WHERE films.id =$1", filmID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []*models.Actor{}
	for rows.Next() {
		actor := models.Actor{Birthdate: &models.CustomDate{}}
		err := rows.Scan(&actor.ID, &actor.FirstName, &actor.LastName, &actor.Sex, &actor.Birthdate.Time)
		if err != nil {
			return nil, err
		}
		res = append(res, &actor)
	}
	return res, nil
}

func (db *DBProvider) SearchForFilmByStringFragment(fragment string) ([]*models.Film, error) {
	rows, err := db.db.Query("SELECT films.* FROM films	JOIN films_actors ON films.id = films_actors.film_id JOIN actors ON films_actors.actor_id = actors.id WHERE LOWER(films.name) LIKE '%' || $1 || '%' OR LOWER(actors.first_name) LIKE '%' || $1 || '%';", fragment)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []*models.Film{}
	for rows.Next() {
		film := models.Film{ReleaseDate: &models.CustomDate{}}
		err := rows.Scan(&film.ID, &film.Name, &film.Description, &film.ReleaseDate.Time, &film.Rating)
		if err != nil {
			return nil, err
		}
		res = append(res, &film)
	}
	return res, nil
}
