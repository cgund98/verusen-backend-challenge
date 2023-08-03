package v1

import "photolib/internal/data/album"

// CreateAlbumRequest is the body for a CreateAlbum request.
type CreateAlbumRequest struct {
	Name string `json:"name" binding:"required"`
}

// UpdateAlbumRequest is the body for a UpdateParty request.
type UpdateAlbumRequest struct {
	Name     string   `json:"name" binding:"required"`
	PhotoIDs []string `json:"photoIDs" binding:"required"`
}

// ListAlbumsResponse is the body for the response to a ListAlbums request.
type ListAlbumsResponse struct {
	Albums []album.Album `json:"albums"`
}

// GetAlbumResponse is the body for the response to a GetAlbum request.
type GetAlbumResponse struct {
	Album  album.Album   `json:"album"`
	Photos []PhotoSimple `json:"photos"`
}

// UpdateAlbumResponse is the body for the response to a UpdateAlbum request.
type UpdateAlbumResponse = GetAlbumResponse
