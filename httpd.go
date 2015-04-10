package main

import (
	"fmt"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"log"
)

type Httpd struct {
	Host string
	Port int
}

type Param struct {
	Channel    string `form:"channel" binding:"required"`
	Message    string `form:"message"`
	Name       string `form:"name"`
	Icon       string `form:"icon"`
	Color      string `form:"color"`
	Pretext    string `form:"pretext"`
	AuthorName string `form:"author_name"`
	AuthorLink string `form:"author_link"`
	AuthorIcon string `form:"author_icon"`
	Title      string `form:"title"`
	TitleLink  string `form:"title_link"`
	Text       string `form:"text"`
	ImageURL   string `form:"image_url"`
	Parse      string `form:"parse"`
	Manual     bool   `form:"manual"`
}

func NewHttpd(host string, port int) *Httpd {
	return &Httpd{
		Host: host,
		Port: port,
	}
}

func (h *Httpd) Run() {
	m := martini.Classic()
	m.Get("/", func() string { return "Hello, I'm Takosan!!1" })
	m.Post("/notice", binding.Bind(Param{}), messageHandler)
	m.Post("/privmsg", binding.Bind(Param{}), messageHandler)
	m.RunOnAddr(fmt.Sprintf("%s:%d", h.Host, h.Port))
}

func messageHandler(p Param) (int, string) {
	ch := make(chan error, 1)

	newMessage := NewMessage(p.Channel, p.Message, ch)
	newMessage.Name = p.Name
	newMessage.Icon = p.Icon
	newMessage.Color = p.Color
	newMessage.Pretext = p.Pretext
	newMessage.AuthorName = p.AuthorName
	newMessage.AuthorLink = p.AuthorLink
	newMessage.AuthorIcon = p.AuthorIcon
	newMessage.Title = p.Title
	newMessage.TitleLink = p.TitleLink
	newMessage.Text = p.Text
	newMessage.ImageURL = p.ImageURL
	newMessage.Parse = p.Parse
	newMessage.Manual = p.Manual

	go MessageBus.Publish(newMessage)
	err := <-ch

	if err != nil {
		message := fmt.Sprintf("Failed to send message to %s: %s\n", p.Channel, err)
		log.Printf(fmt.Sprintf("[error] %s", message))
		return 400, message
	} else {
		return 200, fmt.Sprintf("Message sent successfully to %s", p.Channel)
	}
}
