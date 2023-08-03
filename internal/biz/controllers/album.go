package controllers

import (
	"context"
	"fmt"

	v1 "photolib/api/v1"
	"photolib/internal/data/album"
	"time"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/gofrs/uuid"
)

// AlbumController performs business-logic operations relating to Album entities.
type AlbumController struct {
	alRepo *album.Repo
	phCtrl *PhotoController
}

// NewAlbumController will construct a new instance of AlbumController
func NewAlbumController(alRepo *album.Repo, phCtrl *PhotoController) *AlbumController {
	return &AlbumController{
		alRepo,
		phCtrl,
	}
}

// Create a new album entity
func (c *AlbumController) Create(ctx context.Context, name string) (album.Album, error) {
	// Create unique ID
	albumID, err := uuid.NewV4()
	if err != nil {
		return album.Album{}, fmt.Errorf("unable to create new uuID: %v", err)
	}

	alb := album.Album{
		AlbumID:      albumID.String(),
		CreationTime: time.Now(),
		Name:         name,
	}

	// Persist entity
	return c.alRepo.Create(ctx, alb)
}

// List all albums
func (c *AlbumController) List(ctx context.Context) ([]album.Album, error) {
	return c.alRepo.ListAlbums(ctx)
}

// Get a specific album
func (c *AlbumController) Get(ctx context.Context, albumID string) (v1.GetAlbumResponse, error) {
	alb, err := c.alRepo.GetByID(ctx, albumID)
	if err != nil {
		return v1.GetAlbumResponse{}, err
	}

	photos, err := c.phCtrl.List(ctx, alb.AlbumID)
	if err != nil {
		return v1.GetAlbumResponse{}, err
	}

	res := v1.GetAlbumResponse{
		Album:  alb,
		Photos: photos,
	}

	return res, nil
}

// Update a specific album
func (c *AlbumController) Update(ctx context.Context, albumID, name string,
	photoIDs []string) (v1.UpdateAlbumResponse, error) {
	// Verify that the album exists
	al, err := c.Get(ctx, albumID)
	if err != nil {
		return al, err
	}

	// Fetch a list of existing photo memberships
	photos, err := c.phCtrl.List(ctx, al.Album.AlbumID)
	if err != nil {
		return al, fmt.Errorf("phCtrl.List: %v", err)
	}

	// Create sets of photo IDs
	requestIDs := mapset.NewSet(photoIDs...)

	existingIDs := mapset.NewSet[string]()
	for _, photo := range photos {
		existingIDs.Add(photo.PhotoID)
	}

	// Add memberships
	for photoID := range requestIDs.Difference(existingIDs).Iter() {
		err = c.alRepo.CreateMembership(ctx, albumID, photoID)
		if err != nil {
			return al, fmt.Errorf("alRepo.CreateMembership: %v", err)
		}
	}

	// Remove memberships
	err = c.alRepo.DeleteMemberships(ctx, albumID, existingIDs.Difference(requestIDs).ToSlice())
	if err != nil {
		return al, fmt.Errorf("alRepo.DeleteMemberships: %v", err)
	}

	// Update album fields
	alb := al.Album
	alb.Name = name

	_, err = c.alRepo.Update(ctx, alb)
	if err != nil {
		return al, fmt.Errorf("repo.Update: %v", err)
	}

	return c.Get(ctx, albumID)
}

// Delete a specific album
func (c *AlbumController) Delete(ctx context.Context, albumID string) error {
	// Verify that the album exists
	_, err := c.Get(ctx, albumID)
	if err != nil {
		return err
	}
	return c.alRepo.Delete(ctx, albumID)
}
