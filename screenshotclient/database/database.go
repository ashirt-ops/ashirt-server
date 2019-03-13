// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package database

import (
	"database/sql"
	"log"
	"strconv"
	"time"

	"github.com/theparanoids/ashirt/campaign"
	"github.com/theparanoids/ashirt/evidence"

	"github.com/getlantern/errors"
)

// ErrNoDatabase is returned when the database does not exist
var (
	ErrNoDatabase         = errors.New("database not found")
	ErrScreenshotNotFound = errors.New("screenshot not found")
)

// SetupDB creates and initializes the database for the persistent queue if one
// does not exist
func SetupDB() error {
	database, err := sql.Open("sqlite3", "ashirt.db")
	if err != nil {
		return err
	}
	defer database.Close()

	statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS screenshots (id INTEGER PRIMARY KEY, FileName TEXT, FullPath TEXT, Description TEXT, CampaignID TEXT, CampaignName TEXT, SHA256HASH TEXT, OccurrenceTimestamp INTEGER, SubmittedTimestamp INTEGER,Error INTEGER)")
	if err != nil {
		return err
	}

	_, err = statement.Exec()
	if err != nil {
		return err
	}

	return nil
}

// InsertScreenshot inserts a screenshot into the persistent queue
func InsertScreenshot(screenshot evidence.Screenshot) error {
	database, err := sql.Open("sqlite3", "ashirt.db")
	if err != nil {
		return err
	}
	defer database.Close()

	statement, err := database.Prepare("INSERT INTO screenshots (FileName, FullPath, Description,CampaignID,CampaignName,SHA256HASH,OccurrenceTimestamp) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}

	_, err = statement.Exec(screenshot.FileName, screenshot.FullPath, screenshot.Description, screenshot.Campaign.ID, screenshot.Campaign.Name, screenshot.FileHash, screenshot.OccurrenceTimestamp)
	if err != nil {
		return err
	}

	return nil
}

//UpdateScreenshots updates a screenshot object
func UpdateScreenshots(screenshot evidence.Screenshot) error {
	database, err := sql.Open("sqlite3", "ashirt.db")
	if err != nil {
		return err
	}
	statement, err := database.Prepare("UPDATE screenshots set FileName = ?, FullPath =?, Description =?, CampaignID = ?, CampaignName = ? WHERE ID= ?")
	if err != nil {
		return err
	}
	_, err = statement.Exec(screenshot.FileName, screenshot.FullPath, screenshot.Description, screenshot.Campaign.ID, screenshot.Campaign.Name, screenshot.ID)
	if err != nil {
		return err
	}

	return nil
}

// ListScreenshots returns a list of all screenshots in the queue
func ListScreenshots() ([]evidence.Screenshot, error) {
	database, err := sql.Open("sqlite3", "ashirt.db")
	if err != nil {
		return nil, err
	}
	defer database.Close()

	rows, err := database.Query("SELECT id,FileName,FullPath,Description,OccurrenceTimestamp FROM screenshots")
	if err != nil {
		return nil, err
	}

	var screenshots []evidence.Screenshot

	for rows.Next() {
		var screenshot evidence.Screenshot

		err = rows.Scan(&screenshot.ID, &screenshot.FileName, &screenshot.FullPath, &screenshot.Description, &screenshot.OccurrenceTimestamp)
		if err != nil {
			log.Println(err)
			continue
		}

		screenshots = append(screenshots, screenshot)
	}

	return screenshots, nil
}

// GetUnsubmitted returns a list of all screenshots that have not yet been
// submitted
func GetUnsubmitted() ([]evidence.Screenshot, error) {
	database, err := sql.Open("sqlite3", "ashirt.db")
	if err != nil {
		return nil, err
	}
	defer database.Close()

	rows, err := database.Query("SELECT id,FileName,FullPath,Description,CampaignID,CampaignName,SHA256HASH,OccurrenceTimestamp FROM screenshots where SubmittedTimestamp is null and Error is null")
	if err != nil {
		return nil, err
	}

	var screenshots []evidence.Screenshot

	for rows.Next() {
		var screenshot evidence.Screenshot
		var campaign campaign.Campaign

		err = rows.Scan(&screenshot.ID, &screenshot.FileName, &screenshot.FullPath, &screenshot.Description, &campaign.ID, &campaign.Name, &screenshot.FileHash, &screenshot.OccurrenceTimestamp)
		if err != nil {
			log.Println(err)
			continue
		}

		screenshot.Campaign = campaign
		screenshots = append(screenshots, screenshot)
	}

	return screenshots, nil
}

// MarkSubmitted marks a specific screenshot as having been submitted to the
// ASHIRT api server
func MarkSubmitted(id int) error {
	database, err := sql.Open("sqlite3", "ashirt.db")
	if err != nil {
		return err
	}
	defer database.Close()

	statement, err := database.Prepare("update screenshots set SubmittedTimestamp =?  where id=?")
	if err != nil {
		return err
	}

	var time = time.Now().Unix()
	_, err = statement.Exec(strconv.FormatInt(time, 10), id)
	if err != nil {
		return err
	}

	return nil
}

// MarkError marks that a specific screenshot encountered and error when
// submitting to the ASHIRT api
func MarkError(id int) error {
	database, err := sql.Open("sqlite3", "ashirt.db")
	if err != nil {
		return err
	}
	defer database.Close()

	statement, err := database.Prepare("update screenshots set ERROR =1  where id=?")
	if err != nil {
		return err
	}

	_, err = statement.Exec(id)
	if err != nil {
		return err
	}

	return nil
}
