package db

// SessionUpdateOneTime updates token for one-time access
func SessionUpdateOneTime(email string, accessToken string, refreshToken string) (err error) {
	err = SessionDeleteOneTime(email)
	if err != nil {
		return
	}

	err = SessionCreate(email, accessToken, refreshToken, true)

	return
}
