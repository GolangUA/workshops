create table if not exists users
(
    id      uuid primary key default gen_random_uuid(),
    token   text
);

create unique index if not exists users_id_uindex
    on users (id);

create table if not exists wallets
(
    id      uuid      primary key default gen_random_uuid(),
    balance bigint    not null,
    user_id uuid      not null,

    constraint wallets_users_id_fk
        foreign key (user_id) references users
            on update cascade on delete cascade
);

create table if not exists transactions
(
    id               uuid primary key default gen_random_uuid(),
    credit_wallet_id uuid      not null,
    debit_wallet_id  uuid      not null,
    amount           bigint    not null,
    type             smallint  not null,
    fee_amount       bigint    not null,
    fee_wallet_id    uuid      not null,
    credit_user_id   uuid      not null,
    debit_user_id    uuid      not null,


    constraint transactions_users_id_fk
        foreign key (credit_user_id) references users
            on update cascade on delete cascade,
    constraint transactions_users_id_fk_2
        foreign key (debit_user_id) references users
            on update cascade on delete cascade,
    constraint transactions_wallets_address_fk
        foreign key (credit_wallet_id) references wallets
            on update cascade on delete cascade,
    constraint transactions_wallets_address_fk_2
        foreign key (debit_wallet_id) references wallets
            on update cascade on delete cascade,
    constraint transactions_wallets_address_fk_3
        foreign key (fee_wallet_id) references wallets
            on update cascade on delete cascade
);

create index if not exists transactions_credit_address_index
    on transactions (credit_wallet_id);

create index if not exists transactions_credit_user_id_index
    on transactions (credit_user_id);

create index if not exists transactions_debit_address_index
    on transactions (debit_wallet_id);

create index if not exists transactions_debit_user_id_index
    on transactions (debit_user_id);

create unique index if not exists transactions_id_uindex
    on transactions (id);


INSERT INTO users (id) VALUES ('66aeb414-335a-4d1d-9dd9-6622b9c179a9');
INSERT INTO wallets (id, balance, user_id) VALUES ('85aa7525-4fdb-4436-a600-66ffc55e0f65', 0, '66aeb414-335a-4d1d-9dd9-6622b9c179a9');