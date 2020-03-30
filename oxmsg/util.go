package oxmsg

import (
	"regexp"
)

var propertyByID map[PropertyID]*Property
var propertyByName map[string]*Property

const sPrefix = "Received:"

var rContinuedLine = regexp.MustCompile(`^(\s)+`)
var rMultipleSpaces = regexp.MustCompile(`\s+`)

func init() {
	propertyByID = make(map[PropertyID]*Property)
	propertyByName = make(map[string]*Property)
	for i := range Properties {
		p := Properties[i]
		if p.ID < PsetLAST {
			propertyByID[p.ID] = &p
		}
		propertyByName[p.Name] = &p
	}
}

func GetPropertyByID(id PropertyID) *Property {
	return propertyByID[id]
}

func GetPropertyByName(name string) *Property {
	return propertyByName[name]
}

func TrimMultipleSpaces(s string) string {
	return rMultipleSpaces.ReplaceAllString(s, " ")
}

// func ExtractIP(s string) []net.IP {
// 	return nil
// }
