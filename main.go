package main

import (
	"net/http"
	"sync"

	goutils "github.com/jkaninda/go-utils"
	"github.com/jkaninda/okapi"
	"github.com/jkaninda/okapi/okapicli"
)

var (
	fileMutex sync.Mutex
	filename  = goutils.Env("DATA_FILE_NAME", "data/notifications.json")
)

func main() {
	o := okapi.Default().WithRendererFromDirectory("views", ".html", ".templ")
	cli := okapicli.New(o, "Uzaraka").Int("port", "p", 8080, "HTTP server port")
	// Parse flags
	if err := cli.Parse(); err != nil {
		panic(err)
	}
	o.WithPort(cli.GetInt("port"))
	o.Get("/", func(c okapi.Ctx) error {
		title := "Uzaraka - Bientôt Disponible"
		desciption := "Uzaraka arrive bientôt ! La nouvelle plateforme de petites annonces moderne et intuitive. Lancement le 15 décembre 2025."
		return c.Render(http.StatusOK, "home", okapi.M{
			"title":      title,
			"desciption": desciption,
		})
	})
	o.Get("/healthz", func(c okapi.Ctx) error {
		return c.OK(okapi.M{
			"status": "healthy"})
	})
	o.Post("/notify", func(c okapi.Ctx) error {
		request := &EmailRequest{}
		if err := c.Bind(request); err != nil {
			return c.AbortBadRequest("Bad request", err)
		}

		if err := saveEmail(request.Email); err != nil {
			return c.AbortInternalServerError("Internal server error", err)

		}
		return c.OK(okapi.M{
			"success": true,
			"message": "Email saved successfully",
		})
	}).WithInput(&EmailRequest{})

	// Start the server
	err := cli.Run()
	if err != nil {
		panic(err)
	}
}
