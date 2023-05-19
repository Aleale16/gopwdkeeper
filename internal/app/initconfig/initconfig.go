package initconfig

import (
	"flag"
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)
var PostgresDBURLflag *string
var PostgresDBURL string

var ServerKey = []byte("StrongPhrase_BIuaeruvlkjasdiu%2jl")
var Salt = []byte("StrongSalt_BIuaeruvlkjasdiu%2jl")

func InitFlags() {
	PostgresDBURLflag = flag.String("d", "postgres://postgres:1@localhost:5432/pwdkeeper", "DATABASE_URI flag")
}

// SetinitVars - init global vars according to ENV vars and flags passed.
func SetinitVars() {

	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "15:04:05"})
	fmt.Print("Logger params is set.")

	postgresDBURLENV, postgresDBURLexists := os.LookupEnv("DATABASE_URI")
	if !postgresDBURLexists{		
		PostgresDBURL = *PostgresDBURLflag
		fmt.Println("Set from flag: PostgresDBURL:", PostgresDBURL)
	} else {
		PostgresDBURL = postgresDBURLENV
		fmt.Println("Set from ENV: PostgresDBURL:", PostgresDBURL)
	}
	
}

func SetinitclientVars() {

	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "15:04:05"})
	fmt.Print("Logger params is set.")

}