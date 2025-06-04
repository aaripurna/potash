/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/aaripurna/go-web-template/go-fullstack-template/core"
	"github.com/aaripurna/go-web-template/go-fullstack-template/web"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starting the server",
	Long: `Starting the http server
		-p to change port
		-b to change ip binding
	`,
	Run: func(cmd *cobra.Command, args []string) {
		flags := cmd.Flags()

		port, err := flags.GetInt16("port")

		if err != nil {
			panic("Failed to get port")
		}

		bind, err := flags.GetString("bind")

		if err != nil {
			panic("Failed to get bind")
		}

		if err := Container.Invoke(func(engine *html.Engine) {
			core.AssetHtml(engine)
		}); err != nil {
			log.Fatal(err)
		}

		if err := Container.Invoke(pagesRoute); err != nil {
			log.Fatal(err)
		}

		if err := Container.Invoke(func(app *fiber.App) {
			app.Static("/", "./public")

			log.Fatal(app.Listen(fmt.Sprintf("%s:%v", bind, port)))
		}); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().Int16P("port", "p", 8000, "to set the port")
	serveCmd.Flags().StringP("bind", "b", "127.0.0.1", "to set the IP binding")
}

func pagesRoute(app *fiber.App, pagesWeb *web.PagesWeb) {
	app.Get("/", core.HandleReq(pagesWeb.Index))
}
