package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"html/template"

	"github.com/gorilla/mux"
)

type matchup struct {
	Type     string
	Matchups []string
}

var weaknesses = []matchup{
	{"grass", []string{"fire", "flying", "poison", "bug", "ice"}},
	{"fire", []string{"water", "ground", "rock"}},
	{"water", []string{"grass", "electric"}},
	{"normal", []string{"fighting"}},
	{"flying", []string{"ice", "electric", "rock"}},
	{"poison", []string{"ground", "psychic"}},
	{"bug", []string{"fire", "rock", "flying"}},
	{"ice", []string{"fire", "rock", "fighting", "steel"}},
	{"ground", []string{"grass", "water", "ice"}},
	{"rock", []string{"fighting", "water", "grass", "ground", "steel"}},
	{"electric", []string{"ground"}},
	{"psychic", []string{"ghost", "dark", "bug"}},
	{"ghost", []string{"ghost", "dark"}},
	{"dark", []string{"fairy", "bug", "fighting"}},
	{"fighting", []string{"flying", "psychic", "fairy"}},
	{"dragon", []string{"ice", "dragon", "fairy"}},
	{"fairy", []string{"poison", "steel"}},
	{"steel", []string{"fighting", "ground", "fire"}},
}

var resistances = []matchup{
	{"grass", []string{"grass", "electric"}},
	{"fire", []string{"fire", "fairy", "bug"}},
	{"water", []string{"water", "steel", "fire", "ice"}},
	{"normal", []string{"ghost"}},
	{"flying", []string{"bug", "fighting", "grass", "ground"}},
	{"poison", []string{"poison", "fighting", "fairy"}},
	{"bug", []string{"grass", "ground", "fighting"}},
	{"ice", []string{"ice"}},
	{"ground", []string{"poison", "rock", "electric"}},
	{"rock", []string{"fire", "flying", "normal", "poison"}},
	{"electric", []string{"electric", "flying", "steel"}},
	{"psychic", []string{"fighting", "psychic"}},
	{"ghost", []string{"normal", "fighting", "bug", "poison"}},
	{"dark", []string{"dark", "ghost", "psychic"}},
	{"fighting", []string{"bug", "dark", "rock"}},
	{"dragon", []string{"water", "fire", "grass", "electric"}},
	{"fairy", []string{"bug", "dark", "fighting"}},
	{"steel", []string{"poison", "bug", "dragon", "fairy", "flying", "grass", "ice", "normal", "psychic", "rock", "steel"}},
}

type response struct {
	Name    string  `json:"name"`
	Types   []Type  `json:"types"`
	Sprites Sprites `json:"sprites"`
}

type Pokemon struct {
	Name   string `json:"name"`
	Types  []Type `json:"types"`
	Sprite string `json:"sprite"`
}

type Webdata struct {
	Name       string   `json:"name"`
	Weaknesses []string `json:"weaknesses"`
	Sprite     Sprites  `json:"sprite"`
}

type Type struct {
	Name TypeInfo `json:"type"`
}

type Sprites struct {
	Sprite string `json:"front_default"`
}

type TypeInfo struct {
	Name string `json:"name"`
}

type PokemonPage struct {
	Name       string
	Weaknesses []string
	Image      string
}

func (t TypeInfo) String() string {
	return fmt.Sprintf(t.Name)
}

func (t Type) String() string {
	return fmt.Sprintf("%v", t.Name)
}

func pokeapiCall(name string) response {
	resp, _ := http.Get("https://pokeapi.co/api/v2/pokemon/" + strings.ToLower(name))
	bytes, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	var responseObject response
	json.Unmarshal(bytes, &responseObject)
	// var typeStrings = make(map[Types]string)

	//return responseObject.Types.String()
	//var pokemonJSON = Pokemon{responseObject.Name, responseObject.Types, responseObject.Sprites}
	// var typeStrings = []string{}
	// for i := 0; i < len(responseObject.Types); i++ {
	// 	typeStrings = append(typeStrings, responseObject.Types[i].String())
	// }
	return responseObject
}

func contains(slice []string, target string) bool {
	doesItContain := false
	for i := 0; i < len(slice); i++ {
		if slice[i] == target {
			doesItContain = true
		}
	}
	return doesItContain
}

func switchTypes(zeroOrOne int) int {
	if zeroOrOne == 0 {
		return 1
	}
	return 0
}

func getTypeWeaknesses(name string) []string {
	for i := 0; i < len(weaknesses); i++ {
		if weaknesses[i].Type == name {
			return weaknesses[i].Matchups
		}
	}
	return []string{""}
}

func getTypeResistances(name string) []string {
	for i := 0; i < len(resistances); i++ {
		if resistances[i].Type == name {
			return resistances[i].Matchups
		}
	}
	return []string{""}
}

func calculateWeaknesses(name string) Webdata {
	name = strings.ToLower(name)
	var targetPoke = pokeapiCall(name)
	var answer = []string{}
	if len(targetPoke.Types) == 1 {
		return Webdata{targetPoke.Name, getTypeWeaknesses(targetPoke.Types[0].String()), targetPoke.Sprites}
	}
	for i := 0; i < len(targetPoke.Types); i++ {
		var list = getTypeWeaknesses(targetPoke.Types[i].String())
		for j := 0; j < len(list); j++ {
			if contains(answer, list[j]) == false && contains(getTypeResistances(targetPoke.Types[switchTypes(i)].String()), list[j]) == false {
				answer = append(answer, list[j])
			}
		}
	}
	var respJSON = Webdata{targetPoke.Name, answer, targetPoke.Sprites}
	return respJSON
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("indextemplate.html")
	t.Execute(w, nil)
}

func getPokemon(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	data := calculateWeaknesses(params["name"])
	p := PokemonPage{Name: data.Name, Weaknesses: data.Weaknesses, Image: data.Sprite.Sprite}
	t, _ := template.ParseFiles("pokemontemplate.html")
	t.Execute(w, p)
	// params := mux.Vars(r)
	// w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(calculateWeaknesses(params["name"]))
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/", indexHandler).Methods("GET")
	router.HandleFunc("/pokemon/{name}", getPokemon).Methods("GET")
	http.ListenAndServe(":8000", router)
}
