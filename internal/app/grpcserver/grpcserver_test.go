package grpcserver_test

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"pwdkeeper/internal/app/grpcserver"
	"pwdkeeper/internal/app/initconfig"
	pb "pwdkeeper/internal/app/proto"
	"pwdkeeper/internal/app/storage"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
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

	client:= pb.NewActionsClient(conn)

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

	tests := map[string]struct {
		in       *pb.StoreUserRequest
		expected expectation
	}{
		"User_created": {
			in: &pb.StoreUserRequest{
				Login: "TestUser1",
				Password: "11111111",
				Fek: "5fa06d0b64facf315275f740d850e36bad092368b54116a4346f97306c92d82f6c1236eef780b3b81f0d6a239e9c01a12a47af277afddac3f18c0a9d",
			},
			expected: expectation{
				out: &pb.StoreUserResponse{
					Status:     "200",
					Fek: "authToken",
				},
				err: nil,
			},
		},
		"User_alreadyexists": {
			in: &pb.StoreUserRequest{
				Login: "TestUser1",
				Password: "11111111",
				Fek: "5fa06d0b64facf315275f740d850e36bad092368b54116a4346f97306c92d82f6c1236eef780b3b81f0d6a239e9c01a12a47af277afddac3f18c0a9d",
			},
			expected: expectation{
				out: &pb.StoreUserResponse{
					Status:     "409",
					Fek: "",
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
					tt.expected.out.Fek != out.Fek  {
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
					Status:     "200",
					Fek: "5fa06d0b64facf315275f740d850e36bad092368b54116a4346f97306c92d82f6c1236eef780b3b81f0d6a239e9c01a12a47af277afddac3f18c0a9d",
				},
				err: nil,
			},
		},
		"User_doesntexists": {
			in: &pb.GetUserRequest{
				Login: "BadTestUser1",
			},
			expected: expectation{
				out: &pb.GetUserResponse{
					Status:     "401",
					Fek: "",
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
					if tt.expected.out.Status != out.Status ||
						tt.expected.out.Fek != out.Fek  {
						t.Errorf("Out -> \nWant: %q\nGot : %q", tt.expected.out, out)
					}
				}
	
			})
		}
	}

func TestGetAuthUser(t *testing.T) {

ctx := context.Background()
client, closer := server(ctx)
defer closer()

type expectation struct {
	out *pb.GetUserAuthResponse
	err error
}

tests := map[string]struct {
	in       *pb.GetUserAuthRequest
	expected expectation
}{
	"User_authenticated": {
		in: &pb.GetUserAuthRequest{
			Login: "TestUser1",
			Password: "Password",
		},
		expected: expectation{
			out: &pb.GetUserAuthResponse{
				Status:     "200",
				Fek: "5fa06d0b64facf315275f740d850e36bad092368b54116a4346f97306c92d82f6c1236eef780b3b81f0d6a239e9c01a12a47af277afddac3f18c0a9d",
			},
			err: nil,
		},
	},
	"User_notauthenticated": {
		in: &pb.GetUserAuthRequest{
			Login: "TestUser1",
			Password: "BadPassword",
		},
		expected: expectation{
			out: &pb.GetUserAuthResponse{
				Status:     "401",
				Fek: "",
			},
			err: nil,
		},
	},
}
	for scenario, tt := range tests {
		t.Run(scenario, func(t *testing.T) {
			out, err := client.GetUserAuth(ctx, tt.in)
			if err != nil {
				if tt.expected.err.Error() != err.Error() {
					t.Errorf("Err -> \nWant: %q\nGot: %q\n", tt.expected.err, err)
				}
			} else {
				if tt.expected.out.Status != out.Status ||
					tt.expected.out.Fek != out.Fek  {
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
				Login: "TestUser1",
			},
			expected: expectation{
				out: &pb.GetUserRecordsResponse{
					Status:     "200",
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
				Login: "BadTestUser1",
			},
			expected: expectation{
				out: &pb.GetUserRecordsResponse{
					Status:     "500",
					UserRecordsJSON: "",
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
					if tt.expected.out.Status != out.Status ||
						tt.expected.out.UserRecordsJSON != out.UserRecordsJSON  {
						t.Errorf("Out -> \nWant: %q\nGot : %q", tt.expected.out, out)
					}
				}
	
			})
		}
	}