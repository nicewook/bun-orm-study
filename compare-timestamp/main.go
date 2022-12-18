package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

type MyTime struct {
	bun.BaseModel `bun:"table:my_time,alias:u"`
	ID            int64     `bun:"id,pk,autoincrement"`
	Name          string    `bun:"name"`
	TStamp        time.Time `bun:"t_stamp"`
	TStampTZ      time.Time `bun:"t_stamptz"`
}

func printTime(t MyTime) {
	fmt.Printf("timestamp:   %v\n", t.TStamp)
	fmt.Printf("timestamptz: %v\n", t.TStampTZ)
	fmt.Println("---")
}

func main() {

	// connect
	dsn := "postgres://postgres:@localhost:5432/hsjeong?sslmode=disable"
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqldb, pgdialect.New())

	// prepare data
	ctime := time.Now()

	name := "test1"
	myTime := MyTime{
		Name:     name,
		TStamp:   ctime,
		TStampTZ: ctime,
	}

	printTime(myTime)
	// timestamp:   2022-12-18 11:14:42.443372 +0900 KST m=+0.003532209
	// timestamptz: 2022-12-18 11:14:42.443372 +0900 KST m=+0.003532209

	if _, err := db.NewInsert().Model(&myTime).
		Where("name = ?", name).
		Exec(context.Background()); err != nil {
		fmt.Println(err)
		return
	}

	if err := db.NewSelect().Model(&myTime).Scan(context.Background()); err != nil {
		fmt.Println(err)
	}

	printTime(myTime)
	// timestamp:   2022-12-18 01:41:49.921354 +0000 UTC
	// timestamptz: 2022-12-18 10:41:49.921354 +0900 KST

	// check timezone in Postgres
	var tz string
	db.NewRaw("SHOW TIMEZONE").Scan(context.TODO(), &tz)
	fmt.Println("current timezone:", tz)
	fmt.Println("---")
	// current timezone: Asia/Seoul

	// location test
	locat, _ := time.LoadLocation("Asia/Kolkata")
	ctime = time.Now().In(locat)

	locatName := "test2"
	myTime = MyTime{
		Name:     locatName,
		TStamp:   ctime,
		TStampTZ: ctime,
	}

	fmt.Printf("Asia/Kolkata:%s\n", ctime)
	// Asia/Kolkata:2022-12-18 07:44:42.468642 +0530 IST

	if _, err := db.NewInsert().
		Model(&myTime).
		Exec(context.Background()); err != nil {
		fmt.Println(err)
	}

	if err := db.NewSelect().Model(&myTime).
		Where("name = ?", locatName).
		Scan(context.Background()); err != nil {
		fmt.Println(err)
	}
	printTime(myTime)
	// timestamp:   2022-12-18 02:09:21.439701 +0000 UTC
	// timestamptz: 2022-12-18 11:09:21.439701 +0900 KST
}
