CREATE TABLE yuser (
    id text primary key,
    username text not null
);

CREATE TABLE namespace (
    id text primary key,
    owner text references yuser(id),
    name text not null
);

CREATE TABLE document (
    id text primary key,
    namespace_id text not null,
    name text not null,
    b64_content text not null
);
