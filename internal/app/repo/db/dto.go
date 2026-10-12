package db

import (
	"time"
	"uuid"

	"syspodhot/internal/app/entity"
)

type album struct {
	ID         uuid.UUID `bson:"_id"`
	Photo      string    `bson:"photo"`
	Title      string    `bson:"title"`
	Type       string    `bson:"type"`
	Annotation string    `bson:"annotation"`
	ReleasedAt time.Time `bson:"released_at"`
	Songs      []song    `bson:"songs"`
}

type song struct {
	ID       uuid.UUID `bson:"_id"`
	Title    string    `bson:"title"`
	Duration int       `bson:"duration"`
}

func newAlbum(albumInst *entity.Album) *album {
	dto := &album{
		ID:         albumInst.ID,
		Photo:      albumInst.Photo,
		Title:      albumInst.Title,
		Type:       albumInst.Type,
		Annotation: albumInst.Annotation,
		ReleasedAt: albumInst.ReleasedAt,
		Songs:      make([]song, len(albumInst.Songs)),
	}
	for idx, entitySong := range albumInst.Songs {
		dto.Songs[idx] = song{
			ID:       entitySong.ID,
			Title:    entitySong.Title,
			Duration: entitySong.Duration,
		}
	}
	return dto
}

func (a *album) toEntity() *entity.Album {
	inst := &entity.Album{
		ID:         a.ID,
		Photo:      a.Photo,
		Title:      a.Title,
		Type:       a.Type,
		Annotation: a.Annotation,
		ReleasedAt: a.ReleasedAt,
		Songs:      make([]entity.Song, len(a.Songs)),
	}
	for idx, entitySong := range a.Songs {
		inst.Songs[idx] = entity.Song{
			ID:       entitySong.ID,
			Title:    entitySong.Title,
			Duration: entitySong.Duration,
		}
	}
	return inst
}

func (album) collectionName() string {
	return "albums"
}

type user struct {
	ID       uuid.UUID `bson:"_id"`
	Username string    `bson:"username"`
	Password []byte    `bson:"password"`
}

func newUser(userInst *entity.User) *user {
	return &user{
		ID:       userInst.ID,
		Username: userInst.Username,
		Password: userInst.Password,
	}
}

func (u *user) toEntity() *entity.User {
	return &entity.User{
		ID:       u.ID,
		Username: u.Username,
		Password: u.Password,
	}
}

func (user) collectionName() string {
	return "users"
}
