package cmd

import (
	"fin-auth/config"
	"fin-auth/routes"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/cobra"
)

var (
	port string
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the API server",
	Long:  `Start the fin-auth API server on the specified port`,
	Run: func(cmd *cobra.Command, args []string) {
		startServer()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().StringVarP(&port, "port", "p", "", "Port to run the server on")
}

func startServer() {
	if err := config.LoadConfig(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	cfg := config.GetConfig()
	if port == "" {
		port = cfg.Server.Port
	}
	if port == "" {
		port = "8080"
	}

	e := echo.New()

	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	routes.SetupRoutes(e)

	log.Printf("Starting server [%s] on port %s", cfg.AppEnv, port)
	e.Logger.Fatal(e.Start(":" + port))
}
