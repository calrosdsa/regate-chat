create table if not exists chat(
    id serial primary key,
    parent_id int,
    second_parent_id int,
    type_chat smallint
);

create table if not exists deleted_message(
    id int,
    chat_id int,
    created_at timestamp default current_timestamp
);

create index idx_chat_id on deleted_message(chat_id);