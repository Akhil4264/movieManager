CREATE DATABASE moviemanagercached;

\c moviemanagercached;

CREATE TYPE access AS ENUM (
  'PUBLIC_READ',
  'PUBLIC_WRITE',
  'RESTRICTED'
);

CREATE TABLE "User" (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(50),
    email VARCHAR(50) UNIQUE NOT NULL,
    created_on TIMESTAMP NOT NULL,
    last_login TIMESTAMP NOT NULL
);

CREATE TABLE Playlist (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL,
    shareAccess access NOT NULL,
    owner_id INTEGER NOT NULL,
    created_on TIMESTAMP NOT NULL,
    modified_at TIMESTAMP NOT NULL,
    CONSTRAINT Playlist_owner_id_fkey FOREIGN KEY (owner_id)
      REFERENCES "User"(id) MATCH SIMPLE
      ON UPDATE NO ACTION ON DELETE NO ACTION
);


CREATE TABLE Movie (
    id SERIAL PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    year VARCHAR(4) NOT NULL,
    rated VARCHAR(10),
    released DATE,
    runtime VARCHAR(10),
    genre VARCHAR(100),
    director VARCHAR(100),
    writer VARCHAR(100),
    actors VARCHAR(255),
    plot TEXT,
    language VARCHAR(100),
    country VARCHAR(100),
    awards VARCHAR(255),
    poster VARCHAR(255),
    metascore INTEGER,
    imdb_rating VARCHAR(10),
    imdb_votes VARCHAR(20),
    imdb_id VARCHAR(20) UNIQUE,
    type VARCHAR(20),
    dvd VARCHAR(10),
    box_office VARCHAR(20),
    production VARCHAR(100),
    website VARCHAR(255),
    response BOOLEAN
);

CREATE TABLE Ratings (
    id SERIAL PRIMARY KEY,
    movie_id INTEGER NOT NULL,
    source VARCHAR(100),
    value VARCHAR(10),
    CONSTRAINT fk_movie
        FOREIGN KEY (movie_id)
        REFERENCES Movie(imdb_id)
        ON UPDATE CASCADE
        ON DELETE CASCADE
);


CREATE TABLE moviePlaylist (
  movie_id INTEGER NOT NULL,
  playlist_id INTEGER NOT NULL,
  CONSTRAINT moviePlaylist_movie_id_fkey FOREIGN KEY (movie_id)
    REFERENCES Movie(imdb_id) MATCH SIMPLE
    ON UPDATE NO ACTION ON DELETE CASCADE
  CONSTRAINT moviePlaylist_playlist_id_fkey FOREIGN KEY (playlist_id)
    REFERENCES Playlist(id) MATCH SIMPLE
    ON UPDATE NO ACTION ON DELETE CASCADE
);

CREATE TABLE userPlaylist (
  user_id INTEGER NOT NULL,
  playlist_id INTEGER NOT NULL,
  writeAccess BOOLEAN DEFAULT false,
  CONSTRAINT userPlaylist_user_id_fkey FOREIGN KEY (user_id)
    REFERENCES "User"(id) MATCH SIMPLE
    ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT userPlaylist_playlist_id_fkey FOREIGN KEY (playlist_id)
    REFERENCES Playlist(id) MATCH SIMPLE
    ON UPDATE NO ACTION ON DELETE CASCADE
);

-- DROP DATABASE moviemanager;
