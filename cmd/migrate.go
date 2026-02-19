package cmd

import (
	"fin-auth/config"
	"fin-auth/database"
	"log"

	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations",
	Long:  `Run GORM AutoMigrate to create/update database tables based on models`,
	Run: func(cmd *cobra.Command, args []string) {
		runMigration()
	},
}

var migrateUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Run migrations (create/update tables)",
	Long:  `Run GORM AutoMigrate to create or update database tables`,
	Run: func(cmd *cobra.Command, args []string) {
		runMigration()
	},
}

var migrateStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check migration status",
	Long:  `Check which tables exist in the database`,
	Run: func(cmd *cobra.Command, args []string) {
		checkMigrationStatus()
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
	migrateCmd.AddCommand(migrateUpCmd)
	migrateCmd.AddCommand(migrateStatusCmd)
}

func runMigration() {
	if err := config.LoadConfig(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	log.Println("Connecting to database...")
	db, err := config.InitGormDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database instance: %v", err)
	}
	defer sqlDB.Close()

	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println("Migration completed successfully!")
}

func checkMigrationStatus() {
	if err := config.LoadConfig(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	log.Println("Connecting to database...")
	db, err := config.InitGormDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database instance: %v", err)
	}
	defer sqlDB.Close()

	models := database.GetAllModels()
	log.Println("\nChecking migration status:")
	log.Println("----------------------------")

	for _, model := range models {
		tableName := ""
		if err := db.Raw("SELECT ?::regclass", db.Statement.Table).Scan(&tableName).Error; err == nil {
			if db.Migrator().HasTable(model) {
				log.Printf("✓ Table exists: %s", db.Statement.Table)
			} else {
				log.Printf("✗ Table missing: %s", db.Statement.Table)
			}
		}
	}

	log.Println("----------------------------")
}
