package main

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func openDB() *sql.DB {
	setup, err := sql.Open("mysql", "root:@tcp(localhost:3306)/?charset=utf8mb4")
	if err != nil {
		log.Fatal("mysql open:", err)
	}
	if _, err := setup.Exec("CREATE DATABASE IF NOT EXISTS stockya CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci"); err != nil {
		log.Fatal("create database:", err)
	}
	setup.Close()

	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/stockya?parseTime=true&charset=utf8mb4")
	if err != nil {
		log.Fatal("mysql open:", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal("mysql ping — ¿está corriendo MySQL?:", err)
	}
	applySchema(db)
	return db
}

func applySchema(db *sql.DB) {
	for _, stmt := range schemaStmts {
		if _, err := db.Exec(stmt); err != nil {
			log.Fatalf("schema: %v", err)
		}
	}
}

var schemaStmts = []string{
	`CREATE TABLE IF NOT EXISTS providers (
		id          VARCHAR(50)  PRIMARY KEY,
		name        VARCHAR(255) NOT NULL,
		type        VARCHAR(20)  NOT NULL,
		src         VARCHAR(255),
		text_label  VARCHAR(50),
		badge       VARCHAR(20)
	) CHARACTER SET utf8mb4`,
	`CREATE TABLE IF NOT EXISTS products (
		id    VARCHAR(50)  PRIMARY KEY,
		name  VARCHAR(255) NOT NULL,
		price BIGINT       NOT NULL,
		emoji VARCHAR(50),
		image VARCHAR(255)
	) CHARACTER SET utf8mb4`,
	`CREATE TABLE IF NOT EXISTS product_providers (
		product_id  VARCHAR(50) NOT NULL,
		provider_id VARCHAR(50) NOT NULL,
		PRIMARY KEY (product_id, provider_id)
	)`,
	`CREATE TABLE IF NOT EXISTS branches (
		id      VARCHAR(50)  PRIMARY KEY,
		name    VARCHAR(255) NOT NULL,
		address VARCHAR(255),
		emoji   VARCHAR(50)
	) CHARACTER SET utf8mb4`,
	`CREATE TABLE IF NOT EXISTS businesses (
		id                VARCHAR(32)  PRIMARY KEY,
		name              VARCHAR(255) NOT NULL,
		cuit              VARCHAR(20)  NOT NULL,
		phone             VARCHAR(50),
		address           VARCHAR(255),
		type              VARCHAR(30),
		credit_limit      BIGINT       NOT NULL DEFAULT 50000,
		credit_used       BIGINT       NOT NULL DEFAULT 0,
		assessment_status VARCHAR(30)  NOT NULL DEFAULT 'pending',
		created_at        DATETIME     NOT NULL
	) CHARACTER SET utf8mb4`,
	`CREATE TABLE IF NOT EXISTS advances (
		id              VARCHAR(32)  PRIMARY KEY,
		business_id     VARCHAR(32)  NOT NULL,
		total           BIGINT       NOT NULL,
		logistics       VARCHAR(20)  NOT NULL,
		storage_cost    BIGINT       NOT NULL DEFAULT 0,
		advance_date    VARCHAR(10),
		due_date        VARCHAR(10),
		status          VARCHAR(20)  NOT NULL DEFAULT 'ok',
		delivery_status VARCHAR(30)  NOT NULL DEFAULT '',
		target_branch   VARCHAR(255) NOT NULL DEFAULT '',
		created_at      DATETIME     NOT NULL
	)`,
	`CREATE TABLE IF NOT EXISTS advance_items (
		id         INT AUTO_INCREMENT PRIMARY KEY,
		advance_id VARCHAR(32)  NOT NULL,
		product_id VARCHAR(50)  NOT NULL,
		name       VARCHAR(255) NOT NULL,
		emoji      VARCHAR(50),
		image      VARCHAR(255),
		qty        INT          NOT NULL,
		price      BIGINT       NOT NULL
	) CHARACTER SET utf8mb4`,
}
