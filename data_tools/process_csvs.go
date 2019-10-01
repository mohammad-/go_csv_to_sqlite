package main

import (
	"data_tools/load_data"
	"database/sql"
	"flag"
	"fmt"
	"os"

	// "time"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	mySet := flag.NewFlagSet("", flag.ExitOnError)
	dirPath := mySet.String("dir", "", "Path where CSV is stored")
	dbName := mySet.String("db", "", "Dirrectory where CSV is stored")

	startDate := mySet.String("start", "", "start date YYYY-MM-DD")
	endDate := mySet.String("end", "", "end date YYYY-MM-DD")

	if os.Args[1] == "loaddata" {
		mySet.Parse(os.Args[2:])
		dbpath := *dbName
		db, err := sql.Open("sqlite3", "file:"+dbpath+"?cache=shared&mode=rwc")
		checkErr(err)
		defer func() {
			db.Close()
		}()
		load_data.LoadData(*dirPath, dbpath, db)
	} else if os.Args[1] == "usercount" {
		mySet.Parse(os.Args[2:])
		dbpath := *dbName
		if dbpath == "" {
			panic("Invalid DB Path")
		}
		db, err := sql.Open("sqlite3", "file:"+dbpath+"?cache=shared&mode=rwc")
		checkErr(err)
		defer func() {
			db.Close()
		}()
		s := *startDate + " 00:00:00"
		e := *endDate + " 23:59:59"
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
		checkErr(err)
		var count int
		for rows.Next() {
			err = rows.Scan(&count)
			checkErr(err)
			fmt.Println("New User Count", count)
		}
		rows.Close()
	}

}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

// // insert
// stmt, err := db.Prepare("INSERT INTO userinfo(username, departname, created) values(?,?,?)")
// checkErr(err)

// res, err := stmt.Exec("astaxie", "研发部门", "2012-12-09")
// checkErr(err)

// id, err := res.LastInsertId()
// checkErr(err)

// fmt.Println(id)
// // update
// stmt, err = db.Prepare("update userinfo set username=? where uid=?")
// checkErr(err)

// res, err = stmt.Exec("astaxieupdate", id)
// checkErr(err)

// affect, err := res.RowsAffected()
// checkErr(err)

// fmt.Println(affect)

// // query
// rows, err := db.Query("SELECT * FROM userinfo")
// checkErr(err)
// var uid int
// var username string
// var department string
// var created time.Time

// for rows.Next() {
// 	err = rows.Scan(&uid, &username, &department, &created)
// 	checkErr(err)
// 	fmt.Println(uid)
// 	fmt.Println(username)
// 	fmt.Println(department)
// 	fmt.Println(created)
// }

// rows.Close() //good habit to close

// // delete
// stmt, err = db.Prepare("delete from userinfo where uid=?")
// checkErr(err)

// res, err = stmt.Exec(id)
// checkErr(err)

// affect, err = res.RowsAffected()
// checkErr(err)

// fmt.Println(affect)

// db.Close()
