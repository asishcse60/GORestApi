package main

import (
	"example.com/Service/product"
	"example.com/Service/receipt"
	"fmt"
	"log"
	"net/http"
)


const basePath = "/api"

func main() {
	fmt.Println("Go Test!")
	receipt.SetupRoutes(basePath)
	product.SetupRoutes(basePath)
	var err = http.ListenAndServe(":5000", nil)
	if err != nil{
		log.Fatal(err)
	}
}

