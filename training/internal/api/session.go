package api

import (
	"net/http"
	"strconv"
	"strings"
	"training/internal/db"
	s "training/internal/structs"

	loger "github.com/dredfort42/tools/logprinter"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// parseSessionDataInt is a function to parse a session data for int
func parseSessionDataInt(data []s.JSONLastSessionTypeData) (dbData []s.DBSessionDataInt, err error) {
	dbData = make([]s.DBSessionDataInt, len(data))
	for i, d := range data {
		var intValue int
		intValue, err = strconv.Atoi(strings.Split(d.Quantity, " ")[0])
		if err != nil {
			return
		}

		dbData[i].Timestamp = d.Timestamp
		dbData[i].Data = intValue
	}

	return
}

// parseSessionDataFloat32 is a function to parse a session data for int
func parseSessionDataFloat32(data []s.JSONLastSessionTypeData) (dbData []s.DBSessionDataFloat32, err error) {
	dbData = make([]s.DBSessionDataFloat32, len(data))
	for i, d := range data {
		var float64Value float64
		float64Value, err = strconv.ParseFloat(strings.Split(d.Quantity, " ")[0], 32)
		if err != nil {
			return
		}

		dbData[i].Timestamp = d.Timestamp
		dbData[i].Data = float32(float64Value)
	}

	return
}

// getAverageFromInt is a function to get the average from an int array
func getAverageFromInt(data []s.DBSessionDataInt) (average float32) {
	var sum int
	for _, d := range data {
		sum += d.Data
	}

	average = float32(sum) / float32(len(data))
	return
}

// getTotalFromInt is a function to get the total from an int array
func getTotalFromInt(data []s.DBSessionDataInt) (total int) {
	for _, d := range data {
		total += d.Data
	}

	return
}

// getAverageFromFloat32 is a function to get the average from a float32 array
func getAverageFromFloat32(data []s.DBSessionDataFloat32) (average float32) {
	var sum float32
	for _, d := range data {
		sum += d.Data
	}

	average = sum / float32(len(data))
	return
}

// getTotalFromFloat32 is a function to get the total from a float32 array
func getTotalFromFloat32(data []s.DBSessionDataFloat32) (total float32) {
	for _, d := range data {
		total += d.Data
	}

	return
}

// parseSession is a function to parse a session
func parseSession(jsonSession s.JSONLastSessionData) (dbSession s.DBSession, err error) {
	dbSession = s.DBSession{
		SessionUUID:      uuid.New().String(),
		SessionStartTime: jsonSession.Session.StartTime,
		SessionEndTime:   jsonSession.Session.EndTime,
		Email:            jsonSession.Session.Email,
	}

	dbSession.RouteData = append(dbSession.RouteData, jsonSession.RouteData...)

	dbSession.StepCount, err = parseSessionDataInt(jsonSession.StepCount)
	if err != nil {
		return
	}

	dbSession.RunningPower, err = parseSessionDataInt(jsonSession.RunningPower)
	if err != nil {
		return
	}

	dbSession.VerticalOscillation, err = parseSessionDataFloat32(jsonSession.VerticalOscillation)
	if err != nil {
		return
	}

	dbSession.EnergyBurned, err = parseSessionDataFloat32(jsonSession.EnergyBurned)
	if err != nil {
		return
	}

	dbSession.HeartRate, err = parseSessionDataFloat32(jsonSession.HeartRate)
	if err != nil {
		return
	}

	dbSession.StrideLength, err = parseSessionDataFloat32(jsonSession.StrideLength)
	if err != nil {
		return
	}

	dbSession.GroundContactTime, err = parseSessionDataInt(jsonSession.GroundContactTime)
	if err != nil {
		return
	}

	dbSession.Speed, err = parseSessionDataFloat32(jsonSession.Speed)
	if err != nil {
		return
	}

	dbSession.Distance, err = parseSessionDataFloat32(jsonSession.Distance)
	if err != nil {
		return
	}

	vo2max, err := strconv.ParseFloat(strings.Split(jsonSession.VO2Max[0].Quantity, " ")[0], 32)
	if err != nil {
		return
	}
	dbSession.VO2MaxMLPerMinPerKg = float32(vo2max)

	dbSession.AvrSpeedMPerS = getAverageFromFloat32(dbSession.Speed)
	dbSession.AvrHeartRateCountPerS = getAverageFromFloat32(dbSession.HeartRate)
	dbSession.AvrPowerW = getAverageFromInt(dbSession.RunningPower)
	dbSession.AvrVerticalOscillationCm = getAverageFromFloat32(dbSession.VerticalOscillation)
	dbSession.AvrGroundContactTimeMs = getAverageFromInt(dbSession.GroundContactTime)
	dbSession.AvrStrideLengthM = getAverageFromFloat32(dbSession.StrideLength)

	dbSession.TotalDistanceM = getTotalFromFloat32(dbSession.Distance)
	dbSession.TotalStepsCount = getTotalFromInt(dbSession.StepCount)
	dbSession.TotalEnergyBurnedKcal = int(getTotalFromFloat32(dbSession.EnergyBurned))

	return
}

// CreateSession is a function to create a new session
func CreateSession(c *gin.Context) {
	var email string
	var session s.JSONLastSessionData
	var errorResponse s.ResponseError
	var err error

	if email = VerifyDevice(c); email == "" {
		return
	}

	if err = c.ShouldBindJSON(&session); err != nil {
		loger.Error("Error binding JSON", err)
		errorResponse.Error = "invalid_request"
		errorResponse.ErrorDescription = err.Error()
		c.IndentedJSON(http.StatusBadRequest, errorResponse)
		return
	}

	if session.Session.Email != email {
		errorResponse.Error = "invalid_request"
		errorResponse.ErrorDescription = "Invalid email in session data"
		c.IndentedJSON(http.StatusBadRequest, errorResponse)
		return
	}

	if session.Session.StartTime == 0 || session.Session.EndTime == 0 || session.Session.StartTime > session.Session.EndTime {
		errorResponse.Error = "invalid_request"
		errorResponse.ErrorDescription = "Invalid session time"
		c.IndentedJSON(http.StatusBadRequest, errorResponse)
		return
	}

	dbSession, err := parseSession(session)
	if err != nil {
		loger.Error("Error parsing session", err)
		errorResponse.Error = "invalid_request"
		errorResponse.ErrorDescription = err.Error()
		c.IndentedJSON(http.StatusBadRequest, errorResponse)
		return
	}

	if err = db.CreateSession(dbSession); err != nil {
		errorResponse.Error = "server_error"
		errorResponse.ErrorDescription = "Error creating session"
		c.IndentedJSON(http.StatusInternalServerError, errorResponse)
		return
	}

	loger.Debug("Session created successfully for an ID: ", dbSession.SessionUUID)

	c.JSON(http.StatusOK, session)
}
