package urlmanagement

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"log"
	"regexp"
	"time"

	"github.com/esankhan/project-x/database"
	"github.com/redis/go-redis/v9"
)

func FindUserByEmail(email string) (int, error) {
	var db, err = database.CreateDatabase()
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	defer db.Close()
	var id int
	sqlStatement := `SELECT id FROM register WHERE email=$1`
	err2 := db.QueryRow(sqlStatement, email).Scan(&id)
	if err2 != nil || err2 == sql.ErrNoRows {
		log.Println("No rows were returned x")
		return -1, err2
	}
	return id, nil
}

func ShortenUrl(url string) string {
	h := sha256.New()
	h.Write([]byte(url))
	hashBytes := h.Sum(nil)
	encoded := base64.URLEncoding.EncodeToString(hashBytes)
	return encoded[:8]
}

func IsValidUrl(url string) bool {
	// check if the url is valid
	const urlPattern = `(?i)^(https?:\/\/)` + // validate protocol
		`((([a-z\d]([a-z\d-]*[a-z\d])*)\.)+[a-z]{2,}|` + // validate domain name
		`((\d{1,3}\.){3}\d{1,3}))` + // validate OR ip (v4) address
		`(\:\d+)?(\/(?:[-a-z\d%_.~+'#:,=]|%[a-fA-F0-9]{2})*)*` + // validate port and path
		`(\?[;&a-z\d%_.~+=-]*)?` + // validate query string
		`(\#[-a-z\d_]*)?$`

	// Compile the regex pattern
	re, err := regexp.Compile(urlPattern)
	if err != nil {
		log.Fatal("failed to compile regex", err)
	}

	// Check if the URL matches the pattern
	return re.MatchString(url)
}

func SaveUrl(id int, url string, shortUrl string) {
	saveToPostgres(url, shortUrl, id)
	saveToRedis(shortUrl, url)

}

func saveToPostgres(url string, shortUrl string, id int) {
	var db, err = database.CreateDatabase()
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	defer db.Close()
	sqlStatement := `INSERT INTO short_urls (url, short_url,user_id) VALUES ($1, $2, $3)`
	_, err2 := db.Exec(sqlStatement, url, shortUrl, id)
	if err2 != nil {
		log.Fatalf("Error inserting the URL: %v", err2)
	}
}

func saveToRedis(shortUrl string, url string) {
	rdb := database.CreateRedisConnection(0)
	if err := rdb.Ping(database.Ctx).Err(); err != nil {
		log.Fatal(err)
	}
	_ = rdb.Set(database.Ctx, shortUrl, url, time.Second*100)
}

func getFromRedis(shortUrl string) (string, string, error) {
	rdb := database.CreateRedisConnection(0)
	if err := rdb.Ping(database.Ctx).Err(); err != nil {
		log.Fatal(err)
	}
	value, err3 := rdb.Get(database.Ctx, shortUrl).Result()
	if err3 == redis.Nil {
		return "", "", err3
	}
	_ = rdb.Expire(database.Ctx, shortUrl, time.Second*100).Err()
	return value, "REDIS DB", nil
}

func getFromPostgres(shortUrl string) (string, string, error) {
	db, err := database.CreateDatabase()
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	defer db.Close()
	var url string
	sqlStatement := `SELECT url FROM short_urls WHERE short_url=$1`
	err2 := db.QueryRow(sqlStatement, shortUrl).Scan(&url)
	if err2 != nil || err2 == sql.ErrNoRows {
		log.Println("No rows were returned x")
		return "", "", err2
	}
	saveToRedis(shortUrl, url)
	return url, "POSTGRES DB", nil

}

func ResolveUrl(shortUrl string) (string, string, error) {

	value, dbSource, err := getFromRedis(shortUrl)
	if err == redis.Nil {
		value, dbSource, err = getFromPostgres(shortUrl)
		if err != nil {
			return value, dbSource, err
		}
		return value, dbSource, err

	}

	return value, dbSource, nil

}
