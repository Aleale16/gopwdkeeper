package grpcserver

import (
	// импортируем пакет со сгенерированными protobuf-файлами
	"context"
	"pwdkeeper/internal/app/crypter"
	pb "pwdkeeper/internal/app/proto"
	"pwdkeeper/internal/app/storage"
	"sync"
)

//Authctx контекст авторизации
var Authctx context.Context
// ActionsServer поддерживает все необходимые методы сервера.
type ActionsServer struct {
	// нужно встраивать тип pb.Unimplemented<TypeName>
	// для совместимости с будущими версиями
	pb.UnimplementedActionsServer
}
// Login type for user struct
type Login string
// User type for user struct
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

//IsAuhtorized - checks if user is authorized
func (s *ActionsServer) IsAuhtorized(ctx context.Context, in *pb.IsAuhtorizedRequest) (*pb.IsAuhtorizedResponse, error) {
	var response pb.IsAuhtorizedResponse	

	response.Login = crypter.IsAuhtorized(in.Msg)
	user.login = "auhtorizedLogin"
	Authctx = context.WithValue(Authctx, user.login, response.Login)

	return &response, nil
}

// GetUser - returns user's FEK
func (s *ActionsServer) GetUser(ctx context.Context, in *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	var response pb.GetUserResponse
	
	response.Status, response.Fek = storage.GetUser(in.Login)
	//response.Status, response.Fek = storage.GetUser(fmt.Sprintf(("%s"), Authctx.Value("auhtorizedLogin")))

	return &response, nil
}

// StoreUser - save user in DB, returns FEK
func (s *ActionsServer) StoreUser(ctx context.Context, in *pb.StoreUserRequest) (*pb.StoreUserResponse, error) {
	var response pb.StoreUserResponse

	response.Status, response.Fek = storage.StoreUser(in.Login, in.Password, in.Fek)

	return &response, nil
}

// GetUserAuth - returns FEK to login/password pair
func (s *ActionsServer) GetUserAuth(ctx context.Context, in *pb.GetUserAuthRequest) (*pb.GetUserAuthResponse, error) {
	var response pb.GetUserAuthResponse

	response.Status, response.Fek = storage.AuthenticateUser(in.Login, in.Password)

	return &response, nil
}

// GetUserRecords - returns JSON list of user's records
func (s *ActionsServer) GetUserRecords(ctx context.Context, in *pb.GetUserRecordsRequest) (*pb.GetUserRecordsResponse, error) {
	var response pb.GetUserRecordsResponse
	//log.Debug().Msgf("Authctx auhtorizedLogin= %v", Authctx.Value("auhtorizedLogin"))
	//response.Status, response.UserRecordsJSON = storage.GetUserRecords(fmt.Sprintf(("%v"), Authctx.Value("auhtorizedLogin")))
	response.Status, response.UserRecordsJSON = storage.GetUserRecords(crypter.IsAuhtorized(in.Login))

	return &response, nil
}

// StoreSingleRecord - save record to auth user
func (s *ActionsServer) StoreSingleRecord(ctx context.Context, in *pb.StoreSingleRecordRequest) (*pb.StoreSingleRecordResponse, error) {
	var response pb.StoreSingleRecordResponse
	if crypter.IsAuhtorized(in.Login) != ""{
		response.Status, response.RecordID = storage.StoreRecord(in.DataName, in.SomeData, in.DataType, crypter.IsAuhtorized(in.Login))
	}

	return &response, nil
}

// UpdateRecord - update record to auth user
func (s *ActionsServer) UpdateRecord(ctx context.Context, in *pb.UpdateRecordRequest) (*pb.UpdateRecordResponse, error) {
	var response pb.UpdateRecordResponse
	var m sync.Mutex
	m.Lock()
		response.Status= storage.UpdateRecord(in.RecordID, in.EncryptedData, crypter.IsAuhtorized(in.Login))
	m.Unlock()
	return &response, nil
}

// DeleteRecord - delete record to auth user
func (s *ActionsServer) DeleteRecord(ctx context.Context, in *pb.DeleteRecordRequest) (*pb.DeleteRecordResponse, error) {
	var response pb.DeleteRecordResponse

	response.Status= storage.DeleteRecord(in.RecordID, crypter.IsAuhtorized(in.Login))

	return &response, nil
}

// GetSingleRecord - returns EncryptedData for auth user by Record ID
func (s *ActionsServer) GetSingleRecord(ctx context.Context, in *pb.GetSingleRecordRequest) (*pb.GetSingleRecordResponse, error) {
	var response pb.GetSingleRecordResponse

	//response.EncryptedData, response.DataType = storage.GetRecord(in.RecordID)
	response.EncryptedData, response.DataType = storage.GetRecord(in.RecordID, crypter.IsAuhtorized(in.Login))

	return &response, nil
}

// GetSingleNameRecord returns Record name for auth user by Record ID
func (s *ActionsServer) GetSingleNameRecord(ctx context.Context, in *pb.GetSingleNameRecordRequest) (*pb.GetSingleNameRecordResponse, error) {
	var response pb.GetSingleNameRecordResponse

	//response.EncryptedData, response.DataType = storage.GetRecord(in.RecordID)
	response.DataName = storage.GetNameRecord(in.RecordID, crypter.IsAuhtorized(in.Login))

	return &response, nil
}