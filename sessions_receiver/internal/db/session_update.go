package db

// calculateExtraData calculates extra data for a session
func calculateExtraData(session_uuid string) (err error) {
	query := `
			WITH
				avg_speed AS (
					SELECT COALESCE(AVG(t.speed_m_per_s), 0) AS average_speed
					FROM sessions,
						UNNEST(speed) AS t
					WHERE session_uuid = '` + session_uuid + `'
				),
				total_steps AS (
					SELECT COALESCE(SUM(t.steps_count), 0) AS total_steps
					FROM sessions,
						UNNEST(step_count) AS t
					WHERE session_uuid = '` + session_uuid + `'
				),
				avg_running_power AS (
					SELECT COALESCE(AVG(t.power_w), 0) AS average_running_power
					FROM sessions,
						UNNEST(running_power) AS t
					WHERE session_uuid = '` + session_uuid + `'
				),
				avg_vertical_oscillation AS (
					SELECT COALESCE(AVG(t.vertical_oscillation_cm), 0) AS average_vertical_oscillation
					FROM sessions,
						UNNEST(vertical_oscillation) AS t
					WHERE session_uuid = '` + session_uuid + `'
				),
				total_energy_burned AS (
					SELECT COALESCE(SUM(t.energy_burned_kcal), 0) AS total_energy_burned
					FROM sessions,
						UNNEST(energy_burned) AS t
					WHERE session_uuid = '` + session_uuid + `'
				),
				avg_heart_rate AS (
					SELECT COALESCE(AVG(t.heart_rate_count_per_s), 0) AS average_heart_rate
					FROM sessions,
						UNNEST(heart_rate) AS t
					WHERE session_uuid = '` + session_uuid + `'
				),
				avg_stride_length AS (
					SELECT COALESCE(AVG(t.stride_length_m), 0) AS average_stride_length
					FROM sessions,
						UNNEST(stride_length) AS t
					WHERE session_uuid = '` + session_uuid + `'
				),
				avg_ground_contact_time AS (
					SELECT COALESCE(AVG(t.ground_contact_time_ms), 0) AS average_ground_contact_time
					FROM sessions,
						UNNEST(ground_contact_time) AS t
					WHERE session_uuid = '` + session_uuid + `'
				),
				total_distance AS (
					SELECT COALESCE(SUM(t.distance_m), 0) AS total_distance
					FROM sessions,
						UNNEST(distance) AS t
					WHERE session_uuid = '` + session_uuid + `'
				)
			UPDATE sessions
			SET
				avr_speed_m_per_s = (SELECT average_speed FROM avg_speed),
				total_steps_count = (SELECT total_steps FROM total_steps),
				avr_power_w = (SELECT average_running_power FROM avg_running_power),
				avr_vertical_oscillation_cm = (SELECT average_vertical_oscillation FROM avg_vertical_oscillation),
				total_energy_burned_kcal = (SELECT total_energy_burned FROM total_energy_burned),
				avr_heart_rate_count_per_s = (SELECT average_heart_rate FROM avg_heart_rate),
				avr_stride_length_m = (SELECT average_stride_length FROM avg_stride_length),
				avr_ground_contact_time_ms = (SELECT average_ground_contact_time FROM avg_ground_contact_time),
				total_distance_m = (SELECT total_distance FROM total_distance)
			WHERE session_uuid = '` + session_uuid + `';`

	_, err = db.database.Exec(query)

	return
}
