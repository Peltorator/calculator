drop table if exists history cascade;
create table history
(
    id        serial primary key,
    expression     varchar(255) not null,
    result  varchar(255) not null
);