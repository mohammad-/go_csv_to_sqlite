package load_data

import (
	"bufio"
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func saveRequest(db *sql.DB, query string) (int64, error) {
	stmt, err := db.Prepare(query)
	if err != nil {
		return 0, err
	}
	r, err := stmt.Exec()
	defer stmt.Close()
	if err != nil {
		return 0, err
	}
	c, err := r.RowsAffected()
	if err != nil {
		return 0, err
	}
	return c, nil

}
func readFile(filename string, db *sql.DB) chan string {
	c := make(chan string)
	go func() {
		fmt.Println(filename)
		f, err := os.Open(filename)
		if err != nil {
			panic(err)
		}
		scanner := bufio.NewScanner(f)
		scanner.Scan()
		headers := scanner.Text()
		baseQuery := "INSERT OR IGNORE INTO userinfo(" + headers + ",file_name) values("
		for scanner.Scan() {
			t := scanner.Text()
			query := baseQuery
			for idx, _t := range strings.Split(t, ",") {
				if idx != 0 {
					query += ","
				}
				query += ("'" + _t + "'")
			}
			query = (query + ",'" + filename + "');")
			// fmt.Println(query)
			count, err := saveRequest(db, query)
			if err != nil {
				panic(err)
			}
			if count > 0 {
				c <- t
			}
		}
		go close(c)
	}()
	return c
}

func checkTableOrCreate(db *sql.DB) error {
	r, err := db.Query("SELECT count(name) as count FROM sqlite_master WHERE type='table' and name=='userinfo' ORDER BY name;")
	if err != nil {
		return err
	}
	var count int
	for r.Next() {
		r.Scan(&count)
	}
	r.Close()
	if count == 0 {
		fmt.Println("Creating Table....")
		_, err := db.Exec(`CREATE TABLE userinfo (
			createdAtCompound INTEGER,
			userId TEXT,
			height INTEGER,
			status TEXT,
			createdAt INTEGER,
			id TEXT,
			gender TEXT,
			weight INTEGER,
			age INTEGER,
			updatedAt INTEGER,
			errorCode TEXT,
			created_idx TEXT,
			errorDetail TEXT,
			time_diff TEXT,
			file_name TEXT,
			PRIMARY KEY (id, userId)
			);`)
		if err != nil {
			return err
		}

		_, err = db.Exec(`CREATE INDEX userid_created_data_idx ON userinfo (createdAt, userId);`)
		if err != nil {
			return err
		}
	}
	return nil
}

func LoadData(folderPath, dbPath string, db *sql.DB) error {
	checkTableOrCreate(db)
	files, err := ioutil.ReadDir(folderPath)
	if err != nil {
		return err
	}
	counter := 0
	for _, f := range files {
		if strings.HasSuffix(f.Name(), "csv") {
			c := readFile(folderPath+"/"+f.Name(), db)
			for _ = range c {
				counter += 1
			}
		}
	}
	fmt.Println("Loaded:", counter, "rows")
	return nil
}

// type request struct {
// 	CreatedAtCompound int64
// 	UserId            string
// 	Height            int64
// 	Status            string
// 	CreatedAt         int64
// 	Id                string
// 	Gender            string
// 	Weight            int64
// 	Age               int64
// 	UpdatedAt         int64
// 	ErrorCode         string
// 	CreatedIdx        string
// 	ErrorDetail       string
// 	TimeDiff          string
// }

// func (r *request) save(db *sql.DB) {
// 	stmt, err := db.Prepare("INSERT INTO userinfo values(?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
// 	checkErr(err)
// 	stmt.Exec()
// }

// func makeNewRequest(line string) (*request, error) {
// 	items := strings.Split(line, ",")
// 	r := new(request)

// 	t := reflect.ValueOf(r).Elem()
// 	for i := 0; i < len(items); i++ {
// 		f := t.FieldByIndex([]int{i})
// 		if f.Kind() == reflect.Int64 {
// 			v, err := strconv.ParseFloat(items[i], 64)
// 			checkErr(err)
// 			f.SetInt(int64(v))
// 		} else {
// 			f.SetString(items[i])
// 		}
// 	}
// 	return r, nil
// }
