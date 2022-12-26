create table users
(
    id bigserial primary key,
    name text not null,
    login text unique not null,
    password text not null,
    active bool not null default true,
    created_at timestamptz not null default current_timestamp,
    updated_at timestamptz,
    deleted_at timestamptz
);

create table tokens
(
    id bigserial primary key,
    user_id bigint not null
        references users on delete cascade,
    token text not null
);

create table types
(
    id bigserial primary key,
    name text not null,
    created_at timestamptz not null default current_timestamp,
    updated_at timestamptz,
    deleted_at timestamptz
);

create table accounts
(
    id bigserial primary key,
    name text unique not null,
    user_id bigint not null
        references users on delete cascade,
    balance decimal not null,
    active bool not null default true,
    created_at timestamptz not null default current_timestamp,
    updated_at timestamptz,
    deleted_at timestamptz
);

create table categories
(
    id bigserial primary key,
    type_id bigint not null
        references types on delete restrict,
    name text not null,
    created_at timestamptz not null default current_timestamp,
    updated_at timestamptz,
    deleted_at timestamptz
);

create table operations
(
    id bigserial primary key,
    category_id bigint
        references categories on delete restrict,
    account_id bigint
        references accounts on delete restrict,
    amount decimal not null,
    created_at timestamptz not null default current_timestamp,
    updated_at timestamptz,
    deleted_at timestamptz
);