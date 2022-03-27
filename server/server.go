package server

import (
	"crud/database"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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
	resp := make(map[string]string)
	body, err := ioutil.ReadAll(r.Body) //Lê o body do POST

	if err != nil {
		w.Write([]byte("Falha ao ler a request"))
		return
	}

	var user user

	if err := json.Unmarshal(body, &user); err != nil { //popula a struc com os dados do JSON
		resp["success"] = "false"
		resp["message"] = "Preencha todos os campos"
		log.Println(err)

		respJSON, err := json.Marshal(resp)

		if err != nil {
			log.Println(err)
		}

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)

		w.Write(respJSON)
		return
	}

	db, err := database.Connect()
	if err != nil {
		resp["success"] = "false"
		resp["message"] = "Erro no servidor, por favor, tente novamente mais tarde"
		log.Println(err)

		respJSON, err := json.Marshal(resp)

		if err != nil {
			log.Println(err)
		}

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)

		w.Write(respJSON)
		return
	}
	defer db.Close()

	query := "INSERT INTO users (name, email) VALUES (?, ?)"
	stmt, err := db.Prepare(query)

	if err != nil {
		resp["success"] = "false"
		resp["message"] = "Erro no servidor, por favor, tente novamente mais tarde"
		log.Println(err)

		respJSON, err := json.Marshal(resp)

		if err != nil {
			log.Println(err)
		}

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)

		w.Write(respJSON)
		return
	}

	defer stmt.Close()

	insert, err := stmt.Exec(user.Name, user.Email) //usa os valores da struc populada para por na query

	if err != nil {
		resp["success"] = "false"
		resp["message"] = "Erro no servidor, por favor, tente novamente mais tarde"
		log.Println(err)

		respJSON, err := json.Marshal(resp)

		if err != nil {
			log.Println(err)
		}

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)

		w.Write(respJSON)
		return
	}

	idInsert, err := insert.LastInsertId()

	if err != nil {
		resp["success"] = "false"
		resp["message"] = "Erro no servidor, por favor, tente novamente mais tarde"
		log.Println(err)

		respJSON, err := json.Marshal(resp)

		if err != nil {
			log.Println(err)
		}

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)

		w.Write(respJSON)
		return

	}

	w.Header().Set("content-type", "application/json")

	resp["success"] = "true"
	resp["message"] = fmt.Sprintf("Usuário criado com sucesso! ID: %d", idInsert)

	respJSON, err := json.Marshal(resp)

	w.WriteHeader(http.StatusCreated)
	w.Write(respJSON)

}

//SearchUser procura por todos os usuários no banco
func SearchUsers(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)

	db, err := database.Connect()

	if err != nil {
		resp["success"] = "false"
		resp["message"] = "Ocorreu um erro no servidor, tente novamente mais tarde"

		respJSON, err := json.Marshal(resp)

		if err != nil {
			log.Println(err)
		}

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)

		w.Write(respJSON)
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM users")

	if err != nil {

		resp["success"] = "false"
		resp["message"] = "Usuário não encontrado"
		respJSON, err := json.Marshal(resp)

		if err != nil {
			log.Println(err)
		}

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		w.Write(respJSON)
		return
	}

	var users []user //cria um slice do tipo user (como um array de objetos)

	//da um loop nos rows (linhas do retorno do banco) e preenche cada struc do array com uma linha
	for rows.Next() {
		var user user
		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			resp["success"] = "false"
			resp["message"] = "Ocorreu um erro no servidor, tente novamente mais tarde"
			respJSON, err := json.Marshal(resp)

			if err != nil {
				log.Println(err)
			}

			w.Header().Set("content-type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(respJSON)
			return
		}
		users = append(users, user) //joga esse novo user dentro do slice de users
	}

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(users); err != nil {
		resp["success"] = "false"
		resp["message"] = "Erro no servidor, por favor, tente novamente mais tarde"
		log.Println(err)
		respJSON, err := json.Marshal(resp)

		if err != nil {
			log.Println(err)
		}

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		w.Write(respJSON)

		return

	}

}

//SearchUser procura por um usuário no banco
func SearchUser(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)

	params := mux.Vars(r)

	ID, err := strconv.ParseUint(params["id"], 10, 32)
	if err != nil {
		resp["success"] = "false"
		resp["message"] = "Erro no servidor, por favor, tente novamente mais tarde"
		log.Println(err)

		respJSON, err := json.Marshal(resp)

		if err != nil {
			log.Println(err)
		}

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		w.Write(respJSON)
		return
	}

	db, err := database.Connect()

	if err != nil {
		resp["success"] = "false"
		resp["message"] = "Erro no servidor, por favor, tente novamente mais tarde"
		log.Println(err)
		respJSON, err := json.Marshal(resp)

		if err != nil {
			log.Println(err)
		}

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		w.Write(respJSON)
		return

	}

	row, err := db.Query("SELECT * FROM users WHERE id = ?", ID)

	if err != nil {
		resp["success"] = "false"
		resp["message"] = "Erro no servidor, por favor, tente novamente mais tarde"
		log.Println(err)

		respJSON, err := json.Marshal(resp)

		if err != nil {
			log.Println(err)
		}

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		w.Write(respJSON)
		return

	}

	var user user
	//como só é um usuário específico, popula uma struct
	if row.Next() {
		if err := row.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			resp["success"] = "false"
			resp["message"] = "Erro no servidor, por favor, tente novamente mais tarde"
			log.Println(err)

			respJSON, err := json.Marshal(resp)

			if err != nil {
				log.Println(err)
			}

			w.Header().Set("content-type", "application/json")
			w.WriteHeader(http.StatusNotFound)

			w.Write(respJSON)

			return
		}
	}

	if err := json.NewEncoder(w).Encode(user); err != nil {
		resp["success"] = "false"
		resp["message"] = "Erro no servidor, por favor, tente novamente mais tarde"
		log.Println(err)

		respJSON, err := json.Marshal(resp)

		if err != nil {
			log.Println(err)
		}

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		w.Write(respJSON)
		return
	}

}

//UpdateUser altera os dados de um usuário no banco
func UpdateUser(w http.ResponseWriter, r *http.Request) {

	resp := make(map[string]string)

	params := mux.Vars(r)

	ID, err := strconv.ParseUint(params["id"], 10, 32)

	if err != nil {
		resp["success"] = "false"
		resp["message"] = "Erro no servidor, por favor, tente novamente mais tarde"
		log.Println(err)

		respJSON, err := json.Marshal(resp)

		if err != nil {
			log.Println(err)
		}

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		w.Write(respJSON)

		return
	}

	body, err := ioutil.ReadAll(r.Body) //Lê o body do PUT
	if err != nil {
		resp["success"] = "false"
		resp["message"] = "Erro no servidor, por favor, tente novamente mais tarde"
		log.Println(err)

		respJSON, err := json.Marshal(resp)

		if err != nil {
			log.Println(err)
		}

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		w.Write(respJSON)
		return
	}

	var user user
	if err := json.Unmarshal(body, &user); err != nil { //popula uma struc com os dados do JSON
		resp["success"] = "false"
		resp["message"] = "Erro no servidor, por favor, tente novamente mais tarde"
		log.Println(err)

		respJSON, err := json.Marshal(resp)

		if err != nil {
			log.Println(err)
		}

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		w.Write(respJSON)
		return
	}

	db, err := database.Connect()
	if err != nil {
		resp["success"] = "false"
		resp["message"] = "Erro no servidor, por favor, tente novamente mais tarde"
		log.Println(err)

		respJSON, err := json.Marshal(resp)

		if err != nil {
			log.Println(err)
		}

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		w.Write(respJSON)
		return
	}

	defer db.Close()

	query := "UPDATE users SET name = ?, email = ? WHERE id = ?"

	stmt, err := db.Prepare(query)

	if err != nil {
		resp["success"] = "false"
		resp["message"] = "Erro no servidor, por favor, tente novamente mais tarde"
		log.Println(err)

		respJSON, err := json.Marshal(resp)

		if err != nil {
			log.Println(err)
		}

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		w.Write(respJSON)
		return

	}
	defer stmt.Close()

	if _, err := stmt.Exec(user.Name, user.Email, ID); err != nil {
		w.Write([]byte("Erro ao atualizar user"))
		return
	}

	resp["success"] = "true"
	resp["message"] = "Usuário atualizado com sucesso!"

	respJSON, err := json.Marshal(resp)

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respJSON)

}

//DeleteUser deleta um user do banco
func DeleteUser(w http.ResponseWriter, r *http.Request) {

	resp := make(map[string]string)
	params := mux.Vars(r)

	ID, err := strconv.ParseUint(params["id"], 10, 32)

	if err != nil {
		resp["success"] = "false"
		resp["message"] = "Erro no servidor, por favor, tente novamente mais tarde"
		log.Println(err)

		respJSON, err := json.Marshal(resp)

		if err != nil {
			log.Println(err)
		}

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		w.Write(respJSON)
		return
	}

	db, err := database.Connect()
	if err != nil {
		resp["success"] = "false"
		resp["message"] = "Erro no servidor, por favor, tente novamente mais tarde"
		log.Println(err)

		respJSON, err := json.Marshal(resp)

		if err != nil {
			log.Println(err)
		}

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		w.Write(respJSON)
		return
	}

	defer db.Close()

	query := "DELETE FROM users WHERE id = ?"
	stmt, err := db.Prepare(query)

	if err != nil {
		resp["success"] = "false"
		resp["message"] = "Erro no servidor, por favor, tente novamente mais tarde"
		log.Println(err)

		respJSON, err := json.Marshal(resp)

		if err != nil {
			log.Println(err)
		}

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		w.Write(respJSON)
		return

	}
	defer stmt.Close()

	if _, err := stmt.Exec(ID); err != nil {
		resp["success"] = "false"
		resp["message"] = "Erro no servidor, por favor, tente novamente mais tarde"
		log.Println(err)

		respJSON, err := json.Marshal(resp)

		if err != nil {
			log.Println(err)
		}

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		w.Write(respJSON)
		return
	}

	resp["success"] = "true"
	resp["message"] = "Usuário deletado com sucesso"

	respJSON, err := json.Marshal(resp)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	w.Write(respJSON)

}
