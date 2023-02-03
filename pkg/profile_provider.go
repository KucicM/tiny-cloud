package tinycloud

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func init() {
	err := SetupProfileTabel()
	if err != nil {
		log.Fatalln(err)
	}
}

func ListProfiles() ([]*Profile, error) {
	db, err := sql.Open("sqlite3", DbPath)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT Id, Name, Description FROM profiles")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	profiles := make([]*Profile, 0)
	for rows.Next() {
		profile := &Profile{}
		if err = rows.Scan(&profile.Id, &profile.Name, &profile.Description); err == nil {
			profiles = append(profiles, profile)
		}

	}

	return profiles, nil
}

func AddProfile(newProfile *Profile) error {
	db, err := sql.Open("sqlite3", DbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	query := "INSERT INTO profiles (Name, Description) VALUES (?, ?);"
	_, err = db.Exec(query, newProfile.Name, newProfile.Description)
	return err
}
