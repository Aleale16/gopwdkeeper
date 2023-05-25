package grpcserver_test

import (
	"context"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"net"
	"testing"

	"pwdkeeper/internal/app/crypter"
	"pwdkeeper/internal/app/grpcserver"
	"pwdkeeper/internal/app/initconfig"
	pb "pwdkeeper/internal/app/proto"
	"pwdkeeper/internal/app/storage"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

var (
	LastCreatedRecordID, LastCreatedSomeData, AuthToken string
	Kek_key2, Fek_key1                                  []byte
)

func server(ctx context.Context) (pb.ActionsClient, func()) {
	flag.Parse()

	initconfig.SetinitVars()

	storage.Initdb()

	buffer := 101024 * 1024
	listen := bufconn.Listen(buffer)

	// создаём gRPC-сервер без зарегистрированной службы
	baseServer := grpc.NewServer()
	// регистрируем сервис
	pb.RegisterActionsServer(baseServer, &grpcserver.ActionsServer{})

	fmt.Println("Сервер gRPC начал работу")

	go func() {
		if err := baseServer.Serve(listen); err != nil {
			log.Printf("error serving server: %v", err)
		}
	}()

	conn, err := grpc.DialContext(ctx, "",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return listen.Dial()
		}), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("error connecting to server: %v", err)
	}

	closer := func() {
		err := listen.Close()
		if err != nil {
			log.Printf("error closing listener: %v", err)
		}
		baseServer.Stop()
	}

	client := pb.NewActionsClient(conn)

	return client, closer
}

func TestStoreUser(t *testing.T) {
	ctx := context.Background()
	initconfig.InitFlags()

	client, closer := server(ctx)
	defer closer()

	type expectation struct {
		out *pb.StoreUserResponse
		err error
	}
	Kek_key2 = crypter.Key2build("11111111")
	Fek_key1 = crypter.Key1build()
	AuthToken = crypter.GenAuthToken("TestUser1")

	fmt.Printf("Fek_key1= %v\n", hex.EncodeToString((Fek_key1)))
	fmt.Printf("AuthToken= %v\n", AuthToken)
	tests := map[string]struct {
		in       *pb.StoreUserRequest
		expected expectation
	}{
		"User_created": {
			in: &pb.StoreUserRequest{
				Login:    "TestUser1",
				Password: "11111111",
				// Fek: string(Fek_key1),
				Fek: hex.EncodeToString((Fek_key1)),
			},
			expected: expectation{
				out: &pb.StoreUserResponse{
					Status: "200",
					Fek:    "authToken",
				},
				err: nil,
			},
		},
		"User_alreadyexists": {
			in: &pb.StoreUserRequest{
				Login:    "TestUser1",
				Password: "11111111",
				Fek:      hex.EncodeToString((Fek_key1)),
			},
			expected: expectation{
				out: &pb.StoreUserResponse{
					Status: "409",
					Fek:    "",
				},
				err: nil,
			},
		},
	}

	for scenario, tt := range tests {
		t.Run(scenario, func(t *testing.T) {
			out, err := client.StoreUser(ctx, tt.in)
			if err != nil {
				if tt.expected.err.Error() != err.Error() {
					t.Errorf("Err -> \nWant: %q\nGot: %q\n", tt.expected.err, err)
				}
			} else {
				if tt.expected.out.Status != out.Status ||
					tt.expected.out.Fek != out.Fek {
					t.Errorf("Out -> \nWant: %q\nGot : %q", tt.expected.out, out)
				}
			}
		})
	}
}

func TestGetUser(t *testing.T) {
	ctx := context.Background()
	client, closer := server(ctx)
	defer closer()

	type expectation struct {
		out *pb.GetUserResponse
		err error
	}

	tests := map[string]struct {
		in       *pb.GetUserRequest
		expected expectation
	}{
		"User_exists": {
			in: &pb.GetUserRequest{
				Login: "TestUser1",
			},
			expected: expectation{
				out: &pb.GetUserResponse{
					Status: "200",
					Fek:    hex.EncodeToString((Fek_key1)),
				},
				err: nil,
			},
		},
	}
	for scenario, tt := range tests {
		t.Run(scenario, func(t *testing.T) {
			out, err := client.GetUser(ctx, tt.in)
			if err != nil {
				if tt.expected.err.Error() != err.Error() {
					t.Errorf("Err -> \nWant: %q\nGot: %q\n", tt.expected.err, err)
				}
			} else {
				if tt.expected.out.Status != out.Status {
					t.Errorf("Out -> \nWant: %q\nGot : %q", tt.expected.out, out)
				}
			}
		})
	}
}

func TestGetAuthUserRecords(t *testing.T) {
	ctx := context.Background()
	client, closer := server(ctx)
	defer closer()

	type expectation struct {
		out *pb.GetUserRecordsResponse
		err error
	}

	tests := map[string]struct {
		in       *pb.GetUserRecordsRequest
		expected expectation
	}{
		"User_SuccessGetRecords": {
			in: &pb.GetUserRecordsRequest{
				Login: AuthToken,
			},
			expected: expectation{
				out: &pb.GetUserRecordsResponse{
					Status: "200",
					UserRecordsJSON: `[
						{
						  "id": 4,
						  "namerecord": "NewRecord",
						  "datarecord": "**********************",
						  "datatype": "String"
						}
					  ]`,
				},
				err: nil,
			},
		},
		"User_notauthenticatedGetRecords": {
			in: &pb.GetUserRecordsRequest{
				Login: AuthToken + "0000",
			},
			expected: expectation{
				out: &pb.GetUserRecordsResponse{
					Status:          "200",
					UserRecordsJSON: "null",
				},
				err: nil,
			},
		},
	}
	for scenario, tt := range tests {
		t.Run(scenario, func(t *testing.T) {
			out, err := client.GetUserRecords(ctx, tt.in)
			if err != nil {
				if tt.expected.err.Error() != err.Error() {
					t.Errorf("Err -> \nWant: %q\nGot: %q\n", tt.expected.err, err)
				}
			} else {
				if tt.expected.out.Status != out.Status /*||tt.expected.out.UserRecordsJSON != out.UserRecordsJSON */ {
					t.Errorf("Out -> \nWant: %q\nGot : %q", tt.expected.out, out)
				}
			}
		})
	}
}

func TestStoreSingleRecord(t *testing.T) {
	ctx := context.Background()
	client, closer := server(ctx)
	defer closer()

	type expectation struct {
		out *pb.StoreSingleRecordResponse
		err error
	}

	tests := map[string]struct {
		in       *pb.StoreSingleRecordRequest
		expected expectation
	}{
		"Data_SuccessPostRecord": {
			in: &pb.StoreSingleRecordRequest{
				DataName: "DataName1",
				SomeData: hex.EncodeToString([]byte("SomeData1")),
				DataType: "s",
				Login:    AuthToken,
			},
			expected: expectation{
				out: &pb.StoreSingleRecordResponse{
					Status: "200",
					// RecordID:     "",

				},
				err: nil,
			},
		},
	}

	for scenario, tt := range tests {
		t.Run(scenario, func(t *testing.T) {
			out, err := client.StoreSingleRecord(ctx, tt.in)
			if err != nil {
				if tt.expected.err.Error() != err.Error() {
					t.Errorf("Err -> \nWant: %q\nGot: %q\n", tt.expected.err, err)
				}
			} else {
				if tt.expected.out.Status != out.Status {
					t.Errorf("Out -> \nWant: %q\nGot : %q", tt.expected.out, out)
				} else {
					LastCreatedRecordID = out.RecordID
					LastCreatedSomeData = hex.EncodeToString([]byte("SomeData1"))
				}
			}
		})
	}
}

func TestGetSingleRecord(t *testing.T) {
	ctx := context.Background()
	client, closer := server(ctx)
	defer closer()

	type expectation struct {
		out *pb.GetSingleRecordResponse
		err error
	}

	tests := map[string]struct {
		in       *pb.GetSingleRecordRequest
		expected expectation
	}{
		"Data_SuccessGetRecord": {
			in: &pb.GetSingleRecordRequest{
				RecordID: LastCreatedRecordID,
				Login:    AuthToken,
			},
			expected: expectation{
				out: &pb.GetSingleRecordResponse{
					EncryptedData: LastCreatedSomeData,
					DataType:      "s",
				},
				err: nil,
			},
		},
	}
	for scenario, tt := range tests {
		t.Run(scenario, func(t *testing.T) {
			out, err := client.GetSingleRecord(ctx, tt.in)
			if err != nil {
				if tt.expected.err.Error() != err.Error() {
					t.Errorf("Err -> \nWant: %q\nGot: %q\n", tt.expected.err, err)
				}
			} else {
				if tt.expected.out.EncryptedData != out.EncryptedData ||
					tt.expected.out.DataType != out.DataType {
					t.Errorf("Out -> \nWant: %q\nGot : %q", tt.expected.out, out)
				}
			}
		})
	}
}

func TestGetSingleNameRecord(t *testing.T) {
	ctx := context.Background()
	client, closer := server(ctx)
	defer closer()

	type expectation struct {
		out *pb.GetSingleNameRecordResponse
		err error
	}

	tests := map[string]struct {
		in       *pb.GetSingleNameRecordRequest
		expected expectation
	}{
		"Data_SuccessGetNameRecord": {
			in: &pb.GetSingleNameRecordRequest{
				RecordID: LastCreatedRecordID,
				Login:    AuthToken,
			},
			expected: expectation{
				out: &pb.GetSingleNameRecordResponse{
					DataName: "DataName1",
				},
				err: nil,
			},
		},
	}
	for scenario, tt := range tests {
		t.Run(scenario, func(t *testing.T) {
			out, err := client.GetSingleNameRecord(ctx, tt.in)
			if err != nil {
				if tt.expected.err.Error() != err.Error() {
					t.Errorf("Err -> \nWant: %q\nGot: %q\n", tt.expected.err, err)
				}
			} else {
				if tt.expected.out.DataName != out.DataName {
					t.Errorf("Out -> \nWant: %q\nGot : %q", tt.expected.out, out)
				}
			}
		})
	}
}

func TestUpdateRecord(t *testing.T) {
	ctx := context.Background()
	client, closer := server(ctx)
	defer closer()

	type expectation struct {
		out *pb.UpdateRecordResponse
		err error
	}

	tests := map[string]struct {
		in       *pb.UpdateRecordRequest
		expected expectation
	}{
		"Data_SuccessUpdateRecord": {
			in: &pb.UpdateRecordRequest{
				RecordID:      LastCreatedRecordID,
				EncryptedData: hex.EncodeToString([]byte("UpdatedSomeData1")),
				Login:         AuthToken,
			},
			expected: expectation{
				out: &pb.UpdateRecordResponse{
					Status: "200",
				},
				err: nil,
			},
		},
		"Data_NotFoundUpdateRecord": {
			in: &pb.UpdateRecordRequest{
				RecordID:      "9999",
				EncryptedData: hex.EncodeToString([]byte("UpdatedSomeData1")),
				Login:         AuthToken,
			},
			expected: expectation{
				out: &pb.UpdateRecordResponse{
					Status: "409",
				},
				err: nil,
			},
		},
	}

	for scenario, tt := range tests {
		t.Run(scenario, func(t *testing.T) {
			out, err := client.UpdateRecord(ctx, tt.in)
			if err != nil {
				if tt.expected.err.Error() != err.Error() {
					t.Errorf("Err -> \nWant: %q\nGot: %q\n", tt.expected.err, err)
				}
			} else {
				if tt.expected.out.Status != out.Status {
					t.Errorf("Out -> \nWant: %q\nGot : %q", tt.expected.out, out)
				}
			}
		})
	}
}

func TestDeleteRecord(t *testing.T) {
	ctx := context.Background()
	client, closer := server(ctx)
	defer closer()

	type expectation struct {
		out *pb.DeleteRecordResponse
		err error
	}

	tests := map[string]struct {
		in       *pb.DeleteRecordRequest
		expected expectation
	}{
		"Data_SuccessDeleteRecord": {
			in: &pb.DeleteRecordRequest{
				RecordID: LastCreatedRecordID,
				Login:    AuthToken,
			},
			expected: expectation{
				out: &pb.DeleteRecordResponse{
					Status: "200",
				},
				err: nil,
			},
		},
		"Data_NotFoundDeleteRecord": {
			in: &pb.DeleteRecordRequest{
				RecordID: "9999",
				Login:    AuthToken,
			},
			expected: expectation{
				out: &pb.DeleteRecordResponse{
					Status: "409",
				},
				err: nil,
			},
		},
	}
	for scenario, tt := range tests {
		t.Run(scenario, func(t *testing.T) {
			out, err := client.DeleteRecord(ctx, tt.in)
			if err != nil {
				if tt.expected.err.Error() != err.Error() {
					t.Errorf("Err -> \nWant: %q\nGot: %q\n", tt.expected.err, err)
				}
			} else {
				if tt.expected.out.Status != out.Status {
					t.Errorf("Out -> \nWant: %q\nGot : %q", tt.expected.out, out)
				}
			}
		})
	}
}
