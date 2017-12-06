package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"
	//	"github.com/xdrive/goblueprints/trace"
	//	"os"
	"github.com/spf13/viper"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/google"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("../templates", t.filename)))
	})

	t.templ.Execute(w, r)
}

func init() {
	viper.AddConfigPath("config/")
	viper.AddConfigPath("../config/")
	viper.SetConfigName("config")

	err := viper.ReadInConfig()
	if err != nil {
		panic("Config file not found...")
	}
}

func main() {
	address := viper.GetString("app.host") + ":" + viper.GetString("app.port")
	fmt.Println(callbackURL())

	gomniauth.SetSecurityKey(viper.GetString("auth.secret_key"))
	gomniauth.WithProviders(
		google.New(
			viper.GetString("auth.google.client_id"),
			viper.GetString("auth.google.client_secret"),
			callbackURL(),
		),
	)

	r := newRoom()
	//r.tracer = trace.New(os.Stdout)
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/room", r)
	go r.run()

	// start the server
	log.Println("Starting the app on", address)
	if err := http.ListenAndServe(address, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func callbackURL() string {
	host := viper.GetString("app.host")
	if host == "" {
		host = "localhost"
	}

	return fmt.Sprintf("http://%s:%s%s", host, viper.GetString("app.port"), viper.GetString("auth.google.callback_url"))
}
