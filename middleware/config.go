package middleware

// Configuration holds the middleware settings
type Config struct {
	SecretKey []byte
}

// DefaultConfig returns a default configuration
func DefaultConfig() Config {
	return Config{
		SecretKey: []byte("Shoxrux1801$"),
	}
}

var (
	defaultConfig = DefaultConfig()
)
