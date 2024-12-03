CREATE TABLE users
(
    id            serial       not null unique,
    login         varchar(255) not null unique,
    password_hash varchar(255) not null
);

CREATE TABLE metadata
(
    id      serial       not null unique,
    name    varchar(255) not null,
    file    boolean      not null,
    public  boolean      not null,
    mime    varchar(255) not null,
    created TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE files
(
    id          serial                                         not null unique,
    metadata_id int references metadata (id) on delete cascade not null unique,
    file_data   BYTEA
);

CREATE TABLE json_document
(
    id          serial                                         not null unique,
    metadata_id int references metadata (id) on delete cascade not null unique,
    json_data   JSONB
);

CREATE TABLE users_metadata
(
    id          serial                                         not null unique,
    user_id     int references users (id) on delete cascade    not null,
    metadata_id int references metadata (id) on delete cascade not null
);

CREATE INDEX idx_metadata_name ON metadata (name);
CREATE INDEX idx_users_login ON users (login);
CREATE INDEX idx_users_metadata_metadata_id ON users_metadata (metadata_id);
CREATE INDEX idx_users_metadata_user_id ON users_metadata (user_id);
CREATE INDEX idx_json_document_metadata_id ON json_document (metadata_id);
CREATE INDEX idx_files_metadata_id ON files (metadata_id);