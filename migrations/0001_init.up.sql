create table if not exists klines (
    id serial primary key,
    created_at timestamp not null default current_timestamp,
    symbol text not null,
    open_time bigint not null,
    open_price double precision not null,
    high_price double precision not null,
    low_price double precision not null,
    close_price double precision not null,
    volume double precision not null,
    close_time bigint not null,
    quote_asset_volume double precision not null,
    num_trades bigint not null,
    taker_buy_base_asset_volume double precision not null,
    taker_buy_quote_asset_volume double precision not null
);
