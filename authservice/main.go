package main

import (
	"log"
	"net"
	"strings"

	"github.com/golang-friends/slack-clone/authservice/cmd/authservice"
	"github.com/golang-friends/slack-clone/authservice/models"
	pb "github.com/golang-friends/slack-clone/authservice/protos/authservice"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	LISTEN_INTERFACE    = "listen_interface"
	DATABASE_DBURI      = "database.dburi"
	DATABASE_DBNAME     = "database.dbname"
	DATABASE_DBUSER     = "database.dbuser"
	DATABASE_DBPASSWORD = "database.dbpassword"
	JWT_SECRET          = "jwt.secret"
	JWT_EXPIRY_TIME     = "jwt.expiry_time"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:  "authservice",
	Long: "authentication service for slack clone",
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Println("Starting Server...")

		models.ConnectToDB(models.Database{
			DBURI:      viper.GetString(DATABASE_DBURI),
			DBName:     viper.GetString(DATABASE_DBNAME),
			DBUser:     viper.GetString(DATABASE_DBUSER),
			DBPassword: viper.GetString(DATABASE_DBPASSWORD),
		})
		log.Println("Connected to Database")

		server := grpc.NewServer()
		pb.RegisterAuthServiceServer(server, authservice.NewAuthServer(viper.GetString(JWT_SECRET), viper.GetInt(JWT_EXPIRY_TIME)))

		// Advertises gRPC services for dev env
		// Disable in production
		reflection.Register(server)

		// creates gRPC listener
		listener, err := net.Listen("tcp", viper.GetString(LISTEN_INTERFACE))
		if err != nil {
			return err
		}

		// starts gRPC server
		log.Printf("gRPC server hosted on %s", viper.GetString(LISTEN_INTERFACE))
		return server.Serve(listener)
	},
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file path")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigName("config")
	}

	viper.SetEnvPrefix("SLACK_CLONE")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err == nil {
		return
	}
	if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
		log.Fatal("issue while loading config", err)
	}
}

func main() {
	viper.SetDefault(LISTEN_INTERFACE, "0.0.0.0:9000")
	viper.SetDefault(DATABASE_DBURI, "auth-db")
	viper.SetDefault(DATABASE_DBNAME, "auth")
	viper.SetDefault(DATABASE_DBUSER, "root")
	viper.SetDefault(DATABASE_DBPASSWORD, "password")
	viper.SetDefault(JWT_SECRET, "")
	viper.SetDefault(JWT_EXPIRY_TIME, 30)
	rootCmd.Execute()
}
