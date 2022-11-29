package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/flosch/pongo2"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const tmplPath = "src/template/"

var e = createMux()

func main() {
	connectDB()

	e.GET("/", articleIndex)
	e.GET("/new", articleNew)
	e.GET("/:id", articleShow)
	e.GET("/:id/edit", articleEdit)
	e.Logger.Fatal(e.Start(":8080"))
}

func createMux() *echo.Echo {
	e := echo.New()

	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(middleware.Gzip())

	e.Static("/css", "src/css")
	e.Static("/js", "src/js")

	return e
}

func connectDB() {
	err := godotenv.Load(".env")
	if err != nil {
		panic("Error loading .env file")
	}
	dbconf := os.Getenv("DSN")

	db, err := sql.Open("mysql", dbconf)

	defer db.Close()

	if err != nil {
		log.Println(err.Error())
	}

	err = db.Ping()

	if err != nil {
		log.Println("データベース接続失敗")
		return
	} else {
		log.Println("データベース接続成功")
	}
}

func articleIndex(c echo.Context) error {
	data := map[string]interface{}{
		"Message": "Article Index",
		"Now":     time.Now(),
	}
	return render(c, "article/index.html", data)
}

func articleNew(c echo.Context) error {
	data := map[string]interface{}{
		"Message": "Article New",
		"Now":     time.Now(),
	}

	return render(c, "article/new.html", data)
}

func articleShow(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	data := map[string]interface{}{
		"Message": "Article Show",
		"Now":     time.Now(),
		"ID":      id,
	}

	return render(c, "article/show.html", data)
}

func articleEdit(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	data := map[string]interface{}{
		"Message": "Article Edit",
		"Now":     time.Now(),
		"ID":      id,
	}

	return render(c, "article/edit.html", data)
}

// テンプレートの情報をバイトにして返す
func htmlBlob(file string, data map[string]interface{}) ([]byte, error) {
	return pongo2.Must(pongo2.FromCache(tmplPath + file)).ExecuteBytes(data)
}

func render(c echo.Context, file string, data map[string]interface{}) error {
	b, err := htmlBlob(file, data)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.HTMLBlob(http.StatusOK, b)
}
