create sequence message_seq;

create table if not exists grupo_message(
  id int primary key DEFAULT nextval('message_seq'),
  chat_id int not null,
  reply_to int references grupo_message,
  profile_id int not null,
  type_message smallint DEFAULT 0,
  created_at TIMESTAMP DEFAULT current_timestamp,
  content TEXT NOT NULL,
  is_deleted boolean default false,
  grupo_id int,
  data text,
  CONSTRAINT fk_grupo
  FOREIGN KEY(grupo_id) 
	REFERENCES grupos(grupo_id) on delete cascade);


  create table if not exists user_grupo(
  id serial,
  profile_id int not null,
  grupo_id int not null,
  created_at timestamp DEFAULT current_timestamp,
  updated_at timestamp DEFAULT current_timestamp,
  last_update_messages timestamp DEFAULT current_timestamp, 
  estado int DEFAULT 0,
  is_admin boolean DEFAULT false,
  is_out boolean default false,
  primary key (profile_id,grupo_id),
  CONSTRAINT fk_profile
  FOREIGN KEY(profile_id) 
	REFERENCES profiles(profile_id),
  CONSTRAINT fk_grupo
  FOREIGN KEY(grupo_id) 
	REFERENCES grupos(grupo_id) on delete cascade
);