-- auto-generated definition
create table orders
(
    id           int auto_increment
        primary key,
    user_id      int            not null,
    customer_id  int            not null,
    total_amount decimal(10, 2) not null,
    order_date   date           not null,
    status       varchar(20)    not null
);

INSERT INTO orders (id, user_id, customer_id, total_amount, order_date, status) VALUES (1, 1, 1, 100.00, '2024-04-04', 'pending');

-- auto-generated definition
create table user
(
    id         int                                 not null
        primary key,
    username   varchar(255)                        not null,
    password   varchar(255)                        not null,
    email      varchar(255)                        not null,
    age        int                                 null,
    created_at timestamp default CURRENT_TIMESTAMP null,
    constraint email
        unique (email),
    constraint username
        unique (username)
);

INSERT INTO user (id, username, password, email, age, created_at) VALUES (1, 'user1', 'password1', 'user1@example.com', 30, '2024-03-20 09:50:33');
INSERT INTO user (id, username, password, email, age, created_at) VALUES (2, 'user2', 'password2', 'user2@example.com', 25, '2024-03-20 09:50:33');
INSERT INTO user (id, username, password, email, age, created_at) VALUES (3, 'user3', 'password3', 'user3@example.com', 35, '2024-03-20 09:50:33');
