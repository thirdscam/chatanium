CREATE TABLE attechments (
    id BIGINT PRIMARY KEY,
    message_id BIGINT NOT NULL,
    content TEXT NOT NULL,
    FOREIGN KEY (message_id) REFERENCES messages(message_id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

CREATE TABLE channels (
    id BIGINT PRIMARY KEY,
    guild_id BIGINT NOT NULL,
    name VARCHAR(25) NOT NULL,
    description VARCHAR(1000),
    created_at TIMESTAMP(6) NOT NULL,
    deleted_at TIMESTAMP(6),
    FOREIGN KEY (guild_id) REFERENCES guilds(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

CREATE TABLE guilds (
    id BIGINT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    owner_id BIGINT NOT NULL,
    FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

CREATE TABLE guildusers (
    uuid UUID PRIMARY KEY,
    guild_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    created_at TIMESTAMP(6) NOT NULL,
    quit_at TIMESTAMP(6),
    nickname VARCHAR(32) NOT NULL,
    FOREIGN KEY (guild_id) REFERENCES guilds(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

CREATE TABLE kvstorages (
    user_id BIGINT PRIMARY KEY,
    key TEXT NOT NULL,
    value TEXT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

CREATE TABLE messages (
    message_id BIGINT PRIMARY KEY,
    type SMALLINT NOT NULL,
    guild_id BIGINT,
    channel_id BIGINT,
    user_id BIGINT NOT NULL,
    contents TEXT,
    reference_id BIGINT,
    created_at TIMESTAMP(6) NOT NULL,
    FOREIGN KEY (guild_id) REFERENCES guilds(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    FOREIGN KEY (channel_id) REFERENCES channels(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

CREATE TABLE moduleacl (
    guild_id BIGINT PRIMARY KEY,
    allowed_modules TEXT[] NOT NULL,
    FOREIGN KEY (guild_id) REFERENCES guilds(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

CREATE TABLE users (
    id BIGINT PRIMARY KEY,
    username CHAR(32) NOT NULL,
    created_at TIMESTAMP(6) NOT NULL,
    deleted_at TIMESTAMP(6)
);
