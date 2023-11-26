package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"regexp"
	"strings"

	"log"

	"github.com/gorilla/mux"
	"github.com/russross/blackfriday/v2"
)

var t = template.Must(template.ParseFiles("public/index.html"))

type Link struct {
	Title string
	Link  string
}

type ProjectCard struct {
	Title       string
	Description string
	Image       string
	Stack       string
}

type Project struct {
	Title       string
	Description string
	Image       string
	Stack       string
	Links       []Link
	HTML        template.HTML
}

// Metadata represents the metadata extracted from a Markdown file
type Metadata map[string]string

// ExtractMetadata extracts metadata from a Markdown file
func ExtractMetadata(filePath string) (Metadata, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	contentStr := string(content)

	metadataRegex := regexp.MustCompile(`(?s)---(.*?)---`)

	match := metadataRegex.FindStringSubmatch(contentStr)

	if len(match) > 1 {
		contentStr = contentStr[len(match[0]):]
	}

	lines := strings.Split(match[1], "\n")

	metadata := make(Metadata)

	for _, line := range lines {
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			metadata[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}

	return metadata, nil
}

func main() {

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	router := mux.NewRouter()
	fs := http.FileServer(http.Dir("public"))
	router.PathPrefix("/public/").Handler(http.StripPrefix("/public/", fs))
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := map[string][]ProjectCard{
			"Projects": {},
		}

		files, err := os.ReadDir("./projects")
		if err != nil {
			log.Fatal(err)
		}
		for _, file := range files {

			metadata, err := ExtractMetadata("./projects/" + file.Name())

			if err != nil {
				os.Exit(1)
			}
			project := ProjectCard{metadata["Title"], metadata["Description"], metadata["Image"], metadata["Stack"]}
			data["Projects"] = append(data["Projects"], project)

		}

		t.ExecuteTemplate(w, "index.html", data)
	}).Methods("GET")

	router.HandleFunc("/project", func(w http.ResponseWriter, r *http.Request) {

		ret := template.Must(template.ParseFiles("html/project-modal.html"))
		projectName := r.URL.Query().Get("project")
		filePath := "./projects/" + projectName + ".md"
		content, err := os.ReadFile(filePath)

		if err != nil {
			return
		}

		metadata, err := ExtractMetadata(filePath)

		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		title := metadata["Title"]
		description := metadata["Description"]
		image := metadata["Image"]
		stack := metadata["Stack"]
		linksStr := metadata["Links"]
		contentStr := string(content)

		metadataRegex := regexp.MustCompile(`(?s)---(.*?)---`)

		match := metadataRegex.FindStringSubmatch(contentStr)

		if len(match) > 1 {
			contentStr = contentStr[len(match[0]):]
		}

		rendered := string(blackfriday.Run([]byte(contentStr)))

		var links []Link
		if linksStr != "" {
			linksArr := strings.Split(linksStr, ",")
			for _, link := range linksArr {
				fmt.Println(link)
				linkArr := strings.Split(link, ">")
				links = append(links, Link{linkArr[0], linkArr[1]})
			}
		}

		project := Project{title, description, image, stack, links, template.HTML(rendered)}

		ret.ExecuteTemplate(w, "project-modal.html", project)

	}).Methods("GET")

	log.Println("listening on", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
