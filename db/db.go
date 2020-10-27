package db

import (
	"database/sql"
	"fmt"
	"github.com/hofmann-works/cloudstatus/models"
	"github.com/lib/pq"
	"log"
)

type Database struct {
	Conn *sql.DB
}

func Init(host string, port int, database string, username string, password string) (Database, error) {
	db := Database{}
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, username, password, database)
	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		return db, err
	}
	db.Conn = conn
	err = db.Conn.Ping()
	if err != nil {
		return db, err
	}
	log.Println("Database connection established")
	createTables(db)

	return db, nil
}

func createTables(db Database) {

	_, err := db.Conn.Exec("CREATE TABLE IF NOT EXISTS checks (id SERIAL PRIMARY KEY, Cloud TEXT NOT NULL, lastUpdated timestamp NOT NULL UNIQUE )")
	if err != nil {
		panic(err)
	}
	_, err = db.Conn.Exec("CREATE TABLE IF NOT EXISTS services (check_id INTEGER REFERENCES checks (id), name TEXT NOT NULL)")
	if err != nil {
		panic(err)
	}

	_, err = db.Conn.Exec("CREATE OR REPLACE VIEW latestchecks AS SELECT t.id, t.lastUpdated, t.cloud, s.servicename_array FROM (SELECT DISTINCT on (cloud) * FROM checks ORDER BY cloud, lastUpdated DESC) t LEFT JOIN (SELECT services.check_id AS id,array_agg(services.name) as servicename_array FROM services GROUP BY services.check_id) s USING (id)")
	if err != nil {
		panic(err)
	}

}

func (db Database) AddCheck(check *models.Check) error {
	err := db.Conn.Ping()
	if err != nil {
		log.Printf("Could not reach database: %s\n", err)
	}
	var id int64
	query := `INSERT INTO checks (cloud, lastUpdated) VALUES ($1, $2) ON CONFLICT DO NOTHING RETURNING id`
	err = db.Conn.QueryRow(query, check.Cloud, check.LastUpdated).Scan(&id)
	if err != nil {
		return err
	}

	check.ID = id
	return nil
}

func (db Database) AddService(service *models.Service) error {
	err := db.Conn.Ping()
	if err != nil {
		log.Printf("Could not reach database: %s\n", err)
	}

	query := `INSERT INTO services (check_id, name) VALUES ($1, $2) ON CONFLICT DO NOTHING `
	_, err = db.Conn.Exec(query, service.Check_id, service.Name)
	if err != nil {
		panic(err)
	}
	return nil
}

func (db Database) GetLatestChecks() (models.StatusResponse, error) {
	err := db.Conn.Ping()
	if err != nil {
		log.Printf("Could not reach database: %s\n", err)
	}

	response := models.StatusResponse{}

	rows, err := db.Conn.Query("SELECT cloud,lastupdated,servicename_array FROM latestchecks")
	if err != nil {
		return response, err
	}

	for rows.Next() {
		cloud := models.Cloud{}
		err := rows.Scan(&cloud.Name, &cloud.LastUpdated, pq.Array(&cloud.UnhealthyServices))
		if err != nil {
			return response, err
		}
		response.Clouds = append(response.Clouds, cloud)
	}

	return response, nil
}
