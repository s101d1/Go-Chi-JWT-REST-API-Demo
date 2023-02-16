package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

type Article struct {
	Id      string `json:"id"`
	Title   string `json:"title"`
	Desc    string `json:"desc"`
	Content string `json:"content"`
}

type Login struct {
	UserId   string `json:"userId"`
	Password string `json:"password"`
}

const secretKey = "secret"

var articles []Article

var tokenAuth *jwtauth.JWTAuth

func homePage(w http.ResponseWriter, r *http.Request) {
	respondwithJSON(w, http.StatusOK, map[string]string{"message": "Welcome to home page"})
}

func login(w http.ResponseWriter, r *http.Request) {
	var login Login
	json.NewDecoder(r.Body).Decode(&login)

	// In real life you would query the login data against database, but for sake of simplicity we just use fixed user id and password
	if login.UserId == "admin" && login.Password == "123456" {
		// The token duration
		var duration = 10 * time.Minute

		// Create the claims with expiry by setting the "exp" field
		// Remove the "exp" part if you don't want the token to have expiration date
		var claims = map[string]interface{}{"id": login.UserId, "exp": time.Now().UTC().Unix() + int64(duration.Seconds())}
		// Or alternatively you can set the expiry with jwtauth.SetExpiryIn helper function:
		//var claims = map[string]interface{}{"id": login.UserId}
		//jwtauth.SetExpiryIn(claims, duration)

		// Create the token and return it in the response
		_, tokenString, _ := tokenAuth.Encode(claims)
		respondwithJSON(w, http.StatusOK, map[string]string{"token": tokenString})
		return
	}
	respondWithError(w, http.StatusUnauthorized, "Invalid user ID or password")
}

func handleAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, _, err := jwtauth.FromContext(r.Context())

		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Not authorized")
			return
		}

		if token == nil || jwt.Validate(token) != nil {
			respondWithError(w, http.StatusUnauthorized, "Not authorized")
			return
		}

		// Token is authenticated, pass it through
		next.ServeHTTP(w, r)
	})
}

func getAllArticles(w http.ResponseWriter, r *http.Request) {
	respondwithJSON(w, http.StatusOK, articles)
}

func getSingleArticle(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	for _, article := range articles {
		if article.Id == id {
			respondwithJSON(w, http.StatusOK, article)
			return
		}
	}

	respondWithError(w, http.StatusNotFound, "Article is not found")
}

func createNewArticle(w http.ResponseWriter, r *http.Request) {
	var article Article
	json.NewDecoder(r.Body).Decode(&article)

	article.Id = strconv.Itoa(len(articles) + 1)
	articles = append(articles, article)
	respondwithJSON(w, http.StatusCreated, article)
	//respondwithJSON(w, http.StatusCreated, map[string]string{"message": "New Article is created successfully"})
}

func updateArticle(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var updatedArticle Article
	json.NewDecoder(r.Body).Decode(&updatedArticle)

	for index, article := range articles {
		if article.Id == id {
			article.Title = updatedArticle.Title
			article.Desc = updatedArticle.Desc
			article.Content = updatedArticle.Content
			articles[index] = article
			respondwithJSON(w, http.StatusOK, map[string]string{"message": "Article is updated successfully"})
			return
		}
	}

	respondWithError(w, http.StatusNotFound, "Article is not found")
}

func deleteArticle(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	for index, article := range articles {
		if article.Id == id {
			articles = append(articles[:index], articles[index+1:]...)
			respondwithJSON(w, http.StatusOK, map[string]string{"message": "Article is deleted successfully"})
			return
		}
	}

	respondWithError(w, http.StatusNotFound, "Article is not found")
}

func catch(err error) {
	if err != nil {
		panic(err)
	}
}

// respondwithError return error message
func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondwithJSON(w, code, map[string]string{"message": msg})
}

// respondwithJSON write json response format
func respondwithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	//fmt.Println(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func init() {
	// Initialize articles data
	articles = []Article{
		{Id: "1", Title: "Hello 1", Desc: "Article Description 1", Content: "Article Content 1"},
		{Id: "2", Title: "Hello 2", Desc: "Article Description 2", Content: "Article Content 2"},
		{Id: "3", Title: "Hello 3", Desc: "Article Description 3", Content: "Article Content 3"},
		{Id: "4", Title: "Hello 4", Desc: "Article Description 4", Content: "Article Content 4"},
		{Id: "5", Title: "Hello 5", Desc: "Article Description 5", Content: "Article Content 5"},
	}

	tokenAuth = jwtauth.New("HS256", []byte(secretKey), nil)
}

func router() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", homePage)
	r.Post("/login", login)

	// Protected routes
	r.Group(func(r chi.Router) {
		// Verify and validate JWT token
		r.Use(jwtauth.Verifier(tokenAuth))
		// Handle valid/invalid JWT token
		r.Use(handleAuthentication)

		r.Route("/articles", func(r chi.Router) {
			r.Get("/", getAllArticles)
			r.Get("/{id}", getSingleArticle)
			r.Post("/", createNewArticle)
			r.Put("/{id}", updateArticle)
			r.Delete("/{id}", deleteArticle)
		})
	})

	return r
}

func main() {
	addr := ":3000"
	fmt.Printf("Starting server on %v\n", addr)
	log.Fatal(http.ListenAndServe(addr, router()))
}
