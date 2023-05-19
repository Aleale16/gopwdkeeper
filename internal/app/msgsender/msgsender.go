package msgsender

import (
	"context"
	pb "pwdkeeper/internal/app/proto"

	"github.com/rs/zerolog/log"
)

//SendIsAuhtorizedmsg return user login if sign is valid
func SendIsAuhtorizedmsg(c pb.ActionsClient, authToken string) (login string) {
	resp, err := c.IsAuhtorized(context.Background(), &pb.IsAuhtorizedRequest{
		Msg: authToken,
	})
	if err != nil {
		log.Fatal().Err(err)
	}
	return resp.Login
}

//SendUserGetmsg return user login
func SendUserGetmsg(c pb.ActionsClient, login string) (status string, key1enc string) {
	resp, err := c.GetUser(context.Background(), &pb.GetUserRequest{
		Login: login,
	})
	if err != nil {
		log.Fatal().Err(err)
	}
	return resp.Status, resp.Fek
}

//SendUserAuthmsg Not used
func SendUserAuthmsg(c pb.ActionsClient, login, password string) (status string){
	resp, err := c.GetUserAuth(context.Background(), &pb.GetUserAuthRequest{
		Login: login,
		Password: password,
	})
	if err != nil {
		log.Fatal().Err(err)
	}
	return resp.Status
}

//SendUserStoremsg stores new user
func SendUserStoremsg(c pb.ActionsClient, login, password, fek string) (status string){
	resp, err := c.StoreUser(context.Background(), &pb.StoreUserRequest{
		Login: login,
		Password: password,
		Fek: fek,
	})
	if err != nil {
		log.Fatal().Err(err)
	}
	return resp.Status
}

//SendUserGetRecordsmsg loads all user's records
func SendUserGetRecordsmsg(c pb.ActionsClient, login string) (status string, UserRecordsJSON string){
	resp, err := c.GetUserRecords(context.Background(), &pb.GetUserRecordsRequest{
		Login: login,
	})
	if err != nil {
		log.Fatal().Err(err)
	}
	return resp.Status, resp.UserRecordsJSON
}

//SendUserStoreRecordmsg stores new single user's records
func SendUserStoreRecordmsg(c pb.ActionsClient, namerecord, datarecord, datatype, login string) (status string, RecordID string){
	resp, err := c.StoreSingleRecord(context.Background(), &pb.StoreSingleRecordRequest{
		DataName: namerecord,
		SomeData: datarecord,
		DataType: datatype,
		Login: login,
	})
	if err != nil {
		log.Fatal().Err(err)
	}
	return resp.Status, resp.RecordID
}

//SendUpdateRecordmsg updates single user's record
func SendUpdateRecordmsg(c pb.ActionsClient, recordID, datarecord, login string) (status string){
	resp, err := c.UpdateRecord(context.Background(), &pb.UpdateRecordRequest{
		RecordID: recordID,
		EncryptedData: datarecord,
		Login: login,
	})
	if err != nil {
		log.Fatal().Err(err)
	}
	return resp.Status
}

//SendDeleteRecordmsg deletes single user's record
func SendDeleteRecordmsg(c pb.ActionsClient, recordID, login string) (status string){
	resp, err := c.DeleteRecord(context.Background(), &pb.DeleteRecordRequest{
		RecordID: recordID,
		Login: login,
	})
	if err != nil {
		log.Fatal().Err(err)
	}
	return resp.Status
}

//SendGetSingleRecordmsg loads single user's record
func SendGetSingleRecordmsg(c pb.ActionsClient, recordID, login string) (somedataenc, datatype string){
	resp, err := c.GetSingleRecord(context.Background(), &pb.GetSingleRecordRequest{
		RecordID: recordID,
		Login: login,
	})
	if err != nil {
		log.Fatal().Err(err)
	}
	return resp.EncryptedData, resp.DataType
}

//SendGetSingleRecordmsg loads single user's name record Not used
func SendGetSingleNameRecordmsg(c pb.ActionsClient, recordID, login string) (namerecord string){
	resp, err := c.GetSingleNameRecord(context.Background(), &pb.GetSingleNameRecordRequest{
		RecordID: recordID,
		Login: login,
	})
	if err != nil {
		log.Fatal().Err(err)
	}

	return resp.DataName
}