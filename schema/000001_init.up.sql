CREATE TABLE users
(
    id SERIAL NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL DEFAULT '',
    confirmed BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT (now()),
    updated_at TIMESTAMP NOT NULL DEFAULT (now())
);

CREATE TABLE questions
(
    id SERIAL NOT NULL UNIQUE,
    body text NOT NULL,
    test_id BIGINT NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE answers
(
    id SERIAL NOT NULL UNIQUE,
    title VARCHAR(255) NOT NULL,
    correct BOOLEAN NOT NULL DEFAULT false,
    question_id BIGINT NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE test_passages
(
    id SERIAL NOT NULL UNIQUE,
    user_id BIGINT NOT NULL UNIQUE,
    test_id BIGINT NOT NULL UNIQUE,
    passed BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE tests
(
    id SERIAL NOT NULL UNIQUE,
    title VARCHAR(255) NOT NULL,
    author_id BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT (now()),
    updated_at TIMESTAMP NOT NULL DEFAULT (now())
);

CREATE TABLE refresh_tokens
(
    id SERIAL NOT NULL UNIQUE,
    token VARCHAR(255) NOT NULL UNIQUE,
    user_id BIGINT NOT NULL,
    expires_at TIMESTAMP NOT NULL NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT (now()),
    updated_at TIMESTAMP NOT NULL DEFAULT (now())
);
