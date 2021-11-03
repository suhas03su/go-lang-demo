package routers

import (
	"net/http"

	"github.come/suhas03su/mongoAPI/controllers"
)

func Router() {
	http.HandleFunc("/api/movies", controllers.GetAllMovies)
	http.HandleFunc("/api/movie", controllers.CreateNewMovie)
	http.HandleFunc("/api/update-movie", controllers.MarkAsWatched)
	http.HandleFunc("/api/delete-all", controllers.DeleteAllMovies)
	http.HandleFunc("/api/delete", controllers.DeleteOneMovie)
}
