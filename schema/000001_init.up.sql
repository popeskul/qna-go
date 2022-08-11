CREATE TYPE users_type AS ENUM ('user', 'admin');

CREATE TABLE users
(
    id serial not null unique,
    email varchar(255) NOT NULL unique,
    name varchar(255) NOT NULL DEFAULT '',
    type users_type NOT NULL DEFAULT 'user',
    encrypted_password varchar(255) NOT NULL DEFAULT '',
    created_at timestamp NOT NULL DEFAULT (now()),
    updated_at timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE questions
(
    id serial not null unique,
    body text NOT NULL,
    test_id bigint NOT NULL unique,
    author_id bigint NOT NULL unique,
    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp NOT NULL DEFAULT now()
);

CREATE TABLE answers
(
    id serial not null unique,
    title varchar(255) NOT NULL,
    correct boolean NOT NULL DEFAULT false,
    author_id bigint NOT NULL unique,
    question_id bigint NOT NULL unique,
    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp NOT NULL DEFAULT now()
);

CREATE TABLE test_passages
(
    id serial not null unique,
    user_id bigint NOT NULL unique,
    test_id bigint NOT NULL unique,
    question_ids bigint[] NOT NULL,
    passed_question_ids bigint[] NOT NULL,
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

ALTER TABLE questions ADD CONSTRAINT answers_questions FOREIGN KEY (id) REFERENCES answers (question_id);

ALTER TABLE tests ADD CONSTRAINT questions_tests FOREIGN KEY (id) REFERENCES questions (test_id);

ALTER TABLE tests ADD CONSTRAINT test_passages_tests FOREIGN KEY (id) REFERENCES test_passages (test_id);

ALTER TABLE users ADD CONSTRAINT test_passages_users FOREIGN KEY (id) REFERENCES test_passages (user_id);

ALTER TABLE users ADD CONSTRAINT tests_users FOREIGN KEY (id) REFERENCES tests (author_id);

ALTER TABLE users ADD CONSTRAINT questions_users FOREIGN KEY (id) REFERENCES questions (author_id);

ALTER TABLE users ADD CONSTRAINT answers_users FOREIGN KEY (id) REFERENCES answers (author_id);
