// This file will receive a Pokemon Image and will return the appropriate pokemon.
// This will get called from jokercord.go

package main

import (
	"image"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/corona10/goimagehash"

	"gopkg.in/yaml.v2"
)

// BEGIN structs
// ENDOF structs

// BEGIN function definition
// receive grabs url of Pokemon picture
func receive(url string) string {
	pokemons := make(map[string]string)
	img := Download(url)
	hash := Hash(img)
	readPokemonList(pokemons)
	pokemonName := Compare(hash, pokemons)
	return (pokemonName)
}

// Download grabs Pokemon Picture from receive url
func Download(url string) *image.Image {
	response, err := http.Get(url)
	if err != nil {
		return nil
	}
	defer response.Body.Close()
	// Body is a io ReadCloser, so we can pass it to image.Decode, which receives an io.Reader
	decoded, _, err := image.Decode(response.Body)
	return &decoded
}

// readPokemonList reads hash list
func readPokemonList(pokemonStruct map[string]string) {
	reader, err := ioutil.ReadFile("config/hashes.yaml")
	logErr(err)
	yaml.Unmarshal(reader, pokemonStruct)
}

// Compare checks hash to hash list
func Compare(hash string, pokemonStruct map[string]string) string {
	var name string
	hash = strings.Replace(hash, "p:", "", 1)
	for pokemon, pokemonHash := range pokemonStruct {
		if pokemonHash == hash {
			name = pokemon
		}
	}
	return name

}

// Hash grabs value from Download
func Hash(imageDecoder *image.Image) string {
	if imageDecoder == nil {
		return ""
	}
	hash, err := goimagehash.PerceptionHash(*imageDecoder)
	if err != nil {
		return ""
	}
	return hash.ToString()
}
