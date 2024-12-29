CREATE TABLE users
(
    id             SERIAL PRIMARY KEY,
    chat_id        INT UNIQUE NOT NULL,
    username       VARCHAR(255),
    wallet_address VARCHAR(255),
    balance        DECIMAL    NOT NULL,
    created_at     TIMESTAMP  NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE games
(
    id          SERIAL PRIMARY KEY,
    start_time  TIMESTAMP,
    crash_point DECIMAL      NOT NULL,
    status      VARCHAR(255) NOT NULL,
    created_at  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE bets
(
    id                 SERIAL PRIMARY KEY,
    user_id            INT REFERENCES users (id),
    game_id            INT REFERENCES games (id),
    bet_amount         DECIMAL      NOT NULL,
    status             VARCHAR(255) NOT NULL DEFAULT 'placed',
    cashout_multiplier DECIMAL,
    created_at         TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP
);