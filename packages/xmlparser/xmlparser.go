package xmlparser

import (
	"encoding/xml"
	"go-game/packages/items" // Make sure this import path is correct for your project
	"io"
	"os"
)

// Items struct for holding slices of Weapon and Armor

// ParseItemsXML function to parse the XML file into Items struct
func ParseItemsXML(filePath string) (items.Items, error) {
	var itemsList items.Items

	// Open the XML file
	xmlFile, err := os.Open(filePath)
	if err != nil {
		return itemsList, err
	}
	defer xmlFile.Close()

	// Read the file content into a byte slice
	xmlData, err := io.ReadAll(xmlFile)
	if err != nil {
		return itemsList, err
	}

	// Unmarshal the XML into the Items struct
	err = xml.Unmarshal(xmlData, &itemsList)
	if err != nil {
		return itemsList, err
	}

	// Return the struct containing slices of weapons and armors
	return itemsList, nil
}
