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
	Streaming  streaming `bson:"streaming"`
	Songs      []song    `bson:"songs"`
}

type streaming struct {
	VKLink     string `json:"vk"`
	YandexLink string `json:"yandex"`
}

type song struct {
	Title    string `bson:"title"`
	Duration int    `bson:"duration"`
}

func newAlbum(albumInst *entity.Album) *album {
	dto := &album{
		ID:         albumInst.ID,
		Photo:      albumInst.Photo,
		Title:      albumInst.Title,
		Type:       albumInst.Type,
		Annotation: albumInst.Annotation,
		ReleasedAt: albumInst.ReleasedAt,
		Streaming: streaming{
			VKLink:     albumInst.Streaming.VKLink,
			YandexLink: albumInst.Streaming.YandexLink,
		},
		Songs: make([]song, len(albumInst.Songs)),
	}
	for idx, entitySong := range albumInst.Songs {
		dto.Songs[idx] = song{
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
		Streaming: entity.Streaming{
			VKLink:     a.Streaming.VKLink,
			YandexLink: a.Streaming.YandexLink,
		},
		Songs: make([]entity.Song, len(a.Songs)),
	}
	for idx, entitySong := range a.Songs {
		inst.Songs[idx] = entity.Song{
			Title:    entitySong.Title,
			Duration: entitySong.Duration,
		}
	}
	return inst
}

func (album) collectionName() string {
	return "albums"
}

type concert struct {
	ID         uuid.UUID `bson:"_id"`
	Title      string    `bson:"title"`
	Place      string    `bson:"place"`
	Datetime   time.Time `bson:"datetime"`
	TicketLink string    `bson:"ticket_link"`

	Photos []string `json:"photos"`
	Video  string   `json:"video"`
}

func newConcert(concertInst *entity.Concert) *concert {
	return &concert{
		ID:         concertInst.ID,
		Title:      concertInst.Title,
		Place:      concertInst.Place,
		Datetime:   concertInst.Datetime,
		TicketLink: concertInst.TicketLink,
		Photos:     concertInst.Photos,
		Video:      concertInst.Video,
	}
}

func (c *concert) toEntity() *entity.Concert {
	return &entity.Concert{
		ID:         c.ID,
		Title:      c.Title,
		Place:      c.Place,
		Datetime:   c.Datetime,
		TicketLink: c.TicketLink,
		Photos:     c.Photos,
		Video:      c.Video,
	}
}

func (concert) collectionName() string {
	return "concerts"
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
