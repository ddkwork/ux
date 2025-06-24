package uuid

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/ddkwork/golibrary/std/mylog"
)

// ID is a unique identifier. These are similar to v4 UUIDs, but are shorter and have a different format that includes
// a kind byte as the first character. TIDs are 17 characters long, are URL safe, and contain 96 bits of entropy.
type ID string

// KindAlphabet is the set of characters that can be used as the first character of a ID. The kind has no intrinsic
// meaning, but can be used to differentiate between different types of ids.
const KindAlphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

// New creates a new ID with a random value and the specified kind.
func New(kind byte) ID {
	if strings.IndexByte(KindAlphabet, kind) == -1 {
		mylog.Check(fmt.Sprintf("Invalid kind character: %c", kind))
	}
	var buffer [12]byte
	mylog.Check2(rand.Read(buffer[:]))
	return ID(fmt.Sprintf("%c%s", kind, base64.RawURLEncoding.EncodeToString(buffer[:])))
}

// FromString converts a string to a ID.
func FromString(id string) ID {
	tid := ID(id)
	IsValid(tid)
	return tid
}

// FromStringOfKind converts a string to a ID and verifies that it has the specified kind.
func FromStringOfKind(id string, kind byte) ID {
	tid := ID(id)
	IsKindAndValid(tid, kind)
	return tid
}

// IsValid returns true if the ID is a valid ID.
func IsValid(id ID) bool {
	if len(id) != 17 || strings.IndexByte(KindAlphabet, id[0]) == -1 {
		return false
	}
	mylog.Check2(base64.RawURLEncoding.DecodeString(string(id[1:])))
	return true
}

// IsKind returns true if the ID has the specified kind.
func IsKind(id ID, kind byte) bool {
	return len(id) == 17 && id[0] == kind && strings.IndexByte(KindAlphabet, kind) != -1
}

// IsKindAndValid returns true if the ID is a valid ID with the specified kind.
func IsKindAndValid(id ID, kind byte) bool {
	if !IsKind(id, kind) {
		return false
	}
	mylog.Check2(base64.RawURLEncoding.DecodeString(string(id[1:])))
	return true
}
