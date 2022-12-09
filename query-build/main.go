package main

import (
	"context"
	"database/sql"
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

func main() {
	prepareDB()

	// Query build 방법 1
	query := db.NewSelect().Model((*User)(nil))
	fmt.Printf("memory address %p, query: %v\n", query, query)

	query = query.Where("name =?", "kim") // Where() 메서드의 결과를 왼쪽으로 대입한다
	fmt.Printf("memory address %p, query: %v\n", query, query)

	// Query build 방법 2
	query = db.NewSelect().Model((*User)(nil))
	fmt.Printf("memory address %p, query: %v\n", query, query)

	query.Where("name =?", "kim") // 대입을 시키지 않고 단순히 Where() 만을 실행한다.
	fmt.Printf("memory address %p, query: %v\n", query, query)
}

// 결과를 보면
// 1. query를 추가해주어도 두 경우 모두 메모리 주소가 바뀌지 않고, 생성된 query string 도 동일함을 알 수 있다.
// memory address 0x140002e8000, query: SELECT "u"."id", "u"."name" FROM "users" AS "u"
// memory address 0x140002e8000, query: SELECT "u"."id", "u"."name" FROM "users" AS "u" WHERE (name ='kim')
// memory address 0x140002e81e0, query: SELECT "u"."id", "u"."name" FROM "users" AS "u"
// memory address 0x140002e81e0, query: SELECT "u"."id", "u"."name" FROM "users" AS "u" WHERE (name ='kim')
