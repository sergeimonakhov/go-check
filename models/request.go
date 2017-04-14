package models

import (
	"database/sql"
	"fmt"
)

//Tasks table
type Tasks struct {
	ID       int `json:"id"`
	chekerID int
	userID   int
	Interval int    `json:"interval"`
	Target   string `json:"target"`
	slackID  int
	status   bool
}

//Checker table
type Checker struct {
	CheckerID int `json:"id"`
	hashID    int
	status    bool
}

//GetTasksReq SELECT
func GetTasksReq(db *sql.DB, id int) ([]*Tasks, error) {
	rows, err := db.Query("SELECT * FROM tasks WHERE cheker_id = $1", id)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	bks := make([]*Tasks, 0)
	for rows.Next() {
		bk := new(Tasks)
		err = rows.Scan(&bk.ID, &bk.chekerID, &bk.userID, &bk.Interval, &bk.Target, &bk.slackID, &bk.status)
		if err != nil {
			return nil, err
		}
		bks = append(bks, bk)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return bks, nil
}

//GetCheckerID SELECT
func GetCheckerID(db *sql.DB, hash string) ([]*Checker, error) {
	rows, err := db.Query("SELECT * FROM chekers WHERE hash_id=(SELECT hash_id from hash WHERE hash=$1)", hash)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	bks := make([]*Checker, 0)
	for rows.Next() {
		bk := new(Checker)
		err = rows.Scan(&bk.CheckerID, &bk.hashID, &bk.status)
		if err != nil {
			return nil, err
		}
		bks = append(bks, bk)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return bks, nil
}

//InsertHash INSERT
func InsertHash(db *sql.DB, hash string) error {
	query := `INSERT INTO chekers (hash_id, status)
            VALUES ((SELECT hash_id from hash WHERE hash=$1), true)
            ON CONFLICT DO NOTHING;`

	stmt, err := db.Prepare(query)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	_, err = stmt.Exec(hash)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil
}

//UpdateStatus Update
func UpdateStatus(db *sql.DB, id int, status bool) error {
	query := `UPDATE tasks SET status = $2 WHERE tasks_id = $1;`

	stmt, err := db.Prepare(query)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	_, err = stmt.Exec(id, status)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil
}
