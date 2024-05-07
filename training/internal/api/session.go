package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateSession is a function to create a new session
func CreateSession(c *gin.Context) {
	var session json.RawMessage
	if err := c.BindJSON(&session); err != nil {
		c.JSON(http.StatusBadRequest, ResponseError{Error: "Invalid request", ErrorDescription: err.Error()})
		return
	}

	fmt.Println("Session: ", string(session))

	// if session.ID == "" {
	// 	session.ID = fmt.Sprintf("%d", time.Now().UnixNano())
	// }

	// if err := session.Create(); err != nil {
	// 	c.JSON(http.StatusInternalServerError, ResponseError{Error: "Internal server error", ErrorDescription: err.Error()})
	// 	return
	// }

	c.JSON(http.StatusOK, session)
}
