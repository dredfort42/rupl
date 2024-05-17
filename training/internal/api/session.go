package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
	s "training/internal/structs"

	"github.com/gin-gonic/gin"
)

// CreateSession is a function to create a new session
func CreateSession(c *gin.Context) {

	// var session s.LastSessionData

	var session json.RawMessage

	var email string
	// var errorResponse s.ResponseError
	// var err error

	if email = VerifyDevice(c); email == "" {
		return
	}

	// if err = c.ShouldBindJSON(&session); err != nil {
	// 	errorResponse.Error = "invalid_request"
	// 	errorResponse.ErrorDescription = err.Error()
	// 	c.IndentedJSON(http.StatusBadRequest, errorResponse)
	// 	return
	// }

	if err := c.BindJSON(&session); err != nil {
		c.JSON(http.StatusBadRequest, s.ResponseError{Error: "Invalid request", ErrorDescription: err.Error()})
		return
	}

	fileName := time.Now().Format("2006-01-02-15-04-05") + ".json"

	os.WriteFile("/app/dump/"+fileName, session, 0644)

	fmt.Println(session)

	// if session == nil {
	// 	errorResponse.Error = "invalid_request"
	// 	errorResponse.ErrorDescription = "Missing session"
	// 	c.IndentedJSON(http.StatusBadRequest, errorResponse)
	// 	return
	// }

	// if err = db.CreateSession(email, session); err != nil {
	// 	errorResponse.Error = "server_error"
	// 	errorResponse.ErrorDescription = err.Error()
	// 	c.IndentedJSON(http.StatusInternalServerError, errorResponse)
	// }

	c.JSON(http.StatusOK, session)
}
