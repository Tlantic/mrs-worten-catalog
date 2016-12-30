package main

import (
	"github.com/algolia/algoliasearch-client-go/algoliasearch"

	"fmt"
	"net/http"
	"github.com/gorilla/mux"
	"encoding/json"
	"log"
	"io/ioutil"
	"github.com/gorilla/handlers"
	"strconv"
)

type Item struct {
	ID string `json:"id"`
	Ean string `json:"ean"`
	Status string `json:"status"`
	Name string `json:"name"`
	Brand string `json:"brand"`
	Price float64 `json:"price"`
	Image string `json:"image"`
	Images Images `json:"images"`
	Details []interface{} `json:"details"`
	ConditioningState string `json:"conditioningState"`
	ConditioningEnvironment string `json:"conditioningEnvironment"`
	ConditioningEXP string `json:"conditioningEXP"`
}

type Images struct {
	Ldpi string `json:"ldpi"`
	Mdpi string `json:"mdpi"`
	Hdpi string `json:"hdpi"`
	Xhdpi string `json:"xhdpi"`
}

type ClientReq struct {
	Text  string `json:"text"`
}

type Response struct {
	Status string `json:"status"`
	Result []Item `json:"result"`
}

var index algoliasearch.Index

func main(){

	client := algoliasearch.NewClient("4JF4HATVV9", "d27cc26c973daa4521834f26d391a137")

	index = client.InitIndex("wo-pt-prd")

	r := mux.NewRouter()
	r.HandleFunc("/worten-catalog/products/search", YourHandler).Methods("POST")

	log.Fatal(http.ListenAndServe(":8000", handlers.CORS()(r)))

}


func YourHandler(w http.ResponseWriter, r *http.Request) {

	body, _ := ioutil.ReadAll(r.Body)

	var clientReq ClientReq

	fromParam := r.URL.Query().Get("from")
	var size = 0

	if len(fromParam) > 0{
		size, _ := strconv.Atoi(fromParam)
		size = size/20
	}else{
		size = 0
	}

	json.Unmarshal(body, &clientReq)



	params := algoliasearch.Map{

		"hitsPerPage":          20,
		"page": size,
	}
	res, _ := index.Search(clientReq.Text, params)

	data := []Item{}

	for _, v := range res.Hits{
		x := Item{
			ID:  fmt.Sprintf("%f",v["sku"]),
			Name: v["name"].(string),
			Ean: v["ean"].(string),
			Status: "A",
			Brand: v["brand"].(string),
			Price: v["price"].(float64),
			Image: fmt.Sprintf("http://www.worten.pt%s",v["image_default"].(string)),
			Images: Images{
				fmt.Sprintf("http://www.worten.pt%s",v["image_thumbnail"].(string)),
				fmt.Sprintf("http://www.worten.pt%s",v["image_default"].(string)),
				fmt.Sprintf("http://www.worten.pt%s",v["image_zoom"].(string)),
				fmt.Sprintf("http://www.worten.pt%s",v["image_zoom"].(string)),
			},
		}

		data = append(data, x)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{
		"ok",
		data,
	})
}