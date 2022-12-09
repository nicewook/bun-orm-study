package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

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
		_, err := db.NewInsert().Model(&user).Exec(context.Background())
		if err != nil {
			panic(err)
		}
	}
}

func prepareDB() {
	openDB()
	createTable()
	insertUsers()
}

func prettyJson(object interface{}) string {
	b, _ := json.MarshalIndent(object, "", "  ")
	return string(b)
}

func main() {
	prepareDB()

	var myUsers []User

	// 전체 User를 query 한다. 4명의 User가 모두 Scan되고 Count값도 4가 된다.
	query := db.NewSelect().Model((*User)(nil))
	count, err := query.ScanAndCount(context.Background(), &myUsers)
	fmt.Printf("count: %v, err: %v\n", count, err)
	fmt.Println(prettyJson(myUsers))

	// 이번에는 Limit를 1로 두었다. 1명의 User만이 Scan되지만 Count값이 여전히 전체 개수인 4이다.
	count, err = query.Limit(1).ScanAndCount(context.Background(), &myUsers)
	fmt.Printf("count: %v, err: %v\n", count, err)
	fmt.Println(prettyJson(myUsers))

}
