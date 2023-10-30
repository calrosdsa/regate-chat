create table if not exists conversations(
    conversation_id serial unique,
    establecimiento_id int REFERENCES establecimientos(establecimiento_id) on update cascade on delete cascade,
    profile_id int REFERENCES profiles(profile_id) on update cascade on delete cascade,
    primary key(establecimiento_id,profile_id)
);

create table if not exists conversation_message(
  id int primary key DEFAULT nextval('message_seq'),
  chat_id int not null,
  reply_to int,
  profile_id int,
  type_message smallint DEFAULT 0,
  created_at TIMESTAMP DEFAULT current_timestamp,
  content TEXT NOT NULL,
  establecimiento_id int,
  data text,
  is_user boolean,
  is_read boolean,
  is_deleted boolean default false,
  CONSTRAINT fk_conversation
  FOREIGN KEY(chat_id) 
  REFERENCES chat(id)  on delete cascade
);    


insert into conversation_message(conversation_id,sender_id,content)values(1,1,'First Message');
insert into conversation_message(conversation_id,sender_id,content,reply_to)values(1,1,'First Message',1);

insert into conversations(profile_id,establecimiento_id) values (1,1);