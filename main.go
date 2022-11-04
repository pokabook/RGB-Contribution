package main

import (
	"RGBContribution/contribution"
	"github.com/labstack/echo/v4"
	"net/http"
)

func one(c echo.Context) error {
	name := c.Param("name")
	year := c.Param("year")

	user := contribution.User{
		Name: name,
		Year: year}

	result := contribution.Scr(user)

	return c.JSON(http.StatusOK, result)
}

func main() {
	e := echo.New()

	e.GET("/users/:name/:year", one)

	e.Logger.Fatal(e.Start(":8081"))
}
