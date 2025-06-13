package migrate

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/cobra"

	"github.com/tuananhlai/brevity-go/internal/config"
)

const migrationsDir = "db/migrations"

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Database migration commands",
}

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Run all pending migrations",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.MustLoadConfig()
		m, err := migrate.New("file://"+migrationsDir, cfg.Database.URL)
		if err != nil {
			fmt.Println("Failed to create migrator:", err)
			os.Exit(1)
		}
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			fmt.Println("Migration up failed:", err)
			os.Exit(1)
		}
		fmt.Println("Migrations applied successfully.")
	},
}

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Rollback the last migration",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.MustLoadConfig()
		m, err := migrate.New("file://"+migrationsDir, cfg.Database.URL)
		if err != nil {
			fmt.Println("Failed to create migrator:", err)
			os.Exit(1)
		}
		if err := m.Steps(-1); err != nil && err != migrate.ErrNoChange {
			fmt.Println("Migration down failed:", err)
			os.Exit(1)
		}
		fmt.Println("Rolled back the last migration.")
	},
}

var createCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create a new migration (requires NAME)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		timestamp := time.Now().Format("20060102150405")
		upFile := fmt.Sprintf("%s/%s_%s.up.sql", migrationsDir, timestamp, name)
		downFile := fmt.Sprintf("%s/%s_%s.down.sql", migrationsDir, timestamp, name)
		if err := os.WriteFile(upFile, []byte("-- +migrate Up\n"), 0o644); err != nil {
			fmt.Println("Failed to create up migration:", err)
			os.Exit(1)
		}
		if err := os.WriteFile(downFile, []byte("-- +migrate Down\n"), 0o644); err != nil {
			fmt.Println("Failed to create down migration:", err)
			os.Exit(1)
		}
		fmt.Printf("Created migration %s and %s\n", upFile, downFile)
	},
}

var forceCmd = &cobra.Command{
	Use:   "force [version]",
	Short: "Force migration version (requires VERSION)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.MustLoadConfig()
		m, err := migrate.New("file://"+migrationsDir, cfg.Database.URL)
		if err != nil {
			fmt.Println("Failed to create migrator:", err)
			os.Exit(1)
		}
		version := args[0]
		var v int
		_, err = fmt.Sscanf(version, "%d", &v)
		if err != nil {
			fmt.Println("Invalid version number:", err)
			os.Exit(1)
		}
		if err := m.Force(v); err != nil {
			fmt.Println("Force migration failed:", err)
			os.Exit(1)
		}
		fmt.Println("Forced migration version to", v)
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show current migration version",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.MustLoadConfig()
		m, err := migrate.New("file://"+migrationsDir, cfg.Database.URL)
		if err != nil {
			fmt.Println("Failed to create migrator:", err)
			os.Exit(1)
		}
		v, dirty, err := m.Version()
		if err != nil {
			fmt.Println("Failed to get migration version:", err)
			os.Exit(1)
		}
		fmt.Printf("Current migration version: %d (dirty: %v)\n", v, dirty)
	},
}

var gotoCmd = &cobra.Command{
	Use:   "goto [version]",
	Short: "Migrate to specific version (requires VERSION)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.MustLoadConfig()
		m, err := migrate.New("file://"+migrationsDir, cfg.Database.URL)
		if err != nil {
			fmt.Println("Failed to create migrator:", err)
			os.Exit(1)
		}
		version := args[0]
		var v uint
		_, err = fmt.Sscanf(version, "%d", &v)
		if err != nil {
			fmt.Println("Invalid version number:", err)
			os.Exit(1)
		}
		if err := m.Migrate(v); err != nil && err != migrate.ErrNoChange {
			fmt.Println("Goto migration failed:", err)
			os.Exit(1)
		}
		fmt.Println("Migrated to version", v)
	},
}

var dropCmd = &cobra.Command{
	Use:   "drop",
	Short: "Drop everything in the database",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("WARNING: This will drop all tables in the database.")
		fmt.Print("Are you sure? [y/N] ")
		var confirm string
		fmt.Scanln(&confirm)
		if confirm != "y" && confirm != "Y" {
			fmt.Println("Operation cancelled.")
			return
		}
		cfg := config.MustLoadConfig()
		m, err := migrate.New("file://"+migrationsDir, cfg.Database.URL)
		if err != nil {
			fmt.Println("Failed to create migrator:", err)
			os.Exit(1)
		}
		if err := m.Drop(); err != nil {
			fmt.Println("Drop failed:", err)
			os.Exit(1)
		}
		fmt.Println("All tables dropped.")
	},
}

func init() {
	migrateCmd.AddCommand(upCmd)
	migrateCmd.AddCommand(downCmd)
	migrateCmd.AddCommand(createCmd)
	migrateCmd.AddCommand(forceCmd)
	migrateCmd.AddCommand(versionCmd)
	migrateCmd.AddCommand(gotoCmd)
	migrateCmd.AddCommand(dropCmd)
}

func GetMigrateCmd() *cobra.Command {
	return migrateCmd
}
