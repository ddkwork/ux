# SafeMap

`SafeMap` is a thread-safe map in Go with additional features, such as ordered keys and more.

## Installation

To install the package, use the following command:

```sh
go install github.com/hitsumitomo/safemap@latest
```

### Methods

- `New[K cmp.Ordered, V any](ordered ...bool) *M[K, V]`: Initialize a new SafeMap. If `safemap.Ordered` is provided,
  keys will maintain insertion order.
- `Exists(key K) (exists bool)`: Check if a key exists in the map.
- `Load(key K) (value V, ok bool)`: Retrieves the value associated with a key.
- `Store(key K, value V)`: Store a key-value pair in the map.
- `Delete(key K)`: Removes a key from the map.
- `Add(key K, value V)`: Adds a value to the existing value for the given key.
- `Sub(key K, value V)`: Subtracts a value from the existing value for the given key.
- `LoadAndDelete(key K) (value V, ok bool)`: Load and delete a key from the map.
- `LoadOrStore(key K, value V) (actual V, loaded bool)`: Load or store a key-value pair in the map.
- `Swap(key K, value V) (previous V, loaded bool)`: Swap the value for a key and return the previous value.
- `Range(f func(K, V) bool)`: Iterate over the map with a function.
- `RangeDelete(f func(K, V) int8)`: Iterate over the map and delete keys based on a function that returns: -1 to delete
  the key, 0 to break the loop, and 1 to continue iterating.
- `Clear()`: Clear all entries in the map.
- `Len() int`: Returns the number of entries in the map.
- `Keys() []K`: Returns a slice of all keys in the map.
- `Map() map[K]V`: Returns a clone of the underlying map.
- `KeysInRange(from, to K) []K`: Returns all keys within the specified range `[from, to)`. Note: only works for positive
  numbers.

## Basic example

```go
package main

import (
	"fmt"

	"github.com/hitsumitomo/safemap"
)

func main() {
	sm := safemap.New[int, string]() // *safemap.M

	sm.Store(1, "one")
	sm.Store(2, "two")
	sm.Store(3, "three")

	if value, ok := sm.Load(2); ok {
		fmt.Println("Loaded value:", value)
	}

	if sm.Exists(3) {
		fmt.Println("Key 3 exists")
	}

	sm.Delete(1)

	sm.Range(func(key int, value string) bool {
		fmt.Printf("Key: %d, Value: %s\n", key, value)
		return true
	})

	fmt.Println("Length of map:", sm.Len())

	fmt.Println("\nUpdated map:")
	sm.Range(func(key int, value string) bool {
		fmt.Printf("Key: %d, Value: %s\n", key, value)
		return true
	})

	sm.RangeDelete(func(key int, value string) int8 {
		if key == 2 {
			return -1
		}
		return 1
	})

	fmt.Println("\nUpdated map:")
	sm.Range(func(key int, value string) bool {
		fmt.Printf("Key: %d, Value: %s\n", key, value)
		return true
	})

	sm.Clear()
	fmt.Println("Map cleared, length:", sm.Len())

	// ordered in the order of insertion.
	sm = safemap.New[int, string](safemap.Ordered)
	sm.Store(1, "Efficient")
	sm.Store(5, "solutions")
	sm.Store(2, "for")
	sm.Store(4, "complex")
	sm.Store(3, "tasks")
	fmt.Printf("\nOrdered map\n")
	sm.Range(func(key int, value string) bool {
		fmt.Printf("Key: %d, Value: %s\n", key, value)
		return true
	})
}

// Loaded value: two
// Key 3 exists
// Key: 2, Value: two
// Key: 3, Value: three
// Length of map: 2
//
// Updated map:
// Key: 2, Value: two
// Key: 3, Value: three
//
// Updated map:
// Key: 3, Value: three
// Map cleared, length: 0
//
// Ordered map
// Key: 1, Value: Efficient
// Key: 5, Value: solutions
// Key: 2, Value: for
// Key: 4, Value: complex
// Key: 3, Value: tasks
```

### Benchmarks

<pre>
BenchmarkSafeMap_Store-4                 1000000    1053 ns/op        265 B/op     2 allocs/op
BenchmarkSafeMap_Load-4                 48625240      24.63 ns/op       0 B/op     0 allocs/op
BenchmarkSafeMap_Delete-4                4250373     282.5 ns/op       55 B/op     1 allocs/op

BenchmarkSafeMap_Ordered_Store-4         1000000    1038 ns/op        265 B/op     2 allocs/op
BenchmarkSafeMap_Ordered_Load-4         48494152      24.65 ns/op       0 B/op     0 allocs/op
BenchmarkSafeMap_Ordered_Delete-4        4291180     276.9 ns/op       55 B/op     1 allocs/op

BenchmarkSyncMap_Store-4                 1000000    1001 ns/op        171 B/op     4 allocs/op
BenchmarkSyncMap_Load-4                 38536724      31.13 ns/op       0 B/op     0 allocs/op
BenchmarkSyncMap_Delete-4                1652521     728.0 ns/op      335 B/op     7 allocs/op
</pre>
