CREATE TABLE users
(
    id            serial       not null unique,
    login         varchar(255) not null,
    password_hash varchar(255) not null
);

CREATE TABLE metadata
(
    id        serial       not null unique,
    name      varchar(255) not null,
    file      boolean   default true,
    public    boolean   default false,
    mime      varchar(255) not null,
    created   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    json_data JSONB
);

CREATE TABLE files
(
    id          serial                                         not null unique,
    metadata_id int references metadata (id) on delete cascade not null unique,
    file_data   BYTEA                                          not null
);

CREATE TABLE users_metadata
(
    id      serial                                           not null unique,
    user_id int references users (id) on delete cascade      not null,
    metadata_id int references metadata (id) on delete cascade not null
);
