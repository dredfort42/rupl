package db

// usersTableCheck() checks if the users table exists, if not, it creates it
func usersTableCheck() {
	query := `
		CREATE TABLE IF NOT EXISTS ` + db.tableUsers + ` (
			email VARCHAR(255) PRIMARY KEY,
			first_name VARCHAR(255) NOT NULL,
			last_name VARCHAR(255) NOT NULL,
			date_of_birth DATE NOT NULL,
			gender VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`

	_, err := db.database.Exec(query)
	if err != nil {
		panic(err)
	}
}

// devicesTableCheck checks if the devices table exists, if not, it creates it
func devicesTableCheck() {
	query := `
			CREATE TABLE IF NOT EXISTS ` + db.tableDevices + ` (
				id SERIAL PRIMARY KEY,
				email VARCHAR(255) NOT NULL,
				device_model VARCHAR(255) NOT NULL,
				device_name VARCHAR(255) NOT NULL,
				system_name VARCHAR(255) NOT NULL,
				system_version VARCHAR(255) NOT NULL,
				device_id VARCHAR(255) NOT NULL,
				app_version VARCHAR(255) NOT NULL,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			);
		`

	_, err := db.database.Exec(query)
	if err != nil {
		panic(err)
	}
}

// CheckTables checks if the tables exists, if not, it creates it
func tablesCheck() {
	usersTableCheck()
	devicesTableCheck()
}
