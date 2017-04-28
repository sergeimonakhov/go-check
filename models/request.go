package models

import (
	"database/sql"
	"log"
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

//FailOnRequest func
func FailOnRequest(err error) error {
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

//GetTasksReq SELECT all task from checker
func GetTasksReq(db *sql.DB, id int) ([]*Tasks, error) {
	rows, err := db.Query("SELECT * FROM tasks WHERE cheker_id = $1", id)
	FailOnRequest(err)
	defer rows.Close()

	bks := make([]*Tasks, 0)
	for rows.Next() {
		bk := new(Tasks)
		err = rows.Scan(&bk.ID, &bk.chekerID, &bk.userID, &bk.Interval, &bk.Target, &bk.slackID, &bk.status)
		FailOnRequest(err)
		bks = append(bks, bk)
	}
	err = rows.Err()
	FailOnRequest(err)

	return bks, nil
}

//GetCheckerID SELECT cheker_id with hash_id
func GetCheckerID(db *sql.DB, hash string) (Checker, error) {
	stmt, err := db.Prepare("SELECT * FROM chekers WHERE hash_id=(SELECT hash_id from hash WHERE hash=$1)")
	FailOnRequest(err)

	var id Checker
	err = stmt.QueryRow(hash).Scan(&id.CheckerID, &id.hashID, &id.status)
	FailOnRequest(err)

	return id, nil
}

//InsertHash INSERT new clients
func InsertHash(db *sql.DB, hash string) error {
	query := `INSERT INTO chekers (hash_id, status)
            VALUES ((SELECT hash_id from hash WHERE hash=$1), true)
            ON CONFLICT DO NOTHING;`

	stmt, err := db.Prepare(query)
	FailOnRequest(err)

	_, err = stmt.Exec(hash)
	FailOnRequest(err)

	return nil
}

//UpdateStatus Update targer status
func UpdateStatus(db *sql.DB, id int, status bool) error {
	stmt, err := db.Prepare("UPDATE tasks SET status = $2 WHERE tasks_id = $1;")
	FailOnRequest(err)

	_, err = stmt.Exec(id, status)
	FailOnRequest(err)

	return nil
}
