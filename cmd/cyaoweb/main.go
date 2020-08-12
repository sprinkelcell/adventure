package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/sprinkelcell/adventure"
)

func main() {
	fileName := flag.String("file", "gopher.json", "The Json file name")
	port := flag.Int("port", 3000, "Port to start the web app")
	flag.Parse()

	file, err := os.Open(*fileName)

	if err != nil {
		panic("line no 19")
	}

	story, err := adventure.JsonParser(file)
	if err != nil {
		panic(err)
	}
	//tpl1 := template.Must(template.New("").Parse("Hello"))

	handler := adventure.NewHandler(story)
	addr := fmt.Sprintf(":%d", *port)
	fmt.Printf("Starting the web app on port %s \n", addr)

	log.Fatal(http.ListenAndServe(addr, handler))
}
