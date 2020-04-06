package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
)

const (
	templateDir = "tmpl/"
)

var validPath = regexp.MustCompile("^/(action|edit|save)/([a-zA-Z0-9]+)$")
var templates = template.Must(template.ParseFiles(templateDir+"index.html", templateDir+"edit.html"))
var actions = make(map[string]BoyfriendActor)

func loadActions() (map[string]BoyfriendActor, error) {
	actions := make(map[string]BoyfriendActor)
	fileInfo, err := ioutil.ReadDir(dataDir)
	if err != nil {
		log.Print("dataDir Could Not Be found")
		return nil, err
	}

	for _, file := range fileInfo {
		name := file.Name()
		name = name[:len(name)-4]

		action, err := loadAction(name)
		if err != nil {
			log.Print("dataDir Could Not Be found")
			return nil, err
		}
		actions[name] = action
	}

	return actions, nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, action BoyfriendActor) {
	err := templates.ExecuteTemplate(w, tmpl+".html", action)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

func editHandler(w http.ResponseWriter, r *http.Request, actionName string) {
	action, err := loadAction(actionName)
	if err != nil {
		action = NewGenericAction(actionName, "")
	}
	renderTemplate(w, "edit", action)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	text := r.FormValue("textMsg")
	action := NewGenericAction(title, text)
	err := action.Save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func actionHandler(w http.ResponseWriter, r *http.Request, actionName string) {
	action, err := loadAction(actionName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = action.NotifyBoyfriend()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func homePageHandler(w http.ResponseWriter, r *http.Request) {
	actions, err := loadActions()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	err = templates.ExecuteTemplate(w, "index.html", actions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {

	accountSid := os.Getenv("TWILIO_ACCOUNT_ID")
	authToken := os.Getenv("TWILIO_AUTH_TOKEN")
	bfNumber := os.Getenv("BFMC_RECIEVER_NUMBER")
	fromNumber := os.Getenv("BFMC_FROM_NUMBER")

	twilioClient := NewTwilioClient(
		accountSid,
		authToken,
		bfNumber,
		fromNumber,
	)

	setClient(twilioClient)
	http.HandleFunc("/save/", makeHandler(saveHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/action/", makeHandler(actionHandler))
	http.HandleFunc("/", homePageHandler)
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", nil))
}
