package storage

import (
	"encoding/hex"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

// PGdb - DB pool
var PGdb *pgxpool.Pool

type authUsers struct{ login string; password string; fek string }
type dataRecords struct{ idrecord string; namerecord, datarecord, datatype string; login string }
var (	authUser authUsers;
		record dataRecords;
	)
type storagerUser interface {
	storeuser() (status string, authToken string)
	getuser() (status string, fek string)
	authenticateuser() (status string, fek string)
	getuserrecords() (status string, rowsDataRecordJSON string)
}
type storagerData interface {
	storerecord() (status string, recordID string)
	updaterecord() (status string)
	deleterecord() (status string)
	getrecord() (datarecord string, datatype string)
	getnamerecord() (namerecord string)
	
}

// SUser - user if type
// SData - data if type
var (	SUser storagerUser;
		SData storagerData;
	)

//структура выводимого JSON	 
type rowDataRecord struct {
		IDrecord		int32 	`json:"id,omitempty"`
		Namerecord 		string 	`json:"namerecord"`    
		Datarecord 		string 	`json:"datarecord"`       
		Datatype 		string 	`json:"datatype"`       
	}


// StoreUser - store user in DB 
func StoreUser(login string, password string, fek string) (status string, authToken string) {
	log.Debug().Msg("func StoreUser")
	authUser.login = login
	authUser.password = password
	authUser.fek = fek
	SUser = authUser
	return SUser.storeuser()
}

// GetUser - get user from DB
func GetUser(login string) (status string, publickey string){
	log.Debug().Msg("func GetUser")
	authUser.login = login
	SUser = authUser
	return SUser.getuser()
}

// AuthenticateUser - authenticate User
func AuthenticateUser(login, password string) (status string, publickey string){
	log.Debug().Msg("func GetUser")
	authUser.login = login
	authUser.password = password
	SUser = authUser
	return SUser.authenticateuser()
}

// GetUserRecords - get all user's records
func GetUserRecords(login string) (status string, rowsDataRecordJSON string){
	log.Debug().Msg("func GetUserRecords")
	authUser.login = login
	SUser = authUser
	return SUser.getuserrecords()
}

// StoreRecord - save record in DB
func StoreRecord(namerecord, datarecord, datatype, login string) (status string, recordID string){
	log.Debug().Msg("func StoreRecord")
	record.namerecord = namerecord
	record.datarecord = hex.EncodeToString([]byte(datarecord))
	record.datatype = datatype
	record.login = login
	SData = record
	return SData.storerecord()
}


// UpdateRecord - update record in DB
func UpdateRecord(recordID string, datarecord, login string) (status string){
	log.Debug().Msg("func UpdateRecord")
	record.idrecord = recordID
	record.datarecord = hex.EncodeToString([]byte(datarecord))
	record.login = login
	SData = record
	return SData.updaterecord()
}

//DeleteRecord - delete record from DB
func DeleteRecord(recordID, login string) (status string){
	log.Debug().Msg("func DeleteRecord")
	record.idrecord = recordID
	record.login = login
	SData = record
	return SData.deleterecord()
}

// GetRecord - get record row from DB
func GetRecord(idrecord, login string) (datarecord string, datatype string){
	record.idrecord = idrecord
	record.login = login
	SData = record
	return SData.getrecord()
}

// GetNameRecord - get single record name from DB
func GetNameRecord(idrecord, login string) (namerecord string){
	log.Debug().Msg("func GetNameRecord")
	record.idrecord = idrecord
	record.login = login
	SData = record
	return SData.getnamerecord()
}

