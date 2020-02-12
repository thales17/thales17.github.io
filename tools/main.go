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

type Content struct {
	Kind   string
	Body   string
	Source string
	Alt    string
	Items  []string
	Title  string
}

type Post struct {
	Date        string
	Title       string
	Filename    string
	Description string
	URL         string
	Image       string
	Content     []Content
}

type Home struct {
	Title       string
	Description string
	Body        string
	URL         string
	Image       string
	Posts       []Post
	About       Post
	Resume      Post
}

func main() {

	homeTOML := "data/home.toml"
	aboutTOML := "data/about.toml"
	resumeTOML := "data/resume.toml"
	postDir := "data/posts"
	outputDir := ".."
	htmlSafe := func(html string) template.HTML {
		return template.HTML(html)
	}

	// Parse the Input
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

	// Write the Output
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

	outputHome, err := os.OpenFile(path.Join(outputDir, "index.html"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Creating index.html...")
	err = tmpl.Execute(outputHome, home)
	if err != nil {
		log.Fatal(err)
	}

	tmpl, err = template.ParseFiles(
		"templates/rss.tmpl",
	)

	if err != nil {
		log.Fatal(err)
	}

	outputRSS, err := os.OpenFile(path.Join(outputDir, "rss.xml"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Creating rss.xml...")
	err = tmpl.Execute(outputRSS, home)
	if err != nil {
		log.Fatal(err)
	}

	tmpl, err = template.New("post.tmpl").Funcs(template.FuncMap{
		"htmlSafe": htmlSafe,
	}).ParseFiles(
		"templates/post.tmpl",
		"templates/header.tmpl",
		"templates/footer.tmpl",
	)

	if err != nil {
		log.Fatal(err)
	}

	outputAbout, err := os.OpenFile(path.Join(outputDir, "about.html"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Creating about.html...")
	err = tmpl.Execute(outputAbout, home.About)
	if err != nil {
		log.Fatal(err)
	}

	outputResume, err := os.OpenFile(path.Join(outputDir, "resume.html"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Creating resume.html...")
	err = tmpl.Execute(outputResume, home.Resume)
	if err != nil {
		log.Fatal(err)
	}

	for _, p := range home.Posts {
		fmt.Printf("Creating %s...\n", p.Filename)
		outputPost, err := os.OpenFile(path.Join(outputDir, p.Filename), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			log.Fatal(err)
		}

		err = tmpl.Execute(outputPost, p)
		if err != nil {
			log.Fatal(err)
		}
	}
}
