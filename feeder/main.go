package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/brianvoe/gofakeit"
	_ "github.com/go-sql-driver/mysql"
	"github.com/rs/xid"
)

type Author struct {
	FirstName string
	LastName  string
	Email     string
	Birthdate time.Time
	Added     time.Time
}

type Post struct {
	AuthorID    int
	Title       string
	Description string
	Content     string
	Date        time.Time
}

func main() {

	host := "aurora-mysql.cluster-cprfdqkte67x.ap-southeast-1.rds.amazonaws.com"
	user := "admin"
	password := "Labmysql1!"
	dbName := "dummy"
	connection := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", user, password, host, dbName)

	db, err := sql.Open("mysql", connection)
	if err != nil {
		panic(err)
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	gofakeit.Seed(0)

	feeder := Feeder{
		Db: db,
	}
	for {
		authorCounter := feeder.countAuthor()
		authorLimit := authorCounter + 1
		for authorCounter < authorLimit {
			feeder.generateAuthor()
			authorCounter++
			fmt.Printf("author counter: %d\n", authorCounter)
		}
		break
		// postCounter := feeder.countPost()
		// postLimit := postCounter + 10000
		// for postCounter < postLimit {
		// 	feeder.generatePost(authorLimit)
		// 	postCounter++
		// 	fmt.Printf("post counter: %d\n", postCounter)
		// }
	}
}

type Feeder struct {
	Db *sql.DB
}

func (f *Feeder) generateAuthor() {
	insertStmt, args, _ := sq.
		Insert("authors").Columns(
		"first_name",
		"last_name",
		"email",
		"birthdate",
		"added").
		Values(
			gofakeit.FirstName(),
			gofakeit.LastName(),
			fmt.Sprintf("%s-%s", xid.New().String(), gofakeit.Email()),
			gofakeit.Date(),
			time.Now()).
		Values(
			gofakeit.FirstName(),
			gofakeit.LastName(),
			fmt.Sprintf("%s-%s", xid.New().String(), gofakeit.Email()),
			gofakeit.Date(),
			time.Now()).
		Values(
			gofakeit.FirstName(),
			gofakeit.LastName(),
			fmt.Sprintf("%s-%s", xid.New().String(), gofakeit.Email()),
			gofakeit.Date(),
			time.Now()).
		Values(
			gofakeit.FirstName(),
			gofakeit.LastName(),
			fmt.Sprintf("%s-%s", xid.New().String(), gofakeit.Email()),
			gofakeit.Date(),
			time.Now()).
		Values(
			gofakeit.FirstName(),
			gofakeit.LastName(),
			fmt.Sprintf("%s-%s", xid.New().String(), gofakeit.Email()),
			gofakeit.Date(),
			time.Now()).
		ToSql()

	_, errInsert := f.Db.Exec(insertStmt, args...)
	if errInsert != nil {
		log.Println(errInsert)
	}
}

func (f *Feeder) generatePost(totalAuthor int) {
	insertStmt, args, _ := sq.
		Insert("posts").Columns(
		"author_id",
		"title",
		"description",
		"content",
		"date").
		Values(
			gofakeit.Number(0, totalAuthor),
			gofakeit.Sentence(200),
			gofakeit.Sentence(400),
			gofakeit.Sentence(800),
			time.Now()).
		Values(
			gofakeit.Number(0, totalAuthor),
			gofakeit.Sentence(200),
			gofakeit.Sentence(400),
			gofakeit.Sentence(800),
			time.Now()).
		Values(
			gofakeit.Number(0, totalAuthor),
			gofakeit.Sentence(200),
			gofakeit.Sentence(400),
			gofakeit.Sentence(800),
			time.Now()).
		Values(
			gofakeit.Number(0, totalAuthor),
			gofakeit.Sentence(200),
			gofakeit.Sentence(400),
			gofakeit.Sentence(800),
			time.Now()).
		Values(
			gofakeit.Number(0, totalAuthor),
			gofakeit.Sentence(200),
			gofakeit.Sentence(400),
			gofakeit.Sentence(800),
			time.Now()).
		Values(
			gofakeit.Number(0, totalAuthor),
			gofakeit.Sentence(200),
			gofakeit.Sentence(400),
			gofakeit.Sentence(800),
			time.Now()).
		ToSql()

	_, errInsert := f.Db.Exec(insertStmt, args...)
	if errInsert != nil {
		log.Println(errInsert)
	}
	return
}

func (f *Feeder) countAuthor() (authorTotal int) {
	row := f.Db.QueryRow("SELECT COUNT(*) FROM authors")
	err := row.Scan(&authorTotal)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func (f *Feeder) countPost() (postTotal int) {
	row := f.Db.QueryRow("SELECT COUNT(*) FROM posts")
	err := row.Scan(&postTotal)
	if err != nil {
		log.Fatal(err)
	}
	return
}
