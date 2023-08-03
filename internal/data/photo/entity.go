package photo

import "time"

// Photo is the business-logic representation of an Photo entity.
type Photo struct {
	PhotoID      string    `json:"photoID"`
	CreationTime time.Time `json:"creationTime"`
	Data         []byte    `json:"data"`
	Name         string    `json:"name"`
}
