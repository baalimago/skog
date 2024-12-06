package construct

import (
	"reflect"
	"testing"

	"github.com/baalimago/skog/internal/models"
)

func TestCurrentLevel(t *testing.T) {
	tests := []struct {
		name      string
		data      models.JSONLike
		position  index
		expected  models.JSONLike
		expectErr bool
	}{
		{
			name: "Successful traversal",
			data: models.JSONLike{
				"level1": models.JSONLike{
					"level2": models.JSONLike{
						"level3": "value",
					},
				},
			},
			position: index{
				root: []string{"level1", "level2"},
			},
			expected: models.JSONLike{
				"level3": "value",
			},
			expectErr: false,
		},
		{
			name: "Traverse beyond current level",
			data: models.JSONLike{
				"level1": models.JSONLike{
					"level2": models.JSONLike{
						"level3": "value",
					},
				},
			},
			position: index{
				root: []string{"level1", "level2", "level3"},
			},
			expected:  nil,
			expectErr: true,
		},
		{
			name: "Invalid root path",
			data: models.JSONLike{
				"level1": models.JSONLike{
					"level2": "not a JSONLike",
				},
			},
			position: index{
				root: []string{"level1", "level2"},
			},
			expected:  nil,
			expectErr: true,
		},
		{
			name: "Empty data",
			data: models.JSONLike{},
			position: index{
				root: []string{"level1"},
			},
			expected:  nil,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := builder{
				data:     tt.data,
				position: tt.position,
			}
			result, err := b.CurrentLevel()

			if (err != nil) != tt.expectErr {
				t.Errorf("CurrentLevel() error = %v, expectErr %v", err, tt.expectErr)
				return
			}
			if result != nil && !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("CurrentLevel() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestSet(t *testing.T) {
	// Create a new builder instance
	b := NewBuilder()

	// Test setting a new key-value pair
	key := "testKey"
	value := "testValue"
	b.Set(key, value)

	// Validate the value was set correctly
	if b.data[key] != value {
		t.Errorf("expected %v for key %v, got %v", value, key, b.data[key])
	}

	// Test setting an existing key to a new value
	newValue := "newValue"
	b.Set(key, newValue)

	// Validate the value was updated correctly
	if b.data[key] != newValue {
		t.Errorf("expected %v for key %v, got %v", newValue, key, b.data[key])
	}

	// Test setting another key
	anotherKey := "anotherKey"
	anotherValue := 123
	b.Set(anotherKey, anotherValue)

	// Validate this new key-value pair was added correctly
	if b.data[anotherKey] != anotherValue {
		t.Errorf("expected %v for key %v, got %v", anotherValue, anotherKey, b.data[anotherKey])
	}
}

func TestDel(t *testing.T) {
	// Create a new builder instance
	b := NewBuilder()

	// Test deleting a non-existent key
	key := "nonExistentKey"
	b.Del(key)
	if _, exists := b.data[key]; exists {
		t.Errorf("key %v should not exist after deletion", key)
	}

	// Test deleting an existing key
	existingKey := "existingKey"
	value := "testValue"
	b.Set(existingKey, value)
	b.Del(existingKey)
	if _, exists := b.data[existingKey]; exists {
		t.Errorf("key %v should not exist after deletion", existingKey)
	}

	// Test deleting one key doesn't affect others
	key1 := "key1"
	key2 := "key2"
	value1 := "value1"
	value2 := "value2"

	b.Set(key1, value1)
	b.Set(key2, value2)
	b.Del(key1)

	if _, exists := b.data[key1]; exists {
		t.Errorf("key %v should not exist after deletion", key1)
	}
	if val, exists := b.data[key2]; !exists || val != value2 {
		t.Errorf("key %v should still exist with value %v", key2, value2)
	}
}

func TestNewBuilder(t *testing.T) {
	// Create a new builder instance
	b := NewBuilder()

	// Test that the data field is initialized as an empty JSONLike
	if b.data == nil {
		t.Error("expected data field to be initialized, got nil")
	}

	if len(b.data) != 0 {
		t.Errorf("expected empty data map, got length %v", len(b.data))
	}

	// Test that the position field is initialized with zero values
	if b.position.root != nil {
		t.Errorf("expected nil position root, got %v", b.position.root)
	}

	if b.position.key != "" {
		t.Errorf("expected empty position key, got %v", b.position.key)
	}

	// Test that we can add data to the initialized builder
	b.Set("test", "value")
	if len(b.data) != 1 {
		t.Error("failed to add data to initialized builder")
	}
}

func TestTraverse(t *testing.T) {
	tests := []struct {
		name      string
		data      models.JSONLike
		path      []string
		expected  models.JSONLike
		expectErr bool
	}{
		{
			name: "Successful traversal",
			data: models.JSONLike{
				"level1": models.JSONLike{
					"level2": models.JSONLike{
						"level3": "value",
					},
				},
			},
			path: []string{"level1", "level2"},
			expected: models.JSONLike{
				"level3": "value",
			},
			expectErr: false,
		},
		{
			name: "Empty path",
			data: models.JSONLike{
				"level1": models.JSONLike{},
			},
			path:      []string{},
			expected:  models.JSONLike{"level1": models.JSONLike{}},
			expectErr: false,
		},
		{
			name: "Invalid path",
			data: models.JSONLike{
				"level1": models.JSONLike{
					"level2": "not a JSONLike",
				},
			},
			path:      []string{"level1", "level2"},
			expected:  nil,
			expectErr: true,
		},
		{
			name:      "Empty data",
			data:      models.JSONLike{},
			path:      []string{"level1"},
			expected:  nil,
			expectErr: true,
		},
		{
			name: "Non-existent path",
			data: models.JSONLike{
				"level1": models.JSONLike{},
			},
			path:      []string{"nonexistent"},
			expected:  nil,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := builder{
				data: tt.data,
			}
			result, err := b.Traverse(tt.path)

			if (err != nil) != tt.expectErr {
				t.Errorf("Traverse() error = %v, expectErr %v", err, tt.expectErr)
				return
			}
			if result != nil && !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Traverse() = %v, expected %v", result, tt.expected)
			}
		})
	}
}
