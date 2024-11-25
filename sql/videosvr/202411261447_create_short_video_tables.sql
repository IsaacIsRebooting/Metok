-- Active: 1732610378770@@127.0.0.1@3308
-- +goose Up
CREATE DATABASE IF NOT EXISTS short_video_db;
USE short_video_db;

CREATE TABLE IF NOT EXISTS template (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
    create_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX `create_time_idx` (`create_time`),
    INDEX `update_time_idx` (`update_time`),
    INDEX `title_idx` (`title`)
);

CREATE TABLE IF NOT EXISTS ile (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    domain_name VARCHAR(100) NOT NULL,
    biz_name VARCHAR(100) NOT NULL,
    hash VARCHAR(255) NOT NULL UNIQUE,
    file_size BIGINT NOT NULL DEFAULT 0,
    file_type VARCHAR(255) NOT NULL,
    uploaded BOOLEAN NOT NULL DEFAULT FALSE,
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
    create_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX `create_time_idx` (`create_time`),
    INDEX `update_time_idx` (`update_time`),
    INDEX `hash_idx` (`hash`)
);

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

-- +goose Down
DROP TABLE IF EXISTS video;
DROP TABLE IF EXISTS file;
DROP TABLE IF EXISTS template;
DROP DATABASE IF EXISTS short_video_db;