// This file will receive a Pokemon Image and will return the appropriate pokemon.
package main

import (
	"image"
	"io/ioutil"
	"strconv"
	"fmt"
	"net/http"
	"math"
	"strings"

	"github.com/corona10/goimagehash"
	"github.com/oliamb/cutter"

	"gopkg.in/yaml.v2"
)

// BEGIN structs
// ENDOF structs

// BEGIN function definition
// receive grabs url of Pokemon picture
func receive(url string) string {
	pokemons := make(map[string]string)
	img := Download(url)
	CropUselessArea(img)
	hash := Hash(img)
	readPokemonList(&pokemons)
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
func readPokemonList(pokemonStruct *map[string]string) {
	reader, err := ioutil.ReadFile("config/hashes.yaml")
	if err != nil {
		fmt.Println(err)
	}
	err = yaml.Unmarshal(reader, pokemonStruct)
	if err != nil {
		fmt.Println(err)
	}
}

// Compare checks hash to hash list
func Compare(hash *goimagehash.ImageHash, pokemonStruct map[string]string) string {
	lowestHamming := 100
	var lowestHammingPokemon string
	var similar []string
	var similarLastDistance int
	for pokemon, pokemonHash := range pokemonStruct {
		distance := HammingDistance(strings.Replace(hash.ToString(),"p:","",1), pokemonHash)
		if distance == -1 {
			continue
		}
		if distance < lowestHamming {
			lowestHamming = distance
			lowestHammingPokemon = pokemon
			fmt.Printf("Current lowest: %s :%d\n",pokemon,distance)
		} else if distance == lowestHamming {
			fmt.Println("Same distance: "+ pokemon)
			if distance == similarLastDistance {
				similar = append(similar, pokemon)
			} else {
				similarLastDistance = distance
				similar = nil
				similar = append(similar, pokemon)
			}
		}
	}
	fmt.Printf("Lowest hamming distance: %d, Pokemon: %s\n", lowestHamming, lowestHammingPokemon)
	fmt.Println(similar)
	return lowestHammingPokemon
}

func SplitInPairs(s string) []int64 {
	var pairs []int64
	for i, char := range s  {
		if (i+1) % 2 == 0 {
			pair, _ := strconv.ParseInt(fmt.Sprintf("%d%d", char, s[i-1]), 10,64)
			pairs = append(pairs,pair)
		}
	}
	return pairs
}

func CropUselessArea(img *image.Image) {
	topLeft, size := FindVisibleVertexes(*img)
	newImg, _ := cutter.Crop(*img, cutter.Config{
		Width: size.X,
		Height: size.Y,
		Anchor: topLeft,
		Mode: cutter.TopLeft,
	})
	img = &newImg
}

func FindVisibleVertexes(img image.Image) (image.Point, image.Point) {
	// Bounds returns a rectangle with the shape fit inside
	bounds := img.Bounds()
	// Get the top left pixel
	topLeft := image.Point{X:bounds.Min.X,Y: bounds.Min.Y}
	size := bounds.Size()
	return topLeft, size
}

func HammingDistance(originalHash , probableHash string) int {

	originalHash = strings.Replace(originalHash,"\n","",1)
	probableHash = strings.Replace(probableHash,"\n","",1)
	pairsOriginal := SplitInPairs(originalHash)
	pairsProbable := SplitInPairs(probableHash)

	var hamming int
	for index := 0; index < len(pairsOriginal); index++ {
		currentPairOriginal := pairsOriginal[index]
		currentPairProbable := pairsProbable[index]
		for bit := 0; bit < int(math.Log2(float64(currentPairOriginal))); bit++ {
			// Perform bitwise AND to check each bit
			checkParticularBitOriginal := currentPairOriginal & int64(math.Pow(2, float64(bit)))
			checkParticularBitProbable := currentPairProbable & int64(math.Pow(2, float64(bit)))
			if checkParticularBitProbable ^ checkParticularBitOriginal != 0 {
				hamming += 1
			}

		}
	}
	return hamming
	/*originalHex, err := strconv.ParseInt(strings.Replace(originalHash,"\n","",1),16,64)
	probableHex, err := strconv.ParseInt(strings.Replace(probableHash,"\n","",1),16,64)
	if err != nil {
		fmt.Println(err)
		return -1
	}
	// Loop over both byte arrays at the same time
	var hammingCount int
	for originalByte := len(originalHex)-1; originalByte >= 0; originalByte -= 1 {
		// Perform XOR operation for each byte in the array, add to hammingCount if they differ.
		currentOriginalBit := originalHex & math.Pow(2,originalByte)
		currentProbableBit:= probableHex & math.Pow(2,originalByte)
		if currentOriginalBit ^ currentProbableBit {
			hammingCount += 1
		}
	}
	fmt.Println(hammingCount)
	return hammingCount*/
}

// Hash grabs value from Download
func Hash(imageDecoder *image.Image) *goimagehash.ImageHash {
	if imageDecoder == nil {
		return nil
	}
	hash, err := goimagehash.PerceptionHash(*imageDecoder)
	if err != nil {
		return nil
	}
	return hash
}
