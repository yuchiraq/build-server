create table if not exists tracks
(
    id             bigint unsigned auto_increment
        primary key,
    title          text                          not null,
    album_id       bigint unsigned   default '0' not null,
    author_id      bigint unsigned               not null,
    music_style    smallint unsigned default '0' not null,
    tags           text                          not null,
    price          int               default 0   not null,
    price_currency tinytext                      not null,
    duration       smallint unsigned             not null,
    realise_date   date                          not null,
    sale_type      tinyint unsigned              not null,
    bpm            smallint unsigned             not null,
    for_sale       tinyint(1)        default 0   not null,
    kbps           smallint unsigned default '0' not null
)
    collate = utf8mb4_0900_ai_ci;