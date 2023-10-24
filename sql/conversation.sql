create table if not exists conversations(
    conversation_id serial unique,
    establecimiento_id int REFERENCES establecimientos(establecimiento_id) on update cascade on delete cascade,
    profile_id int REFERENCES profiles(profile_id) on update cascade on delete cascade,
    primary key(establecimiento_id,profile_id)
);

create table if not exists conversation_message(
  id bigint unique,
  conversation_id int not null,
  reply_to bigint,
  sender_id int,
  created_at TIMESTAMP DEFAULT current_timestamp,
  content TEXT NOT NULL,
  primary key(id),
  CONSTRAINT fk_conversation
  FOREIGN KEY(conversation_id) 
  REFERENCES conversations(conversation_id)  on delete cascade
);

-- 9223372036854775807
-- 642830528364498454
-- CONSTRAINT fk_establecimiento
--     FOREIGN KEY(establecimiento_id) 
-- 	REFERENCES establecimientos(establecimiento_id)

insert into conversation_message(conversation_id,sender_id,content)values(1,1,'First Message');
insert into conversation_message(conversation_id,sender_id,content,reply_to)values(1,1,'First Message',1);

insert into conversations(profile_id,establecimiento_id) values (1,1);