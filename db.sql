CREATE DATABASE moviemanager;

\c moviemanager;

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

CREATE TABLE moviePlaylist (
  movie_id INTEGER NOT NULL,
  playlist_id INTEGER NOT NULL,
  PRIMARY KEY (movie_id,playlist_id),
  CONSTRAINT moviePlaylist_playlist_id_fkey FOREIGN KEY (playlist_id)
    REFERENCES Playlist(id) MATCH SIMPLE
    ON UPDATE NO ACTION ON DELETE CASCADE
);

CREATE TABLE userPlaylist (
  user_id INTEGER NOT NULL,
  playlist_id INTEGER NOT NULL,
  writeAccess BOOLEAN DEFAULT false,
  PRIMARY KEY (user_id,playlist_id),
  CONSTRAINT userPlaylist_user_id_fkey FOREIGN KEY (user_id)
    REFERENCES "User"(id) MATCH SIMPLE
    ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT userPlaylist_playlist_id_fkey FOREIGN KEY (playlist_id)
    REFERENCES Playlist(id) MATCH SIMPLE
    ON UPDATE NO ACTION ON DELETE CASCADE
);

-- DROP DATABASE moviemanager;
