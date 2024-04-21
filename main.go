package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/russross/blackfriday/v2"
)

type Post struct {
	Title    string
	Content  template.HTML
	Date     string
	Category string
	FilePath string
}

var templates = template.Must(template.ParseGlob("templates/*.html"))
var posts []Post

func main() {
	// Load posts from markdown files
	var err error
	posts, err = loadPostsFromMarkdown()
	if err != nil {
		log.Fatalf("could not load posts: %v", err)
	}

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/category/", categoryHandler)

	// Before the server starts, open the browser
	go func() {
		url := "http://localhost:8080"
		fmt.Println("Server started on", url)
		openBrowser(url)
	}()

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if err := templates.ExecuteTemplate(w, "index.html", posts); err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
	}
}

func categoryHandler(w http.ResponseWriter, r *http.Request) {
	category := strings.TrimPrefix(r.URL.Path, "/category/")
	var filteredPosts []Post
	for _, post := range posts {
		if post.Category == category {
			filteredPosts = append(filteredPosts, post)
		}
	}
	if err := templates.ExecuteTemplate(w, "category.html", filteredPosts); err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
	}
}

func openBrowser(url string) {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Printf("Failed to open browser: %v", err)
	}
}

func loadPostsFromMarkdown() ([]Post, error) {
	files, err := filepath.Glob("content/posts/*.md")
	if err != nil {
		return nil, err
	}

	var loadedPosts []Post
	for _, file := range files {
		mdContent, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, err
		}

		htmlContent := blackfriday.Run(mdContent)
		// Here you should parse the title, date, and category from the file's content or filename
		// For demonstration, I'll use placeholders
		loadedPosts = append(loadedPosts, Post{
			Title:    "Post Title", // Placeholder - parse from file
			Content:  template.HTML(htmlContent),
			Date:     "Post Date",     // Placeholder - parse from file
			Category: "Post Category", // Placeholder - parse from file
			FilePath: file,
		})
	}

	return loadedPosts, nil
}
