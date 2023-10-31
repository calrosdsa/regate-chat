create table if not exists sala_message(
  id int primary key DEFAULT nextval('message_seq'),
  chat_id int not null,
  reply_to bigint references sala_message,
  profile_id int not null,
  type_message smallint DEFAULT 0,
  created_at TIMESTAMP DEFAULT current_timestamp,
  content TEXT NOT NULL,
  is_deleted boolean default false,
  sala_id int,
  data text,
  CONSTRAINT fk_sala
  FOREIGN KEY(sala_id) 
	REFERENCES salas(sala_id) on delete cascade
);



create table if not exists users_sala(
  id serial,
  profile_id int not null,
  precio decimal(12,2) not null,
  sala_id int not null,
  created_at timestamp not null,
  updated_at timestamp DEFAULT current_timestamp,
  estado int DEFAULT 0,
  is_admin boolean DEFAULT false,
  is_out boolean default false,
  primary key (id),
  CONSTRAINT fk_profile
  FOREIGN KEY(profile_id) 
	REFERENCES profiles(profile_id),
  CONSTRAINT fk_sala
  FOREIGN KEY(sala_id) 
	REFERENCES salas(sala_id) on delete cascade
);