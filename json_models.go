package majordomo_ai

type AWSS3Info struct {
	Region      string `json:"region"`
	AccessKey   string `json:"access_key"`
	SecretToken string `json:"secret_token"`
}

type AzureBlobInfo struct {
	ClientId     string `json:"client_id"`
	TenantId     string `json:"tenant_id"`
	ClientSecret string `json:"client_secret"`
	AccountURL   string `json:"account_url"`
	Container    string `json:"container"`
	Blob         string `json:"blob"`
}

type WebpageInfo struct {
	URL string `json:"url"`
}

type SQLDataStore struct {
	URL    string `json:"url"`
	DbName string `json:"db_name"`
	Table  string `json:"table"`
}

type IngestionOptions struct {
	ChunkingSize   int    `json:"chunking_size"`
	EmbeddingModel string `json:"embedding_model"`
	LLMModel       string `json:"llm_model"`
	Extractor      string `json:"extractor"`
}

type WorkspaceRegistration struct {
	Account   int    `json:"account_id"`
	ClientURL string `json:"client_url"`
	Workspace string `json:"workspace"`
}

type UserInfoMessage struct {
	MdApiKey string `json:"md_api_key"`
	UserName string `json:"user_name"`
}

type WorkerInfo struct {
	ClientURL string
}

type Credentials struct {
	AccountId int    `json:"account_id"`
	Workspace string `json:"workspace"`
	MdApiKey  string `json:"md_api_key"`
	ExtraTags string `json:"extra_tags"`
}

type CredentialsMessage struct {
	Credentials Credentials
}
