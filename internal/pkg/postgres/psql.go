package postgres

import (
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Host     string
	Database string
	Port     uint32
	User     string
	Password string

	MaxConnections uint64
	AfterTimeout   uint64
}

type Psql struct {
	Primary   *pgx.ConnPool
	Migration *pgx.ConnPool
}

func NewPsql() *Psql {

	primarySql := NewPrimarySql()
	migrationSql := NewMigrationSql()

	return &Psql{
		Primary:   primarySql,
		Migration: migrationSql,
	}

}

func NewPrimarySql() *pgx.ConnPool {
	username := os.Getenv("POSTGRES_USER")
	if username == "" {
		username = "postgres"
	}
	password := os.Getenv("POSTGRES_PASSWORD")
	if password == "" {
		password = "z20339"
	}

	host := os.Getenv("POSTGRES_HOST")

	if host == "" {
		host = "127.0.0.1"
	}
	port := os.Getenv("POSTGRES_PORT")
	if port == "" {
		port = "5432"
	}

	db := os.Getenv("POSTGRES_DB")
	if db == "" {
		db = "pickreward_v2"
	}

	portInt, err := strconv.ParseUint(port, 10, 64)

	logrus.Info("[psql.host]", host)
	logrus.Info("[psql.db]", db)

	pgxConfig := pgx.ConnConfig{
		Host:     host, //host.docker.internal
		Database: db,
		Port:     uint16(portInt),
		User:     username,
		Password: password,
	}

	conn, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig:     pgxConfig,
		MaxConnections: 5,
		AfterConnect:   nil,
		AcquireTimeout: time.Duration(30) * time.Second,
	})

	if err != nil {
		logrus.Fatal(err)
		panic(-1)
	}

	return conn
}

func NewMigrationSql() *pgx.ConnPool {
	username := os.Getenv("POSTGRES_MIGRATION_USER")
	if username == "" {
		username = "postgres"
	}
	password := os.Getenv("POSTGRES_MIGRATION_PASSWORD")
	if password == "" {
		password = "z20339"
	}

	host := os.Getenv("POSTGRES_MIGRATION_HOST")

	if host == "" {
		host = "127.0.0.1"
	}
	port := os.Getenv("POSTGRES_MIGRATION_PORT")
	if port == "" {
		port = "5432"
	}

	db := os.Getenv("POSTGRES_MIGRATION_DB")
	if db == "" {
		db = "pickreward"
	}

	portInt, err := strconv.ParseUint(port, 10, 64)

	logrus.Info("[psql] host:", host)

	pgxConfig := pgx.ConnConfig{
		Host:     host, //host.docker.internal
		Database: db,
		Port:     uint16(portInt),
		User:     username,
		Password: password,
	}

	conn, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig:     pgxConfig,
		MaxConnections: 5,
		AfterConnect:   nil,
		AcquireTimeout: time.Duration(30) * time.Second,
	})

	if err != nil {
		logrus.Fatal(err)
		panic(-1)
	}

	return conn
}
