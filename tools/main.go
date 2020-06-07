package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sort"
	"time"

	"github.com/BurntSushi/toml"
)

// Content represents a single content item in a post
type Content struct {
	Kind    string
	Body    string
	Source  string
	Alt     string
	Items   []string
	Title   string
	Left    bool
	Feature bool
}

// Post represents the entire contents of a post
type Post struct {
	Date        string
	Title       string
	Filename    string
	Description string
	URL         string
	Image       string
	Content     []Content
}

// Project represents a single project
type Project struct {
	Name   string
	Image  string
	URL    string
	Github string
	RunURL string
}

// Projects represent the contents of the project page
type Projects struct {
	Title       string
	Filename    string
	Description string
	URL         string
	Image       string
	Projects    []Project
}

// Home represents the contents of the home page
type Home struct {
	Title       string
	Description string
	Body        string
	URL         string
	Image       string
	Posts       []Post
	About       Post
	Resume      Post
	Projects    Projects
}

var outputDir = ".."
var htmlSafe = func(html string) template.HTML {
	return template.HTML(html)
}

func parseTOML() *Home {
	homeTOML := "data/home.toml"
	aboutTOML := "data/about.toml"
	resumeTOML := "data/resume.toml"
	projectTOML := "data/projects.toml"
	postDir := "data/posts"

	home := Home{}
	if _, err := toml.DecodeFile(homeTOML, &home); err != nil {
		log.Fatal(err)
	}

	files, err := ioutil.ReadDir(postDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		post := Post{}
		if _, err := toml.DecodeFile(path.Join(postDir, f.Name()), &post); err != nil {
			log.Fatal(err)
		}
		post.URL = home.URL + "/" + post.Filename
		home.Posts = append(home.Posts, post)
	}

	if _, err := toml.DecodeFile(aboutTOML, &home.About); err != nil {
		log.Fatal(err)
	}

	if _, err := toml.DecodeFile(resumeTOML, &home.Resume); err != nil {
		log.Fatal(err)
	}

	if _, err := toml.DecodeFile(projectTOML, &home.Projects); err != nil {
		log.Fatal(err)
	}

	// Sort the Posts by date
	sort.Slice(home.Posts[:], func(i, j int) bool {
		t1, err := time.Parse("01/02/2006", home.Posts[i].Date)
		if err != nil {
			return false
		}
		t2, err := time.Parse("01/02/2006", home.Posts[j].Date)
		if err != nil {
			return false
		}

		return t1.After(t2)
	})

	return &home
}

func render(data interface{}, tmpl *template.Template, filename string) error {
	output, err := os.OpenFile(path.Join(outputDir, filename), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	fmt.Printf("Creating %s...\n", filename)
	err = tmpl.Execute(output, data)
	if err != nil {
		return err
	}

	return nil
}

func renderHome(home *Home) {
	tmpl, err := template.New("home.tmpl").Funcs(template.FuncMap{
		"htmlSafe": htmlSafe,
	}).ParseFiles(
		"templates/home.tmpl",
		"templates/header.tmpl",
		"templates/footer.tmpl",
	)

	if err != nil {
		log.Fatal(err)
	}
	err = render(home, tmpl, "index.html")
	if err != nil {
		log.Fatal(err)
	}
}

func renderRSS(home *Home) {
	tmpl, err := template.ParseFiles(
		"templates/rss.tmpl",
	)

	if err != nil {
		log.Fatal(err)
	}

	err = render(home, tmpl, "rss.xml")
	if err != nil {
		log.Fatal(err)
	}
}

func renderProjects(home *Home) {
	tmpl, err := template.ParseFiles(
		"templates/projects.tmpl",
		"templates/header.tmpl",
		"templates/footer.tmpl",
	)

	if err != nil {
		log.Fatal(err)
	}

	err = render(home.Projects, tmpl, "projects.html")
	if err != nil {
		log.Fatal(err)
	}
}

func renderPosts(home *Home) {
	tmpl, err := template.New("post.tmpl").Funcs(template.FuncMap{
		"htmlSafe": htmlSafe,
	}).ParseFiles(
		"templates/post.tmpl",
		"templates/header.tmpl",
		"templates/footer.tmpl",
	)

	if err != nil {
		log.Fatal(err)
	}

	err = render(home.About, tmpl, "about.html")
	if err != nil {
		log.Fatal(err)
	}

	err = render(home.Resume, tmpl, "resume.html")
	if err != nil {
		log.Fatal(err)
	}

	for _, p := range home.Posts {
		err = render(p, tmpl, p.Filename)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	home := parseTOML()

	renderHome(home)
	renderRSS(home)
	renderProjects(home)
	renderPosts(home)
}
