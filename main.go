package main

import (
	"net/http"
	"regexp"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var deployments sync.Map

type Deployment struct {
	ID       uuid.UUID         `json:"id"`
	Labels   map[string]string `json:"labels"`
	Replicas int               `json:"replicas"`
	Image    string            `json:"image"`
	Name     string            `json:"name"`
	Ports    []Port            `json:"ports"`
	CreateAt time.Time         `json:"createAt"`
}

type Port struct {
	Name   string `json:"name"`
	Number uint   `json:"port"`
}

type Error struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func main() {
	r := gin.Default()
	r.POST("/deployments", func(c *gin.Context) {
		var deployment Deployment
		c.ShouldBindJSON(&deployment)

		if deployment.ID == uuid.Nil {
			c.Writer.Header().Set("Content-Type", "application/json")
			c.JSON(http.StatusBadRequest, Error{
				Message: "ID must be specified",
				Code:    6,
			})
			return
		}

		if deployment.Replicas == 0 {
			c.Writer.Header().Set("Content-Type", "application/json")
			c.JSON(http.StatusBadRequest, Error{
				Message: "Replicas must be greater than 0",
				Code:    1,
			})
			return
		}

		nameMatch, _ := regexp.MatchString(`^[a-zA-Z0-9_]+$`, deployment.Image)
		if !nameMatch {
			c.Writer.Header().Set("Content-Type", "application/json")
			c.JSON(http.StatusBadRequest, Error{
				Message: "Image name must be alphanumeric",
				Code:    2,
			})
			return
		}

		if len(deployment.Ports) == 0 {
			c.Writer.Header().Set("Content-Type", "application/json")
			c.JSON(http.StatusBadRequest, Error{
				Message: "Ports must be specified",
				Code:    3,
			})
			return
		}

		for _, port := range deployment.Ports {
			if port.Number == 0 || port.Number > uint(65535) {
				c.Writer.Header().Set("Content-Type", "application/json")
				c.JSON(http.StatusBadRequest, Error{
					Message: "Port number must be between 1 and 65535",
					Code:    4,
				})
				return
			}
		}

		_, exists := deployments.Load(deployment.ID.String())
		if exists {
			c.Writer.Header().Set("Content-Type", "application/json")
			c.JSON(http.StatusConflict, Error{
				Message: "Deployment already found with this ID",
				Code:    5,
			})
			return
		}

		deployments.Store(deployment.ID.String(), deployment)

		c.Writer.Header().Set("Content-Type", "application/json")
		c.JSON(http.StatusCreated, deployment)
	})

	r.GET("/deployments/:id", func(c *gin.Context) {
		id := c.Param("id")

		d, exists := deployments.Load(id)
		if !exists {
			c.Writer.Header().Set("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, Error{
				Message: "Deployment not found",
				Code:    5,
			})
			return
		}

		c.JSON(http.StatusOK, d)
	})

	r.DELETE("/deployments/:id", func(c *gin.Context) {
		id := c.Param("id")

		_, exists := deployments.Load(id)
		if !exists {
			c.Writer.Header().Set("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, Error{
				Message: "Deployment not found",
				Code:    5,
			})
			return
		}

		deployments.Delete(id)

		c.Status(http.StatusNoContent)
	})

	r.Run()
}
