package data

import (
	"database/sql"
	"fmt"

	tinycloud "github.com/kucicm/tiny-cloud/pkg"
)

// save new profile to db
// returns an error if profile with the same name already exists
func CreateProfile(profile *tinycloud.Profile) error {
	if err := profile.Valid(); err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = insertNewProfile(tx, profile)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func insertNewProfile(tx *sql.Tx, profile *tinycloud.Profile) error {
	query := "INSERT INTO Profiles (Name, Description) VALUES (?, ?)"
	res, err := tx.Exec(query, profile.Name, profile.Description)
	if err != nil {
		return err
	}

	profileId, err := res.LastInsertId()
	if err != nil {
		return err
	}

	settigns := profile.Settings
	cloud := settigns.ResolveCloudName()
	switch cloud {
	case "aws":
		aws := settigns.Aws
		query = `INSERT INTO AwsSettings (ProfileId, Region, AccessKey, SecretAccessKey)
			VALUES (?, ?, ?, ?);`
		_, err = tx.Exec(
			query,
			profileId,
			aws.AwsRegion,
			aws.AwsAccessKeyId,
			aws.AwsSeacretAccessKey,
		)
		return err
	default:
		return fmt.Errorf("unsupported cloud '%s'", cloud)
	}
}

// return all profiles from database
func GetProfiles() (tinycloud.Profiles, error) {
	rows, err := db.Query(`
	SELECT 
		Name, Description, Active,
		AwsRegion, AwsAccessKey, AwsSecretAccessKey
	FROM v_profiles;`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	profiles := make([]*tinycloud.Profile, 0)
	for rows.Next() {
		settings := &tinycloud.CloudSettings{}
		aws := &tinycloud.AwsSettings{}
		profile := &tinycloud.Profile{Settings: settings}
		err = rows.Scan(
			&profile.Name, &profile.Description, &profile.Active,
			&aws.AwsRegion, &aws.AwsAccessKeyId, &aws.AwsSeacretAccessKey,
		)
		if aws.Valid() == nil {
			settings.Aws = aws
		}
		if err != nil {
			return nil, err
		}
		profiles = append(profiles, profile)
	}

	return profiles, nil
}

// if profile dose not exist, returns error
// else update active of all profiles to false except profile with the given name
func UpdateProfileToActive(profileName string) error {
	var id int
	err := db.QueryRow("SELECT Id FROM v_profiles WHERE Name = ?;", profileName).Scan(&id)
	if err == sql.ErrNoRows {
		return fmt.Errorf("no profile with name '%s'", profileName)
	}

	if err != nil {
		return err
	}

	_, err = db.Exec("UPDATE Profiles SET Active = (Id = ?);", id)
	return err
}

// returns account which is marked as active or
// if non returns newest account
func GetActiveProfile() (*tinycloud.Profile, error) {
	query := `
	SELECT 
		Name, Description, Active,
		AwsRegion, AwsAccessKey, AwsSecretAccessKey
	FROM v_profiles
	WHERE Active = 1`

	settings := &tinycloud.CloudSettings{}
	aws := &tinycloud.AwsSettings{}
	profile := &tinycloud.Profile{Settings: settings}
	err := db.QueryRow(query).Scan(
		&profile.Name, &profile.Description, &profile.Active,
		&aws.AwsRegion, &aws.AwsAccessKeyId, &aws.AwsSeacretAccessKey,
	)

	if aws.Valid() == nil {
		settings.Aws = aws
	}
	// else if gcp etc

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("there are no profiles, create new one")
	}

	if err != nil {
		return nil, err
	}
	return profile, nil
}

// deletes profile and settings,
// if profile dose not exists it returns an error
func DeleteProfile(profileName string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	err = deleteProfile(tx, profileName)

	if err == sql.ErrNoRows {
		return fmt.Errorf("profile with name '%s' does not exists", profileName)
	}

	if err != nil {
		return err
	}

	return tx.Commit()
}

func deleteProfile(tx *sql.Tx, profileName string) error {
	query := "SELECT Id FROM Profiles WHERE Name = ?;"
	var profileId int
	if err := tx.QueryRow(query, profileName).Scan(&profileId); err != nil {
		return err
	}

	if _, err := tx.Exec("DELETE FROM Profiles WHERE Id = ?", profileId); err != nil {
		return err
	}

	if _, err := tx.Exec("DELETE FROM AwsSettings WHERE ProfileId = ?", profileId); err != nil {
		return err
	}
	return nil
}

func DoseProfileExists(profileName string) (bool, error) {
	var id int
	err := db.QueryRow("SELECT Id FROM v_profiles WHERE Name = ?;", profileName).Scan(&id)
	if err == sql.ErrNoRows {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}
