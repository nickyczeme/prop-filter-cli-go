package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"strconv"
	"math"

	"github.com/manifoldco/promptui"
	"github.com/pkg/errors"
	"github.com/olekukonko/tablewriter"

)

type Property struct {
	SquareFootage int               `json:"squareFootage"`
	Lighting      string            `json:"lighting"`
	Price         float64           `json:"price"`
	Rooms         int               `json:"rooms"`
	Bathrooms     int               `json:"bathrooms"`
	Location      [2]float64        `json:"location"`
	Description   string            `json:"description"`
	Amenities     map[string]bool   `json:"ammenities"`
}

func readProperties() ([]Property, error) {
	file, err := os.Open("properties.json")
	if err != nil {
		return nil, errors.Wrap(err, "Error opening properties.json")
	}
	defer file.Close()

	var properties []Property
	if err := json.NewDecoder(file).Decode(&properties); err != nil {
		return nil, errors.Wrap(err, "Error decoding properties.json")
	}

	return properties, nil
}

func displayFilteredProperties(filtered []Property) {
	if len(filtered) > 0 {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Sq ft", "Lighting", "Price", "Rooms", "WC", "Location", "Description", "Amenities"})
		table.SetBorders(tablewriter.Border{Left: true, Top: true, Right: true, Bottom: true})
		table.SetHeaderLine(true)

		// Iterate through filtered properties
		for _, property := range filtered {
			amenities := []string{}
			for amenity, hasAmenity := range property.Amenities {
				if hasAmenity {
					amenities = append(amenities, amenity)
				}
			}

			amenitiesStr := strings.Join(amenities, ", ")

			// Add row to the table
			table.Append([]string{
				strconv.Itoa(property.SquareFootage),
				property.Lighting,
				fmt.Sprintf("$%.2f", property.Price),
				strconv.Itoa(property.Rooms),
				strconv.Itoa(property.Bathrooms),
				fmt.Sprintf("(%.4f, %.4f)", property.Location[0], property.Location[1]),
				property.Description,
				amenitiesStr,
			})
		}

		// Render the table
		table.Render()
	} else {
		fmt.Println("No properties match your criteria.")
	}
}

func main() {
	properties, err := readProperties()
	if err != nil {
		fmt.Println("Error reading properties:", err)
		return
	}

	fmt.Println("Welcome to the Prop Filter CLI!")
	for {
		prompt := promptui.Select{
			Label: "What would you like to do?",
			Items: []string{
				"Filter by price", 
				"Filter by amenities", 
				"Filter by description",
				"Filter by distance", 
				"Filter by rooms", 
				"Filter by bathrooms", 
				"Filter by square footage",
				"Filter by lighting", 
				"Exit",
			},
		}

		_, action, err := prompt.Run()
		if err != nil {
			fmt.Println("Prompt failed:", err)
			return
		}

		switch action {
		case "Filter by price":
			filterByPrice(properties)
		case "Filter by amenities":
			filterByAmenities(properties)
		case "Filter by description":
			filterByDescription(properties)
		case "Filter by distance":
			filterByDistance(properties)
		case "Filter by rooms":
			filterByRooms(properties)
		case "Filter by bathrooms":
			filterByBathrooms(properties)
		case "Filter by square footage":
			filterBySquareFootage(properties)
		case "Filter by lighting":
			filterByLighting(properties)
		case "Exit":
			fmt.Println("Goodbye!")
			return
		}
	}
}

func filterByPrice(properties []Property) {
	promptOperator := promptui.Select{
		Label: "Choose a filter operator for price",
		Items: []string{"greaterThan", "lessThan", "equal"},
	}

	_, operator, err := promptOperator.Run()
	if err != nil {
		fmt.Println("Prompt failed:", err)
		return
	}

	promptValue := promptui.Prompt{
		Label: "Enter the price",
	}

	valueStr, err := promptValue.Run()
	if err != nil {
		fmt.Println("Prompt failed:", err)
		return
	}

	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		fmt.Println("Invalid price input")
		return
	}

	filtered := filterByNumber(properties, "price", operator, value)
	displayFilteredProperties(filtered)
}

func filterByAmenities(properties []Property) {
	prompt := promptui.Prompt{
		Label: "Enter the amenity to filter by (e.g., garage, pool)",
	}

	amenity, err := prompt.Run()
	if err != nil {
		fmt.Println("Prompt failed:", err)
		return
	}

	filtered := filterByAmenity(properties, amenity)
	displayFilteredProperties(filtered)
}

func filterByDescription(properties []Property) {
	prompt := promptui.Prompt{
		Label: "Enter a keyword to search in the description",
	}

	keyword, err := prompt.Run()
	if err != nil {
		fmt.Println("Prompt failed:", err)
		return
	}

	filtered := filterByDesc(properties, keyword)
	displayFilteredProperties(filtered)
}

func filterByDistance(properties []Property) {
	promptLat := promptui.Prompt{
		Label: "Enter your latitude",
	}

	latStr, err := promptLat.Run()
	if err != nil {
		fmt.Println("Prompt failed:", err)
		return
	}

	latitude, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		fmt.Println("Invalid latitude input")
		return
	}

	promptLon := promptui.Prompt{
		Label: "Enter your longitude",
	}

	lonStr, err := promptLon.Run()
	if err != nil {
		fmt.Println("Prompt failed:", err)
		return
	}

	longitude, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		fmt.Println("Invalid longitude input")
		return
	}

	promptDist := promptui.Prompt{
		Label: "Enter the maximum distance (in miles)",
	}

	distStr, err := promptDist.Run()
	if err != nil {
		fmt.Println("Prompt failed:", err)
		return
	}

	distance, err := strconv.ParseFloat(distStr, 64)
	if err != nil {
		fmt.Println("Invalid distance input")
		return
	}

	filtered := filterByDistanceFunc(properties, [2]float64{latitude, longitude}, distance)
	displayFilteredProperties(filtered)
}

func filterByRooms(properties []Property) {
	promptOperator := promptui.Select{
		Label: "Choose a filter operator for rooms",
		Items: []string{"greaterThan", "lessThan", "equal"},
	}

	_, operator, err := promptOperator.Run()
	if err != nil {
		fmt.Println("Prompt failed:", err)
		return
	}

	promptValue := promptui.Prompt{
		Label: "Enter the number of rooms",
	}

	valueStr, err := promptValue.Run()
	if err != nil {
		fmt.Println("Prompt failed:", err)
		return
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		fmt.Println("Invalid number of rooms input")
		return
	}

	filtered := filterByNumber(properties, "rooms", operator, float64(value))
	displayFilteredProperties(filtered)
}

func filterByBathrooms(properties []Property) {
	promptOperator := promptui.Select{
		Label: "Choose a filter operator for bathrooms",
		Items: []string{"greaterThan", "lessThan", "equal"},
	}

	_, operator, err := promptOperator.Run()
	if err != nil {
		fmt.Println("Prompt failed:", err)
		return
	}

	promptValue := promptui.Prompt{
		Label: "Enter the number of bathrooms",
	}

	valueStr, err := promptValue.Run()
	if err != nil {
		fmt.Println("Prompt failed:", err)
		return
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		fmt.Println("Invalid number of bathrooms input")
		return
	}

	filtered := filterByNumber(properties, "bathrooms", operator, float64(value))
	displayFilteredProperties(filtered)
}

func filterBySquareFootage(properties []Property) {
	// Ask for square footage criteria
	promptOperator := promptui.Select{
		Label: "Choose a filter operator for square footage",
		Items: []string{"greaterThan", "lessThan", "equal"},
	}

	_, operator, err := promptOperator.Run()
	if err != nil {
		fmt.Println("Prompt failed:", err)
		return
	}

	promptValue := promptui.Prompt{
		Label: "Enter the square footage",
	}

	valueStr, err := promptValue.Run()
	if err != nil {
		fmt.Println("Prompt failed:", err)
		return
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		fmt.Println("Invalid square footage input")
		return
	}

	filtered := filterByNumber(properties, "squareFootage", operator, float64(value))
	displayFilteredProperties(filtered)
}

func filterByLighting(properties []Property) {
	// Ask for lighting level
	prompt := promptui.Select{
		Label: "Select the lighting level",
		Items: []string{"low", "medium", "high"},
	}

	_, lighting, err := prompt.Run()
	if err != nil {
		fmt.Println("Prompt failed:", err)
		return
	}

	filtered := filterByLight(properties, lighting)
	displayFilteredProperties(filtered)
}

func filterByNumber(properties []Property, field string, operator string, value float64) []Property {
	var filtered []Property
	for _, property := range properties {
		var fieldValue float64
		switch field {
		case "price":
			fieldValue = property.Price
		case "rooms":
			fieldValue = float64(property.Rooms)
		case "bathrooms":
			fieldValue = float64(property.Bathrooms)
		case "squareFootage":
			fieldValue = float64(property.SquareFootage)
		}

		switch operator {
		case "equal":
			if fieldValue == value {
				filtered = append(filtered, property)
			}
		case "lessThan":
			if fieldValue < value {
				filtered = append(filtered, property)
			}
		case "greaterThan":
			if fieldValue > value {
				filtered = append(filtered, property)
			}
		}
	}
	return filtered
}

func filterByAmenity(properties []Property, amenity string) []Property {
	var filtered []Property
	for _, property := range properties {
		if property.Amenities[amenity] {
			filtered = append(filtered, property)
		}
	}
	return filtered
}

func filterByDesc(properties []Property, keyword string) []Property {
	var filtered []Property
	for _, property := range properties {
		if strings.Contains(strings.ToLower(property.Description), strings.ToLower(keyword)) {
			filtered = append(filtered, property)
		}
	}
	return filtered
}

func filterByLight(properties []Property, lighting string) []Property {
	var filtered []Property
	for _, property := range properties {
		if property.Lighting == lighting {
			filtered = append(filtered, property)
		}
	}
	return filtered
}

func getDistance(lat1, lon1, lat2, lon2 float64) float64 {
	radlat1 := (math.Pi * lat1) / 180
	radlat2 := (math.Pi * lat2) / 180
	theta := lon1 - lon2
	radtheta := (math.Pi * theta) / 180

	dist := math.Sin(radlat1)*math.Sin(radlat2) + math.Cos(radlat1)*math.Cos(radlat2)*math.Cos(radtheta)
	dist = math.Acos(dist)
	dist = (dist * 180) / math.Pi
	dist = dist * 60 * 1.1515 // Convert to miles
	return dist
}

func filterByDistanceFunc(properties []Property, userLocation [2]float64, maxDistance float64) []Property {
	var filtered []Property
	userLat, userLon := userLocation[0], userLocation[1]

	for _, property := range properties {
		propLat, propLon := property.Location[0], property.Location[1]
		distance := getDistance(userLat, userLon, propLat, propLon)

		if distance <= maxDistance {
			filtered = append(filtered, property)
		}
	}

	return filtered
}