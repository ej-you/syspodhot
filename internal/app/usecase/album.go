package usecase

import (
	"context"
	"uuid"

	"syspodhot/internal/app/entity"
	repodb "syspodhot/internal/app/repo/db"
)

// AlbumUC represents an album usecase.
type AlbumUC struct {
	albumRepoDB *repodb.AlbumRepoDB
}

// NewAlbumUC returns a new instance of AlbumUC.
func NewAlbumUC(albumRepoDB *repodb.AlbumRepoDB) *AlbumUC {
	return &AlbumUC{
		albumRepoDB: albumRepoDB,
	}
}

func (u *AlbumUC) Create(album *entity.Album) error {
	album.ID = uuid.New()
	return u.albumRepoDB.Create(context.Background(), album) // err OR nil
}

func (u *AlbumUC) Get(id uuid.UUID) (*entity.Album, error) {
	album, err := u.albumRepoDB.GetByID(context.Background(), id)
	if err != nil {
		return nil, err
	}
	return album, nil
}

func (u *AlbumUC) GetAll() ([]*entity.Album, error) {
	albums, err := u.albumRepoDB.GetAll(context.Background())
	if err != nil {
		return nil, err
	}
	return albums, nil
}

func (u *AlbumUC) Update(album *entity.Album) error {
	return u.albumRepoDB.UpdateFull(context.Background(), album) // err OR nil
}

func (u *AlbumUC) Delete(id uuid.UUID) error {
	return u.albumRepoDB.DeleteByID(context.Background(), id) // err OR nil
}
