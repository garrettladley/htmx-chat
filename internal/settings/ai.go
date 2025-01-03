package settings

type AI struct {
	APIKey    string `env:"API_KEY"`
	MaxTokens uint   `env:"MAX_TOKENS" envDefault:"64"`
}
