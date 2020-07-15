package product

type Product struct {
	ProductId      int    `json:"productId" bson:"productid"`
	Manufacturer   string `json:"manufacturer" bson:"manufacturer"`
	Sku            string `json:"sku" bson:"sku"`
	Upc            string `json:"upc" bson:"upc"`
	PricePerUnit   string `json:"pricePerUnit" bson:"priceperunit"`
	QuantityOnHand int    `json:"quantityOnHand" bson:"quantityonhand"`
	ProductName    string `json:"productName" bson:"productname"`
}
