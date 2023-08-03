// Package v1 defines the API with all input and output types.
package v1

import (
	"photolib/internal/data/photo"
	"time"
)

// CreatePhotoRequest is the body for a CreatePhoto request.
type CreatePhotoRequest struct {
	Name string `form:"name"`
}

// CreatePhotoResponse is the body for the response to a CreatePhoto request.
type CreatePhotoResponse = photo.Photo

// PhotoSimple is a simplified representation of a photo entity.
type PhotoSimple struct {
	PhotoID      string    `json:"photoID"`
	CreationTime time.Time `json:"creationTime"`
	Name         string    `json:"name"`
}

// ListPhotosResponse is the body for the response to a ListPhotos request.
type ListPhotosResponse struct {
	Photos []PhotoSimple `json:"photos"`
}

// UpdatePhotoRequest is the body for a UpdatePhoto request.
type UpdatePhotoRequest struct {
	Name string `json:"name" binding:"required"`
}
