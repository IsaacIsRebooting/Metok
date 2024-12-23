-- Active: 1732610357622@@127.0.0.1@3307
-- +goose Up
CREATE DATABASE IF NOT EXISTS base_db;
USE base_db;

CREATE TABLE IF NOT EXISTS account (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    mobile VARCHAR(20) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(64) NOT NULL,
    salt VARCHAR(64) NOT NULL,
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
    create_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX `account_mobile_idx` (`mobile`),
    INDEX `account_email_idx` (`email`)
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
    FOREIGN KEY (account_id) REFERENCES account(id) ON DELETE CASCADE,
    UNIQUE KEY `idx_users_account_id` (`account_id`),
    UNIQUE KEY `idx_users_email` (`email`),
    UNIQUE KEY `idx_users_mobile` (`mobile`)
);

CREATE TABLE IF NOT EXISTS config (
    key_name VARCHAR(100) PRIMARY KEY,
    value TEXT NOT NULL,
    description VARCHAR(255),
    create_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS config;
DROP TABLE IF EXISTS user;
DROP TABLE IF EXISTS account;
DROP DATABASE IF EXISTS base_db;