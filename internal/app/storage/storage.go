package storage

import (
	"context"
	"encoding/hex"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

// PGdb - DB pool
var PGdb *pgxpool.Pool

type authUsers struct {
	ctx      context.Context
	login    string
	password string
	fek      string
}
type dataRecords struct {
	ctx                              context.Context
	idrecord                         string
	namerecord, datarecord, datatype string
	login                            string
}

var (
	authUser authUsers
	record   dataRecords
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
var (
	SUser storagerUser
	SData storagerData
)

// структура выводимого JSON
type rowDataRecord struct {
	IDrecord   int32  `json:"id,omitempty"`
	Namerecord string `json:"namerecord"`
	Datarecord string `json:"datarecord"`
	Datatype   string `json:"datatype"`
}

// StoreUser - store user in DB
func StoreUser(ctx context.Context, login string, password string, fek string) (status string, authToken string) {
	log.Debug().Msg("func StoreUser")
	authUser.ctx = ctx
	authUser.login = login
	authUser.password = password
	authUser.fek = fek
	SUser = authUser
	return SUser.storeuser()
}

// GetUser - get user from DB
func GetUser(ctx context.Context, login string) (status string, publickey string) {
	log.Debug().Msg("func GetUser")
	authUser.ctx = ctx
	authUser.login = login
	SUser = authUser
	return SUser.getuser()
}

// AuthenticateUser - authenticate User
func AuthenticateUser(ctx context.Context, login, password string) (status string, publickey string) {
	log.Debug().Msg("func GetUser")
	authUser.ctx = ctx
	authUser.login = login
	authUser.password = password
	SUser = authUser
	return SUser.authenticateuser()
}

// GetUserRecords - get all user's records
func GetUserRecords(ctx context.Context, login string) (status string, rowsDataRecordJSON string) {
	log.Debug().Msg("func GetUserRecords")
	authUser.ctx = ctx
	authUser.login = login
	SUser = authUser
	return SUser.getuserrecords()
}

// StoreRecord - save record in DB
func StoreRecord(ctx context.Context, namerecord, datarecord, datatype, login string) (status string, recordID string) {
	log.Debug().Msg("func StoreRecord")
	record.ctx = ctx
	record.namerecord = namerecord
	record.datarecord = hex.EncodeToString([]byte(datarecord))
	record.datatype = datatype
	record.login = login
	SData = record
	return SData.storerecord()
}

// UpdateRecord - update record in DB
func UpdateRecord(ctx context.Context, recordID string, datarecord, login string) (status string) {
	log.Debug().Msg("func UpdateRecord")
	record.ctx = ctx
	record.idrecord = recordID
	record.datarecord = hex.EncodeToString([]byte(datarecord))
	record.login = login
	SData = record
	return SData.updaterecord()
}

// DeleteRecord - delete record from DB
func DeleteRecord(ctx context.Context, recordID, login string) (status string) {
	log.Debug().Msg("func DeleteRecord")
	record.ctx = ctx
	record.idrecord = recordID
	record.login = login
	SData = record
	return SData.deleterecord()
}

// GetRecord - get record row from DB
func GetRecord(ctx context.Context, idrecord, login string) (datarecord string, datatype string) {
	record.ctx = ctx
	record.idrecord = idrecord
	record.login = login
	SData = record
	return SData.getrecord()
}

// GetNameRecord - get single record name from DB
func GetNameRecord(ctx context.Context, idrecord, login string) (namerecord string) {
	log.Debug().Msg("func GetNameRecord")
	record.ctx = ctx
	record.idrecord = idrecord
	record.login = login
	SData = record
	return SData.getnamerecord()
}
