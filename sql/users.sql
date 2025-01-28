create table if not exists users
(
    id             bigint unsigned auto_increment primary key,
    login          text                          not null,
    password       text                         not null,
    companyID      bigint unsigned default '0' not null,
    firstName       text,
    secondName      text,
    lastName        text
)
    collate = utf8mb4_0900_ai_ci;