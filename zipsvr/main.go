//Package main is used for your program's main() function,
//plus any other code you want to include in the main package.
//All types, functions, and globals in a package are visible
//throughout the package, but only exported identifiers
//to code in other packages. To export something, make its
//name start with a capital letter (I know, it's kind of goofy
//but that's just the way Go works).
package main

//If you want to use functions or types defined in other
//packages, you need to import them. For standard library
//packages, you just use the package name here. After you
//import the package, you can refer to things in the package
//by using the package name like an object that exposes all
//of its exported types and functions as properties and
//methods of that object. See below for examples.
import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

type zip struct {
	Zip   string `json:"zip"`
	City  string `json:"city"`
	State string `json:"state"`
}

type zipSlice []*zip
type zipIndex map[string]zipSlice

//helloHandler handles requests made to the /hello path.
//Every HTTP handler has this same signature:
//  func (w http.ResponseWriter, r *http.Request)
//The `w` parameter allows you to set response headers
//and status codes, as well as write the response body.
//The `r` parameter gives you access to all of the request
//information and any content in the request body.
//For more details, see:
// - https://golang.org/pkg/net/http/#ResponseWriter
// - https://golang.org/pkg/net/http/#Request
//or just put your cursor on the type name of these
//parameters and hit F12 (Go to Definition command)
func helloHandler(w http.ResponseWriter, r *http.Request) {
	//get the `name` query string parameter
	name := r.URL.Query().Get("name")

	//if it's zero-length, set name to "World"
	if len(name) == 0 {
		name = "World"
	}

	//set the Content-Type header to "text/plain"
	//as we are just writing plain text in the response
	w.Header().Add("Content-Type", "text/plain")

	//write the response body
	//w.Write() accepts a byte slice so that you can
	//write either text or binary data (e.g., images).
	//To convert a string to a byte slice, just do a
	//type conversion: []byte(myString)
	//This works for converting any variable to another
	//type, provided the conversion is deterministic
	w.Write([]byte("Hello " + name))
}

func (zi zipIndex) zipsForCityHandler(w http.ResponseWriter, r *http.Request) {
	// /zips/city/seattle
	_, city := path.Split(r.URL.Path)
	lcity := strings.ToLower(city)

	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.Header().Add("Access-Control-Allow-Origin", "*")

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(zi[lcity]); err != nil {
		http.Error(w, "error encoding json: "+err.Error(), http.StatusInternalServerError)
	}
}

//main is the entry-point for all go programs
//program execution starts with this function
func main() {
	//get the ADDR envrionment variable
	//to set this, execute the following in your terminal
	//before running this program:
	//  export ADDR=localhost:8000
	//Here we use the `os` package from the standard library.
	//We imported it above. Once you import it, you can access
	//all of it's exported types and functions use `os.`
	addr := os.Getenv("ADDR")
	if len(addr) == 0 {
		//log.Fatal() writes the message to stdout and
		//exits with a code of 1, indicating an error
		log.Fatal("please set ADDR environment variable")
	}

	f, err := os.Open("../data/zips.json")
	if err != nil {
		log.Fatal("error opening zips file: " + err.Error())
	}

	zips := make(zipSlice, 0, 43000)
	decoder := json.NewDecoder(f)
	if err := decoder.Decode(&zips); err != nil {
		log.Fatal("error decoding zips json: " + err.Error())
	}
	fmt.Printf("loaded %d zips\n", len(zips))

	zi := make(zipIndex)

	//build a map of city name (lower-case)
	//to a zipSlice ( []*zip )
	for _, z := range zips {
		lower := strings.ToLower(z.City)
		zi[lower] = append(zi[lower], z)
	}

	fmt.Printf("there are %d zips in Seattle\n", len(zi["seattle"]))

	//Register our helloHandler as the handler for
	//the `/hello` resource path. Whenever a request
	//is made to this path, the Go web server will
	//call our helloHandler function.
	http.HandleFunc("/hello", helloHandler)

	http.HandleFunc("/zips/city/", zi.zipsForCityHandler)

	//Let the client know what address the server is
	//listening on. The `fmt` package lets you write
	//messages to stdout. It can also format messages
	//by replacing tokens like %s with strings you
	//pass as additional parameters. For more details see:
	//https://golang.org/pkg/fmt/
	fmt.Printf("server is listening at %s...\n", addr)

	//Start the web server on the address, and use the
	//default router. The default router is what you
	//configured above when you called http.HandleFunc().
	//http.ListenAndServe() is a blocking function so
	//it won't return until the web server is stopped,
	//but if it can't actually start (e.g., can't bind)
	//to the port number you gave it), it will return
	//and error, which we will log using log.Fatal().
	log.Fatal(http.ListenAndServe(addr, nil))
}
