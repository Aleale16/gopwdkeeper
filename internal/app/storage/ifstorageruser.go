package storage

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/jackc/pgx"
	"github.com/rs/zerolog/log"
)

func (user authUsers) storeuser() (status string, authToken string) {
	result, err := PGdb.Exec(context.Background(), `INSERT into users(login, password, fek) values ($1, $2, $3) on conflict (login) DO NOTHING`, user.login, user.password, user.fek)
	if err != nil {
		log.Fatal().Err(err)
		return
	}
	if result.RowsAffected() == 0 {
		log.Warn().Msgf("NEW user was not created! Login %v already exists!", user.login)
		status = "409"
	} else {
		log.Info().Msgf("NEW user %v registered successfully.", user.login)
		status = "200"
		authToken = "authToken"
	}
	return status, authToken
}

func (user authUsers) getuser() (status string, fek string) {
	err := PGdb.QueryRow(context.Background(), `SELECT users.fek FROM users WHERE login=$1`, user.login).Scan(&fek)
	if err != nil {		
		if errors.Is(err, pgx.ErrNoRows){
			log.Error().Msg("User doesn't exist")
			status = "401"
		} else {
			log.Debug().Msg(err.Error())
			status = "500"
		}
	} else {
		log.Info().Msg("User exists")
		status = "200"
		}
	return status, fek
}

func (user authUsers) authenticateuser() (status string, fek string) {
	err := PGdb.QueryRow(context.Background(), `SELECT users.fek FROM users WHERE login=$1 AND password=$2`, user.login, user.password).Scan(&fek)
	if err != nil {		
		if errors.Is(err, pgx.ErrNoRows){
			log.Error().Msg("User login or password is invalid")
			status = "401"
		} else {
			log.Debug().Msg(err.Error())
			status = "500"
		}
	} else {
		log.Info().Msg("User login and password are OK")
		status = "200"
		}
	return status, fek
}

func (user authUsers) getuserrecords() (status string, DataRecordsJSON string) {
	var id int32
	var namerecord, datarecord, datatype, datarecordmask string
	var rowsDataRecordJSON []rowDataRecord
	rows, err := PGdb.Query(context.Background(), `SELECT data.id, data.namerecord, encode(data.datarecord,'hex'), data.datatype FROM data WHERE login_fkey=$1`, user.login)
	if err != nil {
		log.Error().Msgf(err.Error())
		}
	//defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&id, &namerecord, &datarecord, &datatype)
		if err != nil {
			log.Debug().Str("Query","SELECT data.id, data.namerecord, data.datarecord, data.datatype FROM data WHERE login_fkey=$1").Msg(err.Error())
			log.Error().Msgf(err.Error())
			status = "500"
			return
		}
		//datarecordbyte, _ := hex.DecodeString(datarecord)
		switch datatype{
			case "String":
				datarecordmask = "**********************"
			case "File":
				datarecordmask = "*.*"
			case "Bankcard":
				datarecordmask = "**** **** **** ****, **/**, ***"
		}
		rowsDataRecordJSON = append(rowsDataRecordJSON, rowDataRecord{
			IDrecord:			id,
			Namerecord:			namerecord,
			//Datarecord:			string(datarecordbyte),
			Datarecord:			datarecordmask,
			Datatype:			datatype,
		})
	}
	JSONdata, err := json.MarshalIndent(rowsDataRecordJSON, "", "  ")
	if err != nil {
		log.Fatal().Str("JSONdata","rowsDataRecordJSON").Msg(err.Error())
	}
	log.Info().Msgf("Data rows %v extracted successfully.", string(JSONdata))
	status = "200"
	return status, string(JSONdata)
}