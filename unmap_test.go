package flat

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUMapWithNil(t *testing.T) {
	result := UMap(nil)
	assert.Nil(t, result)
}

func TestUListWithNil(t *testing.T) {
	result := UMap(nil)
	assert.Nil(t, result)
}

func TestUMapWithEmptyDelimiter(t *testing.T) {
	flat := Map(readMap(datasetMap), nil)
	result := UMap(flat).(map[string]any)

	assert.Equal(t, "MapSingleField", result["Field"])
}

func TestUListWithEmptyDelimiter(t *testing.T) {
	flat := Map(readList(datasetList), nil)
	result := UMap(flat)
	assert.Equal(t, "0", result.([]any)[0].([]any)[0])
}

func TestUnflatMapSingleField(t *testing.T) {
	flat := Map(readMap(datasetMap), nil)
	result := UMap(flat, nil).(map[string]any)

	assert.Equal(t, "MapSingleField", result["Field"])
}

func TestUnflatMapNestedFields(t *testing.T) {
	opts := &Options{
		Delimiter: "<>",
		Fold:      strings.ToLower,
	}

	flat := Map(readMap(datasetMap), WithOptions(opts))
	result := UMap(flat, WithOptions(opts)).(map[string]any)

	assert.Equal(t, "MapNestedField", result["nested"].(map[string]any)["nested"].(map[string]any)["field"])
	assert.Equal(t, "AnotherValue", result["nested"].(map[string]any)["nested"].(map[string]any)["antotherfield"])
	assert.Nil(t, result["nested"].(map[string]any)["emptyobject"])
}

func TestUnflatMapHazard(t *testing.T) {
	flat := Map(readMap(datasetMap), nil)
	result := UMap(flat, nil).(map[string]any)

	assert.Equal(t, "0", result["Hazard"].(map[string]any)["0"])
	assert.Equal(t, "1", result["Hazard"].(map[string]any)["1"])
	assert.Equal(t, "Value", result["Hazard"].(map[string]any)["Nested"].(map[string]any)["Field"])
}

func TestUnflatMapList(t *testing.T) {
	flat := Map(readMap(datasetMap), nil)
	result := UMap(flat, nil).(map[string]any)

	assert.Equal(t, "0", result["List"].([]any)[0])
	assert.Equal(t, "1", result["List"].([]any)[1])
	assert.Equal(t, "2", result["List"].([]any)[2])
}

func TestUnflatMapListVariousTypes(t *testing.T) {
	flat := Map(readMap(datasetMap), nil)
	result := UMap(flat, nil).(map[string]any)

	assert.Equal(t, "0", result["ListVariousTypes"].([]any)[0])
	assert.Equal(t, 1, int(result["ListVariousTypes"].([]any)[1].(float64)))
	assert.Equal(t, 2.2, result["ListVariousTypes"].([]any)[2])
	assert.Equal(t, false, result["ListVariousTypes"].([]any)[3])
}

func TestUnflatMapNestedList(t *testing.T) {
	flat := Map(readMap(datasetMap), nil)
	result := UMap(flat, nil).(map[string]any)

	assert.Equal(t, "0", result["NestedList"].([]any)[0].([]any)[0])
	assert.Equal(t, "0", result["NestedList"].([]any)[1].([]any)[0])
	assert.Equal(t, "1", result["NestedList"].([]any)[1].([]any)[1])
	assert.Equal(t, "0", result["NestedList"].([]any)[2].([]any)[0])
	assert.Equal(t, "1", result["NestedList"].([]any)[2].([]any)[1])
	assert.Equal(t, "2", result["NestedList"].([]any)[2].([]any)[2])
}

func TestUnflatListNestedListsWithObjects(t *testing.T) {
	flat := Map(readList(datasetList), nil)
	result := UMap(flat, nil).([]any)

	assert.Equal(t, 0, int(result[2].(map[string]any)["List"].([]any)[0].(float64)))
	assert.Equal(t, false, result[2].(map[string]any)["List"].([]any)[1])
	assert.Equal(t, "Value", result[2].(map[string]any)["List"].([]any)[2].([]any)[0])
	assert.Equal(t, "Value", result[2].(map[string]any)["List"].([]any)[2].([]any)[1].(map[string]any)["Field"])
}

func TestUnflatListNestedObjects(t *testing.T) {
	flat := Map(readList(datasetList), nil)
	result := UMap(flat, nil).([]any)

	assert.Equal(t, "Value", result[2].(map[string]any)["Field"])
	assert.Equal(t, "Value", result[2].(map[string]any)["Nested"].(map[string]any)["Nested"].(map[string]any)["Field"])
}

func TestUnflatListNestedLists(t *testing.T) {
	flat := Map(readList(datasetList), nil)
	result := UMap(flat, nil).([]any)

	assert.Equal(t, "0", result[0].([]any)[0])
	assert.Equal(t, "0", result[1].([]any)[0])
	assert.Equal(t, "1", result[1].([]any)[1])
}
