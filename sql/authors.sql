create table if not exists authors
(
    id             bigint unsigned auto_increment
        primary key,
    name           text                          not null,
    music_style    smallint unsigned default '0' not null,
    about_info     text                          not null
)
    collate = utf8mb4_0900_ai_ci;