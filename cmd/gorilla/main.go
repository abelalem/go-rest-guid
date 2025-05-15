package main

import (
	"net/http"

	"github.com/abelalem/go-rest-guid/pkg/recipes"
	"github.com/gorilla/mux"
)

type homeHandler struct{}

func (h *homeHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	rw.Write([]byte("This is my home page"))
}

type recipeStore interface {
	Add(name string, recipe recipes.Recipe) error
	Get(name string) (recipes.Recipe, error)
	List() (map[string]recipes.Recipe, error)
	Update(name string, recipe recipes.Recipe) error
	Remove(name string) error
}

type RecipesHandler struct {
	store recipeStore
}

func NewRecipesHandler(s recipeStore) *RecipesHandler {
	return &RecipesHandler{
		store: s,
	}
}

func (h RecipesHandler) CreateRecipe(rw http.ResponseWriter, r *http.Request) {}
func (h RecipesHandler) ListRecipes(rw http.ResponseWriter, r *http.Request)  {}
func (h RecipesHandler) GetRecipe(rw http.ResponseWriter, r *http.Request)    {}
func (h RecipesHandler) UpdateRecipe(rw http.ResponseWriter, r *http.Request) {}
func (h RecipesHandler) DeleteRecipe(rw http.ResponseWriter, r *http.Request) {}

func main() {
	// Create the store and Recipe Handler
	store := recipes.NewMemStore()
	recipesHandler := NewRecipesHandler(store)
	home := homeHandler{}

	// Create the router
	router := mux.NewRouter()

	// Register the routes
	router.HandleFunc("/", home.ServeHTTP).Methods("GET")
	router.HandleFunc("/recipes", recipesHandler.ListRecipes).Methods("GET")
	router.HandleFunc("/recipes", recipesHandler.CreateRecipe).Methods("POST")
	router.HandleFunc("/recipes/{id}", recipesHandler.GetRecipe).Methods("GET")
	router.HandleFunc("/recipes/{id}", recipesHandler.UpdateRecipe).Methods("PUT")
	router.HandleFunc("/recipes/{id}", recipesHandler.DeleteRecipe).Methods("DELETE")

	// Start the server
	http.ListenAndServe(":8010", router)
}
