CREATE DATABASE moviemanagercached;

\c moviemanagercached;

CREATE TYPE access AS ENUM (
  'PUBLIC_READ',
  'PUBLIC_WRITE',
  'RESTRICTED'
);

CREATE TABLE Users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(255),
    email VARCHAR(50) UNIQUE NOT NULL,
    created_on BIGINT NOT NULL,
    last_login BIGINT NOT NULL
);

CREATE TABLE Playlist (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL,
    shareAccess access NOT NULL,
    ownerId INTEGER NOT NULL,
    createdOn BIGINT NOT NULL,
    modifiedAt BIGINT NOT NULL,
    CONSTRAINT Playlist_owner_id_fkey FOREIGN KEY (ownerId)
      REFERENCES Users(id) MATCH SIMPLE
      ON UPDATE NO ACTION ON DELETE NO ACTION
);

CREATE TABLE Movie (
    Title VARCHAR(100) NOT NULL,
    Year VARCHAR(4) NOT NULL,
    Rated VARCHAR(10),
    Released VARCHAR(255),
    Runtime VARCHAR(10),
    Genre VARCHAR(100),
    Director VARCHAR(100),
    Writer VARCHAR(100),
    Actors VARCHAR(255),
    Plot TEXT,
    Language VARCHAR(100),
    Country VARCHAR(100),
    Awards VARCHAR(255),
    Poster VARCHAR(255),
    Metascore VARCHAR(10),
    ImdbRating VARCHAR(10),
    ImdbVotes VARCHAR(20),
    Imdbid VARCHAR(20) PRIMARY KEY UNIQUE NOT NULL,
    Type VARCHAR(20),
    Dvd VARCHAR(10),
    BoxOffice VARCHAR(20),
    Production VARCHAR(100),
    Website VARCHAR(255),
    Response VARCHAR(20)
);



CREATE TABLE Rating (
    id SERIAL PRIMARY KEY,
    movieId VARCHAR(20) NOT NULL,
    source VARCHAR(100),
    value VARCHAR(10),
    CONSTRAINT fk_movie
        FOREIGN KEY (movie_id)
        REFERENCES Movie(imdb_id)
        ON UPDATE CASCADE
        ON DELETE CASCADE
);


CREATE TABLE moviePlaylist (
  movie_id VARCHAR(20) NOT NULL,
  playlist_id INTEGER NOT NULL,
  PRIMARY KEY (movie_id, playlist_id),
  CONSTRAINT fk_movie
    FOREIGN KEY (movie_id)
    REFERENCES Movie(imdb_id) 
    ON UPDATE NO ACTION 
    ON DELETE CASCADE,
  CONSTRAINT fk_playlist
    FOREIGN KEY (playlist_id)
    REFERENCES Playlist(id) 
    ON UPDATE NO ACTION 
    ON DELETE CASCADE
);

CREATE TABLE userPlaylist (
  user_id INTEGER NOT NULL,
  playlist_id INTEGER NOT NULL,
  writeAccess BOOLEAN DEFAULT false,
  -- add primary key 
  PRIMARY KEY (user_id,playlist_id)
  --
  CONSTRAINT userPlaylist_user_id_fkey FOREIGN KEY (user_id)
    REFERENCES Users(id) MATCH SIMPLE
    ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT userPlaylist_playlist_id_fkey FOREIGN KEY (playlist_id)
    REFERENCES Playlist(id) MATCH SIMPLE
    ON UPDATE NO ACTION ON DELETE CASCADE
);

-- DROP DATABASE moviemanager;



select movie.*,ratings.id as ratingId,ratings.source as source,ratings.value as value from movie 
join
ratings
on 
ratings.movieId = movie.Imdbid


-- tt0499549

SELECT row_to_json(movieratings)
FROM 
(SELECT m.*
, COALESCE((SELECT json_agg(r) FROM ratings r WHERE r.movieId = 'tt0499549')
, '[]') AS ratings movie m WHERE m.imdbID = 'tt0499549') 
AS 
movieratings;




