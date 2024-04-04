package api

import (
	"net/http"
	"training/internal/db"

	"github.com/gin-gonic/gin"
)

func GetTask(c *gin.Context) {
	// TODO: Implement GetTask
	var task db.Task

	task.ID = 1
	task.Description = "Test task"
	task.Intervals = []db.Interval{
		{
			ID:        1,
			Speed:     10,
			PulseZone: 3,
			Distance:  1000,
			Duration:  600,
		},
		{
			ID:        2,
			Speed:     12,
			PulseZone: 4,
			Distance:  2000,
			Duration:  1200,
		},
		{
			ID:        3,
			Speed:     14,
			PulseZone: 5,
			Distance:  3000,
			Duration:  1800,
		},
	}

	c.IndentedJSON(http.StatusOK, task)
}
