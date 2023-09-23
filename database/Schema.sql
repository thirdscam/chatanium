CREATE TABLE "public".attechments
(
    "id"       bigint NOT NULL,
    message_id bigint NOT NULL,
    content    text NOT NULL,
    CONSTRAINT PK_1 PRIMARY KEY ( "id" ),
    CONSTRAINT FK_9 FOREIGN KEY ( message_id ) REFERENCES "public".messages ( message_id )
);

CREATE INDEX FK_1 ON "public".attechments
(
    message_id
);

CREATE TABLE "public".channels
(
    "id"            bigint NOT NULL,
    guild_id        bigint NOT NULL,
    name            varchar(25) NOT NULL,
    description     varchar(1000) NULL,
    created_at      timestamp NOT NULL,
    deleted_at      timestamp NULL,
    CONSTRAINT PK_1 PRIMARY KEY ( "id" ),
    CONSTRAINT FK_7 FOREIGN KEY ( guild_id ) REFERENCES "public".guilds ( "id" )
);

CREATE INDEX FK_1 ON "public".channels
(
    guild_id
);

CREATE TABLE "public".guilds
(
    "id"       bigint NOT NULL,
    name       varchar(100) NOT NULL,
    owner_id   bigint NOT NULL,
    CONSTRAINT PK_1 PRIMARY KEY ( "id" ),
    CONSTRAINT FK_2 FOREIGN KEY ( owner_id ) REFERENCES "public".users ( "id" )
);

CREATE INDEX FK_1 ON "public".guilds
(
    owner_id
);

CREATE TABLE "public".guildUsers
(
    uuid       uuid NOT NULL,
    guild_id   bigint NOT NULL,
    user_id    bigint NOT NULL,
    created_at timestamp NOT NULL,
    quit_at    timestamp NULL,
    CONSTRAINT PK_1 PRIMARY KEY ( uuid ),
    CONSTRAINT FK_3 FOREIGN KEY ( guild_id ) REFERENCES "public".guilds ( "id" )
);

CREATE INDEX FK_1 ON "public".guildUsers
(
    user_id
);

CREATE INDEX FK_2 ON "public".guildUsers
(
    guild_id
);

CREATE TABLE "public".kvStorages
(
    user_id     bigint NOT NULL,
    key         text NOT NULL,
    value       text NOT NULL,
    CONSTRAINT  PK_1 PRIMARY KEY ( user_id ),
    CONSTRAINT  FK_6 FOREIGN KEY ( user_id ) REFERENCES "public".users ( "id" )
);

CREATE INDEX FK_1 ON "public".kvStorages
(
    user_id
)

CREATE TABLE "public".messages
(
    message_id   bigint NOT NULL,
    type         smallint NOT NULL,
    guild_id     bigint NULL,
    channel_id   bigint NULL,
    user_id      bigint NOT NULL,
    contents     text NULL,
    reference_id bigint NULL,
    created_at   timestamp NOT NULL,
    CONSTRAINT   PK_1 PRIMARY KEY ( message_id ),
    CONSTRAINT   FK_4 FOREIGN KEY ( guild_id ) REFERENCES "public".guilds ( "id" ),
    CONSTRAINT   FK_5 FOREIGN KEY ( user_id ) REFERENCES "public".users ( "id" ),
    CONSTRAINT   FK_10 FOREIGN KEY ( channel_id ) REFERENCES "public".channels ( "id" )
);

CREATE INDEX FK_1 ON "public".messages
(
    guild_id
);

CREATE INDEX FK_2 ON "public".messages
(
    user_id
);

CREATE INDEX FK_3 ON "public".messages
(
    channel_id
);

CREATE TABLE "public".users
(
    "id"       bigint NOT NULL,
    username   char(32) NOT NULL,
    created_at timestamp NOT NULL,
    deleted_at timestamp NULL,
    CONSTRAINT PK_1 PRIMARY KEY ( "id" )
);

CREATE TABLE "public".moduleACL
(
    "id"             bigint NOT NULL,
    allowed_modules  text[] NOT NULL,
    CONSTRAINT       PK_1 PRIMARY KEY ( "id" ),
    CONSTRAINT       FK_10_1 FOREIGN KEY ( "id" ) REFERENCES "public".channels ( "id" )
);

CREATE INDEX FK_1 ON "public".moduleACL
(
    "id"
);