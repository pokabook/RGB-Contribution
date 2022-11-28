package main

import (
	"RGBContribution/contribution"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

func one(c echo.Context) error {
	name := c.Param("name")
	year := c.Param("year")
	var resultChannel = make(chan contribution.Result)
	user := contribution.User{
		Name: name,
		Year: year}

	go contribution.Scr(user, resultChannel)
	result := <-resultChannel
	return c.JSON(http.StatusOK, result)
}

func main() {
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	e.GET("/users/:name/:year", one)

	e.Logger.Fatal(e.Start(":8090"))
}
