package server

import (
	"context"
	"crypto/subtle"
	"database/sql"
	"desprit/hammertime/src/config"
	"io"
	"text/template"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

type Server struct {
	db *sql.DB
	e  *echo.Echo
}

func NewServer(d *sql.DB) *Server {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		cfg := config.GetConfig()
		if subtle.ConstantTimeCompare([]byte(username), []byte(cfg.WEB_USER)) == 1 &&
			subtle.ConstantTimeCompare([]byte(password), []byte(cfg.WEB_PASS)) == 1 {
			return true, nil
		}
		return false, nil
	}))
	t := &Template{
		templates: template.Must(template.ParseGlob("src/templates/*.html")),
	}
	e.Renderer = t
	e.Static("/static", "src/assets")

	e.GET("/", func(c echo.Context) error {
		return c.Redirect(302, "/schedule")
	})

	return &Server{db: d, e: e}
}

func (s *Server) RegisterHandlers() {
	NewScheduleServer(s.e).RegisterHandlers(s.db)
	NewSubscriptionServer(s.e).RegisterHandlers(s.db)
}

func (s *Server) Logger() echo.Logger {
	return s.e.Logger
}

func (s *Server) Start() error {
	return s.e.Start(":8080")
}

func (s *Server) Shutdown() error {
	return s.e.Shutdown(context.Background())
}
