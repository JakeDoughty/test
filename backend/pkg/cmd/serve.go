package cmd

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"

	"github.com/JakeDoughty/customer-io-homework-backend/pkg/db/dbal"
	dbutils "github.com/JakeDoughty/customer-io-homework-backend/pkg/db/utils"
	"github.com/JakeDoughty/customer-io-homework-backend/pkg/rest"
)

var listenAddress = ":3000"
var generatorId uint16 = 0
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start processing requests",
	Run: func(cmd *cobra.Command, args []string) {
		dbutils.SetGeneratorId(generatorId)
		if db, err := openDatabase(cmd); err != nil {
			cobra.CheckErr(err)
		} else if err = dbal.AutoMigrateDB(db); err != nil {
			cobra.CheckErr(fmt.Errorf("failed to migrate DB: %w", err))
		} else {
			ctx := dbal.UseDatabase(context.Background(), db)
			engine := gin.Default()
			engine.Use(cors.New(cors.Config{
				AllowAllOrigins: true,
				AllowMethods:    []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
				AllowHeaders:    []string{"*"},
			}))
			engine.Static("static", "public")
			rest.RestServer_v1(engine.Group("/api/v1"))

			webServer := &http.Server{
				Addr:    listenAddress,
				Handler: engine.Handler(),
				// pass the context that contains database reference to the web server
				BaseContext: func(l net.Listener) context.Context { return ctx },
			}

			breakChan := make(chan os.Signal, 1)
			signal.Notify(breakChan, os.Interrupt)
			go func() {
				<-breakChan
				fmt.Println("Received Interrupt singal(Ctrl+C), closing webserver")
				webServer.Shutdown(context.Background())
			}()
			webServer.ListenAndServe()
			fmt.Println("web server stopped")
		}
	},
}

func init() {
	serveCmd.PersistentFlags().StringVar(&listenAddress, "addr", ":3000",
		"the address that server should listen on")
	serveCmd.PersistentFlags().Uint16Var(&generatorId, "instance-id", 1,
		"index of this server in list of listening servers")
	rootCmd.AddCommand(serveCmd)
}
