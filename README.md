# Flat - Go Package for Flattening and Unflattening Complex Structures

FlatMapper is a Go package that provides functions for flattening and unflattening complex structures, allowing you to easily transform nested objects and lists into a flat map, and vice versa.

## Installation
```bash
go get github.com/andreasisnes/go-flat
```

## Usage

### Map
The Map function is used to flatten a complex structure. It takes a complex structure as input and returns a flat map, with keys representing the path to each value.
```go
package main

import (
	"fmt"
	"github.com/andreasisnes/go-flat"
)

type Person struct {
	Name struct {
		First string
		Last  string
	}
	Age int
}

func main() {
	// Create a complex structure
	complexData := Person{
		Name: struct {
			First string
			Last  string
		}{
			First: "John",
			Last:  "Doe",
		},
		Age: 30,
	}

	// Flatten the complex structure
	flatMap := flat.Map(complexData, flat.WithDelimiter("_"))

	// Print the flattened map
	fmt.Println("Flattened Map:")
	for key, value := range flatMap {
		fmt.Printf("%s: %v\n", key, value)
	}
    // Output
    // Name_First: John
    // Name_Last: Doe
    // Age: 30
}
```
### UMap
The UMap function allows you to recreate nested structures from a flattened map. This is particularly useful when you need to reverse the flattening process.

```go
flatMap := map[string]interface{}{
  "person.name.first": "John",
  "person.name.last":  "Doe",
  "person.age":        30,
}
nestedStructure := flat.UMap(flatMap, flat.WithDelimiter("."))
```

## Contributing
If you have suggestions, bug reports, or would like to contribute, feel free to open an issue or create a pull request.