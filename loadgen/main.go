package main

import (
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
)

const frontend = "http://localhost:8080"

func main() {
	for {
		log.Println("requesting index")
		resp, err := http.Get(frontend + "/")
		if err != nil {
			log.Println("ERROR:", err)
			continue
		}
		_, _ = io.ReadAll(resp.Body)
		num := getInt(5, 500000)
		log.Println("sending form with number =", num)
		resp, err = http.PostForm(frontend+"/submit", url.Values{"number": []string{strconv.Itoa(num)}})
		if err != nil {
			log.Println("ERROR:", err)
			continue
		}
		_, _ = io.ReadAll(resp.Body)
	}
}

func getInt(min, max int) int {
	mean := float64(min) + float64(max)/2
	dev := float64(max) - mean
	for {
		candidate := int(rand.NormFloat64()*dev + mean)
		if candidate >= min && candidate <= max {
			return candidate
		}
	}
}
