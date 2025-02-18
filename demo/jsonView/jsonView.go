package main

import (
	_ "embed"
)

//go:generate go build -x .
//go:generate go install .

//go:embed *.go
var icon []byte

// https://faststone-photo-resizer.en.lo4d.com/windows
