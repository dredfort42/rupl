package db

// checkUsersTable checks if the users table exists, if not, it creates it
func checkUsersTable() {
	_, err := db.database.Exec("CREATE EXTENSION IF NOT EXISTS pgcrypto;")
	if err != nil {
		panic(err)
	}

	query := `
		CREATE TABLE IF NOT EXISTS ` + db.tableUsers + ` (
			email VARCHAR(255) PRIMARY KEY,
			password_hash VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
		);
	`
	_, err = db.database.Exec(query)
	if err != nil {
		panic(err)
	}
}

// checkSessionsTable checks if the sessions table exists, if not, it creates it
func checkSessionsTable() {
	query := `
		CREATE TABLE IF NOT EXISTS ` + db.tableSessions + ` (
			email VARCHAR(255) NOT NULL,
			access_token VARCHAR(255) NOT NULL,
			refresh_token VARCHAR(255) NOT NULL,
			is_one_time BOOLEAN DEFAULT FALSE NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
			CONSTRAINT sessions_unique UNIQUE (email, access_token)
		);
	`
	_, err := db.database.Exec(query)
	if err != nil {
		panic(err)
	}
}

// checkDevicesTable checks if the devices table exists, if not, it creates it
func checkDevicesTable() {
	query := `
		CREATE TABLE IF NOT EXISTS ` + db.tableDevices + ` (
			device_uuid VARCHAR(255) NOT NULL,
			device_access_token VARCHAR(255) NOT NULL,
			device_refresh_token VARCHAR(255) NOT NULL,
			email VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
			CONSTRAINT devices_unique UNIQUE (device_uuid, email)
		);
	`
	_, err := db.database.Exec(query)
	if err != nil {
		panic(err)
	}
}

// tablesCheck checks if the tables exist, if not, it creates them
func tablesCheck() {
	checkUsersTable()
	checkSessionsTable()
	checkDevicesTable()
}
