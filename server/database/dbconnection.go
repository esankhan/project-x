package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/redis/go-redis/v9"
)

func CreateDatabase()(*sql.DB, error ){
    // Connection string
    connStr := "user=postgres password=postgres dbname=postgres host=localhost port=5432 sslmode=disable"

    // Open the database connection
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        return nil,err
    }
    // Ping the database to check the connection
    err = db.Ping()
    if err != nil {
        log.Fatal("Error pinging the database:", err)
        return nil,err
    }
    fmt.Println("Successfully connected to the database!")
    return db,nil
}

var Ctx = context.Background()

func CreateRedisConnection(dbNum int) *redis.Client{
    rdb:= redis.NewClient(&redis.Options{
		Addr:   "0.0.0.0:6379",
		Password: "",
		DB: dbNum,
        ReadTimeout:  5 * time.Second,  // Increase read timeout
    WriteTimeout: 5 * time.Second,
	  })
	  return rdb;
}