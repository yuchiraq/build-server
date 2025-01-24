create table if not exists albums
(
    id             bigint unsigned auto_increment
        primary key,
    title          text                          not null,
    author_id      bigint unsigned               not null,
    quantity       smallint unsigned             not null,
    music_style    smallint unsigned default '0' not null,
    tags           text                          not null,
    price          int               default 0   not null,
    price_currency tinytext                      not null,
    duration       smallint unsigned             not null,
    realise_date   date                          not null,
    sale_type      tinyint unsigned              not null,
    for_sale       tinyint(1)        default 0   not null,
    preview        bool                          not null
)
    collate = utf8mb4_0900_ai_ci;