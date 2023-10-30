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

