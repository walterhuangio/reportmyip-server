package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo"
)

type config struct {
	Port           int
	ConfigLocation string
}

var realConfig config

func main() {
	argsWithoutProg := os.Args[1:]

	if len(argsWithoutProg) == 0 {
		log.Fatal("Please provide config file as the first argument")
	}

	readConfig(argsWithoutProg[0])
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, I'm Up!")
	})

	e.POST("/hi", hi)

	// write to a file
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", realConfig.Port)))
}

func hi(c echo.Context) error {
	// Get name and email
	name := c.FormValue("name")
	ip := c.FormValue("ip")

	if name == "" || ip == "" {
		return c.String(http.StatusBadRequest, "please provide both name and ip")
	}

	err := ioutil.WriteFile(realConfig.ConfigLocation, []byte(ip), 0644)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusOK, "name:"+name+", ip:"+ip)
}

func readConfig(file string) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal("failed to read config file", err)
	}

	err = json.Unmarshal([]byte(content), &realConfig)
	if err != nil {
		log.Fatal("failed to unmarshal config", err)
	}
}
