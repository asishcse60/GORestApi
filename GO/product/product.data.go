package product

import (
	"context"
	"encoding/json"
	"example.com/Service/database"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"sort"
	"strings"
	"sync"
)

var productMap = struct {
	sync.Mutex
	m map[int]Product
}{m: make(map[int]Product)}

func init() {
	fmt.Println("Loading Product data....")
	prodMap, err := loadingProductMap()
	productMap.m =  prodMap
	if err != nil{
		log.Fatal(err)
	}
	fmt.Printf("%d products loaded...\n", len(productMap.m))
}

func loadingProductMap() (map[int]Product, error) {

	fileName := "products.json"
	_,err:= os.Stat(fileName)
if os.IsNotExist(err){
	return nil, fmt.Errorf("file [%s] does not exist", fileName)
}

	file,_ := ioutil.ReadFile(fileName)
	productList := make([]Product, 0)

	err = json.Unmarshal([]byte(file), &productList)
	if err != nil{
		log.Fatal(err)
	}
	prodMap := make(map[int]Product)
	//Get the database Context
	ClientContext, _ := database.GetDatabaseClientStorage()
	collection := getCollectionDB(ClientContext, productList[0])
	var documents = GetProductDocuments(collection)

	if len(documents) > 0{
		bsonBytes, _ := bson.Marshal(documents)
		bson.Unmarshal(bsonBytes, &productList)
		for i:=0;i<len(productList);i++ {
			prodMap[productList[i].ProductId] = productList[i]
		}
	} else{
		var multipleProducts []interface{}

		for i:=0;i<len(productList);i++ {
			prodMap[productList[i].ProductId] = productList[i]
			multipleProducts=append(multipleProducts, productList[i])
		}
		_, err2 := collection.InsertMany(context.Background(), multipleProducts)
		if err2 != nil{
			log.Print("Multiple documents Insert failed...")
		}
		fmt.Printf("%d Insert Success!\n", len(productList))
	}

	return prodMap, nil
}

func searchForProductData(report ReportFilter) ([]Product, error) {
	andQuery := []bson.M{}
	if len(report.NameFilter) > 0{
		likeValue := strings.ToLower(report.NameFilter)+".*"
		nameFilter := bson.M{"productname": primitive.Regex{likeValue, ""}} //
		andQuery = append(andQuery, nameFilter)
	}
	if len(report.ManufactureFilter) > 0{
		likeValue := report.ManufactureFilter+".*"
        manuFilter := bson.M{"manufacturer": primitive.Regex{likeValue, ""}}
		andQuery = append(andQuery, manuFilter)
	}
	if len(report.SKUFilter) > 0{
		likeValue := report.SKUFilter+".*"
		skuFilter := bson.M{"sku": primitive.Regex{likeValue, ""}}
		andQuery = append(andQuery, skuFilter)
	}
	var product Product
	var products [] Product
	ClientContext, _ := database.GetDatabaseClientStorage()
	collection := getCollectionDB(ClientContext, product)
	filter := bson.M{"$and": andQuery}
	cursorDoc, err := collection.Find(context.TODO(), filter)
	if err != nil{
		log.Print(err)
		return nil, err
	}
	err = cursorDoc.All(context.TODO(), &products)
	if err != nil{
		log.Print(err)
		return nil, err
	}
	return products, nil
}

func getProduct(productId int)	*Product  {
	productMap.Lock()
	defer productMap.Unlock()
	ClientContext, _ := database.GetDatabaseClientStorage()
	var product Product
	collection := getCollectionDB(ClientContext, product)

	//var objID, _ = primitive.ObjectIDFromHex("5ef7a7ed34c50e6556dcf7e4")
     filter:=bson.M{"productid": productId}

	 err := collection.FindOne(context.TODO(), filter).Decode(&product)
	if err != nil{
		log.Printf("ProductId %d Item is not found", productId)
	}
	fmt.Println(product)
	return &product
	/*
	if product, ok:=productMap.m[productId]; ok{
			return &product
		}
		return nil
	*/
}
func GetTopTenProducts() ([] Product, error) {
	var product Product
	ClientContext, _ := database.GetDatabaseClientStorage()
	collection := getCollectionDB(ClientContext, product)
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"quantityonhand", -1}})
	findOptions.SetLimit(10)
	cursor, err := collection.Find(context.TODO(), bson.D{}, findOptions)
	if err != nil{
		fmt.Println("Top 10 document is not found")
		return nil, err
	}

	var multipleProducts []Product

	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		if err = cursor.Decode(&product); err != nil {
			log.Fatal(err)
		}
		multipleProducts = append(multipleProducts, product)
	}
	return multipleProducts, nil

}

func removeProduct(productId int){
	productMap.Lock()
	defer productMap.Unlock()
	var product Product
	ClientContext, _ := database.GetDatabaseClientStorage()
	collection := getCollectionDB(ClientContext, product)
	filter:=bson.M{"productid": productId}
	deleteResult, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil{
		fmt.Println("Document remove in unsuccess!")
		return
	}
	fmt.Println(deleteResult)
	//delete(productMap.m, productId)
}

func getProductList() [] Product {
	//productMap.Lock()
	ClientContext, _ := database.GetDatabaseClientStorage()
	var product Product
	collection := getCollectionDB(ClientContext, product)
	var documents = GetProductDocuments(collection)
	if len(documents) > 0{
		var productList [] Product
		bsonBytes, _ := json.Marshal(documents)
		json.Unmarshal(bsonBytes, &productList)

		return productList
	}else {
		products := make([]Product, 0, len(productMap.m))
		for _, value := range productMap.m	{
			products = append(products, value)
		}
	//	productMap.Unlock()
		return products
	}
}

func getProductIds() [] int {
	//productMap.Lock()
	var productIds []int
	for key:= range productMap.m{
		productIds = append(productIds, key)
	}
	//productMap.Unlock()
	sort.Ints(productIds)
	return productIds
}

func getNextProductId() int {
	productsId := getProductIds()
	return productsId[len(productsId) - 1] + 1
}

func addOrUpdateProduct(product Product) (int, error)  {
	// if the product id is set, update, otherwise add
	addOrUpdateID := -1
	if product.ProductId > 0{
		oldProduct := getProduct(product.ProductId)
		if oldProduct == nil{
			return 0, fmt.Errorf("product id [%d] doesn't exist", product.ProductId)
		}
		addOrUpdateID = oldProduct.ProductId
	}else{
		addOrUpdateID =  getNextProductId()
		product.ProductId = addOrUpdateID
	}

	//productMap.Lock()
	ClientContext, _ := database.GetDatabaseClientStorage()
	collection := getCollectionDB(ClientContext, product)
	_, err := collection.InsertOne(context.TODO(), product)
	if err != nil{
		fmt.Println("Add or Update failed...")
		return addOrUpdateID, err
	}
	fmt.Printf("New or updated document id is %d\n", addOrUpdateID)
	productMap.m[addOrUpdateID] = product
	//productMap.Unlock()
	return addOrUpdateID, nil
}

func getType(genericObject interface{}) string {
	return reflect.TypeOf(genericObject).Name()
}
//data base utility called for function
func GetProductDocuments(collection *mongo.Collection) []bson.M {
	/*
		var m bson.M
		var s Struct1

		// convert m to s
		bsonBytes, _ := bson.Marshal(m)
	    bson.Unmarshal(bsonBytes, &s)
			or
		bsonBytes, _ := json.Marshal(m)
	    json.Unmarshal(bsonBytes, &s)

	*/

	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	var multipleProducts []bson.M

	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var product bson.M
		if err = cursor.Decode(&product); err != nil {
			log.Fatal(err)
		}
		multipleProducts = append(multipleProducts, product)
	}
	return multipleProducts
}

func getCollectionDB(dbClientStorage *database.MongoDBClientStorage, genericObject interface{})  *mongo.Collection {
	collectionName := getType(genericObject)
	dataBaseContext := dbClientStorage.Client.Database(dbClientStorage.DatabaseName)
	collection := dataBaseContext.Collection(collectionName)
	return collection
}
// end of database utility