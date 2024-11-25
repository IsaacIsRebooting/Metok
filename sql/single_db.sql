-- +goose Up
-- +goose StatementBegin
-- 创建 account 表，用于存储用户账户信息
CREATE TABLE IF NOT EXISTS account (
    id BIGINT PRIMARY KEY, -- 主键，唯一标识每个账户
    mobile VARCHAR(20) NOT NULL, -- 手机号，不能为空
    email VARCHAR(100) NOT NULL, -- 邮箱，不能为空
    password VARCHAR(64) NOT NULL, -- 密码，不能为空
    salt VARCHAR(64) NOT NULL, -- 盐值，用于密码加密，不能为空
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE, -- 是否已删除，默认为未删除
    create_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- 创建时间，默认为当前时间
    update_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, -- 更新时间，默认为当前时间，并在更新记录时自动更新
    INDEX `account_mobile_idx` (`mobile`), -- 建立手机号索引
    INDEX `account_email_idx` (`email`), -- 建立邮箱索引
    INDEX `create_time_idx` (`create_time`), -- 建立创建时间索引
    INDEX `update_time_idx` (`update_time`) -- 建立更新时间索引
);
-- +goose StatementEnd

-- +goose StatementBegin
-- 创建 user 表，用于存储用户详细信息
CREATE TABLE `user` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT, -- 主键，自增
    `account_id` bigint(20) DEFAULT NULL, -- 关联 account 表的外键
    `mobile` varchar(20) DEFAULT NULL, -- 手机号
    `email` varchar(100) DEFAULT NULL, -- 邮箱
    `name` varchar(50) DEFAULT NULL, -- 用户名
    `avatar` varchar(255) DEFAULT NULL, -- 头像 URL
    `background_image` varchar(255) DEFAULT NULL, -- 背景图 URL
    `signature` varchar(255) DEFAULT NULL, -- 个人签名
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- 创建时间，默认为当前时间
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, -- 更新时间，默认为当前时间，并在更新记录时自动更新
    PRIMARY KEY (`id`), -- 设置主键
    UNIQUE KEY `idx_users_account_id` (`account_id`) USING BTREE, -- 建立 account_id 唯一索引
    UNIQUE KEY `idx_users_email` (`email`) USING BTREE, -- 建立邮箱唯一索引
    UNIQUE KEY `idx_users_mobile` (`mobile`) USING BTREE -- 建立手机号唯一索引
) ENGINE=InnoDB AUTO_INCREMENT=7240204103809514232 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
-- +goose StatementEnd

-- +goose StatementBegin
-- 创建 template 表，用于存储模板信息
CREATE TABLE IF NOT EXISTS template (
    id BIGINT PRIMARY KEY, -- 主键，唯一标识每个模板
    title VARCHAR(255) NOT NULL, -- 模板标题，不能为空
    content TEXT NOT NULL, -- 模板内容，不能为空
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE, -- 是否已删除，默认为未删除
    create_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- 创建时间，默认为当前时间
    update_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, -- 更新时间，默认为当前时间，并在更新记录时自动更新
    INDEX `create_time_idx` (`create_time`), -- 建立创建时间索引
    INDEX `update_time_idx` (`update_time`), -- 建立更新时间索引
    INDEX `title_idx` (`title`) -- 建立标题索引
);
-- +goose StatementEnd

-- +goose StatementBegin
-- 创建 file 表，用于存储文件信息
CREATE TABLE IF NOT EXISTS file (
    id BIGINT PRIMARY KEY, -- 主键，唯一标识每个文件
    domain_name VARCHAR(100) NOT NULL, -- 文件所属域名称，不能为空
    biz_name VARCHAR(100) NOT NULL, -- 业务名称，不能为空
    hash VARCHAR(255) NOT NULL, -- 文件哈希值，不能为空
    file_size BIGINT NOT NULL DEFAULT 0, -- 文件大小，默认为 0
    file_type VARCHAR(255) NOT NULL, -- 文件类型，不能为空
    uploaded BOOLEAN NOT NULL DEFAULT FALSE, -- 是否已上传，默认为未上传
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE, -- 是否已删除，默认为未删除
    create_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- 创建时间，默认为当前时间
    update_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, -- 更新时间，默认为当前时间，并在更新记录时自动更新
    INDEX `create_time_idx` (`create_time`), -- 建立创建时间索引
    INDEX `update_time_idx` (`update_time`), -- 建立更新时间索引
    INDEX `hash_idx` (`hash`) -- 建立哈希值索引
);
-- +goose StatementEnd

-- +goose StatementBegin
-- 创建 video 表，用于存储视频信息
CREATE TABLE `video` (
     `id` bigint(20) NOT NULL AUTO_INCREMENT, -- 主键，自增
     `user_id` bigint(20) DEFAULT NULL, -- 关联 user 表的外键
     `title` varchar(255) DEFAULT NULL, -- 视频标题
     `description` varchar(255) DEFAULT NULL, -- 视频描述
     `video_url` varchar(255) DEFAULT NULL, -- 视频 URL
     `cover_url` varchar(255) DEFAULT NULL, -- 封面图 URL
     `like_count` bigint(20) DEFAULT 0, -- 点赞数，默认为 0
     `comment_count` bigint(20) DEFAULT 0, -- 评论数，默认为 0
     `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- 创建时间，默认为当前时间
     `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, -- 更新时间，默认为当前时间，并在更新记录时自动更新
     PRIMARY KEY (`id`) -- 设置主键
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- 删除 account 表
DROP TABLE IF EXISTS account;
-- +goose StatementEnd

-- +goose StatementBegin
-- 删除 user 表
DROP TABLE IF EXISTS `user`;
-- +goose StatementEnd

-- +goose StatementBegin
-- 删除 video 表
DROP TABLE IF EXISTS video;
-- +goose StatementEnd

-- +goose StatementBegin
-- 删除 template 表
DROP TABLE IF EXISTS template;
-- +goose StatementEnd

-- +goose StatementBegin
-- 删除 file 表
DROP TABLE IF EXISTS file;
-- +goose StatementEnd