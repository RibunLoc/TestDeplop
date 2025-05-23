// models/user.go
package models


type Music struct {
	ID        uint      `gorm:"primaryKey"`
	Name	  string    `gorm:"not null"`
	Artist    string    `gorm:"not null"`
	// Lyrics    string    `gorm:"not null"` // Lyrics of the song
	AccessURL string    `gorm:"not null"` // URL to access the song
	// ImageURL  string    `gorm:"not null"` // URL to access the image of the song
	Genre	string    `gorm:"not null"` // Genre of the song
	PlayCount uint      `gorm:"default:0"` // Number of times the song has been played
}

type UserPlaylist struct {
	PlaylistID        uint      `gorm:"primaryKey"`
	UserID    int      `gorm:"not null"`
	Topic	 string    `gorm:"not null"`
}

type Playlist struct {
	PlaylistID        uint      `gorm:"not null"`
	MusicID    uint      `gorm:"not null"`
}