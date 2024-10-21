package majordomo_ai

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IngestTypes int

const (
	IngestTypeText IngestTypes = iota + 1
	IngestTypeTextWithImages
	IngestTypeCustom
)

type InputTypes int

const (
	AWSS3 InputTypes = iota + 1
	AzureBlob
	Webpage
	Local
)

// Used only for validation.
type IngestParams struct {
	ChunkSize     int               `json:"chunk_size"`
	ChunkOverlap  int               `json:"chunk_overlap"`
	LlmModel      string            `json:"llm_model"`
	Summarize     bool              `json:"summarize"`
	CustomScript  string            `json:"custom_script"`
	FileExtractor map[string]string `json:"file_extractor"`
	ApiKey        string            `json:"api_key"`
	MetaData      map[string]string `json:"metadata"`
}

type IngestPipeline struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Workspace string             `json:"workspace" bson:"workspace"`
	UserName  string             `json:"user_name" bson:"user_name"`
	DataStore string             `json:"data_store" bson:"data_store"`
	Name      string             `json:"name" bson:"name"`

	InputType   InputTypes `json:"input_type"`
	InputFilter string     `json:"input_filter"`
	InputKeys   string     `json:"input_keys"`

	IngestType    IngestTypes `json:"ingest_type" bson:"ingest_type"`
	IngestParams  string      `json:"ingest_params" bson:"ingest_params"`
	TimerInterval int         `json:"timer_interval" bson:"timer_interval"`
	TimerOn       bool        `json:"timer_on" bson:"timer_on"`

	CreatedAt  int64 `json:"created_at" bson:"created_at"`
	LastUpdate int64 `json:"last_update" bson:"last_update"`
}

type IngestPipelineList struct {
	IngestPipelines []IngestPipeline `json:"ingest_pipelines"`
}

type IngestPipelineMessage struct {
	Credentials Credentials `json:"credentials"`
	DataStore   string      `json:"data_store"`
	Name        string      `json:"name"`

	InputType   InputTypes `json:"input_type"`
	InputFilter string     `json:"input_filter"`
	InputKeys   string     `json:"input_keys"`

	IngestType    IngestTypes `json:"ingest_type"`
	IngestParams  string      `json:"ingest_params"`
	TimerInterval int         `json:"timer_interval"`
	TimerOn       bool        `json:"timer_on"`
}

type IngestPipelineQuickRunMessage struct {
	Credentials Credentials `json:"credentials"`
	DataStore   string      `json:"data_store"`
	Name        string      `json:"name"`

	InputType   InputTypes `json:"input_type"`
	InputFilter string     `json:"input_filter"`
	InputKeys   string     `json:"input_keys"`

	IngestType   IngestTypes `json:"ingest_type"`
	IngestParams string      `json:"ingest_params"`
}

type IngestPipelineRunMessage struct {
	Credentials Credentials `json:"credentials"`
	DataStore   string      `json:"data_store"`
	Name        string      `json:"name"`
}

// Creates a new file upload http request with optional extra params
func NewfileUploadRequest(url string, params map[string]string, paramName, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		return err
	}
	_, err = io.Copy(part, file)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)

	if err == nil {
		body := &bytes.Buffer{}
		_, err := body.ReadFrom(resp.Body)
		if err != nil {
			return err
		}
		resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return errors.New("File upload request failed")
		}
	}
	return err
}
func IngestPipelineRun(dataStore string, name string) (*http.Response, error) {

	var v IngestPipelineRunMessage

	err := setCredentialsFromEnv(&v.Credentials)
	if err != nil {
		log.Fatal(err)
	}
	v.DataStore = dataStore
	v.Name = name
	jsonData, err := json.Marshal(v)
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("POST", os.Getenv("MAJORDOMO_AI_DIRECTOR")+"/ingest_pipeline_run", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	return (client.Do(req))
}

func CreateOrUpdateIngestPipeline(
	create bool,
	dataStore string,
	name string,
	sourceType InputTypes,
	sourceFilter string,
	ingestType IngestTypes,
	ingestParams string,
	timerInterval int,
	timerOn bool) (*http.Response, error) {

	var v IngestPipelineMessage

	err := setCredentialsFromEnv(&v.Credentials)
	if err != nil {
		log.Fatal(err)
	}
	v.Name = name
	v.InputType = sourceType
	v.InputFilter = sourceFilter
	v.DataStore = dataStore
	v.IngestType = ingestType
	v.IngestParams = ingestParams
	v.TimerInterval = timerInterval
	v.TimerOn = timerOn

	jsonData, err := json.Marshal(v)
	if err != nil {
		log.Fatal(err)
	}

	op := "POST"
	if !create {
		op = "PUT"
	}
	req, err := http.NewRequest(op, os.Getenv("MAJORDOMO_AI_DIRECTOR")+"/ingest_pipeline", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	return (client.Do(req))
}

func DeleteIngestPipeline(dataStore string, name string) (*http.Response, error) {

	var v IngestPipelineMessage

	err := setCredentialsFromEnv(&v.Credentials)
	if err != nil {
		log.Fatal(err)
	}
	v.DataStore = dataStore
	v.Name = name

	jsonData, err := json.Marshal(v)
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("DELETE", os.Getenv("MAJORDOMO_AI_DIRECTOR")+"/ingest_pipeline", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	return (client.Do(req))
}

func GetIngestPipeline(dataStore string, name string) (*http.Response, error) {

	var v IngestPipelineMessage

	err := setCredentialsFromEnv(&v.Credentials)
	if err != nil {
		log.Fatal(err)
	}
	v.DataStore = dataStore
	v.Name = name

	jsonData, err := json.Marshal(v)
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("GET",
		os.Getenv("MAJORDOMO_AI_DIRECTOR")+"/ingest_pipeline",
		bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	return (client.Do(req))
}

func UploadToDirector(mdApiKey string, url string, f string) {
	// Request the controller for the client URL for this group.
	x := UserInfoMessage{MdApiKey: mdApiKey, UserName: os.Getenv("MAJORDOMO_AI_USER")}

	jsonData, err := json.Marshal(x)
	req, err := http.NewRequest("GET", url+"/worker_info", bytes.NewBuffer(jsonData))

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		if resp.StatusCode != http.StatusOK {
			log.Fatal(errors.New(string(bodyBytes)))
		}

		var s WorkerInfo
		if err := json.Unmarshal(bodyBytes, &s); err != nil {
			log.Fatal(err)
		}

		extraParams := map[string]string{
			"md_api_key": "",
		}
		err = NewfileUploadRequest(s.ClientURL+"/file_upload", extraParams, "file", f)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func DataStoreIngest(dataStore string,
	sourceType InputTypes,
	sourceFilter string,
	ingestType IngestTypes,
	ingestParams string) (*http.Response, error) {

	var v IngestPipelineMessage

	err := setCredentialsFromEnv(&v.Credentials)
	if err != nil {
		log.Fatal(err)
	}
	v.DataStore = dataStore
	v.InputType = sourceType
	v.InputFilter = sourceFilter
	v.DataStore = dataStore
	v.IngestType = ingestType
	v.IngestParams = ingestParams

	if v.InputType == Local {
		var infoMap map[string]string
		json.Unmarshal([]byte(sourceFilter), &infoMap)

		f, ok := infoMap["files"]
		if !ok {
			return nil, errors.New("Specify local files")
		}

		extraParams := map[string]string{
			"md_api_key": "",
		}
		split := strings.Split(f, ",")
		for _, ff := range split {
			err := NewfileUploadRequest(
				os.Getenv("MAJORDOMO_AI_DIRECTOR")+"/file_upload",
				extraParams,
				"file",
				ff)
			if err != nil {
				log.Fatal(err)
			}
		}

	}

	jsonData, err := json.Marshal(v)
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("POST",
		os.Getenv("MAJORDOMO_AI_DIRECTOR")+"/data_store_ingest",
		bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	return (client.Do(req))
}
