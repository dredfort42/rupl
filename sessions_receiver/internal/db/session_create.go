package db

import (
	"errors"
	s "sessions_receiver/internal/structs"
)

// CreateSession creates a new session in the database
func CreateSession(session s.DBSession) (err error) {
	if CheckSessionExists(session.SessionStartTime, session.SessionEndTime, session.Email) {
		return errors.New("session already exists")
	}

	query := `
		INSERT INTO ` + db.tableSessions + ` (
				session_uuid,
				session_start_time,
				session_end_time,
				email,
				route_data,
				step_count,
				running_power,
				vertical_oscillation,
				energy_burned,
				heart_rate,
				stride_length,
				ground_contact_time,
				speed,
				distance,
				vo2max_mL_per_min_per_kg
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
		)`

	_, err = db.database.Exec(query,
		session.SessionUUID,
		session.SessionStartTime,
		session.SessionEndTime,
		session.Email,
		session.RouteData,
		session.StepCount,
		session.RunningPower,
		session.VerticalOscillation,
		session.EnergyBurned,
		session.HeartRate,
		session.StrideLength,
		session.GroundContactTime,
		session.Speed,
		session.Distance,
		session.VO2MaxMLPerMinPerKg)
	if err != nil {
		return
	}

	err = calculateExtraData(session.SessionUUID)

	return
}
