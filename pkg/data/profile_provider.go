package data

import (
	tinycloud "github.com/kucicm/tiny-cloud/pkg"
	_ "github.com/mattn/go-sqlite3"
)

func ListProfiles() (tinycloud.Profiles, error) {
	rows, err := db.Query("SELECT Id, Name, Description FROM profiles")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	profiles := make([]*tinycloud.Profile, 0)
	for rows.Next() {
		profile := &tinycloud.Profile{}
		if err = rows.Scan(&profile.Id, &profile.Name, &profile.Description); err == nil {
			profiles = append(profiles, profile)
		}

	}

	return profiles, nil
}

func AddProfile(newProfile *tinycloud.Profile) error {
	query := "INSERT INTO profiles (Name, Description) VALUES (?, ?);"
	_, err := db.Exec(query, newProfile.Name, newProfile.Description)
	return err
}
