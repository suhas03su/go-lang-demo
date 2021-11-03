package main

import (
	"fmt"
	"log"
	"net/http"

	"github.come/suhas03su/mongoAPI/routers"
)

const port = ":3000"

func main() {
	fmt.Println("MONGO DB API")
	fmt.Println("SERVER IS GETTING STARTED...")

	routers.Router()

	fmt.Println("Welcome to creating APIs with GO")
	log.Fatal(http.ListenAndServe(port, nil))
	fmt.Println("LISTENING AT PORT 3000")
}
