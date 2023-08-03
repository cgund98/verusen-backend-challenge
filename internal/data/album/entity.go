package album

import "time"

// Album is the business-logic representation of an Album entity.
type Album struct {
	AlbumID      string    `json:"albumID"`
	CreationTime time.Time `json:"creationTime"`
	Name         string    `json:"name"`
}
