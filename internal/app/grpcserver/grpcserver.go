package grpcserver

import (
	// импортируем пакет со сгенерированными protobuf-файлами
	"context"
	"pwdkeeper/internal/app/crypter"
	pb "pwdkeeper/internal/app/proto"
	"pwdkeeper/internal/app/storage"
)

var Authctx context.Context
// ActionsServer поддерживает все необходимые методы сервера.
type ActionsServer struct {
	// нужно встраивать тип pb.Unimplemented<TypeName>
	// для совместимости с будущими версиями
	pb.UnimplementedActionsServer
}
type Login string
type User struct {
	login Login
}
var user User

//IsAuhtorized хорошо бы запустить аналогично как если бы были хттп хэндлеры. 
//	r.Route("/api/user", func(r chi.Router) {
//		r.Use(isAuhtorized)
//		r.Post("/orders", handler.PostUserOrders)
//		...
//	})
// Не получается. Не могу передать Authctx = context.WithValue(Authctx, user.login, response.Login)
func (s *ActionsServer) IsAuhtorized(ctx context.Context, in *pb.IsAuhtorizedRequest) (*pb.IsAuhtorizedResponse, error) {
	var response pb.IsAuhtorizedResponse	

	response.Login = crypter.IsAuhtorized(in.Msg)
	user.login = "auhtorizedLogin"
	Authctx = context.WithValue(Authctx, user.login, response.Login)

	return &response, nil
}

func (s *ActionsServer) GetUser(ctx context.Context, in *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	var response pb.GetUserResponse
	
	response.Status, response.Fek = storage.GetUser(in.Login)
	//response.Status, response.Fek = storage.GetUser(fmt.Sprintf(("%s"), Authctx.Value("auhtorizedLogin")))

	return &response, nil
}

func (s *ActionsServer) StoreUser(ctx context.Context, in *pb.StoreUserRequest) (*pb.StoreUserResponse, error) {
	var response pb.StoreUserResponse

	response.Status, response.Fek = storage.StoreUser(in.Login, in.Password, in.Fek)

	return &response, nil
}

func (s *ActionsServer) GetUserAuth(ctx context.Context, in *pb.GetUserAuthRequest) (*pb.GetUserAuthResponse, error) {
	var response pb.GetUserAuthResponse

	response.Status, response.Fek = storage.AuthenticateUser(in.Login, in.Password)

	return &response, nil
}

func (s *ActionsServer) GetUserRecords(ctx context.Context, in *pb.GetUserRecordsRequest) (*pb.GetUserRecordsResponse, error) {
	var response pb.GetUserRecordsResponse
	//log.Debug().Msgf("Authctx auhtorizedLogin= %v", Authctx.Value("auhtorizedLogin"))
	//response.Status, response.UserRecordsJSON = storage.GetUserRecords(fmt.Sprintf(("%v"), Authctx.Value("auhtorizedLogin")))
	response.Status, response.UserRecordsJSON = storage.GetUserRecords(crypter.IsAuhtorized(in.Login))

	return &response, nil
}

func (s *ActionsServer) StoreSingleRecord(ctx context.Context, in *pb.StoreSingleRecordRequest) (*pb.StoreSingleRecordResponse, error) {
	var response pb.StoreSingleRecordResponse
	if crypter.IsAuhtorized(in.Login) != ""{
		response.Status, response.RecordID = storage.StoreRecord(in.DataName, in.SomeData, in.DataType, crypter.IsAuhtorized(in.Login))
	}

	return &response, nil
}

func (s *ActionsServer) UpdateRecord(ctx context.Context, in *pb.UpdateRecordRequest) (*pb.UpdateRecordResponse, error) {
	var response pb.UpdateRecordResponse

		response.Status= storage.UpdateRecord(in.RecordID, in.EncryptedData, crypter.IsAuhtorized(in.Login))

	return &response, nil
}

func (s *ActionsServer) DeleteRecord(ctx context.Context, in *pb.DeleteRecordRequest) (*pb.DeleteRecordResponse, error) {
	var response pb.DeleteRecordResponse

	response.Status= storage.DeleteRecord(in.RecordID, crypter.IsAuhtorized(in.Login))

	return &response, nil
}

func (s *ActionsServer) GetSingleRecord(ctx context.Context, in *pb.GetSingleRecordRequest) (*pb.GetSingleRecordResponse, error) {
	var response pb.GetSingleRecordResponse

	//response.EncryptedData, response.DataType = storage.GetRecord(in.RecordID)
	response.EncryptedData, response.DataType = storage.GetRecord(in.RecordID, crypter.IsAuhtorized(in.Login))

	return &response, nil
}

func (s *ActionsServer) GetSingleNameRecord(ctx context.Context, in *pb.GetSingleNameRecordRequest) (*pb.GetSingleNameRecordResponse, error) {
	var response pb.GetSingleNameRecordResponse

	//response.EncryptedData, response.DataType = storage.GetRecord(in.RecordID)
	response.DataName = storage.GetNameRecord(in.RecordID, crypter.IsAuhtorized(in.Login))

	return &response, nil
}