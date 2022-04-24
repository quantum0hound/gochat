CREATE TABLE users
(
    id            serial       not null unique,
    username      varchar(255) not null unique,
    password_hash varchar(255) not null
    
);

CREATE TABLE channels
(
    id          serial       not null unique,
    name        varchar(255) not null unique,
    creator     int references users (id) not null,
    description varchar(255) not null
);

CREATE TABLE users_channels
(
    id          bigserial       not null unique,
    user_id     int references users (id) on delete cascade not null,
    channel_id  int references channels (id) on delete cascade not null,
    unique(user_id,channel_id)
);

CREATE TABLE messages
(
    id         bigserial    not null unique,
    content    text      not null,
    channel_id int references channels (id) on delete cascade not null,
    user_id    int references users (id) on delete cascade not null,
    posted     timestamp not null,
    modified   timestamp not null
);



