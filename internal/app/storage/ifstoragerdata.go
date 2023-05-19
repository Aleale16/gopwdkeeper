package storage

import (
	"context"
	"encoding/hex"
	"errors"

	"github.com/jackc/pgx"
	"github.com/rs/zerolog/log"
)

func (data dataRecords) storerecord() (status string, recordID string) {
	
err := PGdb.QueryRow(context.Background(), `INSERT into data(namerecord, datarecord, datatype, login_fkey) values($1,decode($2,'hex'),$3,$4) RETURNING (id)`, data.namerecord, data.datarecord, data.datatype, data.login ).Scan(&recordID)
	//result, err := PGdb.Exec(context.Background(), `INSERT into data(namerecord, datarecord, login_fkey) values($1,$2,$3)`, data.namerecord, data.datarecord, data.login)
	if err != nil {
		log.Error().Msg(err.Error())
		return "500", ""
	}
	log.Info().Msgf("NEW data record %v inserted successfully.", recordID)
	status = "200"
	return status, recordID
}

func (data dataRecords) updaterecord() (status string) {	
	result, err := PGdb.Exec(context.Background(), `UPDATE data SET datarecord = decode($1,'hex') where data.id = $2 and data.login_fkey=$3`, data.datarecord, data.idrecord, data.login )
	if err != nil {
		log.Error().Msg(err.Error())
		return "500"
	}
	if result.RowsAffected() == 0 {
		log.Error().Msgf("Data record is %v is not availible (deleted)", data.idrecord)
		status = "409"
		return status
	}
	log.Info().Msgf("Data record %v updated successfully.", data.idrecord)
	status = "200"
	return status
}

func (data dataRecords) deleterecord() (status string) {	
	result, err := PGdb.Exec(context.Background(), `DELETE FROM data WHERE data.id = $1 and data.login_fkey=$2`, data.idrecord, data.login )
	if err != nil {
		log.Error().Msg(err.Error())
		return "500"
	}
	if result.RowsAffected() == 0 {
		log.Error().Msgf("Data record is %v is not availible (deleted)", data.idrecord)
		status = "409"
		return status
	}
	log.Info().Msgf("Data record %v deleted successfully.", data.idrecord)
	status = "200"
	return status
}

func (data dataRecords) getrecord() (datarecord string, datatype string) {
	var namerecord string
	err := PGdb.QueryRow(context.Background(), `SELECT data.namerecord, encode(data.datarecord,'hex'), data.datatype FROM data WHERE id=$1 and data.login_fkey=$2`, data.idrecord, data.login).Scan(&namerecord, &datarecord, &datatype)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows){
			log.Error().Msgf("Data record is %v is not availible (deleted)", data.idrecord)
			return
		} else {
			log.Error().Msg(err.Error())
			return
		}
	}

	datarecordbyte, _ := hex.DecodeString(datarecord)
	log.Info().Msgf("Data record row %v extracted successfully.", data.idrecord)
	return string(datarecordbyte), datatype
}

func (data dataRecords) getnamerecord() (namerecord string) {
	err := PGdb.QueryRow(context.Background(), `SELECT data.namerecord FROM data WHERE id=$1 and data.login_fkey=$2`, data.idrecord, data.login).Scan(&namerecord)
	if err != nil {
		log.Error().Msgf(err.Error())
		}

	log.Info().Msgf("Data name record row %v extracted successfully.", data.idrecord)
	return namerecord
}

