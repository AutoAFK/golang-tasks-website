// This project is a simple to do list app.
// It will save all the tasks in a txt file.
package main

import (
	"bufio"
	"bytes"
	"errors"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"
	"todo/logger"
)

type Task string

var (
	templates   = template.Must(template.ParseFiles("home.html"))
	infoLogger  = logger.InfoLogger()
	debugLogger = logger.DebugLogger()
	fatalLogger = logger.FatalLogger()
)

func renderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, data any) {
	err := templates.ExecuteTemplate(w, tmpl+".html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func readData(filename string) []string {
	if _, err := os.Stat(filename); err != nil {
		debugLogger.Println("File does not exists")
		return []string{}
	}
	file, err := os.Open(filename)
	if err != nil {
		infoLogger.Println("Couldn't open the file")
		return []string{}
	}
	defer file.Close()
	fileData, err := io.ReadAll(file)
	if err != nil {
		log.Println("Error", log.Lshortfile, ":Couldn't read the file")
		return []string{}
	}
	lines := strings.Split(string(fileData), "\n")
	return lines
}

func writeDataToFile(data string, filename string) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE, fs.FileMode(fs.ModeAppend))
	if err != nil {
		return errors.New("file not found")
	}
	defer file.Close()
	_, err = file.WriteString(data + "\n")
	if err != nil {
		return errors.New("can't append to file")
	}
	return nil
}

func deleteDataFromFile(data string, filename string) error {
	file, _ := os.OpenFile(filename, os.O_RDWR, 0644)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var bs []byte
	buffer := bytes.NewBuffer(bs)
	var text string
	for scanner.Scan() {
		text = scanner.Text()
		if text != data {
			_, err := buffer.WriteString(text + "\n")
			if err != nil {
				return errors.New(err.Error())
			}
		}
	}
	file.Truncate(0)
	file.Seek(0, 0)
	buffer.WriteTo(file)
	return nil
}

func home(w http.ResponseWriter, r *http.Request) {
	data := readData("tasks.txt")
	renderTemplate(w, r, "home", data)
}

func newTask(w http.ResponseWriter, r *http.Request) {
	infoLogger.Println("Adding a task...")
	formData := r.FormValue("new_task")
	debugLogger.Println("Form data:", formData)
	err := writeDataToFile(formData, "tasks.txt")
	if err != nil {
		debugLogger.Println(err.Error())
		http.Redirect(w, r, "/home/", http.StatusInternalServerError)
	}
	infoLogger.Println("Redirect to home page")
	http.Redirect(w, r, "/home/", http.StatusFound)
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/delete/"):]
	err := deleteDataFromFile(id, "tasks.txt")
	if err != nil {
		debugLogger.Println(err.Error())
		http.Redirect(w, r, "/home/", http.StatusInternalServerError)
	}
	infoLogger.Println("Redirect to home page")
	http.Redirect(w, r, "/home/", http.StatusFound)
}

func main() {
	http.HandleFunc("/home/", home)
	http.HandleFunc("/new/", newTask)
	http.HandleFunc("/delete/", deleteTask)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
