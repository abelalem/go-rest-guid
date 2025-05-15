package main

import (
	"encoding/json"
	"net/http"

	"github.com/abelalem/go-rest-guid/pkg/recipes"
	"github.com/gorilla/mux"
	"github.com/gosimple/slug"
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

func InternalServerErrorHandler(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusInternalServerError)
	rw.Write([]byte("500 Internal Server Error"))
}

func NotFoundHandler(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusNotFound)
	rw.Write([]byte("404 Not Found"))
}

type RecipesHandler struct {
	store recipeStore
}

func NewRecipesHandler(s recipeStore) *RecipesHandler {
	return &RecipesHandler{
		store: s,
	}
}

func (h RecipesHandler) CreateRecipe(rw http.ResponseWriter, r *http.Request) {
	// Recipe object that will be populated from json payload
	var recipe recipes.Recipe

	if err := json.NewDecoder(r.Body).Decode(&recipe); err != nil {
		InternalServerErrorHandler(rw, r)
		return
	}

	resourceID := slug.Make(recipe.Name)

	if err := h.store.Add(resourceID, recipe); err != nil {
		InternalServerErrorHandler(rw, r)
		return
	}

	rw.WriteHeader(http.StatusOK)
}
func (h RecipesHandler) ListRecipes(rw http.ResponseWriter, r *http.Request) {
	recipesList, err := h.store.List()
	if err != nil {
		InternalServerErrorHandler(rw, r)
		return
	}

	jsonBytes, err := json.Marshal(recipesList)
	if err != nil {
		InternalServerErrorHandler(rw, r)
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write(jsonBytes)
}
func (h RecipesHandler) GetRecipe(rw http.ResponseWriter, r *http.Request) {
	// Get recipeId from the request
	recipeId := mux.Vars(r)["id"]

	recipe, err := h.store.Get(recipeId)
	if err != nil {
		if err == recipes.NotFoundErr {
			NotFoundHandler(rw, r)
			return
		}

		InternalServerErrorHandler(rw, r)
		return
	}

	jsonBytes, err := json.Marshal(recipe)
	if err != nil {
		InternalServerErrorHandler(rw, r)
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write(jsonBytes)
}
func (h RecipesHandler) UpdateRecipe(rw http.ResponseWriter, r *http.Request) {
	// Get recipeId from the request
	recipeId := mux.Vars(r)["id"]

	// Recipe object that will be populated from JSON payload
	var recipe recipes.Recipe
	if err := json.NewDecoder(r.Body).Decode(&recipe); err != nil {
		InternalServerErrorHandler(rw, r)
		return
	}

	if err := h.store.Update(recipeId, recipe); err != nil {
		if err == recipes.NotFoundErr {
			NotFoundHandler(rw, r)
			return
		}
		InternalServerErrorHandler(rw, r)
		return
	}

	rw.WriteHeader(http.StatusOK)
}
func (h RecipesHandler) DeleteRecipe(rw http.ResponseWriter, r *http.Request) {
	// Get recipeId from the request
	recipeId := mux.Vars(r)["id"]

	if err := h.store.Remove(recipeId); err != nil {
		InternalServerErrorHandler(rw, r)
		return
	}

	rw.WriteHeader(http.StatusOK)
}

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
