package types

type Config struct {
	Env       string
	Port      string
	MongoURI  string
	JWTSecret string
	PublicURL string
	APIKey    string
}
