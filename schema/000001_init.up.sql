CREATE TABLE users
(
    id serial not null unique,
    email varchar(255) NOT NULL unique,
    name varchar(255) NOT NULL,
    encrypted_password varchar(255) NOT NULL DEFAULT '',
    created_at timestamp NOT NULL DEFAULT (now()),
    updated_at timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE questions
(
    id serial not null unique,
    body text NOT NULL,
    test_id bigint NOT NULL unique,
    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp NOT NULL DEFAULT now()
);

CREATE TABLE answers
(
    id serial not null unique,
    title varchar(255) NOT NULL,
    correct boolean NOT NULL DEFAULT false,
    question_id bigint NOT NULL unique,
    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp NOT NULL DEFAULT now()
);

CREATE TABLE test_passages
(
    id serial not null unique,
    user_id bigint NOT NULL unique,
    test_id bigint NOT NULL unique,
    passed boolean NOT NULL DEFAULT false,
    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp NOT NULL DEFAULT now()
);

CREATE TABLE tests
(
    id serial not null unique,
    title varchar(255),
    author_id bigint NOT NULL unique,
    created_at timestamp NOT NULL DEFAULT (now()),
    updated_at timestamp NOT NULL DEFAULT (now())
);
