package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"entgo.io/ent/dialect/sql"
	_ "github.com/lib/pq"

	"api_server/ent"
	"api_server/ent/migrate"
)

var DEFAULT_MAX_IDLE_CONNECTIONS = 30
var DEFAULT_MAX_OPEN_CONNECTIONS = 30

type KAISConfigModel struct {
	KAISPORT  string `json:"port"`
	DBIP      string `json:"db_ip"`
	DBPORT    int    `json:"db_port"`
	DB        string `json:"db"`
	BROKERIP  string `json:"broker_ip"`
	WORKSPACE string `json:"workspace"`
}

var kaisConfig *KAISConfigModel

func GetDBConfig() string {
	fileconfig := GetFileConfig()

	d_ip := fileconfig.DBIP
	if d_ip == "" {
		d_ip = "localhost"
	}

	d_port := fileconfig.DBPORT
	if d_port == 0 {
		d_port = 5432
	}

	d_db := fileconfig.DB
	if d_db == "" {
		d_db = "kaier_backend_poc"
	}

	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", d_ip, d_port, DATABASE_DEFAULT_USER, DATABASE_DEFAULT_PWD, d_db)
}

var entClient *ent.Client

func initEntClient() {
	driver, err := sql.Open("postgres", GetDBConfig())
	if err != nil {
		println("DBMS Connection Error : ", err.Error())
		return
	}

	db := driver.DB()
	db.SetMaxIdleConns(DEFAULT_MAX_IDLE_CONNECTIONS)
	db.SetMaxOpenConns(DEFAULT_MAX_OPEN_CONNECTIONS)
	db.SetConnMaxLifetime(time.Second * 1)

	entClient = ent.NewClient(ent.Driver(driver))
}

func GetEntClient() *ent.Client {
	if entClient == nil {
		initEntClient()
	}

	return entClient
}

func InitDBMS() {
	if entClient == nil {
		initEntClient()
	}

	if err := entClient.Schema.Create(
		context.Background(),
		migrate.WithDropIndex(true),
		migrate.WithDropColumn(true),
	); err != nil {
		log.Print(err)
	}
}

func ReadConfigFromFile() *KAISConfigModel {
	defaultconfig := &KAISConfigModel{
		KAISPORT: ":8970",
		DBIP:     "localhost",
		DBPORT:   5432,
		DB:       "kaier_backend_poc",
		BROKERIP: "http://localhost:8880/auth",
	}

	bytejson, err := os.ReadFile("./config_kais.json")
	if err != nil {
		println("Read file Error : ", err.Error())
		return defaultconfig
	}

	// remove \r\n
	d := bytes.Replace(bytejson, []byte{13, 10}, []byte{}, -1)
	// remove \r
	d = bytes.Replace(d, []byte{13}, []byte{}, -1)
	// remove \n
	d = bytes.Replace(d, []byte{10}, []byte{}, -1)
	// remove sapce
	d = bytes.Replace(d, []byte{32}, []byte{}, -1)
	err1 := json.Unmarshal(d, &kaisConfig)
	if err1 != nil {
		println(err1.Error())
	}

	return kaisConfig
}

func GetFileConfig() *KAISConfigModel {
	if kaisConfig.DBIP == "" {
		return ReadConfigFromFile()
	}
	return kaisConfig
}
