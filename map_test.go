package flat

import (
	"encoding/json"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	datasetMap    = "map.json"
	datasetList   = "list.json"
	datasetRandom = "random.json"
	datasetDir    = "data"
)

func readMap(dataset string) map[string]any {
	return readfile(dataset, make(map[string]any)).(map[string]any)
}

func readList(dataset string) []any {
	return readfile(dataset, make([]any, 0)).([]any)
}

func readfile(dataset string, result any) any {
	if content, err := os.ReadFile(path.Join(datasetDir, dataset)); err != nil {
		panic(err)
	} else {
		result = make(map[string]any)
		if err = json.Unmarshal(content, &result); err != nil {
			panic(err)
		}
	}

	return result
}

type testStructOuter struct {
	Nested *testStructInner
}

type testStructInner struct {
	FieldChan    chan string
	FieldPointer *string
	Field        string
}

func TestFlatWithUnestedValues(t *testing.T) {
	res := Map(5)
	assert.Len(t, res, 0)

	pValue := 5
	res = Map(&pValue)
	assert.Len(t, res, 0)

	res = Map("test")
	assert.Len(t, res, 0)

	res = Map(false)
	assert.Len(t, res, 0)
}

func TestListWithCustomValues(t *testing.T) {
	innerString := "TestPointer"
	outerString := &innerString
	c := make(chan string)
	value := []any{
		[]int{1, 2, 3, 4, 5},
		[]int{2, 3, 4},
		testStructOuter{
			Nested: &testStructInner{
				FieldChan:    c,
				FieldPointer: outerString,
				Field:        "Test",
			},
		},
	}

	result := Map(value)

	assert.Equal(t, c, result["2.Nested.FieldChan"])
	assert.Equal(t, outerString, result["2.Nested.FieldPointer"])
	assert.Equal(t, "Test", result["2.Nested.Field"])

	assert.Equal(t, 1, result["0.0"])
	assert.Equal(t, 5, result["0.4"])
	assert.Equal(t, 2, result["1.0"])
	assert.Equal(t, 4, result["1.2"])
}

func TestMapWithNil(t *testing.T) {
	result := Map(nil)
	assert.Nil(t, result)
}

func TestListWithNil(t *testing.T) {
	result := Map(nil)
	assert.Nil(t, result)
}

func TestMapWithDifferentDelimiter(t *testing.T) {
	result := Map(readMap(datasetMap), func(options *Options) {
		options.Delimiter = "<>"
	})

	assert.Equal(t, "MapNestedField", result["Nested<>Nested<>Field"])
}

func TestMapWithUpperCaseFolding(t *testing.T) {
	result := Map(readMap(datasetMap), func(options *Options) {
		options.Fold = strings.ToUpper
	})

	assert.Equal(t, "MapSingleField", result["FIELD"])
}

func TestMapWithEmptyDelimiter(t *testing.T) {
	result := Map(readMap(datasetMap))
	assert.Equal(t, "MapSingleField", result["Field"])
}

func TestListWithEmptyDelimiter(t *testing.T) {
	result := Map(readList(datasetList))
	assert.Equal(t, "0", result["0.0"])
}

func TestFlatMapSingleField(t *testing.T) {
	data := readMap(datasetMap)

	result := Map(data, nil)
	assert.Equal(t, "MapSingleField", result["Field"])
}

func TestFlatMapNestedFields(t *testing.T) {
	data := readMap(datasetMap)

	result := Map(data, nil)
	assert.Equal(t, "MapNestedField", result["Nested.Nested.Field"])
	assert.Equal(t, "AnotherValue", result["Nested.Nested.AntotherField"])
}

func TestFlatMapList(t *testing.T) {
	data := readMap(datasetMap)

	result := Map(data, nil)
	assert.Equal(t, "0", result["List.0"])
	assert.Equal(t, "1", result["List.1"])
	assert.Equal(t, "2", result["List.2"])
}

func TestFlatMapListVariousTypes(t *testing.T) {
	data := readMap(datasetMap)

	result := Map(data, nil)
	assert.Equal(t, "0", result["ListVariousTypes.0"])
	assert.Equal(t, 1, int(result["ListVariousTypes.1"].(float64)))
	assert.Equal(t, 2.2, result["ListVariousTypes.2"])
	assert.Equal(t, false, result["ListVariousTypes.3"])
}

func TestFlatMapNestedList(t *testing.T) {
	data := readMap(datasetMap)

	result := Map(data, nil)
	assert.Equal(t, "0", result["NestedList.0.0"])
	assert.Equal(t, "0", result["NestedList.1.0"])
	assert.Equal(t, "1", result["NestedList.1.1"])
	assert.Equal(t, "0", result["NestedList.2.0"])
	assert.Equal(t, "1", result["NestedList.2.1"])
	assert.Equal(t, "2", result["NestedList.2.2"])
}

func TestFlatListNestedListsWithObjects(t *testing.T) {
	data := readList(datasetList)

	result := Map(data, nil)

	assert.Equal(t, 0, int(result["2.List.0"].(float64)))
	assert.Equal(t, false, result["2.List.1"])
	assert.Equal(t, "Value", result["2.List.2.0"])
	assert.Equal(t, "Value", result["2.List.2.1.Field"])
}

func TestFlatListNestedObjects(t *testing.T) {
	data := readList(datasetList)

	result := Map(data, nil)

	assert.Equal(t, "Value", result["2.Field"])
	assert.Equal(t, "Value", result["2.Nested.Nested.Field"])
}

func TestFlatListNestedLists(t *testing.T) {
	data := readList(datasetList)

	result := Map(data, nil)

	assert.Equal(t, "0", result["0.0"])
	assert.Equal(t, "0", result["1.0"])
	assert.Equal(t, "1", result["1.1"])
}

func TestRandomJsonData(t *testing.T) {
	data := readList(datasetRandom)
	result := Map(data, nil)

	assert.Equal(t, result["0.id"], "0001")
	assert.Equal(t, result["0.batters.batter.2.id"], "1003")
}
