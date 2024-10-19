package application

import (
	"main/handler"
	"main/repository/order"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (a *App) loadRoutes() {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.Writer.WriteHeader(http.StatusOK)
	})

	v1 := router.Group("/v1")
	ordersGroup := v1.Group("/orders")

	a.loadOrderRoutes(ordersGroup)

	a.router = router
}

func (a *App) loadOrderRoutes(group *gin.RouterGroup) {
	orderHandler := &handler.Order{
		Repo: &order.RedisRepo{
			Client: a.rdb,
		},
	}

	group.POST("/", orderHandler.Create)
	group.GET("/", orderHandler.List)
	group.GET("/:id", orderHandler.GetByID)
	group.PUT("/:id", orderHandler.UpdateByID)
	group.DELETE("/:id", orderHandler.DeleteByID)
}
