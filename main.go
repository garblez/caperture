package main

import (
  "log"
  "net/http"
  "html/template"
)

var templates map[string]*template.Template

func main() {
  templates = make(map[string]*template.Template)
  templates["home"] = template.Must(template.ParseFiles("templates/base.html", "templates/home.html"))
  templates["contact"] = template.Must(template.ParseFiles("templates/base.html", "templates/contact.html"))
  templates["gallery"] = template.Must(template.ParseFiles("templates/base.html", "templates/gallery.html"))
  templates["thanks"] = template.Must(template.ParseFiles("templates/thanks.html"))

  http.HandleFunc("/", homeHandler)
  http.HandleFunc("/contact", contactHandler)
  http.HandleFunc("/gallery", galleryHandler)

  log.Println("Starting the server on :8080")
  if err := http.ListenAndServe(":8080", nil); err != nil {
    log.Fatalf("Could not start the server: %v", err)
  }
}


func homeHandler(w http.ResponseWriter, r *http.Request) {
  data := struct {
    Title string
    Heading string
    Message string
  }{
    Title: "Home Page",
    Heading: "Glasgow Photographer and Alumnus of The Garage Nightclub",
    Message: "Who I am will go here.",
  }

  if err := templates["home"].ExecuteTemplate(w, "base.html", data); err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
  if r.Method == http.MethodPost {
    err := r.ParseForm()
    if err != nil {
      http.Error(w, "Error parsing the form", http.StatusBadRequest)
      return
    }

    name := r.Form.Get("name")
    email := r.Form.Get("email")
    details := r.Form.Get("details")

    submission := struct{
      Name string
      Email string
      Details string
    }{
      Name: name,
      Email: email,
      Details: details,
    }

    log.Printf("Contact form sent with following details: %v", submission)
    if err := templates["thanks"].Execute(w, submission); err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
    }

  } else {
    data := struct {
      Title string
      Heading string
    }{
      Title: "Contact Page",
      Heading: "Contact me using this form and I will get back to you as soon as I can.",
    }

    if err := templates["contact"].ExecuteTemplate(w, "base.html", data); err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
    }
  }
}

func galleryHandler(w http.ResponseWriter, r *http.Request) {
  data := struct {
    Title string
    Heading string
  }{
    Title: "Gallery Page",
    Heading: "Place some photos here and we should be good to go. Likely just link to B2",
  }


  if err := templates["gallery"].ExecuteTemplate(w, "base.html", data); err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
}
