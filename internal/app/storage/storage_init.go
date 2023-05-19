package storage

import (
	"context"
	"fmt"
	"pwdkeeper/internal/app/initconfig"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

func Initdb() {
	//----------------------------//
	//Подключаемся к СУБД postgres
	//----------------------------//
	//urlExample := "postgres://postgres:1@localhost:5432/pwdkeeper"
	//os.Setenv("DATABASE_DSN", urlExample)
	//initconfig.PostgresDBURL = urlExample
	if initconfig.PostgresDBURL != "" {
		poolConfig, err := pgxpool.ParseConfig(initconfig.PostgresDBURL)
		if err != nil {
			log.Error().Err(err)
			log.Fatal().Msg("Unable to parse DATABASE_DSN")
		}
		//log.Debug().Msgf("poolConfig: %v", poolConfig)

		PGdb, err = pgxpool.NewWithConfig(context.Background(), poolConfig)
		if err != nil {
			fmt.Println("ERROR! PGdbOpened = false")
			panic(err)
		} else {
			_, err := PGdb.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS public.users
			(
				login character varying(20) NOT NULL,
				password character varying(400) NOT NULL,
				fek character varying(400),
				CONSTRAINT login PRIMARY KEY (login),
				CONSTRAINT login UNIQUE (login)
			);
			
			CREATE TABLE IF NOT EXISTS public.data
			(
				id serial NOT NULL,
				namerecord character varying(50),
				datarecord bytea,
				datatype character varying(10),
				login_fkey character varying(20) NOT NULL,
				PRIMARY KEY (id)
			);
			
			ALTER TABLE IF EXISTS public.data
				ADD FOREIGN KEY (login_fkey)
				REFERENCES public.users (login) MATCH SIMPLE
				ON UPDATE NO ACTION
				ON DELETE NO ACTION
				NOT VALID;
			
			END;`)
			if err != nil {
				log.Error().Err(err)
			}
			log.Info().Msg("PGdbOpened = TRUE")
		}
	} else {
		log.Info().Msg("PGdbOpened = FALSE")
	}

}