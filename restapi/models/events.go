package models

import (
	"project/restapi/db"
	"time"
)

type Event struct {
	ID          int64
	Name        string `binding:"required"`
	Description string
	Location    string
	DateTime    time.Time
	UserID      int
}

func (e Event) Save() error {
	query := `INSERT INTO events (name, description, location, dateTime, userID)
	VALUES (?,?,?,?,?)
	`
	sql_smt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer sql_smt.Close()

	result, err := sql_smt.Exec(e.Name, e.Description, e.Location, e.DateTime, e.UserID)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	e.ID = id

	return err
}

func (e Event) Update() error{
	query := `
	UPDATE events
	SET name = ?, description = ?, location = ?, dateTime = ?
	WHERE id = ?
	`

	sql_smt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer sql_smt.Close()

	_, err = sql_smt.Exec(e.Name, e.Description, e.Location, e.DateTime, e.ID)

	return err

}

func GetAllEvents() ([]Event, error) {
	query := `SELECT * FROM events`

	rows, err := db.DB.Query(query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var events []Event

	for rows.Next() {
		var event Event

		err := rows.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &event.DateTime, &event.UserID)

		if err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	return events, nil
}

func GetAllEventByID(id int64) (*Event, error) {
	query := `SELECT * FROM events WHERE id = ?`

	rows := db.DB.QueryRow(query, id)

	var event Event

	err := rows.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &event.DateTime, &event.UserID)

	if err != nil {
		return nil, err
	}

	return &event, nil
}


func DeleteEventByID(id int64) error{
	query := `DELETE FROM events WHERE id = ?`

	sql_smt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer sql_smt.Close()

	_, err = sql_smt.Exec(id)

	if err != nil {
		return err
	}

	return nil
}