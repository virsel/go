/*

The package works on 2 tables on a PostgreSQL data base server.

The names of the tables are:

	* Users
	* Userdata

The definitions of the tables in the PostgreSQL server are:

	CREATE TABLE Users (
    	ID SERIAL,
    	Username VARCHAR(100) PRIMARY KEY
	);

	CREATE TABLE Userdata (
    	UserID Int NOT NULL,
    	Name VARCHAR(100),
    	Surname VARCHAR(100),
    	Description VARCHAR(200)
	);

	This is rendered as code

This is not rendered as code

*/
package pqfunctions

// BUG(1): Function ListUsers() not working as expected
// BUG(2): Function AddUser() is too slow

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

/*
This block of global variables holds the connection details to the Postgres server
	Hostname: is the IP or the hostname of the server
 	Port: is the TCP port the DB server listens to
	Username: is the username of the database user
	Password: is the password of the database user
	Database: is the name of the Database in PostgreSQL
*/
var (
	Hostname = ""
	Port     = 2345
	Username = ""
	Password = ""
	Database = ""
)

// The Userdata structure is for holding full user data
// from the Userdata table and the Username from the
// Users table
type Userdata struct {
	ID          int
	Username    string
	Name        string
	Surname     string
	Description string
}

// openConnection() is for opening the Postgres connection
// in order to be used by the other functions of the package.
func openConnection() (*sql.DB, error) {
	// connection string
	conn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", Hostname, Port, Username, Password, Database)
	// open database
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// The function returns the User ID of the username
// -1 if the user does not exist
func exists(username string) int {
	username = strings.ToLower(username)
	db, err := openConnection()
	if err != nil {
		fmt.Println(err)
		return -1
	}
	defer db.Close()
	userID := -1
	statement := fmt.Sprintf(`SELECT "id" FROM "users" where username = '%s'`, username)
	rows, err := db.Query(statement)
	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			fmt.Println("Scan", err)
			return -1
		}
		userID = id
	}
	defer rows.Close()
	return userID
}

// AddUser adds a new user to the database
//
// Returns new User ID
// -1 if there was an error
func AddUser(d Userdata) int {
	d.Username = strings.ToLower(d.Username)
	db, err := openConnection()
	if err != nil {
		fmt.Println(err)
		return -1
	}
	defer db.Close()
	userID := exists(d.Username)
	if userID != -1 {
		fmt.Println("User already exists:", Username)
		return -1
	}
	insertStatement := `insert into "users" ("username") values ($1)`
	_, err = db.Exec(insertStatement, d.Username)
	if err != nil {
		fmt.Println(err)
		return -1
	}
	userID = exists(d.Username)
	if userID == -1 {
		return userID
	}
	insertStatement = `insert into "userdata" ("userid", "name", "surname", "description") values ($1, $2, $3, $4)`
	_, err = db.Exec(insertStatement, userID, d.Name, d.Surname, d.Description)
	if err != nil {
		fmt.Println("db.Exec()", err)
		return -1
	}

	return userID
}

/*
	DeleteUser deletes an existing user if the user exists.

	It requires the User ID of the user to be deleted.
*/
func DeleteUser(id int) error {
	db, err := openConnection()
	if err != nil {
		return err
	}
	defer db.Close()
	// Does the ID exist?
	statement := fmt.Sprintf(`SELECT "username" FROM "users" where id = %d`, id)
	rows, err := db.Query(statement)
	var username string
	for rows.Next() {
		err = rows.Scan(&username)
		if err != nil {
			return err
		}
	}
	defer rows.Close()
	if exists(username) != id {
		return fmt.Errorf("User with ID %d does not exist", id)
	}
	// Delete from Userdata
	deleteStatement := `delete from "userdata" where userid=$1`
	_, err = db.Exec(deleteStatement, id)
	if err != nil {
		return err
	}
	// Delete from Users
	deleteStatement = `delete from "users" where id=$1`
	_, err = db.Exec(deleteStatement, id)
	if err != nil {
		return err
	}
	return nil
}

// ListUsers lists all users in the database
// and returns a slice of Userdata.
func ListUsers() ([]Userdata, error) {
	Data := []Userdata{}
	db, err := openConnection()
	if err != nil {
		return Data, err
	}
	defer db.Close()
	rows, err := db.Query(`SELECT  
        "id","username","name","surname","description"
        FROM "users","userdata"
        WHERE users.id = userdata.userid`)
	if err != nil {
		return Data, err
	}
	for rows.Next() {
		var id int
		var username string
		var name string
		var surname string
		var description string
		err = rows.Scan(&id, &username, &name, &surname, &description)
		temp := Userdata{ID: id, Username: username, Name: name, Surname: surname, Description: description}
		Data = append(Data, temp)
		if err != nil {
			return Data, err
		}
	}
	defer rows.Close()
	return Data, nil
}

// UpdateUser is for updating an existing user
// given a Userdata structure.
// The user ID of the user to be updated is found
// inside the function.
func UpdateUser(d Userdata) error {
	db, err := openConnection()
	if err != nil {
		return err
	}
	defer db.Close()
	userID := exists(d.Username)
	if userID == -1 {
		return errors.New("User does not exist")
	}
	d.ID = userID
	updateStatement := `update "userdata" set "name"=$1, "surname"=$2, "description"=$3 where "userid"=$4`
	_, err = db.Exec(updateStatement, d.Name, d.Surname, d.Description, d.ID)
	if err != nil {
		return err
	}
	return nil
}
