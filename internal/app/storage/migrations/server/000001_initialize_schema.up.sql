CREATE TABLE IF NOT EXISTS public.users
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
			
			END;