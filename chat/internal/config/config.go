package config

// Config contains the Application configuration.
// Currently, the Port is the only item as it does not have db repo implemented.
type Config struct {
	Port int
}
