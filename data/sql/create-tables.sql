CREATE TABLE users
(
    id                      UUID        PRIMARY KEY     DEFAULT uuid_generate_v1mc(),
    email                   TEXT        NOT NULL        UNIQUE,
    hashed_password         TEXT        NOT NULL,
    is_admin                BOOL        NOT NULL        DEFAULT false,
    date_created            TIMESTAMP   NOT NULL
);

CREATE TABLE twitter_accounts
(
    id                      UUID        PRIMARY KEY     DEFAULT uuid_generate_v1mc(),
    user_id                 UUID        NOT NULL,
    username                TEXT        NOT NULL        UNIQUE,
    date_created            TIMESTAMP   NOT NULL,
    consumer_key            TEXT        NOT NULL,
    consumer_secret         TEXT        NOT NULL,
    access_token            TEXT        NOT NULL,
    access_token_secret     TEXT        NOT NULL,

    FOREIGN KEY (user_id)
    REFERENCES users(id)
        ON DELETE CASCADE
        ON UPDATE NO ACTION
);

CREATE TABLE tweets
(
    id                      UUID        PRIMARY KEY     DEFAULT uuid_generate_v1mc(),
    twitter_account_id      UUID        NOT NULL,
    tweet                   TEXT        NOT NULL,
    post_on                 TIMESTAMP   NOT NULL,
    is_posted               BOOL        NOT NULL        DEFAULT false,
    date_created            TIMESTAMP   NOT NULL,

    FOREIGN KEY (twitter_account_id)
    REFERENCES twitter_accounts(id)
        ON DELETE CASCADE
        ON UPDATE NO ACTION
);
