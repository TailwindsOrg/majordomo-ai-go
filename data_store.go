package majordomo_ai

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type DataStoreTypes int

const (
	DataStoreVectorDB DataStoreTypes = iota + 1
	DataStoreSQL
	DataStoreMongoDB
)

type DataStoreMessage struct {
	Credentials Credentials    `json:"credentials"`
	Name        string         `json:"name"`
	Type        DataStoreTypes `json:"type"`

	VectorDBProfile string `json:"vectordb_profile"`
	EmbeddingModel  string `json:"embedding_model"`

	DatabaseURL   string `json:"db_url"`
	DatabaseName  string `json:"db_name"`
	DatabaseTable string `json:"db_table"`

	Shared bool `json:"shared"`
}

func CreateOrUpdateVectorDB(
	create bool,
	name string,
	vectorDBProfile string,
	embeddingModel string,
	shared bool) (*http.Response, error) {

	var v DataStoreMessage

	v.Name = name
	err := setCredentialsFromEnv(&v.Credentials)
	if err != nil {
		log.Fatal(err)
	}
	v.Type = DataStoreVectorDB
	v.VectorDBProfile = vectorDBProfile
	v.EmbeddingModel = embeddingModel
	v.DatabaseURL = ""
	v.DatabaseName = ""
	v.DatabaseTable = ""
	v.Shared = shared

	jsonData, err := json.Marshal(v)
	if err != nil {
		log.Fatal(err)
	}

	op := "POST"
	if !create {
		op = "PUT"
	}
	req, err := http.NewRequest(op, os.Getenv("MAJORDOMO_AI_DIRECTOR")+"/data_store", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	return (client.Do(req))
}

func CreateOrUpdateStructedDB(
	create bool,
	name string,
	storeType DataStoreTypes,
	embeddingModel string,
	databaseURL string,
	databaseName string,
	databaseTable string) (*http.Response, error) {

	var v DataStoreMessage

	err := setCredentialsFromEnv(&v.Credentials)
	if err != nil {
		log.Fatal(err)
	}
	v.Name = name
	v.Type = storeType
	v.VectorDBProfile = ""
	v.EmbeddingModel = embeddingModel
	v.DatabaseURL = databaseURL
	v.DatabaseName = databaseName
	v.DatabaseTable = databaseTable
	v.Shared = false

	jsonData, err := json.Marshal(v)
	if err != nil {
		log.Fatal(err)
	}

	op := "POST"
	if !create {
		op = "PUT"
	}
	req, err := http.NewRequest(op, os.Getenv("MAJORDOMO_AI_DIRECTOR")+"/data_store", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	return (client.Do(req))
}

func DeleteDataStore(name string) (*http.Response, error) {

	var v DataStoreMessage

	err := setCredentialsFromEnv(&v.Credentials)
	if err != nil {
		log.Fatal(err)
	}
	v.Name = name

	jsonData, err := json.Marshal(v)
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("DELETE", os.Getenv("MAJORDOMO_AI_DIRECTOR")+"/data_store", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	return (client.Do(req))
}
