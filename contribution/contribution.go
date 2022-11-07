package contribution

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"strconv"
)

var BASEURL = "https://github.com/"

type User struct {
	Name string
	Year string
}

type UserURL struct {
	Name string
	Url  string
}

type Result struct {
	UserName      string         `json:"userName"`
	Contributions []Contribution `json:"contributions"`
}

type Contribution struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
	Level int    `json:"level"`
}

func Scr(user User) Result {
	var c1 = make(chan UserURL)
	var c2 = make(chan []Contribution)
	var c4 = make(chan string)

	go getURL(user, c1)

	githubLink := <-c1
	go getContribution(githubLink, c2, c4)

	con := <-c2
	name := <-c4
	result := Result{
		UserName:      name,
		Contributions: con}
	return result
}

func getContribution(user UserURL, c2 chan []Contribution, c4 chan string) {
	var Contributions []Contribution
	c3 := make(chan Contribution)

	fmt.Println(user.Name + " : " + user.Url)
	res, err := http.Get(user.Url)

	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)

	grass := doc.Find(" g > rect ")

	grass.Each(func(i int, s *goquery.Selection) {
		go extractGrass(s, c3)
	})
	len := grass.Length()
	for i := 0; i < len; i++ {
		contribution := <-c3
		Contributions = append(Contributions, contribution)
	}
	c2 <- Contributions
	c4 <- user.Name
}

func extractGrass(s *goquery.Selection, c3 chan Contribution) {
	date, _ := s.Attr("data-date")
	tmpCount, _ := s.Attr("data-count")
	tmpLevel, _ := s.Attr("data-level")

	count, err := strconv.Atoi(tmpCount)
	checkErr(err)
	level, err := strconv.Atoi(tmpLevel)
	checkErr(err)
	c3 <- Contribution{
		Date:  date,
		Count: count,
		Level: level}
}

func getURL(user User, c1 chan UserURL) {
	url := BASEURL + user.Name + "?tab=overview&from=" + user.Year + "-01-01&to=" + user.Year + "-12-31"
	c1 <- UserURL{
		Name: user.Name,
		Url:  url}
}

func checkCode(res *http.Response) {
	if res.StatusCode >= 400 {
		log.Fatalln(res.StatusCode)
	}
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
