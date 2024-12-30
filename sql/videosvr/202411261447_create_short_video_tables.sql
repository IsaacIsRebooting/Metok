-- Active: 1732610378770@@127.0.0.1@3308
-- +goose Up
CREATE DATABASE IF NOT EXISTS short_video_db;
USE short_video_db;


CREATE TABLE IF NOT EXISTS video (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL,
    title VARCHAR(255),
    description VARCHAR(255),
    video_url VARCHAR(255),
    cover_url VARCHAR(255),
    like_count BIGINT DEFAULT 0,
    comment_count BIGINT DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    -- 注意：这里不使用外键引用，而是通过应用层逻辑处理关联
    INDEX `created_at_idx` (`created_at`),
    INDEX `updated_at_idx` (`updated_at`)
);
CREATE TABLE IF NOT EXISTS user (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    account_id BIGINT NOT NULL,
    name VARCHAR(50),
    email VARCHAR(100),
    mobile VARCHAR(20),
    avatar VARCHAR(255),
    background_image VARCHAR(255),
    signature VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY `idx_users_account_id` (`account_id`),
    UNIQUE KEY `idx_users_email` (`email`),
    UNIQUE KEY `idx_users_mobile` (`mobile`)
);

-- +goose Down
DROP TABLE IF EXISTS video;
DROP TABLE IF EXISTS file;
DROP TABLE IF EXISTS template;
DROP DATABASE IF EXISTS short_video_db;