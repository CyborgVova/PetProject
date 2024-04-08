package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"jwt/types"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	jwtt "github.com/form3tech-oss/jwt-go"
	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
)

type Store interface {
	GetPlaces(limit int, offset int) ([]types.Place, int, error)
}

type Implement struct{}

const limit = 10

func (Implement) GetPlaces(limit int, offset int) ([]types.Place, int, error) {
	var total int
	places := []types.Place{}
	es := types.CreateEsClient()
	count, err := es.Count(es.Count.WithBody(strings.NewReader(`{
		"query": {
		  "match": {
			"_index": "places"
		  }
		}
	  }`)))
	if err != nil {
		log.Fatal(err)
	}
	var rr map[string]interface{}
	if err := json.NewDecoder(count.Body).Decode(&rr); err == nil {
		total = int(rr["count"].(float64))
	} else {
		log.Printf("Error parsing the response body of count documents: %s", err)
	}
	if offset < 0 || offset > total {
		return places, total, errors.New("illegal value of offset")
	}
	res, err := es.Search(es.Search.WithBody(strings.NewReader(fmt.Sprintf(`{
		"size":%d,
		"from":%d,
		"query": {
		  "match": {
			"_index": "places"
		  }
		},
		"sort": [
		{
			"_score": "desc"
		},
		{
			"id": "asc"
		}
		]
	  }`, limit, offset))),
		es.Search.WithPretty())
	if err != nil {
		log.Fatal(err)
	}
	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err == nil {
		for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
			source := hit.(map[string]interface{})["_source"]
			v := reflect.ValueOf(source).MapRange()
			place := types.Place{}
			for v.Next() {
				switch key := v.Key().String(); key {
				case "id":
					id, err := strconv.ParseFloat(fmt.Sprintf("%v", v.Value()), 64)
					if err != nil {
						log.Fatal("unable to get ID", err)
					}
					place.Id = int(id)
				case "name":
					place.Name = fmt.Sprintf("%s", v.Value())
				case "address":
					place.Address = fmt.Sprintf("%s", v.Value())
				case "phone":
					place.Phone = fmt.Sprintf("%s", v.Value())
				case "location":
					place.Location.Latitude = v.Value().Interface().(map[string]interface{})["lat"].(float64)
					place.Location.Longitude = v.Value().Interface().(map[string]interface{})["lon"].(float64)
				}
			}
			places = append(places, place)
		}

	} else {
		log.Printf("Error parsing the response body: %s", err)
	}
	return places, total, nil
}

func head(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	numpage, err := strconv.Atoi(page)
	if err == nil {
		temp, _ := template.ParseFiles("index.html")
		places, total, err := Implement{}.GetPlaces(limit, numpage*limit)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid 'page' value: '%d'", numpage), http.StatusBadRequest)
		} else {
			temp.ExecuteTemplate(w, "index", struct {
				Places []types.Place
				Total  int
				Prev   int
				Next   int
			}{places, total, numpage - 1, numpage + 1})
		}
	} else {
		http.Error(w, fmt.Sprintf("Invalid 'page' value: '%v'", page), http.StatusBadRequest)
	}
}

func api(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	page := r.URL.Query().Get("page")
	numpage, err := strconv.Atoi(page)
	if err == nil {
		places, total, err := Implement{}.GetPlaces(limit, numpage*limit)
		if err == nil {
			prev := numpage - 1
			next := numpage + 1
			b, err := json.Marshal(struct {
				Name   string        `json:"name"`
				Total  int           `json:"total"`
				Places []types.Place `json:"places"`
				Prev   int           `json:"prev_page"`
				Next   int           `json:"next_page"`
				Last   int           `json:"last_page"`
			}{"Places", total, places, prev, next, total / 10})
			if err != nil {
				log.Fatal("unable to create json", err)
			}
			var out bytes.Buffer
			json.Indent(&out, b, "", "    ")
			w.Write(out.Bytes())
		} else {
			b, err := json.Marshal(struct {
				Error string `json:"error"`
			}{
				fmt.Sprintf("Invalid 'page' value: '%d'", numpage),
			})
			if err != nil {
				log.Fatal("Bad marshalling 'error'", err)
			}
			var out bytes.Buffer
			json.Indent(&out, b, "", "    ")
			w.Write(out.Bytes())
		}
	} else {
		b, err := json.Marshal(struct {
			Error string `json:"error"`
		}{
			fmt.Sprintf("Invalid 'page' value: '%v'", page),
		})
		if err != nil {
			log.Fatal("Bad marshalling 'error'", err)
		}
		var out bytes.Buffer
		json.Indent(&out, b, "", "    ")
		w.Write(out.Bytes())
	}

}

func GetPlacesR(lat, lon float64) []types.Place {
	places := []types.Place{}
	es := types.CreateEsClient()
	res, err := es.Search(es.Search.WithBody(strings.NewReader(fmt.Sprintf(`{
		"size":3,
		"query": {
		  "match": {
			"_index": "places"
		  }
		},
		"sort": [
			{
			  "_geo_distance": {
				"location": {
				  "lat": %v,
				  "lon": %v
				},
				"order": "asc",
				"unit": "km",
				"mode": "min",
				"distance_type": "arc",
				"ignore_unmapped": true
			  }
			}
		]
		
	  }`, lat, lon))),
		es.Search.WithPretty())
	if err != nil {
		log.Fatal(err)
	}
	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err == nil {
		for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
			source := hit.(map[string]interface{})["_source"]
			v := reflect.ValueOf(source).MapRange()
			place := types.Place{}
			for v.Next() {
				switch key := v.Key().String(); key {
				case "id":
					id, err := strconv.ParseFloat(fmt.Sprintf("%v", v.Value()), 64)
					if err != nil {
						log.Fatal("unable to get ID", err)
					}
					place.Id = int(id)
				case "name":
					place.Name = fmt.Sprintf("%s", v.Value())
				case "address":
					place.Address = fmt.Sprintf("%s", v.Value())
				case "phone":
					place.Phone = fmt.Sprintf("%s", v.Value())
				case "location":
					place.Location.Latitude = v.Value().Interface().(map[string]interface{})["lat"].(float64)
					place.Location.Longitude = v.Value().Interface().(map[string]interface{})["lon"].(float64)
				}
			}
			places = append(places, place)
		}

	} else {
		log.Printf("Error parsing the response body: %s", err)
	}
	return places
}

var recommend = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	lat := r.URL.Query().Get("lat")
	lon := r.URL.Query().Get("lon")
	numlat, err := strconv.ParseFloat(lat, 64)
	if err != nil || numlat > 90 || numlat < -90 {
		b, err := json.Marshal(struct {
			Error string `json:"error"`
		}{
			fmt.Sprintf("Invalid 'lat' value: '%v'", numlat),
		})
		if err != nil {
			log.Fatal("Bad marshalling 'error'", err)
		}
		var out bytes.Buffer
		json.Indent(&out, b, "", "    ")
		w.Write(out.Bytes())
		return
	}
	numlon, err := strconv.ParseFloat(lon, 64)
	if err != nil || numlon > 180 || numlon < -180 {
		b, err := json.Marshal(struct {
			Error string `json:"error"`
		}{
			fmt.Sprintf("Invalid 'lon' value: '%v'", numlon),
		})
		if err != nil {
			log.Fatal("Bad marshalling 'error'", err)
		}
		var out bytes.Buffer
		json.Indent(&out, b, "", "    ")
		w.Write(out.Bytes())
		return
	}
	places := GetPlacesR(numlat, numlon)
	b, err := json.Marshal(struct {
		Name   string        `json:"name"`
		Places []types.Place `json:"places"`
	}{"Recommendation", places})
	if err != nil {
		log.Fatal("Recommendation ", err)
	}
	var out bytes.Buffer
	json.Indent(&out, b, "", "    ")
	w.Write(out.Bytes())
})

func getToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["username"] = "delilahl"
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()
	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		log.Fatal(err)
	}
	b, err := json.Marshal(struct {
		Token string `json:"token"`
	}{tokenString})
	if err != nil {
		log.Fatal("Token ", err)
	}
	var out bytes.Buffer
	json.Indent(&out, b, "", "    ")
	w.Write(out.Bytes())
}

var jwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
	ValidationKeyGetter: func(token *jwtt.Token) (interface{}, error) {
		return []byte("secret"), nil
	},
	SigningMethod: jwt.SigningMethodHS256,
})

func Handlers() {
	router := mux.NewRouter()
	http.Handle("/", router)
	router.HandleFunc("/", head).Methods("GET")
	router.HandleFunc("/api/places", api).Methods("GET")
	router.Handle("/api/recommend", jwtMiddleware.Handler(recommend)).Methods("GET")
	router.HandleFunc("/api/get_token", getToken).Methods("GET")
}

func main() {
	Handlers()
	if err := http.ListenAndServe(":8888", nil); err != nil {
		log.Fatal(err)
	}
}
