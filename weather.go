package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"github.com/gin-gonic/gin"
)

type ForecastResponse struct { //I create a struct named ForecastResponse where all the data like temp,feels_like is stored under Main object and the Main object is under List array
	List []struct {				// of forecast struct
		Main struct {
			Temp float64 `json:"temp"`
			FeelsLike float64 `json:"feels_like"`
			Humidity int `json:"humidity"`
		} `json:"main"`
		Pop float64 `json:"pop"` //Chances of rain
		DtText string `json:"dt_txt"` // timestamp
	}
}



func fetchWeather(city, apiKey string) (ForecastResponse, error) { //I create a function called fetchWeather to call it in func main(), it has city and apiKey string as input. It returns 
	var forecast ForecastResponse									// ForecastResponse as output and error, I also created a forecast struct of the the ForecastResponse struct type
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/forecast?q=%s&appid=%s&units=metric", city, apiKey)

	resp, err := http.Get(url) // I initialied resp to the incoming data stream path created by the api. I used http.Get(url) method to get the data from the api
	if err != nil { //If error is not empty, it will return this error message
		fmt.Print("Error fetching data")
		return forecast, err //returns forecast where forecast is empty and err is filled
	}
	defer resp.Body.Close() //Close the stream of incoming data to prevent any leaks. This works as the function ends

	body, err := io.ReadAll(resp.Body) //I read all the info from the resp.Body stream to body by using io.ReadAll method
	if err != nil { //If error is not empty, it will return this error message
		fmt.Print("Error reading data")
		return forecast, err //returns forecast where forecast is empty and err is filled

	}

	err = json.Unmarshal(body, &forecast) //Parses the data from body to a Go-friendly structure using json.Unmarshal method
	if err != nil { //If error is not empty, it will return this error message
		fmt.Print("Error parsing JSON")
		return forecast, err //returns forecast where forecast is empty and err is filled
	}
	return forecast, nil

}

func main() {

	apiKey := "ff78855bd201acc9386304ad2f079e37" //Api key
	r := gin.Default() //gin router(engine)

	// serve static CSS/JS files

	r.Static("/static", "./static") //Ia
	r.LoadHTMLGlob("templates/*") //Instructing gin to the html file

	//Home page
r.GET("/", func(c *gin.Context) { //It gives you access to the request (like headers, query params if you want).
	
	c.HTML(http.StatusOK, "index.html", nil) //Sending an html response
})
 //API Endpoint returning JSON
r.GET("/api/weather", func(c *gin.Context) { ////It gives you access to the request (like headers, query params if you want).
	city := c.Query("city") // extract the "city" query parameter from request URL
	if city == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "city is required"}) //if query returns empty, it'll show this
		return
	}
	

	forecast, err := fetchWeather(city, apiKey) //call func fetchWeather() with city and apiKey with inputs and stores it in forecast
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})//if query returns empty, it'll show this 
		return 

	}

	type Item struct {		//Creates item struct type for front end easier to use
		DtTxt string `json:"dt_txt"`
		Temp float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
		Humidity int `json:"humidity"`
		Pop float64 `json:"pop"`

	}

	items := []Item{}		//creates a slice from Item struct

	for i, entry := range forecast.List { //Stores List slice inside entry
		if i >= 8 {   //next 24h
			break
		}
		items = append(items, Item{ //storing entry slice inside items slice for not exposing barebones and flattening out/renaming fields
			DtTxt: entry.DtText,
			Temp: entry.Main.Temp,
			FeelsLike: entry.Main.FeelsLike,
			Humidity: entry.Main.Humidity,
			Pop: entry.Pop,
		})
	}
	c.JSON(http.StatusOK, gin.H{"city": city, "list": items}) //returns httpstatus and the items slice

	
})
r.Run(":8080")
}








































// func main() {
// 	if len(os.Args) < 2 {
// 		fmt.Println("Usage: go run weather.go [City_Name]")
// 		return
// 	}

// 	city := os.Args[1]
// 	apiKey := "ff78855bd201acc9386304ad2f079e37"

// 	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/forecast?q=%s&appid=%s&units=metric", city, apiKey)
// 	resp, err := http.Get(url)
	
// 	if err != nil {
// 		fmt.Print("Error fetching data", err)
// 		return
// 	}
// 	defer resp.Body.Close()

// 	body, err := io.ReadAll(resp.Body)
// 	fmt.Println(string(body)) //raw Json text
// 	if err != nil {
// 		fmt.Println("Error parsing JSON:", err)
// 		return
// 	}

// 	var forecast ForecastResponse
// 	err = json.Unmarshal(body, &forecast)
// 	if err != nil {
// 		fmt.Print("Error parsing JSON:", err)
// 		return
// 	}

// 	fmt.Printf("Weather for %s:\n", city)
// 	for _, item := range forecast.List {
// 		fmt.Printf("Date and Time: %s\n", item.DtText)
// 		fmt.Printf("🌡 Temp: %.2f°C (Feels like %.2f°C)\n", item.Main.Temp, item.Main.FeelsLike)
// 		fmt.Printf("💧 Humidity: %d%%\n", item.Main.Humidity)
// 		fmt.Printf("🌧 Chance of rain: %f\n", item.Pop)
// 		fmt.Println("---------------------------")
// 	}


// }