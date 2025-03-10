package main

import (
	"html/template"
	"log"
	"net/http"
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

	tmpl = template.Must(template.ParseFiles("templates/base.html"))
	data = templateData{
		Title:     portfolio.Title,
		Portfolio: portfolio,
	}

	err = generateOutputFile()
	http.HandleFunc("/", renderPortfolio)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}

// For dev testing prior to generating the index.html
func renderPortfolio(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := tmpl.Execute(w, data)
	if err != nil {
		msg := http.StatusText(http.StatusInternalServerError)
		log.Printf("template.Execute: %v", err)
		http.Error(w, msg, http.StatusInternalServerError)
	}
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

	log.Printf("Generated %s successfully", filename)
	return nil
}
