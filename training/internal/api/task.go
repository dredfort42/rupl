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
	task.Description = "Tempo Intervals"
	task.Compleated = false
	task.Intervals = []training.Interval{
		{
			ID:          1,
			Description: "Easy run",
			Speed:       0,
			PulseZone:   1,
			Distance:    0,
			Duration:    10 * 60,
		},
		{
			ID:          2,
			Description: "Tempo run",
			Speed:       0,
			PulseZone:   3,
			Distance:    0,
			Duration:    5 * 60,
		},
		{
			ID:          3,
			Description: "Easy run",
			Speed:       0,
			PulseZone:   2,
			Distance:    0,
			Duration:    2*60 + 30,
		},
		{
			ID:          4,
			Description: "Tempo run",
			Speed:       0,
			PulseZone:   3,
			Distance:    0,
			Duration:    5 * 60,
		},
		{
			ID:          5,
			Description: "Easy run",
			Speed:       0,
			PulseZone:   2,
			Distance:    0,
			Duration:    2*60 + 30,
		},
		{
			ID:          6,
			Description: "Tempo run",
			Speed:       0,
			PulseZone:   3,
			Distance:    0,
			Duration:    5 * 60,
		},
		{
			ID:          7,
			Description: "Easy run",
			Speed:       0,
			PulseZone:   2,
			Distance:    0,
			Duration:    2*60 + 30,
		},
		{
			ID:          1,
			Description: "Steady run",
			Speed:       0,
			PulseZone:   2,
			Distance:    0,
			Duration:    10 * 60,
		},
	}

	c.IndentedJSON(http.StatusOK, task)
}

// DeclineTask declines a task
func DeclineTask(c *gin.Context) {
	taskDeclined = true
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Task successfully declined"})
}
