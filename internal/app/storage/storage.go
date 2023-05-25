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
	login    string
	password string
	fek      string
}
type dataRecords struct {
	idrecord                         string
	namerecord, datarecord, datatype string
	login                            string
}

var (
	authUser authUsers
	record   dataRecords
)

type storagerUser interface {
	storeuser(ctx context.Context) (status string, authToken string)
	getuser(ctx context.Context) (status string, fek string)
	authenticateuser(ctx context.Context) (status string, fek string)
	getuserrecords(ctx context.Context) (status string, rowsDataRecordJSON string)
}
type storagerData interface {
	storerecord(ctx context.Context) (status string, recordID string)
	updaterecord(ctx context.Context) (status string)
	deleterecord(ctx context.Context) (status string)
	getrecord(ctx context.Context) (datarecord string, datatype string)
	getnamerecord(ctx context.Context) (namerecord string)
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
	authUser.login = login
	authUser.password = password
	authUser.fek = fek
	SUser = authUser
	return SUser.storeuser(ctx)
}

// GetUser - get user from DB
func GetUser(ctx context.Context, login string) (status string, publickey string) {
	log.Debug().Msg("func GetUser")
	authUser.login = login
	SUser = authUser
	return SUser.getuser(ctx)
}

// AuthenticateUser - authenticate User
func AuthenticateUser(ctx context.Context, login, password string) (status string, publickey string) {
	log.Debug().Msg("func GetUser")
	authUser.login = login
	authUser.password = password
	SUser = authUser
	return SUser.authenticateuser(ctx)
}

// GetUserRecords - get all user's records
func GetUserRecords(ctx context.Context, login string) (status string, rowsDataRecordJSON string) {
	log.Debug().Msg("func GetUserRecords")
	authUser.login = login
	SUser = authUser
	return SUser.getuserrecords(ctx)
}

// StoreRecord - save record in DB
func StoreRecord(ctx context.Context, namerecord, datarecord, datatype, login string) (status string, recordID string) {
	log.Debug().Msg("func StoreRecord")
	record.namerecord = namerecord
	record.datarecord = hex.EncodeToString([]byte(datarecord))
	record.datatype = datatype
	record.login = login
	SData = record
	return SData.storerecord(ctx)
}

// UpdateRecord - update record in DB
func UpdateRecord(ctx context.Context, recordID string, datarecord, login string) (status string) {
	log.Debug().Msg("func UpdateRecord")
	record.idrecord = recordID
	record.datarecord = hex.EncodeToString([]byte(datarecord))
	record.login = login
	SData = record
	return SData.updaterecord(ctx)
}

// DeleteRecord - delete record from DB
func DeleteRecord(ctx context.Context, recordID, login string) (status string) {
	log.Debug().Msg("func DeleteRecord")
	record.idrecord = recordID
	record.login = login
	SData = record
	return SData.deleterecord(ctx)
}

// GetRecord - get record row from DB
func GetRecord(ctx context.Context, idrecord, login string) (datarecord string, datatype string) {
	record.idrecord = idrecord
	record.login = login
	SData = record
	return SData.getrecord(ctx)
}

// GetNameRecord - get single record name from DB
func GetNameRecord(ctx context.Context, idrecord, login string) (namerecord string) {
	log.Debug().Msg("func GetNameRecord")
	record.idrecord = idrecord
	record.login = login
	SData = record
	return SData.getnamerecord(ctx)
}
