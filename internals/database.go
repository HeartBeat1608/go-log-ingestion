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
	cm := &DBConnectionManager{
		connections: make(map[string]*sqlx.DB),
	}
	cm.loadConnections()
	return cm
}

func (cm *DBConnectionManager) CloseAll() {
	for k, v := range cm.connections {
		log.Printf("Closing %s", k)
		v.Close()
	}
}

func (cm *DBConnectionManager) getServiceDBPath(service string) string {
	service = strings.Split(service, ".")[0]
	cfg := GetConfig()
	return fmt.Sprintf("%s/%s.db", cfg.DataSources, service)
}

func (cm *DBConnectionManager) initService(db *sqlx.DB) error {
	ctx := context.Background()
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = tx.Exec(queries.CREATE_LOG_TABLE)
	if err != nil {
		log.Println(err)
		return tx.Rollback()
	}

	_, err = tx.Exec(queries.CREATE_META_TABLE)
	if err != nil {
		log.Println(err)
		return tx.Rollback()
	}

	return tx.Commit()
}

func (cm *DBConnectionManager) openConnection(service string) (*sqlx.DB, error) {
	serviceDBName := strings.ToLower(service)
	return sqlx.Open(DB_DRIVER, cm.getServiceDBPath(serviceDBName))
}

func (cm *DBConnectionManager) loadConnections() {
	cfg := GetConfig()
	entries, err := os.ReadDir(cfg.DataSources)
	if err != nil {
		panic(err)
	}

	for _, fd := range entries {
		if fd.IsDir() {
			continue
		}

		serviceName, _ := strings.CutSuffix(fd.Name(), ".")
		_, _ = cm.GetConnection(serviceName)
	}
}

func (cm *DBConnectionManager) GetConnection(service string) (*sqlx.DB, error) {
	conn, ok := cm.connections[service]
	if !ok {
		log.Printf("Creating datastore for %s", service)
		conn, err := cm.openConnection(service)
		if err != nil {
			return nil, ErrConnectionNotFound
		}
		cm.connections[service] = conn
		err = cm.initService(conn)
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
