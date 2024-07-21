CREATE TABLE messageTable
(
    id serial not null unique,
    text varchar(255) not null,
    time timestamp not null
);