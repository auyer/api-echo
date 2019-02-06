package controllers

import (
	"crypto/md5"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/auyer/muxapi/db"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

//ErrorBody structure is used to improve error reporting in a JSON response body
type ErrorBody struct {
	Reason string `json:"reason"`
}

//User structure is used to store data used by this API
type User struct {
	ID                int    `db:"id, primarykey, autoincrement" json:"id"`
	Idpessoa          int    `db:"id_pessoa" json:"id_pessoa"`
	Nome              string `db:"nome" json:"nome"`
	Matriculainterna  int    `db:"matricula_interna" json:"matricula_interna"`
	Nomeidentificacao string `db:"nome_identificacao" json:"nome_identificacao"`
	Datanascimento    string `db:"data_nascimento" json:"data_nascimento"`
	Sexo              string `db:"sexo" json:"sexo"`
}

//Controller is used to export the API handler functions
type Controller struct {
	DB db.Pointer
}

//GetAll funtion returns the full list of document
func (ctrl Controller) GetAll(w http.ResponseWriter, r *http.Request) {
	rows, err := ctrl.DB.GetAll()
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(ErrorBody{
			Reason: err.Error(),
		})
		log.Print("[MUX] " + " | 500 | " + r.Method + "  " + r.URL.Path)
		return
	}
	defer rows.Close()
	var users []User
	var id, idpessoa, matriculainterna, siape int
	var nome, nomeidentificacao, datanascimento, sexo string
	for rows.Next() {
		err := rows.Scan(&id, &siape, &idpessoa, &matriculainterna, &nomeidentificacao, &nome, &datanascimento, &sexo)
		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(ErrorBody{
				Reason: err.Error(),
			})
			log.Print("[MUX] " + " | 500 | " + r.Method + "  " + r.URL.Path)
			return
		}
		date, _ := time.Parse("1969-02-12", datanascimento)
		users = append(users, User{
			ID:                id,
			Idpessoa:          idpessoa,
			Nome:              nome,
			Matriculainterna:  matriculainterna,
			Nomeidentificacao: nomeidentificacao,
			Datanascimento:    date.Format("1969-02-12"),
			Sexo:              sexo,
		})
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(ErrorBody{
			Reason: err.Error(),
		})
		log.Print("[MUX] " + " | 500 | " + r.Method + "  " + r.URL.Path)
		return
	}
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(&users)
	log.Print("[MUX] " + " | 200 | " + r.Method + "  " + r.URL.Path)
	return
}

//GetByID funtion returns document by id
func (ctrl Controller) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	mat := vars["id"]
	rows, err := ctrl.DB.GetByID(mat)
	if err != nil {
		log.Println(err)
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(ErrorBody{
			Reason: err.Error(),
		})
		log.Print("[MUX] " + " | 400 | " + r.Method + "  " + r.URL.Path)
		return
	}
	defer rows.Close()
	var users []User
	var id, idpessoa, matriculainterna, siape int
	var nome, nomeidentificacao, datanascimento, sexo string
	for rows.Next() {
		err := rows.Scan(&id, &siape, &idpessoa, &matriculainterna, &nomeidentificacao, &nome, &datanascimento, &sexo)
		if err != nil {
			log.Println(err)
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(ErrorBody{
				Reason: err.Error(),
			})
			log.Print("[MUX] " + " | 400 | " + r.Method + "  " + r.URL.Path)
			return
		}
		date, _ := time.Parse("1969-02-12", datanascimento)
		users = append(users, User{
			ID:                id,
			Idpessoa:          idpessoa,
			Nome:              nome,
			Matriculainterna:  matriculainterna,
			Nomeidentificacao: nomeidentificacao,
			Datanascimento:    date.Format("1969-02-12"),
			Sexo:              sexo,
		})
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(ErrorBody{
			Reason: err.Error(),
		})
		log.Print("[MUX] " + " | 400 | " + r.Method + "  " + r.URL.Path)
		return
	}
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(&users)
	log.Print("[MUX] " + " | 200 | " + r.Method + "  " + r.URL.Path)
	return
}

// WebsocketHandler handles the websocket connection
func (ctrl Controller) WebsocketHandler(writer http.ResponseWriter, request *http.Request) {
	socket, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		fmt.Println(err)
	}
	for {
		// Vamos ler a mensagem recebida via Websocket
		msgType, msg, err := socket.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		// Logando no console do Webserver
		fmt.Println("Mensagem recebida: ", string(msg))

		// Devolvendo a mensagem recebida de volta para o cliente
		err = socket.WriteMessage(msgType, msg)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

//PostServidor function reads a JSON body and store it in the database
func (ctrl Controller) PostServidor(w http.ResponseWriter, r *http.Request) {
	regexcheck := false
	var ser User
	var Reasons []ErrorBody
	decoder := json.NewDecoder(r.Body)
	errDecode := decoder.Decode(&ser)
	if errDecode != nil {
		log.Println(errDecode)
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(ErrorBody{
			Reason: errDecode.Error(),
		})
		log.Print("[MUX] " + " | 400 | " + r.Method + "  " + r.URL.Path)
		return
	}
	// REGEX CHEKING PHASE
	regex, _ := regexp.Compile(`^(19[0-9]{2}|2[0-9]{3})-(0[1-9]|1[012])-([123]0|[012][1-9]|31)$`)
	if !regex.MatchString(ser.Datanascimento) {
		regexcheck = true
		Reasons = append(Reasons, ErrorBody{
			Reason: "[data_nascimento] failed to match API requirements. It should look like this: 1969-02-12",
		})
	}
	regex, _ = regexp.Compile(`^([A-Z][a-z]+([ ]?[a-z]?['-]?[A-Z][a-z]+)*)$`)
	if !regex.MatchString(ser.Nome) {
		regexcheck = true
		Reasons = append(Reasons, ErrorBody{
			Reason: "[nome] failed to match API requirements. It should look like this: Firstname Middlename(optional) Lastname",
		})
	}
	regex, _ = regexp.Compile(`^([A-Z][a-z]+([ ]?[a-z]?['-]?[A-Z][a-z]+)*)$`)
	if !regex.MatchString(ser.Nomeidentificacao) {
		regexcheck = true
		Reasons = append(Reasons, ErrorBody{
			Reason: "[nome_identificacao] failed to match API requirements. It should look like this: Firstname Middlename(optional) Lastname",
		})
	}
	regex, _ = regexp.Compile(`\b[MF]{1}\b`)
	if !regex.MatchString(ser.Sexo) {
		regexcheck = true
		Reasons = append(Reasons, ErrorBody{
			Reason: "[sexo] failed to match API requirements. It should look like this: M for male, F for female",
		})
	}
	if regexcheck {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(&Reasons)
		log.Print("[MUX] " + " | 400 | " + r.Method + "  " + r.URL.Path)
		return
	}
	// END OF REGEX CHEKING PHASE
	timestamp := time.Now().UTC().Format("2006-01-02T15:04:05-0700")
	b := md5.Sum([]byte(fmt.Sprintf(string(ser.Nome), string(timestamp))))
	bid := binary.BigEndian.Uint64(b[:])
	ser.Matriculainterna = int(bid % 9999999)
	q := fmt.Sprintf(`
		INSERT INTO rh.servidor_tmp(
			nome, nome_identificacao, siape, id_pessoa, matricula_interna, id_foto,
			data_nascimento, sexo)
			VALUES ('%s', '%s', %d, %d, %d, null, '%s', '%s');
			`, ser.Nome, ser.Nomeidentificacao, ser.Idpessoa, ser.Matriculainterna,
		ser.Datanascimento, ser.Sexo) //String formating
	rows, err := ctrl.DB
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(ErrorBody{
			Reason: err.Error(),
		})
		log.Print("[MUX] " + " | 500 | " + r.Method + "  " + r.URL.Path)
		return
	}
	defer rows.Close()
	w.Header().Add("location", r.URL.Host+"/api/servidor/"+strconv.Itoa(ser.Matriculainterna))
	w.WriteHeader(201)
	log.Print("[MUX] " + " | 201 | " + r.Method + "  " + r.URL.Path)
	return
}
