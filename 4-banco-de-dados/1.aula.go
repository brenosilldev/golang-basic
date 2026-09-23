package bancodedados

import (
	"context"
	"database/sql"

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

	// user := NewUser("John Doe", "breno@hotmail.com")

	// result, err := InsertUser(db, user)

	// if err != nil {
	// 	panic(err)
	// }

	// rowsAffected, err := result.RowsAffected()

	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Printf("User inserted successfully: id=%s rowsAffected=%d\n", user.ID, rowsAffected)

	// user.Name = "Jane Doe(Breno Silva)"

	// result, err = UpdateUser(db, user)

	// if err != nil {
	// 	panic(err)
	// }

	// rowsAffected, err = result.RowsAffected()

	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Printf("User updated successfully: id=%s rowsAffected=%d\n", user.ID, rowsAffected)

	// userFromDB, err := GetUserByID(db, context.Background(), user.ID)

	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Printf("User retrieved successfully: id=%s name=%s email=%s\n", userFromDB.ID, userFromDB.Name, userFromDB.Email)

	users, err := GetAllUsers(db, context.Background())

	if err != nil {
		panic(err)
	}

	for _, user := range users {
		println("User retrieved successfully: id=" + user.ID + " name=" + user.Name + " email=" + user.Email)
	}

	userIDToDelete := users[0].ID // Replace with the actual user ID you want to delete

	err = DeleteUser(db, userIDToDelete)

	if err != nil {
		panic(err)
	}

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

func GetUserByID(db *sql.DB, ctx context.Context, id string) (*User, error) {

	stmt, err := db.Prepare("SELECT id, name, email FROM users WHERE id = $1")

	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	var user User

	// QueryRow  = Executa a consulta SQL e retorna uma única linha de resultado. Se a consulta não retornar nenhuma linha, ele retornará um erro sql.ErrNoRows.
	err = stmt.QueryRowContext(ctx, id).Scan(&user.ID, &user.Name, &user.Email)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetAllUsers(db *sql.DB, ctx context.Context) ([]User, error) {

	rows, err := db.Query("SELECT id, name, email FROM users")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []User

	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Email)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil

}

func DeleteUser(db *sql.DB, id string) error {

	stmt, err := db.Prepare("DELETE FROM users WHERE id = $1")

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(id)

	if err != nil {
		return err
	}

	return nil
}
