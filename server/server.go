package server

import (
	"D/Avito/cache"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

//Server ...
type Server struct {
	config *Config
	logger *logrus.Logger
	router *mux.Router
	cache  *cache.Cache
}

//New ...
func New(config *Config) *Server {
	return &Server{
		config: config,
		logger: logrus.New(),
		router: mux.NewRouter(),
		cache:  cache.New(5*time.Minute, 10*time.Minute),
	}
}

//Start ...
func (s *Server) Start() error {
	s.configureRouter()
	if err := s.configureLogger(); err != nil {
		return err
	}
	s.logger.Info("starting server")
	return http.ListenAndServe(s.config.BindAddr, s.router)
}

func (s *Server) configureLogger() error {
	level, err := logrus.ParseLevel(s.config.LogLevel)
	if err != nil {
		return err
	}
	s.logger.SetLevel(level)
	return nil
}

func (s *Server) configureRouter() {
	s.router.HandleFunc("/map", s.Get)
	s.router.HandleFunc("/set", s.Set)
	s.router.HandleFunc("/delete", s.Delete)
	s.router.HandleFunc("/keys", s.Keys)
	s.router.HandleFunc("/save", s.Save)
}

//Set ...
func (s *Server) Set(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("Key")
	value := r.FormValue("Value")
	duration, err := time.ParseDuration(r.FormValue("Duration"))
	if err != nil {
		log.Fatal(err)
	}
	var expiration int64
	// Если продолжительность жизни равна 0 - используется значение по-умолчанию
	if duration == 0 {
		duration = s.cache.DefaultExpiration
	}

	// Устанавливаем время истечения кеша
	if duration > 0 {
		expiration = time.Now().Add(duration).UnixNano()
	}

	s.cache.Lock()

	defer s.cache.Unlock()

	s.cache.Items[key] = cache.Item{
		Value:      value,
		Expiration: expiration,
		Created:    time.Now(),
	}
	http.Redirect(w, r, "/map", 301)

}

//Get ...
func (s *Server) Get(w http.ResponseWriter, r *http.Request) {
	s.cache.RLock()
	defer s.cache.RUnlock()
	tmpl, err := template.ParseFiles("static/map.html")
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	if err := tmpl.Execute(w, s.cache.Items); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
}

//Delete ...
func (s *Server) Delete(w http.ResponseWriter, r *http.Request) {
	s.cache.Lock()
	defer s.cache.Unlock()
	key := r.FormValue("Key1")
	delete(s.cache.Items, key)
	http.Redirect(w, r, "/map", 301)
}

//Keys ...
func (s *Server) Keys(w http.ResponseWriter, r *http.Request) {
	regex := r.FormValue("Key2")
	regexpr, err := regexp.Compile(regex)
	if err != nil {
		log.Fatal(err)
		return
	}
	tmpl, err := template.ParseFiles("static/keys.html")
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	for key := range s.cache.Items {
		if regexpr.MatchString(key) {
			if err := tmpl.Execute(w, key); err != nil {
				http.Error(w, err.Error(), 400)
				return
			}
		}
	}
}

func (s *Server) Save(w http.ResponseWriter, r *http.Request) {
	b, err := json.Marshal(s.cache.Items)
	if err != nil {
		fmt.Println("error:", err)
	}
	myjson := []byte(b)
	error := ioutil.WriteFile("Cache.txt", myjson, 0777)
	if error != nil {
		fmt.Println(error)
	}
	http.Redirect(w, r, "/map", 301)
}
