// Package album encapsulates logic relating to storing Album entities in the datastore.
package album

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repo is an abstraction that handles queries to the database.
type Repo struct {
	pg *pgxpool.Pool
}

// NewRepo constructs a new instance of the albums Repo.
func NewRepo(pg *pgxpool.Pool) *Repo {
	return &Repo{
		pg,
	}
}

// Create will create a new album entity in the database
func (r *Repo) Create(ctx context.Context, album Album) (Album, error) {
	// Define query
	query := "INSERT INTO albums (albumID, creationTime, name) " +
		"VALUES ($1, $2, $3) " +
		"RETURNING albumID, creationTime, name"

	// Execute query
	var res Album
	err := r.pg.QueryRow(ctx, query, album.AlbumID, album.CreationTime, album.Name).
		Scan(&res.AlbumID, &res.CreationTime, &res.Name)

	if err != nil {
		return Album{}, fmt.Errorf("pg.QueryRow: %v", err)
	}

	return res, nil
}

// ListAlbums will query all albums in the database.
func (r *Repo) ListAlbums(ctx context.Context) ([]Album, error) {
	// Define query
	query := "SELECT albumID, creationTime, name FROM albums"

	// Execute query
	rows, err := r.pg.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("pg.Query: %v", err)
	}

	// Parse results
	res := []Album{}
	for rows.Next() {
		album := Album{}
		if err := rows.Scan(&album.AlbumID, &album.CreationTime, &album.Name); err != nil {
			return nil, fmt.Errorf("rows.Scan: %v", err)
		}
		res = append(res, album)
	}

	return res, nil
}

// GetByID will fetch a specific album by its ID
func (r *Repo) GetByID(ctx context.Context, albumID string) (Album, error) {
	// Define query
	query := "SELECT albumID, creationTime, name FROM albums WHERE albumID=$1"

	// Execute query
	var res Album
	err := r.pg.QueryRow(ctx, query, albumID).
		Scan(&res.AlbumID, &res.CreationTime, &res.Name)

	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			err = fmt.Errorf("pg.QueryRow: %v", err)
		}

		return Album{}, err
	}

	return res, nil
}

// Update will persist changes to an existing album
func (r *Repo) Update(ctx context.Context, album Album) (Album, error) {
	// Define query
	query := "UPDATE albums SET name=$2 " +
		"WHERE albumID=$1 " +
		"RETURNING albumID, creationTime, name"

	// Execute query
	var res Album
	err := r.pg.QueryRow(ctx, query, album.AlbumID, album.Name).
		Scan(&res.AlbumID, &res.CreationTime, &res.Name)

	if err != nil {
		return Album{}, fmt.Errorf("pg.QueryRow: %v", err)
	}

	return res, nil
}

// Delete will remove a specified album from the database.
func (r *Repo) Delete(ctx context.Context, albumID string) error {
	// Delete album memberships
	query := "DELETE FROM album_memberships WHERE albumID=$1"
	_, err := r.pg.Query(ctx, query, albumID)
	if err != nil {
		return fmt.Errorf("pg.QueryRow: %v", err)
	}

	// Define album
	query = "DELETE FROM albums WHERE albumID=$1"
	_, err = r.pg.Query(ctx, query, albumID)
	if err != nil {
		return fmt.Errorf("pg.QueryRow: %v", err)
	}

	return nil
}

// CreateMembership will create a new album memberships
func (r *Repo) CreateMembership(ctx context.Context, albumID string, photoID string) error {
	// Create membership
	query := "INSERT INTO album_memberships (albumID, photoID) " +
		"VALUES ($1, $2) "
	_, err := r.pg.Query(ctx, query, albumID, photoID)
	if err != nil {
		return fmt.Errorf("pg.QueryRow: %v", err)
	}

	return nil
}

// DeleteMemberships will delete existing album memberships
func (r *Repo) DeleteMemberships(ctx context.Context, albumID string, photoIDs []string) error {
	// Check that there is a list of IDs
	if len(photoIDs) == 0 {
		return nil
	}

	// Delete memberships
	insertExp := []string{}
	for i := range photoIDs {
		insertExp = append(insertExp, fmt.Sprintf("$%d", i+2))
	}
	args := []any{albumID}
	for _, photoID := range photoIDs {
		args = append(args, photoID)
	}

	query := fmt.Sprintf("DELETE FROM album_memberships WHERE albumID = $1 AND photoID IN (%s)",
		strings.Join(insertExp, ", "))
	_, err := r.pg.Query(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("pg.QueryRow: %v", err)
	}

	return nil
}
