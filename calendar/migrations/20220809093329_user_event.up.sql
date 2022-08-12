create table user_event
(
    username varchar(31) not null
        constraint user_event_username_fkey
            references users (name),
    event_id varchar(36) not null
        constraint user_event_event_fkey
            references event (id),
    constraint user_event_pkey
        primary key (username, event_id)
);

