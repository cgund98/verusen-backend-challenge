package service

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	v1 "photolib/api/v1"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/spf13/viper"
)

func (svc *Service) createPhoto(c *gin.Context) {
	// Parse request body
	var req v1.CreatePhotoRequest
	if err := c.ShouldBind(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// Parse uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		log.Printf("error: %v", err)
	}

	if file == nil {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("must specify an image file with key 'file'"))
		return
	}

	tempFile, err := os.CreateTemp(viper.GetString("files.tmp"), file.Filename)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("os.TempFile: %v", err))
	}
	defer tempFile.Close()

	c.SaveUploadedFile(file, tempFile.Name())

	imgBytes, err := io.ReadAll(tempFile)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("ioutil.ReadAll: %v", err))
	}

	// Persist photo to DB
	photo, err := svc.phtoCtrl.Create(c, imgBytes, req.Name)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, photo)
}

func (svc *Service) listPhotos(c *gin.Context) {
	albumID, _ := c.GetQuery("albumID")
	photos, err := svc.phtoCtrl.List(c, albumID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	res := v1.ListPhotosResponse{
		Photos: photos,
	}

	c.JSON(http.StatusOK, res)
}

func (svc *Service) getPhoto(c *gin.Context) {
	photoID := c.Param("photoID")
	photo, err := svc.phtoCtrl.Get(c, photoID)
	if err != nil {
		// Check error type
		if errors.Is(err, pgx.ErrNoRows) {
			c.AbortWithStatus(http.StatusNotFound)
		} else {
			c.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}

	c.JSON(http.StatusOK, photo)
}

func (svc *Service) updatePhoto(c *gin.Context) {
	photoID := c.Param("photoID")

	var req v1.UpdatePhotoRequest
	if err := c.ShouldBind(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	photo, err := svc.phtoCtrl.Update(c, photoID, req.Name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.AbortWithStatus(http.StatusNotFound)
		} else {
			c.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}

	c.JSON(http.StatusOK, photo)
}

func (svc *Service) deletePhoto(c *gin.Context) {
	photoID := c.Param("photoID")
	err := svc.phtoCtrl.Delete(c, photoID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.AbortWithStatus(http.StatusNotFound)
		} else {
			c.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}

	c.String(http.StatusOK, "Deleted photo.")
}
