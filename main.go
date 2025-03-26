package main

import (
	"encoding/xml"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

type Portfolio struct {
	Name            string        `yaml:"name"`
	Title           string        `yaml:"title"`
	Url             string        `yaml:"url"`
	JobTitle        template.HTML `yaml:"job_title"`
	Blurb           template.HTML `yaml:"blurb"`
	CurrentPosition string        `yaml:"current_position"`
	Contact         Contact       `yaml:"contact"`
	Socials         []Social      `yaml:"socials"`
	Experience      []Experience  `yaml:"experience"`
	Meta            Meta          `yaml:"meta"`
}

type Meta struct {
	Keywords    string `yaml:"keywords"`
	Description string `yaml:"description"`
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
	Meta      Meta
	Domain    string
}

var (
	data      templateData
	indexTmpl *template.Template
)

const distDir = "dist"

func main() {
	if _, err := os.Stat(distDir); os.IsNotExist(err) {
		err = os.MkdirAll(distDir, 0755)
		if err != nil {
			log.Fatalf("Error creating dist directory: %v", err)
		}
		log.Println("Created dist directory")
	}

	portfolio, err := loadPortfolioData("data.yaml")
	if err != nil {
		log.Fatalf("Error loading portfolio data: %v", err)
	}

	funcMap := template.FuncMap{
		"sub": func(a, b int) int {
			return a - b
		},
	}

	indexTmpl = template.Must(template.New("base.html").Funcs(funcMap).ParseFiles("templates/base.html"))

	data = templateData{
		Title:     portfolio.Title,
		Portfolio: portfolio,
		Meta:      portfolio.Meta,
		Domain:    "stevendeleon.com",
	}

	err = generateIndexPage()
	if err != nil {
		log.Fatalf("Error generating output file: %v", err)
	}

	err = generate404Page()
	if err != nil {
		log.Fatalf("Error generating 404 page: %v", err)
	}

	err = generateSitemap("sitemap.xml", portfolio.Url)
	if err != nil {
		log.Fatalf("Error generating sitemap: %v", err)
	}

	err = generateRobotsTXT("robots.txt", "https://stevendeleon.com/sitemap.xml")
	if err != nil {
		log.Fatalf("Error generating robots.txt: %v", err)
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

func generateIndexPage() error {
	filePath := filepath.Join(distDir, "index.html")
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	err = indexTmpl.Execute(file, data)
	if err != nil {
		return err
	}

	return nil
}

func generate404Page() error {
	notFoundTmpl := template.Must(template.ParseFiles("templates/base_404.html"))
	filePath := filepath.Join(distDir, "404.html")
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	err = notFoundTmpl.Execute(file, data)
	if err != nil {
		return err
	}

	log.Println("Generated 404.html page successfully.")
	return nil
}

type SitemapURL struct {
	Loc        string  `xml:"loc"`
	LastMod    string  `xml:"lastmod"`
	ChangeFreq string  `xml:"changefreq"`
	Priority   float32 `xml:"priority"`
}

type Sitemap struct {
	XMLName     xml.Name     `xml:"urlset"`
	Xmlns       string       `xml:"xmlns,attr"`
	SitemapURLS []SitemapURL `xml:"url"`
}

func generateSitemap(filename string, domain string) error {
	filePath := filepath.Join(distDir, filename)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	sitemap := Sitemap{
		Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
		SitemapURLS: []SitemapURL{
			{
				Loc:        domain + "/",
				LastMod:    time.Now().Format(time.RFC3339),
				ChangeFreq: "monthly",
				Priority:   1.0,
			},
			{
				Loc:        domain + "/index.html",
				LastMod:    time.Now().Format(time.RFC3339),
				ChangeFreq: "monthly",
				Priority:   0.8,
			},
			{
				Loc:        domain + "/404.html",
				LastMod:    time.Now().Format(time.RFC3339),
				ChangeFreq: "monthly",
				Priority:   0.5,
			},
		},
	}

	enc := xml.NewEncoder(file)
	enc.Indent("  ", "    ")
	_, err = file.WriteString(xml.Header)
	if err != nil {
		return err
	}

	err = enc.Encode(sitemap)
	if err != nil {
		return err
	}

	fmt.Println("Sitemap generated successfully at", filename)
	return nil
}

func generateRobotsTXT(filename, sitemapURL string) error {
	filePath := filepath.Join(distDir, filename)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	content := fmt.Sprintf("User-agent: *\nAllow: /\n\nSitemap: %s", sitemapURL)

	_, err = file.WriteString(content)
	if err != nil {
		return err
	}

	fmt.Println("robots.txt generated successfully at", filename)
	return nil
}
