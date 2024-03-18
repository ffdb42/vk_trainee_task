CREATE TABLE IF NOT EXISTS users (
  id SERIAL PRIMARY KEY,
  name VARCHAR(100) UNIQUE,
  role VARCHAR(100),
  password VARCHAR(100)
);

CREATE TABLE IF NOT EXISTS films (
    id SERIAL PRIMARY KEY,
    name VARCHAR(150) NOT NULL,
    description VARCHAR(1500),
    release_date DATE NOT NULL,
    rating integer DEFAULT 0 CHECK (rating >= 0 AND rating <= 10)
);

CREATE TABLE IF NOT EXISTS actors (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(20),
    last_name VARCHAR(20),
    sex VARCHAR(1),
    birthdate DATE NOT NULL
);

CREATE TABLE IF NOT EXISTS films_actors(
    id SERIAL PRIMARY KEY,
    film_id INTEGER REFERENCES films (id) ON UPDATE CASCADE ON DELETE CASCADE,
    actor_id INTEGER REFERENCES actors (id) ON UPDATE CASCADE ON DELETE CASCADE,
    UNIQUE (film_id, actor_id)
);
