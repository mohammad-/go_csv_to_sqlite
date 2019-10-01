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
	dbName := mySet.String("db", "", "Path to sqlite file")
	showCount := mySet.Bool("show-count", false, "At the end show number of items displayed")

	startDate := mySet.String("start", "", "start date YYYY-MM-DD")
	endDate := mySet.String("end", "", "end date YYYY-MM-DD")
	if len(os.Args) < 2 {
		fmt.Println("Commands loaddata, usercount, list")
		mySet.Usage()
	} else if os.Args[1] == "loaddata" {
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
		s := *startDate + " 00:00:00"
		e := *endDate + " 23:59:59"
		count, err := load_data.CountUsers(s, e, *dbName)
		checkErr(err)
		fmt.Println("New User Count", count)
	} else if os.Args[1] == "list" {
		mySet.Parse(os.Args[2:])
		s := *startDate + " 00:00:00"
		e := *endDate + " 23:59:59"
		count, err := load_data.ListRequest(s, e, *dbName)
		checkErr(err)
		if *showCount {
			fmt.Println("Total:", count)
		}
	} else {
		mySet.Usage()
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
