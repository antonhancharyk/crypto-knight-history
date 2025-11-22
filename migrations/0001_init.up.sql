CREATE TABLE IF NOT EXISTS klines (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    symbol TEXT NOT NULL,
    open_time BIGINT NOT NULL,
    open_price DOUBLE PRECISION NOT NULL,
    high_price DOUBLE PRECISION NOT NULL,
    low_price DOUBLE PRECISION NOT NULL,
    close_price DOUBLE PRECISION NOT NULL,
    volume DOUBLE PRECISION NOT NULL,
    close_time BIGINT NOT NULL,
    quote_asset_volume DOUBLE PRECISION NOT NULL,
    num_trades BIGINT NOT NULL,
    taker_buy_base_asset_volume DOUBLE PRECISION NOT NULL,
    taker_buy_quote_asset_volume DOUBLE PRECISION NOT NULL
);
