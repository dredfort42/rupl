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
			ID:        1,
			Speed:     10,
			PulseZone: 0,
			Distance:  1000,
			Duration:  0,
		},
		{
			ID:        2,
			Speed:     0,
			PulseZone: 4,
			Distance:  0,
			Duration:  1200,
		},
		{
			ID:        3,
			Speed:     10,
			PulseZone: 0,
			Distance:  1000,
			Duration:  0,
		},
	}

	c.IndentedJSON(http.StatusOK, task)
}
