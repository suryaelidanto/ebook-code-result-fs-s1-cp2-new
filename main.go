package main

import (
	"html/template"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	// Create new Echo instance
	e := echo.New()

	// Serve static files from "public" directory
	e.Static("/public", "public")

	t := &Template{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}

	e.Renderer = t

	// Routing
	e.GET("/hello", helloWorld)
	e.GET("/", home)
	e.GET("/contact", contact)
	e.GET("/blog", blog)

	// Start server
	println("Server running on port 5000")
	e.Logger.Fatal(e.Start("localhost:5000"))
}

func helloWorld(c echo.Context) error {
	return c.String(http.StatusOK, "Hello World")
}

func home(c echo.Context) error {
	return c.Render(http.StatusOK, "index.html", nil)
}

func contact(c echo.Context) error {
	return c.Render(http.StatusOK, "contact.html", nil)
}

func blog(c echo.Context) error {
	return c.Render(http.StatusOK, "blog.html", nil)
}
