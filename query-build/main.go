package main

import (
	"context"
	"database/sql"
	"log"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
)

var db *bun.DB

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID   int64  `bun:"id,pk,autoincrement"`
	Name string `bun:"name,notnull"`
}

var users = []User{
	{Name: "kim"},
	{Name: "song"},
	{Name: "park"},
	{Name: "han"},
}

func openDB() {
	sqldb, err := sql.Open(sqliteshim.ShimName, "file::memory:?cache=shared")
	if err != nil {
		panic(err)
	}

	db = bun.NewDB(sqldb, sqlitedialect.New())
}

func createTable() {
	_, err := db.NewCreateTable().
		Model((*User)(nil)).
		Exec(context.Background())
	if err != nil {
		panic(err)
	}
}

func insertUsers() {
	for _, user := range users {
		res, err := db.NewInsert().Model(&user).Exec(context.Background())
		log.Println("res:", res, "err:", err)
	}
}

func prepareDB() {
	openDB()
	createTable()
	insertUsers()
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	prepareDB()

	query := db.NewSelect().Model((*User)(nil))
	log.Println(query)
	query.Where("name =?", "kim")
	log.Println(query)

}
