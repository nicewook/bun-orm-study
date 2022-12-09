package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
)

var db *bun.DB

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID        int64     `bun:"id,pk,autoincrement"`
	Name      string    `bun:"name,notnull"`
	CreatedAt time.Time `bun:"createdAt,nullzero,notnull,default:current_timestamp"`
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
		_, err := db.NewInsert().Model(&user).Exec(context.Background())
		if err != nil {
			panic(err)
		}
		time.Sleep(time.Second)
	}
}

func prepareDB() {
	openDB()
	createTable()
	insertUsers()
}

func main() {
	prepareDB()

	var createdTimes []time.Time
	if err := db.NewSelect().Model((*User)(nil)).
		Column("createdAt").
		Scan(context.Background(), &createdTimes); err != nil {
		panic(err)
	}

	for i, ctime := range createdTimes {
		fmt.Printf("record #%d, created at %v\n", i+1, ctime)
	}
}

// 아래와 같이 INSERT 된 시간이 출력된다.
// record #1, created at 2022-12-09 17:26:21 +0000 UTC
// record #2, created at 2022-12-09 17:26:22 +0000 UTC
// record #3, created at 2022-12-09 17:26:23 +0000 UTC
// record #4, created at 2022-12-09 17:26:24 +0000 UTC
