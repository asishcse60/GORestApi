package product

import (
	"encoding/json"
	"example.com/Service/cors"
	"fmt"
	"golang.org/x/net/websocket"
	"net/http"
	"strconv"
	"strings"
)
const productsPath = "products"
func SetupRoutes(apiBasePath string) {

	productsHandler := http.HandlerFunc(handleProducts)
	productHandler := http.HandlerFunc(handleProduct)
	reportHandler := http.HandlerFunc(handleProductReport)
	http.Handle("/websocket", websocket.Handler(productSocket))
	http.Handle(fmt.Sprintf("%s/%s", apiBasePath, productsPath), cors.Middleware(productsHandler))
	http.Handle(fmt.Sprintf("%s/%s/", apiBasePath, productsPath), cors.Middleware(productHandler))
	http.Handle(fmt.Sprintf("%s/%s/reports", apiBasePath, productsPath), cors.Middleware(reportHandler))
}


func handleProduct(writer http.ResponseWriter, request *http.Request) {
	urlPathSegments := strings.Split(request.URL.Path, "products/")
	productID, err:= strconv.Atoi(urlPathSegments[len(urlPathSegments) - 1])
	if err != nil{
		writer.WriteHeader(http.StatusNotFound)
		return
	}
	product := getProduct(productID)

	if product == nil{
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	switch request.Method {

	case http.MethodGet:
		productsJson, err := json.Marshal(product)
		if err != nil{
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, err = writer.Write(productsJson)
		if err != nil{
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		writer.WriteHeader(http.StatusOK)

	case http.MethodPut:
		var updateProduct Product
		err := json.NewDecoder(request.Body).Decode(&updateProduct)

		if err != nil{
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
		if updateProduct.ProductId != productID {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
		_, err = addOrUpdateProduct(updateProduct)
		if err != nil{
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
		writer.WriteHeader(http.StatusOK)
		return

	case http.MethodDelete:
		removeProduct(productID)

	case http.MethodOptions:
		return

	default:
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}


func handleProducts(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {

	case http.MethodGet:
		productList := getProductList()
		productsJson, err := json.Marshal(productList)
		if err != nil{
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, err = writer.Write(productsJson)
		if err != nil{
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		writer.WriteHeader(http.StatusOK)

	case http.MethodPost:
		var newProduct Product
		err := json.NewDecoder(request.Body).Decode(&newProduct)

		if err != nil{
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
		_,err = addOrUpdateProduct(newProduct)
		if err != nil{
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
		writer.WriteHeader(http.StatusCreated)
		return

	case http.MethodOptions:
		return

	default:
		writer.WriteHeader(http.StatusMethodNotAllowed)
	}
}

