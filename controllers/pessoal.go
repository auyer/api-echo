package controllers

import (
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/labstack/echo"
	"github.com/restsec/api-gingonic/db"
)

//ErrorBody structure is used to improve error reporting in a JSON response body
type ErrorBody struct {
	Reason string `json:"reason"`
}

//Servidor structure is used to store data used by this API
type Servidor struct {
	ID                int    `db:"id, primarykey, autoincrement" json:"id"`
	Siape             int    `db:"siape" json:"siape"`
	Idpessoa          int    `db:"id_pessoa" json:"id_pessoa"`
	Nome              string `db:"nome" json:"nome"`
	Matriculainterna  int    `db:"matricula_interna" json:"matricula_interna"`
	Nomeidentificacao string `db:"nome_identificacao" json:"nome_identificacao"`
	Datanascimento    string `db:"data_nascimento" json:"data_nascimento"`
	Sexo              string `db:"sexo" json:"sexo"`
}

//ServidorController is used to export the API handler functions
type ServidorController struct{} // THis is used to make functions callable from ServidorCOntroller

//GetServidores funtion returns the full list of "servidores" in the database
func (ctrl ServidorController) GetServidores(c echo.Context) error {
	q := `select s.id_servidor, s.siape, s.id_pessoa, s.matricula_interna, s.nome_identificacao,
		p.nome, p.data_nascimento, p.sexo from rh.servidor s
	inner join comum.pessoa p on (s.id_pessoa = p.id_pessoa)` //Manual query
	rows, err := db.GetDB().Query(q) //Get Database cursor from DB module
	if err != nil {
		log.Println(err)
		return c.JSON(400, ErrorBody{
			Reason: err.Error(),
		})
	}
	defer rows.Close() //will close DB after function GetServidor is over.
	var servidores []Servidor
	var id, idpessoa, siape, matriculainterna int
	var nome, nomeidentificacao, datanascimento, sexo string
	for rows.Next() {
		err := rows.Scan(&id, &siape, &idpessoa, &matriculainterna, &nomeidentificacao, &nome, &datanascimento, &sexo)
		if err != nil {
			log.Println(err)
			return c.JSON(400, ErrorBody{
				Reason: err.Error(),
			})

		}
		// log.Println(id)
		date, _ := time.Parse("1969-02-12", datanascimento)
		servidores = append(servidores, Servidor{
			ID:                id,
			Siape:             siape,
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
		return c.JSON(400, ErrorBody{
			Reason: err.Error(),
		})
	}
	c.JSON(200, servidores)
	if err != nil {
		log.Println(err)
		return c.JSON(400, ErrorBody{
			Reason: err.Error(),
		})
	}
	return nil
}

//GetServidorMat funtion returns the "servidor" matching a given id
func (ctrl ServidorController) GetServidorMat(c echo.Context) error {
	mat := c.Param("matricula") // URL parameter
	// Data security checking to be insterted here
	r, err := regexp.Compile(`\b[0-9]+\b`)
	if !r.MatchString(mat) {
		log.Println(err)
		return c.JSON(404, ErrorBody{
			Reason: err.Error(),
		})
	}
	q := fmt.Sprintf(`select s.id_servidor, s.siape, s.id_pessoa, s.matricula_interna, s.nome_identificacao,
		p.nome, p.data_nascimento, p.sexo from rh.servidor s
	inner join comum.pessoa p on (s.id_pessoa = p.id_pessoa) where s.matricula_interna = %s`, mat) //String formating
	rows, err := db.GetDB().Query(q)
	if err != nil {
		log.Println(err)
		return c.JSON(400, ErrorBody{
			Reason: err.Error(),
		})
	}
	defer rows.Close()
	var servidores []Servidor
	var id, idpessoa, siape, matriculainterna int
	var nome, nomeidentificacao, datanascimento, sexo string
	for rows.Next() {
		err := rows.Scan(&id, &siape, &idpessoa, &matriculainterna, &nomeidentificacao, &nome, &datanascimento, &sexo)
		if err != nil {
			log.Println(err)
			return c.JSON(400, ErrorBody{
				Reason: err.Error(),
			})
		}
		date, _ := time.Parse("1969-02-12", datanascimento)
		servidores = append(servidores, Servidor{
			ID:                id,
			Siape:             siape,
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
		return c.JSON(400, ErrorBody{
			Reason: err.Error(),
		})
	}
	return c.JSON(200, servidores)
}

//PostServidor function reads a JSON body and store it in the database
func (ctrl ServidorController) PostServidor(c echo.Context) error {
	regexcheck := false
	var ser Servidor
	var Reasons []ErrorBody
	err := c.Bind(ser)
	if err != nil {
		log.Println(err)
		return c.JSON(400, ErrorBody{
			Reason: err.Error(),
		})
	}

	// REGEX CHEKING PHASE
	r, _ := regexp.Compile(`^(19[0-9]{2}|2[0-9]{3})-(0[1-9]|1[012])-([123]0|[012][1-9]|31)$`)
	if !r.MatchString(ser.Datanascimento) {
		regexcheck = true
		Reasons = append(Reasons, ErrorBody{
			Reason: "[data_nascimento] failed to match API requirements. It should look like this: 1969-02-12",
		})
	} else {
		now := time.Now()
		now.Format(time.RFC3339)
		time, err := time.Parse("2006-01-02 15:04:05 -0200", fmt.Sprint(ser.Datanascimento, " 00:00:00 -0200"))
		if err != nil {
			log.Println(err)
			return c.JSON(500, ErrorBody{
				Reason: err.Error(),
			})
		}
		if !now.After(time) {
			regexcheck = true
			Reasons = append(Reasons, ErrorBody{
				Reason: "[data_nascimento] failed to match API requirements. It should not be in future",
			})
		}
	}
	r, _ = regexp.Compile(`^([A-Z][a-z]+([ ]?[a-z]?['-]?[A-Z][a-z]+)*)$`)
	if !r.MatchString(ser.Nome) {
		regexcheck = true
		Reasons = append(Reasons, ErrorBody{
			Reason: "[nome] failed to match API requirements. It should look like this: Firstname Middlename(optional) Lastname",
		})
	} else if len(ser.Nome) > 100 {
		regexcheck = true
		Reasons = append(Reasons, ErrorBody{
			Reason: "[nome] failed to match API requirements. It should have a maximum of 100 characters",
		})
	}
	r, _ = regexp.Compile(`^([A-Z][a-z]+([ ]?[a-z]?['-]?[A-Z][a-z]+)*)$`)
	if !r.MatchString(ser.Nomeidentificacao) {
		regexcheck = true
		Reasons = append(Reasons, ErrorBody{
			Reason: "[nome_identificacao] failed to match API requirements. It should look like this: Firstname Middlename(optional) Lastname",
		})
	} else if len(ser.Nomeidentificacao) > 100 {
		regexcheck = true
		Reasons = append(Reasons, ErrorBody{
			Reason: "[nome_identificacao] failed to match API requirements. It should have a maximum of 100 characters",
		})
	}
	r, _ = regexp.Compile(`\b[MF]{1}\b`)
	if !r.MatchString(ser.Sexo) {
		regexcheck = true
		Reasons = append(Reasons, ErrorBody{
			Reason: "[sexo] failed to match API requirements. It should look like this: M for male, F for female",
		})
	}
	if regexcheck {
		return c.JSON(400, Reasons)
	}
	// END OF REGEX CHEKING PHASE
	timestamp := time.Now().UTC().Format("2006-01-02T15:04:05-0700")
	b := md5.Sum([]byte(fmt.Sprintf(string(ser.Nome), string(timestamp))))
	bid := binary.BigEndian.Uint64(b[:])
	ser.Matriculainterna = int(bid % 99999)
	q := fmt.Sprintf(`
		INSERT INTO rh.servidor_tmp(
			nome, nome_identificacao, siape, id_pessoa, matricula_interna, id_foto,
			data_nascimento, sexo)
			VALUES ('%s', '%s', %d, %d, %d, null, '%s', '%s');
			`, ser.Nome, ser.Nomeidentificacao, ser.Siape, ser.Idpessoa, ser.Matriculainterna,
		ser.Datanascimento, ser.Sexo) //String formating
	rows, err := db.GetDB().Query(q)
	if err != nil {
		log.Println(err)
		return c.JSON(400, ErrorBody{
			Reason: err.Error(),
		})
	}
	defer rows.Close()
	c.Response().Header().Set(echo.HeaderLocation, fmt.Sprint("https://", c.Request().Host, "/api/servidor/", strconv.Itoa(ser.Matriculainterna)))
	return c.String(201, "")
}
