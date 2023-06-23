package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

var knownDialects = map[string]func(dsn string) gorm.Dialector{
	"sqlite":     sqlite.Open,
	"postgres":   postgres.Open,
	"postgresql": postgres.Open,
	"mysql":      mysql.Open,
	"sqlserver":  sqlserver.Open,
}

func openDatabase(cmd *cobra.Command) (*gorm.DB, error) {
	dbAddress, err := cmd.Flags().GetString("db")
	if err != nil {
		return nil, err
	}

	n := strings.Index(dbAddress, "://")
	if n == -1 {
		return nil, errors.New("invalid database connection information")
	}

	config := &gorm.Config{}
	dialect, dataSourceName := dbAddress[:n], dbAddress[n+3:]
	if opener, ok := knownDialects[dialect]; ok {
		return gorm.Open(opener((dataSourceName)), config)
	} else {
		return nil, fmt.Errorf("%q is not a supported database dialect", dialect)
	}
}
