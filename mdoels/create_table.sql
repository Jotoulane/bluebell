create table user(
    id bigint(20) not null auto_increment,
    user_id bigint(20) not null,
    username varchar(64) collate utf8mb4_general_ci not null ,
    password varchar(64) collate utf8mb4_general_ci not null ,
    email varchar(64) collate utf8mb4_general_ci,
    gender tinyint(4) not null default '0',
    create_time timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    update_time timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE KEY idx_username( username ) USING BTREE,
    UNIQUE KEY idx_user_id( user_id ) USING BTREE
)ENGINE=InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci;
