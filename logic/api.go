package logic

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/thallesrangel/crud_go/database"
)

type user struct {
	ID    uint32 `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email`
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		w.Write([]byte("Falha ao ler o Body"))
		return
	}

	var user user

	if err = json.Unmarshal(body, &user); err != nil {
		w.Write([]byte("Falha ao converter user para struct"))
		return
	}

	fmt.Println(user)

	db, err := database.Conn()

	if err != nil {
		w.Write([]byte("Erro ao conectar no database"))
		return
	}

	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO users (name, email) VALUES(?, ?)")

	if err != nil {
		w.Write([]byte("Erro ao criar o statement"))
		return
	}

	defer stmt.Close()

	insert, err := stmt.Exec(user.Name, user.Email)

	if err != nil {
		w.Write([]byte("Erro ao executar o statement"))
		return
	}

	lastId, err := insert.LastInsertId()

	if err != nil {
		w.Write([]byte("Erro ao obter o ID inserido"))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("Usuario inserido com sucesso. ID: %d ", lastId)))
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	db, err := database.Conn()

	if err != nil {
		w.Write([]byte("Erro ao conectar no database"))
		return
	}

	defer db.Close()

	line, err := db.Query("SELECT * FROM users")

	if err != nil {
		w.Write([]byte("Erro ao buscar o usuario"))
		return
	}

	defer line.Close()

	var data []user

	for line.Next() {
		var user user

		if err := line.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			w.Write([]byte("Erro ao scanear o usuario"))
			return
		}

		data = append(data, user)
	}

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		w.Write([]byte("Erro ao converter os usuarios para JSON"))
		return
	}
}

func GetUserById(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	ID, err := strconv.ParseUint(params["id"], 10, 32)

	if err != nil {
		w.Write([]byte("Erro ao converter os usuarios para JSON"))
		return
	}

	db, err := database.Conn()

	if err != nil {
		w.Write([]byte("Erro ao conectar no banco de dados"))
		return
	}

	line, err := db.Query("SELECT * FROM users WHERE id = ?", ID)

	if err != nil {
		w.Write([]byte("Erro ao buscar usuario"))
		return
	}

	var data user

	if line.Next() {
		if err := line.Scan(&data.ID, &data.Name, &data.Email); err != nil {
			w.Write([]byte("Erro ao scanear usuario"))
			return
		}
	}

	if data.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Usuario nao existe"))
		return
	}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		w.Write([]byte("Erro ao converter usuario para JSON"))
		return
	}
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	ID, err := strconv.ParseUint(params["id"], 10, 32)

	if err != nil {
		w.Write([]byte("Erro ao ler o parametro ID"))
		return
	}

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		w.Write([]byte("Erro ao ler o Body"))
		return
	}

	var data user

	if err := json.Unmarshal(body, &data); err != nil {
		w.Write([]byte("Falha ao converter user para struct"))
		return
	}

	db, err := database.Conn()

	if err != nil {
		w.Write([]byte("Erro ao conectar no database"))
		return
	}

	defer db.Close()

	stmt, err := db.Prepare("UPDATE users SET name = ?, email = ? WHERE id = ?")

	if err != nil {
		w.Write([]byte("Erro ao criar o statement"))
		return
	}

	defer stmt.Close()

	if _, err := stmt.Exec(data.Name, data.Email, ID); err != nil {
		w.Write([]byte("Erro ao atualizar o usuario"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	ID, err := strconv.ParseUint(params["id"], 10, 32)

	if err != nil {
		w.Write([]byte("Erro ao ler o parametro ID"))
		return
	}

	db, err := database.Conn()

	if err != nil {
		w.Write([]byte("Erro ao conectar no database"))
		return
	}

	defer db.Close()

	stmt, err := db.Prepare("DELETE FROM users WHERE id = ?")

	if err != nil {
		w.Write([]byte("Erro ao criar o statement"))
		return
	}

	defer stmt.Close()

	if _, err := stmt.Exec(ID); err != nil {
		w.Write([]byte("Erro ao deletar o usuario"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
