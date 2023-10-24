create table if not exists grupo_message(
  id bigint not null,
  chat_id bigint not null,
  reply_to bigint references grupo_message,
  profile_id int not null,
  type_message int DEFAULT 0,
  created_at TIMESTAMP DEFAULT current_timestamp,
  content TEXT NOT NULL,
  grupo_id int,
  data text,
  PRIMARY KEY(id),
  CONSTRAINT fk_grupo
  FOREIGN KEY(grupo_id) 
	REFERENCES grupos(grupo_id) on delete cascade);