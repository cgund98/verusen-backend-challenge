// Package photo encapsulates logic relating to storing Photo entities in the datastore.
package photo

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repo is an abstraction that handles queries to the database.
type Repo struct {
	pg *pgxpool.Pool
}

// NewRepo constructs a new instance of the photos Repo.
func NewRepo(pg *pgxpool.Pool) *Repo {
	return &Repo{
		pg,
	}
}

// Create will create a new photo entity in the database
func (r *Repo) Create(ctx context.Context, photo Photo) (Photo, error) {
	// Define query
	query := "INSERT INTO photos (photoID, creationTime, data, name) " +
		"VALUES ($1, $2, $3, $4) " +
		"RETURNING photoID, creationTime, data, name"

	// Execute query
	var res Photo
	err := r.pg.QueryRow(ctx, query, photo.PhotoID, photo.CreationTime, photo.Data, photo.Name).
		Scan(&res.PhotoID, &res.CreationTime, &res.Data, &res.Name)

	if err != nil {
		return Photo{}, fmt.Errorf("pg.QueryRow: %v", err)
	}

	return res, nil
}

// ListPhotos will query all photos in the database.
func (r *Repo) ListPhotos(ctx context.Context) ([]Photo, error) {
	// Define query
	query := "SELECT photoID, creationTime, name FROM photos"

	// Execute query
	rows, err := r.pg.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("pg.Query: %v", err)
	}

	// Parse results
	res := []Photo{}
	for rows.Next() {
		photo := Photo{}
		if err := rows.Scan(&photo.PhotoID, &photo.CreationTime, &photo.Name); err != nil {
			return nil, fmt.Errorf("rows.Scan: %v", err)
		}
		res = append(res, photo)
	}

	return res, nil
}

// ListPhotosByAlbumID will query all photos that belong to a specific album.
func (r *Repo) ListPhotosByAlbumID(ctx context.Context, albumID string) ([]Photo, error) {
	// Define query
	query := `
		SELECT photos.photoID, photos.creationTime, photos.name FROM photos
		INNER JOIN album_memberships as mem
			ON photos.photoID = mem.photoID
		WHERE mem.albumID = $1`

	// Execute query
	rows, err := r.pg.Query(ctx, query, albumID)
	if err != nil {
		return nil, fmt.Errorf("pg.Query: %v", err)
	}

	// Parse results
	res := []Photo{}
	for rows.Next() {
		photo := Photo{}
		if err := rows.Scan(&photo.PhotoID, &photo.CreationTime, &photo.Name); err != nil {
			return nil, fmt.Errorf("rows.Scan: %v", err)
		}
		res = append(res, photo)
	}

	return res, nil
}

// GetByID will fetch a specific photo by its ID
func (r *Repo) GetByID(ctx context.Context, photoID string) (Photo, error) {
	// Define query
	query := "SELECT photoID, creationTime, data, name FROM photos WHERE photoID=$1"

	// Execute query
	var res Photo
	err := r.pg.QueryRow(ctx, query, photoID).
		Scan(&res.PhotoID, &res.CreationTime, &res.Data, &res.Name)

	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			err = fmt.Errorf("pg.QueryRow: %v", err)
		}

		return Photo{}, err
	}

	return res, nil
}

// Update will persist changes to an existing photo
func (r *Repo) Update(ctx context.Context, photo Photo) (Photo, error) {
	// Define query
	query := "UPDATE photos SET name=$2 " +
		"WHERE photoID=$1 " +
		"RETURNING photoID, creationTime, data, name"

	// Execute query
	var res Photo
	err := r.pg.QueryRow(ctx, query, photo.PhotoID, photo.Name).
		Scan(&res.PhotoID, &res.CreationTime, &res.Data, &res.Name)

	if err != nil {
		return Photo{}, fmt.Errorf("pg.QueryRow: %v", err)
	}

	return res, nil
}

// Delete will remove a specified photo from the database.
func (r *Repo) Delete(ctx context.Context, photoID string) error {
	// Define query
	query := "DELETE FROM photos WHERE photoID=$1"

	// Execute query
	_, err := r.pg.Query(ctx, query, photoID)

	if err != nil {
		return fmt.Errorf("pg.QueryRow: %v", err)
	}

	return nil
}
