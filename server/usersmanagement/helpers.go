package usersmanagement

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/esankhan/project-x/database"
	"golang.org/x/crypto/bcrypt"
)

func FindUser(email string) bool{
	db,err := database.CreateDatabase();
if(err != nil){
		panic(err)
	}
	defer db.Close();
	sqlStatement := `SELECT email FROM register WHERE email=$1`;
	err2 := db.QueryRow(sqlStatement, email).Scan(&email)
	if(err2 != nil || err2 == sql.ErrNoRows){
		log.Println("No rows were returned")
		return false
	}
	return true

}


func FindUserByEmail(email string) loginRequest{
	db,err := database.CreateDatabase();
if(err != nil){
		panic(err)
	}
	var req2 loginRequest
	defer db.Close();
	sqlStatement := `SELECT email,password FROM register WHERE email=$1`;
	err2 := db.QueryRow(sqlStatement, email).Scan(&req2.Email, &req2.Password)
	if(err2 != nil || err2 == sql.ErrNoRows){
		log.Println("No rows were returned x")
		return loginRequest{}
	}
	return req2
}



 func InsertUser(req registerRequest) int64{
db,err := database.CreateDatabase();

if(err != nil){
		panic(err)
	}
	defer db.Close()
sqlStatement := `INSERT INTO register (username, email, password) VALUES ($1, $2, $3) RETURNING id`;
var id int64

error := db.QueryRow(sqlStatement, req.Username, req.Email, req.Password).Scan(&id)
if error != nil {
	log.Fatal("Error inserting into the database:", error)
}
// set the key value pair in redis


fmt.Println("New record ID is XXXXXXXXXXXXXXX:", id)
return id
 }


func HashPassword(password string) string {
	bytes ,err := bcrypt.GenerateFromPassword([]byte(password), 14);
	if(err != nil){
		log.Fatal("failed to hash password",err)
	}
	return string(bytes)
}


func ComparePassword (hashedPassword string, password string) bool{
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}