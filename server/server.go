package server

import (
	"crud/database"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type user struct {
	ID    uint32 `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

//CreateUser cria um novo usuário no banco
func CreateUser(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		w.Write([]byte("Falha ao ler a request"))
		return
	}

	var user user

	if err := json.Unmarshal(body, &user); err != nil {
		w.Write([]byte("Erro ao converter user para struct"))
		return
	}

	db, err := database.Connect()
	if err != nil {
		w.Write([]byte("Erro ao conectar no banco!"))
		return
	}
	defer db.Close()

	query := "INSERT INTO users (name, email) VALUES (?, ?)"
	stmt, err := db.Prepare(query)

	if err != nil {
		w.Write([]byte("Erro ao criar o statment!"))
		return
	}

	defer stmt.Close()

	insert, err := stmt.Exec(user.Name, user.Email)

	if err != nil {
		w.Write([]byte("Erro ao executar statemant"))
	}

	idInsert, err := insert.LastInsertId()

	if err != nil {
		w.Write([]byte("Erro ao obter o id inserido"))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("Usuário criado com sucesso! ID: %d", idInsert)))

}

//SearchUser procura por todos os usuários no banco
func SearchUsers(w http.ResponseWriter, r *http.Request) {
	db, err := database.Connect()

	if err != nil {
		w.Write([]byte("Erro ao contectar com o banco"))
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM users")

	if err != nil {
		w.Write([]byte("Erro ao buscar user"))
		return
	}

	var users []user

	for rows.Next() {
		var user user
		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			w.Write([]byte("Erro ao scanear o user"))
			return
		}
		users = append(users, user)
	}

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(users); err != nil {
		w.Write([]byte("Erro ao converter user para json"))
		return
	}

}

//SearchUser procura por um usuário no banco
func SearchUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	ID, err := strconv.ParseUint(params["id"], 10, 32)
	if err != nil {
		w.Write([]byte("Erro ao converter parâmetro para int"))
		return
	}

	db, err := database.Connect()

	if err != nil {
		w.Write([]byte("Erro ao conectar com o banco"))
		return
	}

	row, err := db.Query("SELECT * FROM users WHERE id = ?", ID)

	if err != nil {
		w.Write([]byte("Erro ao conectar buscar usuario"))
		return
	}

	var user user
	if row.Next() {
		if err := row.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			w.Write([]byte("Erro ao conectar escanear usuario"))
			return
		}
	}

	if err := json.NewEncoder(w).Encode(user); err != nil {
		w.Write([]byte("Erro ao converter user para json"))
		return
	}

}

//UpdateUser altera os dados de um usuário no banco
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	ID, err := strconv.ParseUint(params["id"], 10, 32)

	if err != nil {
		w.Write([]byte("Erro ao converter user para id"))
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte("Erro ao ler corpo da request"))
		return
	}

	var user user
	if err := json.Unmarshal(body, &user); err != nil {
		w.Write([]byte("Erro ao converter user para struct"))
		return
	}

	db, err := database.Connect()
	if err != nil {
		w.Write([]byte("Erro ao conectar no banco"))
		return
	}

	defer db.Close()

	query := "UPDATE users SET name = ?, email = ? WHERE id = ?"

	stmt, err := db.Prepare(query)

	if err != nil {
		w.Write([]byte("Erro ao criar o statment"))
		fmt.Printf("%s\n", err)
		return
	}
	defer stmt.Close()

	if _, err := stmt.Exec(user.Name, user.Email, ID); err != nil {
		w.Write([]byte("Erro ao atualizar user"))
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

//DeleteUser deleta um user do banco
func DeleteUser(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	ID, err := strconv.ParseUint(params["id"], 10, 32)

	if err != nil {
		w.Write([]byte("Erro ao converter user para id"))
		return
	}

	db, err := database.Connect()
	if err != nil {
		w.Write([]byte("Erro ao conectar no banco"))
		return
	}

	defer db.Close()

	query := "DELETE FROM users WHERE id = ?"
	stmt, err := db.Prepare(query)

	if err != nil {
		w.Write([]byte("Erro ao criar o statment"))
		fmt.Printf("%s\n", err)
		return
	}
	defer stmt.Close()

	if _, err := stmt.Exec(ID); err != nil {
		w.Write([]byte("Erro ao deletar user"))
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode("Usuário deletado com sucesso")

}
