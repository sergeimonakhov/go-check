package models

import (
	"database/sql"
	"fmt"
)

//Tasks table
type Tasks struct {
	id       int
	chekerID int
	userID   int
	Interval int    `json:"interval"`
	Target   string `json:"target"`
	slackID  int
	Status   bool `json:"status"`
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
		err = rows.Scan(&bk.id, &bk.chekerID, &bk.userID, &bk.Interval, &bk.Target, &bk.slackID, &bk.Status)
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
