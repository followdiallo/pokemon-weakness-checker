package main

import (
	"fmt"
	"strings"
)

type pokemon struct {
	Name   string
	Types  []string
	Sprite string
}

type matchup struct {
	Type     string
	Matchups []string
}

var pokedex = []pokemon{
	{"venusaur", []string{"grass", "poison"}, "url"},
	{"charizard", []string{"fire", "flying"}, "url"},
	{"blastoise", []string{"water"}, "url"},
	{"pidgey", []string{"normal", "flying"}, "url"},
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
	{"water", []string{"water", "steel", "fire"}},
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

func calculateWeaknesses(name string) []string {
	name = strings.ToLower(name)
	var targetPoke pokemon
	for p := 0; p < len(pokedex); p++ {
		if pokedex[p].Name == name {
			targetPoke = pokedex[p]
		}
	}
	var answer = []string{}
	if len(targetPoke.Types) == 1 {
		return getTypeWeaknesses(targetPoke.Types[0])
	}
	for i := 0; i < len(targetPoke.Types); i++ {
		var list = getTypeWeaknesses(targetPoke.Types[i])
		// fmt.Println("LINE 115", list)
		// fmt.Println("LINE 116", getTypeResistances(targetPoke.Types[switchTypes(i)]))
		for j := 0; j < len(list); j++ {
			if contains(answer, list[j]) == false && contains(getTypeResistances(targetPoke.Types[switchTypes(i)]), list[j]) == false {
				// fmt.Println("INSIDE THE IF", list[j])
				answer = append(answer, list[j])
			}
		}
	}
	return answer
}

// var bulbasaur = []string{"for the love of God please work", "two"}
// var charizard = pokemon{"CHARIZARD", []string{"FIRE"}, "PIC"}
// var data = []pokemon{pokemon{"VILEPLUME", []string{"GRASS"}, "url"}}

func main() {
	//fmt.Println(getTypeWeaknesses("fire"))
	fmt.Println(calculateWeaknesses("pidgey"))
	fmt.Println(calculateWeaknesses("charizard"))
	fmt.Println(calculateWeaknesses("blastoise"))
	fmt.Println(calculateWeaknesses("venusaur"))
}
