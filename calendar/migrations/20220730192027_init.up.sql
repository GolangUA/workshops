create table event
(
    id             varchar(36)              not null primary key,
    title          varchar(256)             not null,
    description    text,
    timestamp_from timestamp with time zone not null,
    timestamp_to   timestamp with time zone not null,
    notes          varchar(1024)[]
);
