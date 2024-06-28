package internals

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/HeartBeat1608/logmaster/internals/queries"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

const (
	DB_DRIVER string = "sqlite3"
)

var (
	ErrConnectionNotFound = errors.New("connection not found")
)

type DBConnectionManager struct {
	connections map[string]*sqlx.DB
}

func NewConnectionManager() *DBConnectionManager {
	return &DBConnectionManager{
		connections: make(map[string]*sqlx.DB),
	}
}

func (cm *DBConnectionManager) CloseAll() {
	for k, v := range cm.connections {
		log.Printf("Closing %s", k)
		v.Close()
	}
}

func (cm *DBConnectionManager) getServiceDBPath(service string) string {
	cfg := LoadConfig(service)
	return fmt.Sprintf("%s/%s.db", cfg.DataSources, service)
}

func (cm *DBConnectionManager) initService(db *sqlx.DB) error {
	ctx := context.Background()
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	stmt, err := tx.Preparex(queries.CREATE_LOG_TABLE)
	if err != nil {
		log.Println(err)
		return tx.Rollback()
	}
	res, err := stmt.Exec()
	if err != nil {
		log.Println(err)
		return tx.Rollback()
	}

	log.Printf("Last Exec Result: %v\n", res)

	return tx.Commit()
}

func (cm *DBConnectionManager) openConnection(service string) (*sqlx.DB, error) {
	serviceDBName := strings.ToLower(service)
	return sqlx.Open(DB_DRIVER, cm.getServiceDBPath(serviceDBName))
}

func (cm *DBConnectionManager) GetConnection(service string) (*sqlx.DB, error) {
	conn, ok := cm.connections[service]
	if !ok {
		conn, err := cm.openConnection(service)
		if err != nil {
			return nil, ErrConnectionNotFound
		}
		cm.connections[service] = conn
		return conn, err
	}
	return conn, nil
}

func (cm *DBConnectionManager) AddService(service string) (*sqlx.DB, error) {
	conn, err := cm.openConnection(service)
	if err != nil {
		return nil, err
	}

	if err = cm.initService(conn); err != nil {
		return nil, err
	}

	cm.connections[service] = conn
	return conn, err
}

func (cm *DBConnectionManager) GetOrAddConnection(service string) (*sqlx.DB, error) {
	servicePath := cm.getServiceDBPath(service)

	if _, err := os.Stat(servicePath); err != nil {
		file, err := os.Create(servicePath)
		if err != nil {
			return nil, err
		}
		if err = file.Truncate(0); err != nil {
			return nil, err
		}
		defer file.Close()
		return cm.AddService(service)
	}

	return cm.GetConnection(service)
}
