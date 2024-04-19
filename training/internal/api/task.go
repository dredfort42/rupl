package api

import (
	"net/http"
	"time"

	"training/internal/training"

	"github.com/gin-gonic/gin"
)

// TMP: taskDeclined is a temporary var
var taskDeclined bool

// GetTask returns a task
func GetTask(c *gin.Context) {
	// TMP: taskDeclined is a temporary check
	if time.Now().Minute()%2 == 0 {
		taskDeclined = false
	}

	if taskDeclined {
		c.IndentedJSON(http.StatusNoContent, nil)
		return
	}

	// TODO: Implement GetTask
	var task training.Task

	task.ID = 1
	task.Description = "Test task"
	task.Compleated = false
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

// DeclineTask declines a task
func DeclineTask(c *gin.Context) {
	taskDeclined = true
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Task successfully declined"})
}
