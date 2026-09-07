package bancodedados

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type User struct {
	ID    string
	Name  string
	Email string
}

func NewUser(name, email string) *User {
	return &User{
		ID:    uuid.New().String(),
		Name:  name,
		Email: email,
	}
}

func BancoTeste() {

	db, err := sql.Open("postgres", "host=localhost port=5439 user=root password=root dbname=golang_db sslmode=disable")
	if err != nil {
		panic(err)
	}

	user := NewUser("John Doe", "breno@hotmail.com")

	result, err := InsertUser(db, user)

	if err != nil {
		panic(err)
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		panic(err)
	}

	fmt.Printf("User inserted successfully: id=%s rowsAffected=%d\n", user.ID, rowsAffected)
}

func InsertUser(db *sql.DB, user *User) (sql.Result, error) {

	stmt, err := db.Prepare("INSERT INTO users (id, name, email) VALUES ($1, $2, $3)")

	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	user.Name = "Breno Silva"

	res, err := stmt.Exec(user.ID, user.Name, user.Email)

	if err != nil {
		return nil, err
	}

	return res, nil

}

func UpdateUser(db *sql.DB, user *User) (sql.Result, error) {

	stmt, err := db.Prepare("UPDATE users SET name = $1, email = $2 WHERE id = $3")

	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	res, err := stmt.Exec(user.Name, user.Email, user.ID)

	if err != nil {
		return nil, err
	}

	return res, nil
}
