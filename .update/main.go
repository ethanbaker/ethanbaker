package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/PuerkitoBio/goquery"
)

/** CONSTANTS **/

// Template path
const templatePath = "./readme.tmpl"

// README path
const readmePath = "../README.md"

/** TYPES **/

// Repository type to store the information of a pinned repo
type Repository struct {
	Title       string
	Url         string
	Description string
}

// Data type to provide to template
type Data struct {
	Age       string
	Repos     []Repository
	UpdatedAt string
}

/** HELPER FUNCTIONS **/

// Calculate my age so the repo stays updated
func calculateAge() string {
	// Get the dates to compare
	birthdate := time.Date(2003, time.February, 8, 0, 0, 0, 0, time.UTC)
	now := time.Now()

	// Find the age
	years := now.Year() - birthdate.Year()
	if now.YearDay() < birthdate.YearDay() {
		years--
	}

	return fmt.Sprint(years)
}

// Get a list of all my pinned repos
func getPinnedRepos() []Repository {
	// Request the HTML page
	res, err := http.Get("https://github.com/ethanbaker")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error (code = %d, status = %s)\n", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatalf("error loading html document (err = %v)\n", err)
	}

	var output []Repository

	// Find the pinned repos
	doc.Find(".pinned-item-list-item-content").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the title, link, and description
		title := s.Find("a > span").Text()
		link := s.Find("a").AttrOr("href", "invalid")
		desc := s.Find(".pinned-item-desc").Text()

		// Add the repo to the list
		output = append(output, Repository{
			Title:       title,
			Url:         "https://github.com" + link,
			Description: desc,
		})
	})

	return output
}

/** MAIN **/

func main() {
	// Get the different components of the README and format it
	age := calculateAge()
	repos := getPinnedRepos()
	last := time.Now().Format("Mon Jan 2 15:04 2006")

	// Update template file
	tmpl, err := template.New(templatePath).ParseFiles(templatePath)
	if err != nil {
		log.Fatalf("could not parse template file (err = %v)\n", err)
	}

	// Create a new file for the README
	file, err := os.Create(readmePath)
	if err != nil {
		log.Fatalf("failed to create readme file (err = %v)\n", err)
	}

	// Execute the template
	err = tmpl.Execute(file, Data{
		Age:       age,
		Repos:     repos,
		UpdatedAt: last,
	})
	if err != nil {
		log.Fatalf("could not execute template (err = %v)\n", err)
	}
}
