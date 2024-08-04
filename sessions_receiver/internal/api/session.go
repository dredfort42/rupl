package api

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	s "sessions_receiver/internal/structs"

	loger "github.com/dredfort42/tools/logprinter"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// func parseRouteData(routeData []s.RouteData) (dbRouteData pq.StringArray) {
// 	dbRouteData = make([]string, 0, len(routeData))

// 	for _, rd := range routeData {
// 		dbRouteData = append(dbRouteData, fmt.Sprintf("(%d, %f, %f, %f, %f, %f, %f, %f, %f, %f)",
// 			rd.Timestamp,
// 			rd.Latitude,
// 			rd.Longitude,
// 			rd.HorizontalAccuracy,
// 			rd.Altitude,
// 			rd.VerticalAccuracy,
// 			rd.Speed,
// 			rd.SpeedAccuracy,
// 			rd.Course,
// 			rd.CourseAccuracy))
// 	}

// 	return
// }

// // parseSessionDataInt is a function to parse a session data for int
// func parseSessionDataInt(data []s.JSONLastSessionTypeData) (dbData pq.StringArray) {
// 	dbData = make([]string, 0, len(data))

// 	for _, d := range data {
// 		dbData = append(dbData, fmt.Sprintf("(%d, %s)",
// 			d.Timestamp,
// 			strings.Split(strings.Split(d.Quantity, " ")[0], ".")[0]))
// 	}

// 	return
// }

// // parseSessionDataFloat32 is a function to parse a session data for int
// func parseSessionDataFloat32(data []s.JSONLastSessionTypeData) (dbData pq.StringArray) {
// 	dbData = make([]string, 0, len(data))

// 	for _, d := range data {
// 		dbData = append(dbData, fmt.Sprintf("(%d, %s)",
// 			d.Timestamp,
// 			strings.Split(d.Quantity, " ")[0]))
// 	}

// 	return
// }

// // parseSession is a function to parse a session
// func parseSession(jsonSession s.JSONLastSessionData) (dbSession s.DBSession, err error) {
// 	dbSession = s.DBSession{
// 		SessionUUID:      uuid.New().String(),
// 		SessionStartTime: jsonSession.Session.StartTime,
// 		SessionEndTime:   jsonSession.Session.EndTime,
// 		Email:            jsonSession.Session.Email,
// 	}

// 	dbSession.RouteData = parseRouteData(jsonSession.RouteData)
// 	dbSession.StepCount = parseSessionDataInt(jsonSession.StepCount)
// 	dbSession.RunningPower = parseSessionDataInt(jsonSession.RunningPower)
// 	dbSession.VerticalOscillation = parseSessionDataFloat32(jsonSession.VerticalOscillation)
// 	dbSession.EnergyBurned = parseSessionDataFloat32(jsonSession.EnergyBurned)
// 	dbSession.HeartRate = parseSessionDataFloat32(jsonSession.HeartRate)
// 	dbSession.StrideLength = parseSessionDataFloat32(jsonSession.StrideLength)
// 	dbSession.GroundContactTime = parseSessionDataInt(jsonSession.GroundContactTime)
// 	dbSession.Speed = parseSessionDataFloat32(jsonSession.Speed)
// 	dbSession.Distance = parseSessionDataFloat32(jsonSession.Distance)

// 	vo2max := float64(0)
// 	if len(jsonSession.VO2Max) > 0 {
// 		vo2max, err = strconv.ParseFloat(strings.Split(jsonSession.VO2Max[0].Quantity, " ")[0], 32)
// 		if err != nil {
// 			return
// 		}
// 	}
// 	dbSession.VO2MaxMLPerMinPerKg = float32(vo2max)

// 	return
// }

// CreateSession is a function to create a new session
func CreateSession(c *gin.Context) {
	// var email string
	// var session s.JSONLastSessionData
	var errorResponse s.ResponseError
	var err error

	// if email = VerifyDevice(c); email == "" {
	// 	return
	// }

	// write body to file on disk
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		loger.Error("Error reading body", err)
		errorResponse.Error = "invalid_request"
		errorResponse.ErrorDescription = err.Error()
		c.IndentedJSON(http.StatusBadRequest, errorResponse)
		return
	}

	dir := "app/sessions"
	fileName := fmt.Sprintf("session_%s.json", uuid.New().String())

	filePath := filepath.Join(dir, fileName)

	err = os.MkdirAll(dir, 0755)
	if err != nil {
		loger.Error("Error creating directory", err)
		errorResponse.Error = "invalid_request"
		errorResponse.ErrorDescription = err.Error()
		c.IndentedJSON(http.StatusBadRequest, errorResponse)
		return
	}

	file, err := os.Create(filePath)
	if err != nil {
		loger.Error("Error opening file", err)
		errorResponse.Error = "invalid_request"
		errorResponse.ErrorDescription = err.Error()
		c.IndentedJSON(http.StatusBadRequest, errorResponse)
		return
	}
	defer file.Close()

	_, err = file.Write(body)
	if err != nil {
		loger.Error("Error writing body to file", err)
		errorResponse.Error = "invalid_request"
		errorResponse.ErrorDescription = err.Error()
		c.IndentedJSON(http.StatusBadRequest, errorResponse)
		return
	}

	// if err = c.ShouldBindJSON(&session); err != nil {
	// 	loger.Error("Error binding JSON", err)
	// 	errorResponse.Error = "invalid_request"
	// 	errorResponse.ErrorDescription = err.Error()
	// 	c.IndentedJSON(http.StatusBadRequest, errorResponse)
	// 	return
	// }

	// if session.Session.Email != email {
	// 	errorResponse.Error = "invalid_request"
	// 	errorResponse.ErrorDescription = "Invalid email in session data"
	// 	c.IndentedJSON(http.StatusBadRequest, errorResponse)
	// 	return
	// }

	// if session.Session.StartTime == 0 || session.Session.EndTime == 0 || session.Session.StartTime > session.Session.EndTime {
	// 	errorResponse.Error = "invalid_request"
	// 	errorResponse.ErrorDescription = "Invalid session time"
	// 	c.IndentedJSON(http.StatusBadRequest, errorResponse)
	// 	return
	// }

	// dbSession, err := parseSession(session)
	// if err != nil {
	// 	loger.Error("Error parsing session", err)
	// 	errorResponse.Error = "invalid_request"
	// 	errorResponse.ErrorDescription = err.Error()
	// 	c.IndentedJSON(http.StatusBadRequest, errorResponse)
	// 	return
	// }

	// if err = db.CreateSession(dbSession); err != nil {
	// 	errorResponse.Error = "server_error"
	// 	errorResponse.ErrorDescription = "Error creating session"
	// 	// errorResponse.ErrorDescription = err.Error()
	// 	c.IndentedJSON(http.StatusInternalServerError, errorResponse)
	// 	return
	// }

	// loger.Debug("Session created successfully for an ID: ", dbSession.SessionUUID)

	c.JSON(http.StatusOK, gin.H{"message": "Session created successfully"})
}
