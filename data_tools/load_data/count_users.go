package load_data

import (
	"database/sql"
	"fmt"
)

func CountUsers(startDate, endDate, dbPath string) (int, error) {
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
	query := fmt.Sprintf(`select count(distinct userId)
	from userinfo
	where userId not in
		(
			select distinct userId from userinfo where
			datetime(createdAt, 'unixepoch', 'localtime') < '%s'
		)
	and status=='completed'
	and datetime(createdAt, 'unixepoch', 'localtime') < '%s';`, s, e)
	// fmt.Println(query)
	rows, err := db.Query(query)
	if err != nil {
		return 0, err
	}

	var count int
	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			return 0, err
		}
	}
	rows.Close()
	return count, nil
}
