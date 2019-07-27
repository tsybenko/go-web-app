package feedback

import (
	"database/sql"
	"fmt"
)

type Feedback struct {
	Id    int
	Email string
	Text  string
}

const table string = "feedback"

// TODO realize field selection
func GetAll(db *sql.DB) map[int]Feedback {
	query := fmt.Sprintf("SELECT id, email, text FROM `%s` ORDER BY id DESC", table)
	rows, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
	items := make(map[int]Feedback)
	for rows.Next() {
		var feedback Feedback
		err = rows.Scan(&feedback.Id, &feedback.Email, &feedback.Text)
		if err != nil {
			panic(err.Error())
		}
		items[feedback.Id] = feedback
	}
	return items
}

func Create(db *sql.DB, feedback *Feedback) bool {
	query := fmt.Sprintf("INSERT INTO `%s` (email, text) VALUES (\"%s\", \"%s\")", table, feedback.Email, feedback.Text)
	_, err := db.Exec(query)
	if err != nil {
		panic(err.Error())
	} else {
		return true
	}
	return false
}

func Delete(db *sql.DB, id int) {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = %d", table, id)
	_, err := db.Exec(query)
	if err != nil {
		panic(err.Error())
	}
}
