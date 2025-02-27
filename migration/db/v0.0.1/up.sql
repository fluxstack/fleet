CREATE TABLE `users`
(
    `id`              bigint                           NOT NULL COMMENT '主键，用户 UID',
    `name`            varchar(256) COLLATE utf8mb4_bin NOT NULL COMMENT '用户名',
    `display_name`    varchar(256) COLLATE utf8mb4_bin NOT NULL COMMENT '展示的用户名',
    `avatar_url`      varchar(256) COLLATE utf8mb4_bin          DEFAULT NULL COMMENT '头像',
    `phone`           varchar(32) COLLATE utf8mb4_bin           DEFAULT NULL,
    `email`           varchar(256) COLLATE utf8mb4_bin          DEFAULT NULL,
    `password_hash`   varchar(256) COLLATE utf8mb4_bin NOT NULL,
    `created_at`      bigint                           NOT NULL DEFAULT '0',
    `created_by`      bigint                           NOT NULL DEFAULT '0',
    `last_updated_at` bigint                           NOT NULL DEFAULT '0',
    `last_updated_by` bigint                           NOT NULL DEFAULT '0',
    `role`            tinyint                          NOT NULL DEFAULT '0',
    `is_deleted`      tinyint                          NOT NULL DEFAULT '0',
    `is_locked`       tinyint                          NOT NULL DEFAULT '0',
    `union_id`        varchar(256) COLLATE utf8mb4_bin          DEFAULT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='用户表';