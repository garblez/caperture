package main

import (
  "log"
  "fmt"
  "os"
  "net/http"
  "net/smtp"
  "html/template"

)

var templates map[string]*template.Template
var gsmtpEmail string
var gsmtpPassword string
var gsmtpRecipient string

func main() {
  gsmtpEmail = os.Getenv("GSMTP_EMAIL")
  gsmtpPassword = os.Getenv("GSMTP_PASSWORD")
  gsmtpRecipient = os.Getenv("GSMTP_RECIPIENT")
  fs := http.FileServer(http.Dir("./static/"))

  templates = make(map[string]*template.Template)
  templates["home"] = template.Must(template.ParseFiles("templates/base.html", "templates/home.html"))
  templates["contact"] = template.Must(template.ParseFiles("templates/base.html", "templates/contact.html"))
  templates["gallery"] = template.Must(template.ParseFiles("templates/base.html", "templates/gallery.html"))
  templates["thanks"] = template.Must(template.ParseFiles("templates/thanks.html"))


  http.HandleFunc("/", homeHandler)
  http.HandleFunc("/contact", contactHandler)
  http.HandleFunc("/gallery", galleryHandler)
  http.Handle("/static/", http.StripPrefix("/static/", fs))

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
    SubMessage string
  }{
    Title: "Home Page",
    Heading: "Hi, I'm Jonathan!",
    Message: "I'm an experienced Glasgow based nightlife and events photographer.", 
    SubMessage: "Let's talk photos!",
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
    phonenumber := r.Form.Get("phonenumber")

    submission := struct{
      Name string
      Email string
      Details string
      PhoneNumber string
    }{
      Name: name,
      Email: email,
      Details: details,
      PhoneNumber: phonenumber,
    }

    go sendEmail(name, email, details, phonenumber)

    log.Printf("Contact form sent with following details: %v", submission)
    if err := templates["thanks"].Execute(w, submission); err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
    }

  } else {
    data := struct {
      Title string
      Heading string
      SubHeading string
    }{
      Title: "Contact Page",
      Heading: "Contact",
      SubHeading: "Contact me using this form and I will get back to you as soon as I can.",
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


func sendEmail(name, customerEmail, details, phonenumber string) {
  smtpHost := "smtp.gmail.com"
  smtpPort := "587"

  message := fmt.Sprintf("Subject: Gig Request from %s\r\n\r%s\n\nRequest sent from: %s\nTel: (+44) 0%s", name, details, customerEmail, phonenumber)

  auth := smtp.PlainAuth("", gsmtpEmail, gsmtpPassword, smtpHost)

  err := smtp.SendMail(smtpHost+":"+smtpPort, auth, gsmtpEmail, []string{gsmtpRecipient}, []byte(message))

  if err != nil {
    log.Println("Failed to send email: ", err.Error())
    return
  }

  log.Printf("Email sent to %s from customer %s via %s", gsmtpRecipient, customerEmail, gsmtpEmail)
}
