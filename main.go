package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"text/template"

	goutils "github.com/jkaninda/go-utils"
	"github.com/jkaninda/okapi"
)

type Template struct {
	templates *template.Template
}
type EmailRequest struct {
	Email string `json:"email"`
}

type EmailList struct {
	Emails []string `json:"emails"`
}

var (
	fileMutex sync.Mutex
	filename  = goutils.GetStringEnvWithDefault("DATA_FILE_NAME", "data/notifications.json")
)

func (t *Template) Render(w io.Writer, name string, data interface{}, c okapi.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
func main() {
	o := okapi.New().WithRenderer(&Template{templates: template.Must(template.ParseGlob("views/*.html"))})
	o.Get("/", func(c okapi.Context) error {
		title := "Uzaraka - Bientôt Disponible"
		desciption := "Uzaraka arrive bientôt ! La nouvelle plateforme de petites annonces moderne et intuitive. Lancement le 15 décembre 2025."
		return c.Render(http.StatusOK, "home", okapi.M{
			"title":      title,
			"desciption": desciption,
		})
	})
	o.Get("/healthz", func(c okapi.Context) error {
		return c.OK(okapi.M{
			"status": "healthy"})
	})
	o.Post("/notify", func(c okapi.Context) error {
		request := &EmailRequest{}
		if err := c.Bind(request); err != nil {
			return c.AbortBadGateway("Bad request")
		}
		if request.Email == "" {
			return c.AbortBadGateway("Bad request, email is required")

		}
		if err := saveEmail(request.Email); err != nil {
			return c.AbortInternalServerError("Internal server error")

		}
		return c.OK(okapi.M{
			"message": "Email saved successfully",
		})
	})

	// Start the server
	err := o.Start()
	if err != nil {
		return
	}
}

func saveEmail(email string) error {
	fileMutex.Lock()
	defer fileMutex.Unlock()

	var emailList EmailList

	// Read existing file if it exists
	if _, err := os.Stat(filename); err == nil {
		data, err := os.ReadFile(filename)
		if err != nil {
			return fmt.Errorf("error reading file: %w", err)
		}

		if err := json.Unmarshal(data, &emailList); err != nil {
			return fmt.Errorf("error parsing JSON: %w", err)
		}
	}

	for _, existingEmail := range emailList.Emails {
		if existingEmail == email {
			return nil
		}
	}

	emailList.Emails = append(emailList.Emails, email)

	data, err := json.MarshalIndent(emailList, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %w", err)
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("error writing file: %w", err)
	}
	return nil
}
