create table users
(
    name     varchar(31) not null primary key,
    password varchar(60) not null,
    timezone varchar(32)
);
