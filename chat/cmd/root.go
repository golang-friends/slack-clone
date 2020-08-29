package cmd

import (
	"github.com/golang-friends/slack-clone/chat/internal/application"
	"github.com/golang-friends/slack-clone/chat/internal/config"
	"github.com/golang-friends/slack-clone/chat/internal/repo/inmemoryrepo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "chat",
	Short: "chat is the entry point for chatservice",
	RunE: func(cmd *cobra.Command, args []string) error {
		config := readConfigFromViper()
		app := application.NewApplication(config, inmemoryrepo.NewInmemoryRepo())
		return app.Start()
	},
}

func readConfigFromViper() *config.Config {
	cfg := &config.Config{
		Port: viper.GetInt("port"),
	}
	return cfg
}

func init() {
	viper.SetConfigFile("config")
	viper.SetDefault("port", 9001)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
