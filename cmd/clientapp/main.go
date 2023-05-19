package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"

	"pwdkeeper/internal/app/crypter"
	"pwdkeeper/internal/app/initconfig"
	"pwdkeeper/internal/app/msgsender"

	pb "pwdkeeper/internal/app/proto"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	// buildVersion - global buildVersion value.
	buildVersion string
	// buildDate - global buildDate value.
	buildDate string
	// buildCommit - buildCommit value.
	buildCommit string
)

func main() {

	if buildVersion == "" {
		buildVersion = "N/A"
	}
	if buildDate == "" {
		buildDate = "N/A"
	}
	if buildCommit == "" {
		buildCommit = "N/A"
	}
	_, err := fmt.Printf("Build version: %s", buildVersion)
	if err != nil {
		log.Print(err)
	}
	_, err = fmt.Printf("Build date: %s", buildDate)
	if err != nil {
		log.Print(err)
	}
	_, err = fmt.Printf("Build commit: %s", buildCommit)
	if err != nil {
		log.Print(err)
	}

	initconfig.SetinitclientVars()

	log.Logger.Info().Msg("Starting CLIENT...")
	log.Logger.Info().Msg("Connecting to Server localhost:3200...")
	// устанавливаем соединение с сервером
	conn, err := grpc.Dial(":3200", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal().Err(err)
	}
	//defer conn.Close()
	// получаем переменную интерфейсного типа ActionsClient,
	// через которую будем отправлять сообщения
	c := pb.NewActionsClient(conn)
	log.Logger.Info().Msg("Connected.")
	log.Logger.Info().Msg("Starting UI...")

	StartUI(c)
	// функция, в которой будем отправлять сообщения
}

func StartUI(c pb.ActionsClient) {
	var (
		key1 []byte
		action, userRecordsJSON, status, AuthToken string
		menulevel int32
	)
	AuthToken = ""
	menulevel = 1
	userisNew := true
	login := ""
	loginstatus := ""
	password := ""
	key1enc :=""
	recordisNew := false
	recordIDname :=""			
	somedata :=""			
	datatype :=""			

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Hello! Enter Login: ")

	for {				
		// Scans a line from Stdin(Console)
		scanner.Scan()
		// Holds the string that scanned
		consoleInput := scanner.Text()
		log.Debug().Msgf("consoleInput =%v", consoleInput)
		switch menulevel {
//!Read Login			
			case 1: 				
				login = consoleInput
				loginstatus, key1enc = msgsender.SendUserGetmsg(c, login)
				if (loginstatus == "200") {
					userisNew = false
					fmt.Print("Enter Password: ")
				} else {
					fmt.Print("Create Password for NEW user: ")
				}
				
				menulevel = 2

//!Read Password and Show All encrypted user Records
			case 2: 
				password = consoleInput
				if userisNew {	
					//Key Encryption Key (KEK)
					key2 := crypter.Key2build(password)	
					//File Encryption Key (FEK) 
					key1 = crypter.Key1build()
					
					log.Debug().Msgf("Key2build(password) %v", string(key2))	
					log.Debug().Msgf("key1 %v", hex.EncodeToString(key1))

					EncryptedKey1 := crypter.EncryptKey1(key1, key2)
					
					key2 = crypter.Key2build(password)
					log.Debug().Msgf("Decrypted key1: %v", string(crypter.DecryptKey1((EncryptedKey1), key2)))
					//if msgsender.SendUserStoremsg(c, login, "password", string(EncryptedKey1)) == "200" {
					if msgsender.SendUserStoremsg(c, login, "password", hex.EncodeToString(EncryptedKey1)) == "200" {
						AuthToken = crypter.GenAuthToken(login)
						
						log.Info().Msgf("User %v created and logged in!", login)
						log.Debug().Msgf("User AuthToken = %v", AuthToken)
						fmt.Print("Enter NAME of new record to create: ")
						menulevel = 3
						} else {
							log.Error().Msg("Error!")
							fmt.Print("Enter Login: ")
							menulevel = 1
						}
					} else {
						noncekey1,_ := hex.DecodeString(key1enc)
						key1 = crypter.DecryptKey1([]byte(noncekey1), crypter.Key2build(password))
						log.Debug().Msgf(string(key1))
						if key1 != nil {
							AuthToken = crypter.GenAuthToken(login)
							log.Info().Msgf("Welcome, user %v! Logged in successfully.", login)
							log.Debug().Msgf("User AuthToken = %v", AuthToken)
							//crypter.IsAuhtorized(AuthToken)
							status, userRecordsJSON = msgsender.SendUserGetRecordsmsg(c, AuthToken)
							log.Debug().Msgf("SendUserGetRecordsmsg %v", status)
							log.Info().Msgf("List of user %v records:", login)
							log.Info().Msg(userRecordsJSON)
							fmt.Print("Enter ID of existing record or NAME of new record to create: ")
							menulevel = 3
							} else {
								log.Error().Msg("Error! Wrong Password!")
								fmt.Print("Enter Login: ")
								menulevel = 1
							}
					}

//!Create new or ASK Update/Delete existing Record 1)Read existing ID or new NAME
			case 3:
				recordIDname = consoleInput
				if _, err := strconv.Atoi(recordIDname); err == nil {
					log.Info().Msgf("%q looks like an ID number.\n Decrypting data...\n", recordIDname)
					recordisNew = false
				} else {
					recordisNew = true
				}
				if recordisNew {
					if  len(recordIDname)>1{
						fmt.Print("Enter data type to store [s]tring, [f]ile, [b]ankcard: ")
						menulevel = 31
						} else {
							log.Warn().Msg("Dataname length must be at least 2 symbols!")
							fmt.Print("Enter ID of existing record or NAME of new record to create: ")
							menulevel = 3
						}
					} else {
						log.Debug().Msgf("User AuthToken = %v", AuthToken)
						loadedsomedata, loadeddatatype := msgsender.SendGetSingleRecordmsg(c, recordIDname, AuthToken)
						loadeddataname := msgsender.SendGetSingleNameRecordmsg(c, recordIDname, AuthToken)

						log.Debug().Msgf("loadeddataname: %v",loadeddataname)
						log.Debug().Msgf("loadedsomedata: %v",loadedsomedata)
						log.Debug().Msgf("loadeddatatype: %v",loadeddatatype)
						if loadedsomedata != ""{
							noncedata,_ := hex.DecodeString(loadedsomedata)
							somedataDecrypted := crypter.DecryptData([]byte(noncedata),key1)
							log.Debug().Msgf("somedataDecrypted: %v", string(somedataDecrypted))

							log.Info().Msgf("Decrypted:\n Name=%v\n Data=%q\n Type=%v\n",loadeddataname, string(somedataDecrypted), loadeddatatype) 

							fmt.Printf("[u]pdate, [d]elete, [r]eturn? ")
							//fmt.Printf("Enter somedata to update record ID %v: ", recordIDname)
							menulevel = 41
						} else {
							log.Error().Msgf("Data record with ID %v is not availible (deleted)", recordIDname)
							menulevel = gotoMenulevel3(c, login, AuthToken)
							fmt.Print("Enter ID of existing record or NAME of new record to create: ")
						}

					}				

//!Create new Record 2)Read new datatype
			case 31:
				datatype = consoleInput	
//TODO different input logic for datatypes
				switch datatype{
					case "s":
						datatype = "String"
						fmt.Print("Enter somedata text to store: ")
						menulevel = 32
					case "f":
						datatype = "File"
						fmt.Print("Enter filepath to store: ")
						menulevel = 32
					case "b":
						datatype = "Bankcard"
						fmt.Print("Enter 20 digits, exp, cvc/cvv to store: ")
						menulevel = 32
					default:
						log.Warn().Msgf("Wrong data type %v!", datatype)
						fmt.Print("Enter data type to store [s]tring, [f]ile, [b]ankcard: ")
						menulevel = 31
					}				
				
//!Create new Record 3)Read new someDATA and Store NEW record returning created id
			case 32:
				somedata = consoleInput
				somedataenc := crypter.EncryptData(somedata, key1)
				log.Debug().Msgf("Created somedataenc= %v", hex.EncodeToString(somedataenc))
				status, recordID := msgsender.SendUserStoreRecordmsg(c, recordIDname, hex.EncodeToString(somedataenc), datatype, AuthToken)
				if status == "200" {
					log.Info().Msg("Created new record with ID=")
					log.Info().Msg(recordID)
				} else {
					log.Error().Msg("Error creating NEW data record!")
				}
				menulevel = gotoMenulevel3(c, login, AuthToken)
				fmt.Print("Enter ID of existing record or NAME of new record to create: ")
				
//![u]pdate, [d]elete, [r]eturn? someDATA
			case 41:
				action = consoleInput	
				switch action{
					case "u":
						fmt.Printf("Enter somedata to update record ID %v: ", recordIDname)
						menulevel = 42
//! DELETING record	
					case "d":
						if msgsender.SendDeleteRecordmsg(c, recordIDname, AuthToken) == "200" {
							log.Info().Msgf("Record ID %v deleted successfully", recordIDname)
						} else {
							log.Error().Msgf("Data record with ID %v is not availible (already deleted)", recordIDname)
						}
						menulevel = gotoMenulevel3(c, login, AuthToken)
						fmt.Print("Enter ID of existing record or NAME of new record to create: ")
					case "r":
						menulevel = gotoMenulevel3(c, login, AuthToken)
						fmt.Print("Enter ID of existing record or NAME of new record to create: ")
				}
				 		
//! UPDATING record			
			case 42:
				somedata = consoleInput
				somedataenc := crypter.EncryptData(somedata, key1)
				if msgsender.SendUpdateRecordmsg(c, recordIDname, hex.EncodeToString(somedataenc), AuthToken) == "200" {
					log.Info().Msgf("Record ID %v updated successfully", recordIDname)
				} else {
					log.Error().Msgf("Data record with ID %v is not availible (deleted)", recordIDname)
				}
				menulevel = gotoMenulevel3(c, login, AuthToken)
				fmt.Print("Enter ID of existing record or NAME of new record to create: ")
		}					
	}
}

func gotoMenulevel3(c pb.ActionsClient, login, AuthToken string) (menulevel int32) {
	_, userRecordsJSON := msgsender.SendUserGetRecordsmsg(c, AuthToken)
	log.Info().Msgf("List of user %v records:", login)
	log.Info().Msg(userRecordsJSON)
	menulevel = 3
	return menulevel
}

