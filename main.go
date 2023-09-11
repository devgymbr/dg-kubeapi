package main

import (
	"regexp"
	"sync"

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
		c.BindJSON(&deployment)

		if deployment.ID == uuid.Nil {
			c.JSON(400, Error{
				Message: "ID must be specified",
				Code:    6,
			})
			return
		}

		if deployment.Replicas == 0 {
			c.JSON(400, Error{
				Message: "Replicas must be greater than 0",
				Code:    1,
			})
			return
		}

		nameMatch, _ := regexp.MatchString(`^[a-zA-Z0-9_]+$`, deployment.Image)
		if !nameMatch {
			c.JSON(400, Error{
				Message: "Image name must be alphanumeric",
				Code:    2,
			})
			return
		}

		if len(deployment.Ports) == 0 {
			c.JSON(400, Error{
				Message: "Ports must be specified",
				Code:    3,
			})
			return
		}

		for _, port := range deployment.Ports {
			if port.Number > uint(65535) {
				c.JSON(400, Error{
					Message: "Port number must be less than 65535",
					Code:    4,
				})
				return
			}
		}
		deployments.Store(deployment.ID, deployment)
	})

	r.GET("/deployments/:id", func(c *gin.Context) {
		id := c.Param("id")

		d, exists := deployments.Load(id)
		if !exists {
			c.JSON(404, Error{
				Message: "Deployment not found",
				Code:    5,
			})
			return
		}

		c.JSON(200, d)
	})

	r.Run()
}
