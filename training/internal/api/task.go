package api

import (
	"net/http"

	"training/internal/training"

	"github.com/gin-gonic/gin"
)

func GetTask(c *gin.Context) {
	// TODO: Implement GetTask
	var task training.Task

	task.ID = 1
	task.Description = "Test task"
	task.Intervals = []training.Interval{
		{
			ID:          1,
			Description: "Warm up",
			Speed:       10,
			PulseZone:   0,
			Distance:    100,
			Duration:    0,
		},
		{
			ID:          2,
			Description: "Sprint",
			Speed:       0,
			PulseZone:   4,
			Distance:    0,
			Duration:    60,
		},
		{
			ID:          3,
			Description: "Cool down",
			Speed:       10,
			PulseZone:   0,
			Distance:    100,
			Duration:    0,
		},
	}

	c.IndentedJSON(http.StatusOK, task)
}
