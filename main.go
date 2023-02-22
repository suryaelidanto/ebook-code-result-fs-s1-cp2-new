package main

import (
	"context"
	"ebook-s1-cp2/connection"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

type Blog struct {
	ID         int
	Title      string
	Content    string
	Image      string
	Author     string
	PostDate   time.Time
	FormatDate string
}

var dataBlog = []Blog{
	{
		Title:    "Hallo Title",
		Content:  "Hallo Content",
		Author:   "Surya Elidanto",
		PostDate: time.Now(),
	},
	{
		Title:    "Hallo Title 2",
		Content:  "Hallo Content 2",
		Author:   "Surya Elidanto",
		PostDate: time.Now(),
	},
}

func main() {
	// Called database connection function
	connection.DatabaseConnect()

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
	e.GET("/blog-detail/:id", blogDetail)
	e.GET("/form-blog", formAddBlog)
	e.POST("/blog-delete/:id", deleteBlog)
	e.POST("/add-blog", addBlog)

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
	data, _ := connection.Conn.Query(context.Background(), "SELECT id, title, content, image, post_date FROM tb_blog")

	var result []Blog
	for data.Next() {
		var each = Blog{}

		err := data.Scan(&each.ID, &each.Title, &each.Content, &each.Image, &each.PostDate)
		if err != nil {
			fmt.Println(err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
		}

		each.Author = "Surya Elidanto"
		each.FormatDate = each.PostDate.Format("2 January 2006")

		result = append(result, each)
	}

	blogs := map[string]interface{}{
		"Blogs": result,
	}

	return c.Render(http.StatusOK, "blog.html", blogs)
}

func blogDetail(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	var BlogDetail = Blog{}

	err := connection.Conn.QueryRow(context.Background(), "SELECT id, title, content, image, post_date FROM tb_blog WHERE id=$1", id).Scan(
		&BlogDetail.ID, &BlogDetail.Title, &BlogDetail.Content, &BlogDetail.Image, &BlogDetail.PostDate)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	BlogDetail.Author = "Surya Elidanto"
	BlogDetail.FormatDate = BlogDetail.PostDate.Format("2 January 2006")

	data := map[string]interface{}{
		"Blog": BlogDetail,
	}

	return c.Render(http.StatusOK, "blog-detail.html", data)
}

func formAddBlog(c echo.Context) error {
	return c.Render(http.StatusOK, "add-blog.html", nil)
}

func addBlog(c echo.Context) error {
	title := c.FormValue("inputTitle")
	content := c.FormValue("inputContent")
	// author := "SuryaElz"  you can make this with manually input, and add author column in table, challenge yourself :)

	_, err := connection.Conn.Exec(context.Background(), "INSERT INTO tb_blog (title, content, image, post_date) VALUES ($1, $2, $3, $4)", title, content, "blog-img.png", time.Now())

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.Redirect(http.StatusMovedPermanently, "/blog")
}

func deleteBlog(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	dataBlog = append(dataBlog[:id], dataBlog[id+1:]...)

	return c.Redirect(http.StatusMovedPermanently, "/blog")
}
