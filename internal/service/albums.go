package service

import (
	"errors"
	"net/http"
	v1 "photolib/api/v1"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func (svc *Service) createAlbum(c *gin.Context) {

	// Parse request body
	var req v1.CreateAlbumRequest
	if err := c.ShouldBind(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// Persist album to DB
	album, err := svc.albmCtrl.Create(c, req.Name)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, album)
}

func (svc *Service) listAlbums(c *gin.Context) {
	als, err := svc.albmCtrl.List(c)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	res := v1.ListAlbumsResponse{
		Albums: als,
	}

	c.JSON(http.StatusOK, res)
}

func (svc *Service) getAlbum(c *gin.Context) {
	albumID := c.Param("albumID")
	album, err := svc.albmCtrl.Get(c, albumID)
	if err != nil {
		// Check error type
		if errors.Is(err, pgx.ErrNoRows) {
			c.AbortWithStatus(http.StatusNotFound)
		} else {
			c.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}

	c.JSON(http.StatusOK, album)
}

func (svc *Service) updateAlbum(c *gin.Context) {
	albumID := c.Param("albumID")

	var req v1.UpdateAlbumRequest
	if err := c.ShouldBind(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	album, err := svc.albmCtrl.Update(c, albumID, req.Name, req.PhotoIDs)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.AbortWithStatus(http.StatusNotFound)
		} else {
			c.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}

	c.JSON(http.StatusOK, album)
}

func (svc *Service) deleteAlbum(c *gin.Context) {
	albumID := c.Param("albumID")
	err := svc.albmCtrl.Delete(c, albumID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.AbortWithStatus(http.StatusNotFound)
		} else {
			c.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}

	c.String(http.StatusOK, "Deleted album.")
}
