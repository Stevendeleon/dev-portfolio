package main

import (
	"html/template"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Portfolio struct {
	Name            string        `yaml:"name"`
	Title           string        `yaml:"title"`
	JobTitle        template.HTML `yaml:"job_title"`
	Blurb           template.HTML `yaml:"blurb"`
	CurrentPosition string        `yaml:"current_position"`
	Contact         Contact       `yaml:"contact"`
	Socials         []Social      `yaml:"socials"`
	Experience      []Experience  `yaml:"experience"`
	Projects        []Project     `yaml:"projects"`
}

type Contact struct {
	Message string `yaml:"message"`
	Email   string `yaml:"email"`
}

type Social struct {
	Name string `yaml:"name"`
	URL  string `yaml:"link"`
}

type Experience struct {
	StartDate   string   `yaml:"start_date"`
	EndDate     string   `yaml:"end_date"`
	Position    string   `yaml:"position_title"`
	Company     string   `yaml:"company_title"`
	Description string   `yaml:"description"`
	Skills      []string `yaml:"skills"`
}

type Project struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Image       struct {
		Src string `yaml:"src"`
		Alt string `yaml:"alt"`
	} `yaml:"image"`
	Tags  []string `yaml:"tags"`
	Links []struct {
		Href  string        `yaml:"href"`
		Svg   template.HTML `yaml:"svg"`
		Label string        `yaml:"label"`
	}
}

type templateData struct {
	Title     string
	Portfolio Portfolio
}

var (
	data templateData
	tmpl *template.Template
)

func main() {
	portfolio, err := loadPortfolioData("data.yaml")
	if err != nil {
		log.Fatalf("Error loading portfolio data: %v", err)
	}

	tmpl = template.Must(template.ParseFiles("templates/base.html", "templates/projects.tmpl.html"))
	data = templateData{
		Title:     portfolio.Title,
		Portfolio: portfolio,
	}

	err = generateOutputFile()
	if err != nil {
		log.Fatalf("Error generating output file: %v", err)
	}

	log.Println("Completed generating output file.")
}

func loadPortfolioData(filename string) (Portfolio, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return Portfolio{}, err
	}

	var portfolio Portfolio
	err = yaml.Unmarshal(content, &portfolio)
	if err != nil {
		return Portfolio{}, err
	}

	return portfolio, nil
}

func generateOutputFile() error {
	filename := "index.html"
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	err = tmpl.Execute(file, data)
	if err != nil {
		return err
	}

	return nil
}
