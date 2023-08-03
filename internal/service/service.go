package service

import (
	"photolib/internal/biz/controllers"
	"photolib/internal/data/album"
	"photolib/internal/data/photo"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

const apiPrefix string = "/api/v1"

type Service struct {
	g        *gin.Engine
	albmCtrl *controllers.AlbumController
	phtoCtrl *controllers.PhotoController
}

// NewService constructs a new instance of Service
func NewService(g *gin.Engine, pg *pgxpool.Pool) *Service {
	// Initialize controllers
	phtoCtrl := controllers.NewPhotoController(photo.NewRepo(pg))
	albmCtrl := controllers.NewAlbumController(album.NewRepo(pg), phtoCtrl)

	// MIDdleware
	g.Use(ErrorMIDdleware())

	service := &Service{
		g,
		albmCtrl,
		phtoCtrl,
	}

	// Initialize routes
	service.addRoutes()

	return service
}

// Listen will start the HTTP server and await requests
func (svc *Service) Listen(addr string) {
	svc.g.Run(addr)
}

// addRoutes registers the API endpoints with the Gin engine
func (svc *Service) addRoutes() {
	svc.g.GET(apiPrefix+"/photos", svc.listPhotos)
	svc.g.POST(apiPrefix+"/photos", svc.createPhoto)
	svc.g.GET(apiPrefix+"/photos/:photoID", svc.getPhoto)
	svc.g.PUT(apiPrefix+"/photos/:photoID", svc.updatePhoto)
	svc.g.DELETE(apiPrefix+"/photos/:photoID", svc.deletePhoto)

	svc.g.GET(apiPrefix+"/albums", svc.listAlbums)
	svc.g.POST(apiPrefix+"/albums", svc.createAlbum)
	svc.g.GET(apiPrefix+"/albums/:albumID", svc.getAlbum)
	svc.g.PUT(apiPrefix+"/albums/:albumID", svc.updateAlbum)
	svc.g.DELETE(apiPrefix+"/albums/:albumID", svc.deleteAlbum)
}
