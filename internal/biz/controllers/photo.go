// Package controllers encapsulates most business logic pertaining to entities.
package controllers

import (
	"context"
	"fmt"
	v1 "photolib/api/v1"
	"photolib/internal/data/photo"
	"time"

	"github.com/gofrs/uuid"
)

// PhotoController performs business-logic operations relating to Photo entities.
type PhotoController struct {
	repo *photo.Repo
}

// NewPhotoController will construct a new instance of PhotoController
func NewPhotoController(repo *photo.Repo) *PhotoController {
	return &PhotoController{
		repo,
	}
}

// Create a new photo entity
func (c *PhotoController) Create(ctx context.Context, data []byte, name string) (photo.Photo, error) {
	// Create unique ID
	photoID, err := uuid.NewV4()
	if err != nil {
		return photo.Photo{}, fmt.Errorf("unable to create new uuid: %v", err)
	}

	pho := photo.Photo{
		PhotoID:      photoID.String(),
		CreationTime: time.Now(),
		Data:         data,
		Name:         name,
	}

	// Persist entity
	return c.repo.Create(ctx, pho)
}

// List all photos
func (c *PhotoController) List(ctx context.Context, albumID string) ([]v1.PhotoSimple, error) {
	var photos []photo.Photo
	var err error

	// If albumID is specified, query by that. Otherwise, fetch all photos.
	if albumID == "" {
		photos, err = c.repo.ListPhotos(ctx)
	} else {
		photos, err = c.repo.ListPhotosByAlbumID(ctx, albumID)
	}
	if err != nil {
		return nil, err
	}

	res := []v1.PhotoSimple{}
	for _, pho := range photos {
		simple := v1.PhotoSimple{
			PhotoID:      pho.PhotoID,
			CreationTime: pho.CreationTime,
			Name:         pho.Name,
		}
		res = append(res, simple)
	}

	return res, nil
}

// Get will a specific photo
func (c *PhotoController) Get(ctx context.Context, photoID string) (photo.Photo, error) {
	return c.repo.GetByID(ctx, photoID)
}

// Update a specific photo
func (c *PhotoController) Update(ctx context.Context, photoID, name string) (photo.Photo, error) {
	// Verify that the photo exists
	pho, err := c.Get(ctx, photoID)
	if err != nil {
		return pho, err
	}
	pho.Name = name
	return c.repo.Update(ctx, pho)
}

// Delete a specific photo
func (c *PhotoController) Delete(ctx context.Context, photoID string) error {
	// Verify that the photo exists
	_, err := c.Get(ctx, photoID)
	if err != nil {
		return err
	}
	return c.repo.Delete(ctx, photoID)
}
