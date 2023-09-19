package unicorns

import (
	"bufio"
	"fmt"
	"os"
)

type CapabilitiesList = []string

var Capabilities = CapabilitiesList{
	"super strong", "fullfill wishes", "fighting capabilities", "fly", "swim", "sing", "run", "cry", "change color", "talk", "dance", "code", "design", "drive", "walk", "talk chinese", "lazy",
}

// LIFO stack
type UnicornElement struct {
	Name         string           `json:"name"`
	Capabilities CapabilitiesList `json:"capabilities"`
}

type UnicornList []UnicornElement

func (u *UnicornList) Push(unicorn *UnicornElement) {
	*u = append(*u, *unicorn)
}

func (u *UnicornList) Pop() *UnicornElement {
	uni := *u
	if len(*u) > 0 {
		res := uni[len(uni)-1]
		*u = uni[:len(uni)-1]
		return &res
	}
	return nil
}

func getArrayStringFromTextFile(fileName string) ([]string, error) {
	fn, err := os.Open(fileName)
	defer fn.Close()
	if err != nil {
		return nil, err
	}
	var lines []string
	var scanner = bufio.NewScanner(fn)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, nil
}

type PetNames = []string
type Adjectives = []string

func GetPetNames() (PetNames, error) {
	names, err := getArrayStringFromTextFile("petnames.txt")
	if err != nil {
		fmt.Println("Unicorn names not found")
		return nil, err
	}
	return names, nil
}

func GetAdjectives() (Adjectives, error) {
	adjectives, err := getArrayStringFromTextFile("adj.txt")
	if err != nil {
		fmt.Println("Unicorn adjectives not found")
		return nil, err
	}
	return adjectives, nil
}

func (u *UnicornList) GetUnicornByName(name string) bool {
	for _, item := range *u {
		if item.Name == name {
			return true
		}
	}
	return false
}

func GetCapabilityByName(capabilitiesList CapabilitiesList, name string) bool {
	for _, item := range capabilitiesList {
		if item == name {
			return true
		}
	}
	return false
}
