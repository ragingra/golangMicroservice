package application

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type App struct {
	router *gin.Engine
	rdb    *redis.Client
	Config Config
}

func New(Config Config) *App {
	app := &App{
		rdb:    redis.NewClient(&redis.Options{Addr: Config.RedisAddress}),
		Config: Config,
	}

	app.loadRoutes()

	return app
}

func (app *App) Start(ctx context.Context) error {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", app.Config.ServerPort),
		Handler: app.router,
	}

	err := app.rdb.Ping(ctx).Err()
	if err != nil {
		return fmt.Errorf("failed to ping redis: %w", err)
	}

	defer func() {
		err := app.rdb.Close()
		if err != nil {
			fmt.Println("failed to close redis connection")
		}
	}()

	fmt.Println("Starting server on port", app.Config.ServerPort)
	ch := make(chan error, 1)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			ch <- fmt.Errorf("failed to start server, %w", err)
		}
		close(ch)
	}()

	select {
	case err := <-ch:
		return err
	case <-ctx.Done():
		timeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		return server.Shutdown(timeout)
	}
}
