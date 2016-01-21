package main

import (
	"bufio"

	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", "host=/var/run/postgresql/ dbname=ningyoudb sslmode=disable")
	if err != nil {
		panic(err)
	}

	types := []string{
		"main",
		"syn",
		"short",
		"official",
	}

	resp, err := http.Get("http://anidb.net/api/anime-titles.dat.gz")
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	scanner.Split(bufio.ScanLines)

	var line int
	const query = `INSERT INTO ningyou_titles_anidb (site, show_id, title, type, lang) VALUES ($1, $2, $3, $4, $5)`

	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}
	_, err = tx.Exec("TRUNCATE table ningyou_titles_anidb")
	if err != nil {
		fmt.Println(err)
	}
	for scanner.Scan() {
		line++
		data := strings.Split(scanner.Text(), "|")
		if line < 4 {
			continue
		}
		t, _ := strconv.Atoi(data[1])
		_, err = tx.Exec(query, "anidb", data[0], data[3], types[t-1], data[2])
		if err != nil {
			fmt.Println(err)
		}

	}
	err = tx.Commit()
	if err != nil {
		panic(err)
	}
}
