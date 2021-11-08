create table if not exists wallets
(
    address text      not null,
    balance bigint    not null,
    user_id serial    not null,
    id      bigserial not null,
    constraint wallets_pk
    primary key (address)
    );

create unique index if not exists wallets_address_uindex
    on wallets (address);

create unique index if not exists wallets_id_uindex
    on wallets (id);

create table if not exists users
(
    id    serial not null,
    token text,
    constraint users_pk
    primary key (id)
    );

create unique index if not exists users_id_uindex
    on users (id);

create table if not exists transactions
(
    credit_address text      not null,
    debit_address  text      not null,
    amount         bigint    not null,
    type           smallint  not null,
    fee_amount     bigint    not null,
    fee_address    text      not null,
    credit_user_id integer   not null,
    debit_user_id  integer   not null,
    id             bigserial not null,
    constraint transactions_pk
    primary key (id),
    constraint transactions_users_id_fk
    foreign key (credit_user_id) references users
    on update cascade on delete cascade,
    constraint transactions_users_id_fk_2
    foreign key (debit_user_id) references users
    on update cascade on delete cascade,
    constraint transactions_wallets_address_fk
    foreign key (credit_address) references wallets
    on update cascade on delete cascade,
    constraint transactions_wallets_address_fk_2
    foreign key (debit_address) references wallets
    on update cascade on delete cascade,
    constraint transactions_wallets_address_fk_3
    foreign key (fee_address) references wallets
    on update cascade on delete cascade
    );

create index if not exists transactions_credit_address_index
    on transactions (credit_address);

create index if not exists transactions_credit_user_id_index
    on transactions (credit_user_id);

create index if not exists transactions_debit_address_index
    on transactions (debit_address);

create index if not exists transactions_debit_user_id_index
    on transactions (debit_user_id);

create unique index if not exists transactions_id_uindex
    on transactions (id);

create table if not exists rates
(
    id   bigserial        not null,
    rate double precision not null,
    constraint rates_pk
    primary key (id)
    );


INSERT INTO users (id) VALUES (0);
INSERT INTO wallets (address, balance, user_id) VALUES ('1feeaddressasdasd123asd', 0, 0);