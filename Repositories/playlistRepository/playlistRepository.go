package playlistRepository

import "time"





type Playlist struct {
	Name        string    `json:"name"`
	ShareAccess string    `json:"share_access"`
	OwnerID     int       `json:"owner_id"`
	CreatedOn   time.Time `json:"created_on"`
	ModifiedAt  time.Time `json:"modified_at"`
}

type MoviePlaylist struct {
	MovieID   int `json:"movie_id"`
	PlaylistID int `json:"playlist_id"`
}


type UserPlaylist struct {
	UserID      int  `json:"user_id"`
	PlaylistID  int  `json:"playlist_id"`
	WriteAccess bool `json:"write_access"`
}




