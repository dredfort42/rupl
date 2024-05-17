package db

import (
	"errors"
	s "training/internal/structs"
)

// CreateSession creates a new session in the database
func CreateSession(session s.DBSession) (err error) {
	if CheckSessionExists(session.SessionStartTime, session.SessionEndTime, session.Email) {
		return errors.New("session already exists")
	}

	query := `
		INSERT INTO ` + db.tableSessions + ` (
			session_uuid, session_start_time, session_end_time, email, route_data, step_count, running_power,
			vertical_oscillation, energy_burned, heart_rate, stride_length, ground_contact_time, speed, distance,
			vo2max_mL_per_min_per_kg, avr_speed_m_per_s, avr_heart_rate_count_per_s, avr_power_W,
			avr_vertical_oscillation_cm, avr_ground_contact_time_ms, avr_stride_length_m, total_distance_m,
			total_steps_count, total_energy_burned_kcal
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24
		)`

	_, err = db.database.Exec(query,
		session.SessionUUID, session.SessionStartTime, session.SessionEndTime, session.Email, session.RouteData,
		session.StepCount, session.RunningPower, session.VerticalOscillation, session.EnergyBurned,
		session.HeartRate, session.StrideLength, session.GroundContactTime, session.Speed, session.Distance,
		session.VO2MaxMLPerMinPerKg, session.AvrSpeedMPerS, session.AvrHeartRateCountPerS, session.AvrPowerW,
		session.AvrVerticalOscillationCm, session.AvrGroundContactTimeMs, session.AvrStrideLengthM, session.TotalDistanceM,
		session.TotalStepsCount, session.TotalEnergyBurnedKcal)

	return
}
