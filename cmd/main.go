// Package main is the entrypoint for this application.
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"

	"photolib/internal/config"
	"photolib/internal/service"
)

// Initialize postgres connection pool
func initPostgres(ctx context.Context) (*pgxpool.Pool, error) {
	// Read credentials from environment
	user := viper.GetString("postgres.username")
	password := viper.GetString("postgres.password")
	dbname := viper.GetString("postgres.db")
	host := viper.GetString("postgres.host")

	// Create connection
	connStr := fmt.Sprintf("postgres://%s:%s@%s:5432/%s", user, password, host, dbname)
	pg, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, err
	}

	return pg, nil
}

func main() {
	// Initialize configuration
	err := config.InitConfig()
	if err != nil {
		log.Fatalf("unable to initialize viper configuration: %v", err)
	}

	// Initialize postgres connection
	ctx := context.Background()
	pg, err := initPostgres(ctx)
	if err != nil {
		log.Fatalf("unable to connect to postgres database: %v", err)
	}
	defer pg.Close()

	// Initialize HTTP server
	g := gin.Default()
	svc := service.NewService(g, pg)

	// Start HTTP server
	addr := fmt.Sprintf("0.0.0.0:%d", viper.GetInt("service.port"))
	svc.Listen(addr)

}
