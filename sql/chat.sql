create table if not exists chat(
    id bigserial primary key,
    parent_id int,
    second_parent_id int,
    type_chat smallint
);