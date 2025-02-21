# Prop Filter CLI
## Overview
The Prop Filter CLI is a command-line tool designed to filter a large set of real estate properties based on specific attributes.

This CLI was developed using GO and adheres to clean coding principles, with a focus on usability, extensibility, and performance.

## Features
- Filtering Options:
  - Price, Rooms, Bathrooms, Square Footage: Apply quantifiers (lessThan, greaterThan, or equal) to these numeric attributes.
  - Lighting: Search by lighting levels (low, medium, high).
  - Amenities: Check for specific amenities like yard, garage, or pool.
  - Description Matching: Search for keywords in the property description.
  - Distance Filtering: Find properties within a specific distance.
- Readable Results:
  - Properties are displayed in a table.
- Dynamic Data Loading:
  - The CLI reads property data dynamically from a properties.json file.

## Requirements
- Node.js v16.0.0 or higher
- TypeScript v4.0.0 or higher

## Setup 
1. Clone the Repository

   ```bash
   git clone https://github.com/nickyczeme/prop-filter-cli-go.git
   cd prop-filter-cli-go
   ```
2. Install Dependencies
   ```
   go get
   ```
3. Run the CLI
   ```
   go run main.go
   ```

## Usage 
When you run the CLI, you will be greeted with a menu of options to filter properties. Navigate the menu using the arrow keys and follow the prompts to apply filters.
