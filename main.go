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
	flag "github.com/spf13/pflag"
)

var (
	dbConn string
	dbType string
	reset  bool
)

func init() {
	gofakeit.Seed(0)

	flag.StringVar(&dbConn, "conn", "", "db connection string")
	flag.StringVar(&dbType, "type", "", "db type (mysql/postrgres")
	flag.BoolVar(&reset, "reset", false, "reset (default false")

	flag.Parse()

	log.Printf(`
	______________________
	Initialize Chatter-Box
	DB: %s 
	Type: %s
	Reset mode: %v
	______________________
	`, dbConn, dbType, reset)
}

func main() {
	c := New()
	if reset {
		c.DropTable()
		c.CreateTable()
	}
	for {
		c.CreateData()
		c.ReadData()
	}
}

type chatter struct {
	DbConnString string
	Db           *sql.DB
	MultiThread  bool
}

func New() (c chatter) {
	db, err := sql.Open(dbType, dbConn)
	if err != nil {
		log.Fatalln(err)
	}

	c.Db = db
	return
}

func (c *chatter) CreateTable() {

	_, errCreate := c.Db.Exec(`
	CREATE TABLE authors (
		id INT(11) NOT NULL AUTO_INCREMENT,
		first_name VARCHAR(50) NOT NULL,
		last_name VARCHAR(50) NOT NULL,
		email VARCHAR(100) NOT NULL,
		birthdate DATE NOT NULL,
		added TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (id),
		UNIQUE INDEX email (email)
	)`)
	if errCreate != nil {
		log.Fatalln(errCreate)
	}
}

func (c *chatter) DropTable() {
	_, errDrop := c.Db.Exec(`
	DROP TABLE IF EXISTS authors`)
	if errDrop != nil {
		log.Fatalln(errDrop)
	}
}

func (c *chatter) CreateData() {
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
		ToSql()

	_, errInsert := c.Db.Exec(insertStmt, args...)
	if errInsert != nil {
		log.Println(errInsert)
	}
}

func (c *chatter) ReadData() {
	rows, err := c.Db.Query(`select id from authors order by id desc limit 1`)
	if err != nil {
		log.Fatalln(err)
	}
	defer rows.Close()
	for rows.Next() {
		var (
			id int
		)

		errScan := rows.Scan(&id)
		if errScan != nil {
			log.Fatalln(errScan)
		}
		fmt.Printf("%d, ", id)
	}
}
