package load_data

import (
	"database/sql"
	"fmt"
	"time"
)

func ListRequest(startDate, endDate, dbPath string) (int, error) {
	if dbPath == "" {
		panic("Invalid DB Path")
	}
	db, err := sql.Open("sqlite3", "file:"+dbPath+"?cache=shared&mode=rwc")
	if err != nil {
		return 0, err
	}

	defer func() {
		db.Close()
	}()
	s := startDate + " 00:00:00"
	e := endDate + " 23:59:59"
	query := fmt.Sprintf(`select userId, id, createdAt, height, weight, age, gender
	from userinfo
	where userId not in
		(
			select distinct userId from userinfo where
			datetime(createdAt, 'unixepoch', 'localtime') < '%s'
		)
	and status=='completed'
	and datetime(createdAt, 'unixepoch', 'localtime') < '%s';`, s, e)
	rows, err := db.Query(query)
	if err != nil {
		return 0, err
	}
	var count int = 0
	var userId, id, height, weight, age, gender string
	var createdAt int64
	fmt.Println("userId, id, createdAt, height, weight, age, gender")
	for rows.Next() {
		err = rows.Scan(&userId, &id, &createdAt, &height, &weight, &age, &gender)
		if err != nil {
			return 0, err
		}
		dateString := time.Unix(createdAt, 0).Format("2006-01-02T15:04:05-0700")

		fmt.Printf("%s, %s, %s, %s, %s, %s, %s\n", userId, id, dateString, height, weight, age, gender)
		count += 1
	}
	rows.Close()
	return count, nil
}
