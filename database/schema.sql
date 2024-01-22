create table if not exists users
(
    id         bigint    not null,
    username   char(32)  not null,
    created_at timestamp not null,
    deleted_at timestamp,
    constraint pk_7
        primary key (id)
);

create table if not exists guilds
(
    id       bigint       not null,
    name     varchar(100) not null,
    owner_id bigint       not null,
    constraint pk_3
        primary key (id),
    constraint fk_2
        foreign key (owner_id) references users
);

create table if not exists channels
(
    id          bigint      not null,
    guild_id    bigint      not null,
    name        varchar(25) not null,
    description varchar(1000),
    created_at  timestamp   not null,
    deleted_at  timestamp,
    constraint pk_2
        primary key (id),
    constraint fk_7
        foreign key (guild_id) references guilds
);

create table if not exists guildusers
(
    guild_id   bigint      not null,
    created_at timestamp   not null,
    quit_at    timestamp,
    user_id    bigint      not null,
    nickname   varchar(32) not null,
    uuid       uuid        not null,
    constraint pk_4
        primary key (uuid),
    constraint fk_3
        foreign key (guild_id) references guilds,
    constraint fk_4
        foreign key (user_id) references users
);

create table if not exists kvstorages
(
    user_id bigint not null,
    key     text   not null,
    value   text   not null,
    constraint pk_5
        primary key (user_id),
    constraint fk_6
        foreign key (user_id) references users
);

create table if not exists messages
(
    message_id   bigint    not null,
    type         smallint  not null,
    guild_id     bigint,
    channel_id   bigint,
    user_id      bigint    not null,
    contents     text,
    reference_id bigint,
    created_at   timestamp not null,
    constraint pk_6
        primary key (message_id),
    constraint fk_4
        foreign key (guild_id) references guilds,
    constraint fk_5
        foreign key (user_id) references users,
    constraint fk_10
        foreign key (channel_id) references channels
);

create table if not exists attechments
(
    id         bigint not null,
    message_id bigint not null,
    content    text   not null,
    constraint pk_1
        primary key (id),
    constraint fk_9
        foreign key (message_id) references messages
);

create table if not exists moduleacl
(
    guild_id        bigint not null,
    allowed_modules text[] not null,
    constraint pk_8
        primary key (guild_id),
    constraint fk_10_1
        foreign key (guild_id) references guilds
);

-- comments
comment on column guildusers.user_id is 'ID for User';