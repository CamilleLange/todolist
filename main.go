package main

import (
	"fmt"
	"net/http"
	"time"
)

func helloworld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World! Time : %s", time.Now().Format(time.RFC3339Nano))
}

func main() {
	http.HandleFunc("/", helloworld)

	fmt.Println("API ready !")
	fmt.Println("Listening at : http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
