package data

func GetNewRunId(profileName string) (string, error) {
	var id int
	query := "SELECT Id FROM Profiles WHERE Name = ?;"
	if err := db.QueryRow(query, profileName).Scan(&id); err != nil {
		return "", err
	}

	var runId int
	query = "INSERT INTO RunLogs (ProfileId) VALUES (?) RETURNING Id;"
	if err := db.QueryRow(query, id).Scan(&runId); err != nil {
		return "", err
	}

	var runIdHuman string
	query = "SELECT RunIdHuman FROM v_runIds WHERE RunId = ?;"
	if err := db.QueryRow(query, runId).Scan(&runIdHuman); err != nil {
		return "", err
	}

	return runIdHuman, nil
}
