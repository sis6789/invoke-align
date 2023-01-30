package main

import (
	"embed"
	"fmt"
	"log"
	"os/exec"
	"path"
	"time"

	"github.com/foolin/goview"
	"github.com/foolin/goview/supports/echoview-v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"

	flag "github.com/spf13/pflag"

	"invoke-align/internal/app/web/handler"
	"invoke-align/internal/pkg/env_log"
)

//go:embed views
var viewFS embed.FS

func embeddedFH(config goview.Config, tmpl string) (string, error) {
	wPath := path.Join(config.Root, tmpl)
	bytes, err := viewFS.ReadFile(wPath + config.Extension)
	if err != nil {
		log.Printf("%v", err)
	}
	return string(bytes), err
}

var servicePort = flag.StringP("webPort", "s", "9099", "web service port number")

func main() {

	go func() {
		<-time.After(100 * time.Millisecond)
		err := exec.Command("explorer", "http://127.0.0.1:9099").Run()
		if err != nil {
			log.Println(err)
		}
	}()

	err := env_log.EnvLog()
	if err != nil {
		log.Fatalf("%v", err)
	}
	log.Printf("web_genome.xeq -h 192.168.0.6 -p 27017 -s 9090\n")

	// parameter
	flag.Parse()
	err = viper.BindPFlags(flag.CommandLine)
	if err != nil {
		log.Printf("%v", err)
	}

	// Echo instance
	e := echo.New()
	e.Logger.SetOutput(env_log.GetFile())

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Set Renderer using embedded FS
	viewEngine := echoview.New(goview.Config{
		Root:      "views",
		Extension: ".html",
	})
	viewEngine.SetFileHandler(embeddedFH)
	e.Renderer = viewEngine // echoview.Default()

	// Routes
	e.GET("/", handler.InvokeMain)
	e.POST("/submit", handler.Submit)

	// Start server
	webAddress := fmt.Sprintf(":%s", viper.GetString("webPort")) // 8001 default
	e.Logger.Fatal(e.Start(webAddress))

}
