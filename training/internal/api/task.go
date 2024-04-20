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

	task.Intervals = append(task.Intervals, training.Interval{
		ID:          1,
		Description: "Easy run",
		Speed:       training.Speed{Min: 0, Max: 0},
		Pulse:       training.Pulse{Min: 121, Max: 131},
		Distance:    0,
		Duration:    10 * 60,
	})

	for i := 0; i < 3; i++ {
		task.Intervals = append(task.Intervals, training.Interval{
			ID:          i + 2,
			Description: "Tempo run",
			Speed:       training.Speed{Min: secPKmToMPS(300), Max: secPKmToMPS(297)},
			Pulse:       training.Pulse{Min: 0, Max: 0},
			Distance:    0,
			Duration:    5 * 60,
		})

		task.Intervals = append(task.Intervals, training.Interval{
			ID:          i + 3,
			Description: "Easy run",
			Speed:       training.Speed{Min: 0, Max: 0},
			Pulse:       training.Pulse{Min: 121, Max: 131},
			Distance:    0,
			Duration:    2*60 + 30,
		})
	}

	task.Intervals = append(task.Intervals, training.Interval{
		ID:          8,
		Description: "Steady run",
		Speed:       training.Speed{Min: 0, Max: 0},
		Pulse:       training.Pulse{Min: 121, Max: 131},
		Distance:    0,
		Duration:    10 * 60,
	})

	c.IndentedJSON(http.StatusOK, task)
}

// DeclineTask declines a task
func DeclineTask(c *gin.Context) {
	taskDeclined = true
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Task successfully declined"})
}

// Converter sec/km to m/s
func secPKmToMPS(secPerKm float32) float32 {
	return 1000 / secPerKm
}
