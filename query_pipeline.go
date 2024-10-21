package majordomo_ai

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type QueryType int

const (
	QueryTypeText QueryType = iota + 1
	QueryTypeTextAndImage
	QueryTypeSQL
	QueryTypeChainOfThought
)

type QueryPipeline struct {
	ID             primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Workspace      string             `json:"workspace" bson:"workspace"`
	UserName       string             `json:"user_name" bson:"user_name"`
	Name           string             `json:"name" bson:"name"`
	Type           QueryType          `json:"query_type" bson:"query_type"`
	EmbeddingModel string             `json:"embedding_model" bson:"embedding_model"`
	LLMModel       string             `json:"llm_model" bson:"llm_model"`
	DataStores     []string           `json:"data_stores" bson:"data_stores"`
	QueryParams    string             `json:"query_params" bson:"query_params"`
	CreatedAt      int64              `json:"created_at" bson:"created_at"`
	LastUpdate     int64              `json:"last_update" bson:"last_update"`
}

type QueryPipelineList struct {
	QueryPipelines []QueryPipeline `json:"query_pipelines"`
}

type QueryParams struct {
	TopK          int               `json:"top_k"`
	Temperature   float64           `json:"temperature"`
	Metadata      map[string]string `json:"metadata"`
	MetadataMatch string            `json:"metadata_match"`
}

type QueryPipelineMessage struct {
	Credentials    Credentials `json:"credentials"`
	Name           string      `json:"name"`
	Type           QueryType   `json:"query_type"`
	EmbeddingModel string      `json:"embedding_model" bson:"embedding_model"`
	LLMModel       string      `json:"llm_model"`
	DataStores     []string    `json:"data_stores"`
	QueryParams    string      `json:"query_params"`
}

type QueryPipelineQuickRunMessage struct {
	Credentials    Credentials `json:"credentials"`
	EmbeddingModel string      `json:"embedding_model" bson:"embedding_model"`
	LLMModel       string      `json:"llm_model"`
	DataStore      string      `json:"data_store"`
	Type           QueryType   `json:"query_type"`
	QueryParams    string      `json:"query_params"`
	QueryStr       string      `json:"query_str"`
}

type QueryPipelineRunMessage struct {
	Credentials Credentials `json:"credentials"`
	Name        string      `json:"name"`
	QueryStr    string      `json:"query_str"`
}

func QueryPipelineRun(name string, queryStr string) (*http.Response, error) {

	var v QueryPipelineRunMessage

	err := setCredentialsFromEnv(&v.Credentials)
	if err != nil {
		log.Fatal(err)
	}
	v.Name = name
	v.QueryStr = queryStr
	jsonData, err := json.Marshal(v)
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("POST",
		os.Getenv("MAJORDOMO_AI_DIRECTOR")+"/query_pipeline_run",
		bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	return (client.Do(req))
}

func CreateOrUpdateQueryPipeline(
	create bool,
	name string,
	queryType QueryType,
	embeddingModel string,
	LLMModel string,
	dataStores []string,
	queryParams string) (*http.Response, error) {

	var v QueryPipelineMessage

	v.Name = name
	err := setCredentialsFromEnv(&v.Credentials)
	if err != nil {
		log.Fatal(err)
	}
	v.Type = queryType
	v.EmbeddingModel = embeddingModel
	v.LLMModel = LLMModel
	v.QueryParams = queryParams
	v.DataStores = dataStores

	jsonData, err := json.Marshal(v)
	if err != nil {
		log.Fatal(err)
	}

	op := "POST"
	if !create {
		op = "PUT"
	}
	req, err := http.NewRequest(op,
		os.Getenv("MAJORDOMO_AI_DIRECTOR")+"/query_pipeline",
		bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	return (client.Do(req))
}

func DeleteQueryPipeline(name string) (*http.Response, error) {

	var v QueryPipelineMessage

	v.Name = name
	err := setCredentialsFromEnv(&v.Credentials)
	if err != nil {
		log.Fatal(err)
	}

	jsonData, err := json.Marshal(v)
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("DELETE",
		os.Getenv("MAJORDOMO_AI_DIRECTOR")+"/query_pipeline",
		bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	return (client.Do(req))
}

func DataStoreQuery(queryType QueryType,
	embeddingModel string,
	LLMModel string,
	dataStore string,
	queryParams string,
	queryStr string) (*http.Response, error) {

	var v QueryPipelineQuickRunMessage

	err := setCredentialsFromEnv(&v.Credentials)
	if err != nil {
		log.Fatal(err)
	}
	v.DataStore = dataStore
	v.Type = queryType
	v.LLMModel = LLMModel
	v.QueryParams = queryParams
	v.QueryStr = queryStr

	jsonData, err := json.Marshal(v)
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("POST",
		os.Getenv("MAJORDOMO_AI_DIRECTOR")+"/data_store_query",
		bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	return (client.Do(req))
}
