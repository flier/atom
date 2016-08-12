package atom

import (
	"bytes"
	"hash/fnv"
	"sync"
)

// Atom is an integer code for a string.
//
// The zero value maps to "".
type Atom uint32

// Empty is an empty atom without name
const Empty = Atom(0)

// MaxAtomLen is the maximum length of an atom name
const MaxAtomLen = 127

// Cache hash based atom caches
type Cache map[uint64]Atom

// Value return the atom value
func (a Atom) Value() uint32 {
	return uint32(a)
}

// Hash return the hash value of atom name
func (a Atom) Hash() uint64 {
	return hashAtom(a.Bytes())
}

// IsEmpty reports whether the atom is empty
func (a Atom) IsEmpty() bool {
	return a == Empty
}

// IsEmbedded reports whether the atom contains an embedded string
func (a Atom) IsEmbedded() bool {
	return (uint32(a) & embeddedTag) != 0
}

// Len return the atom name length
func (a Atom) Len() int {
	n := int((uint32(a) >> atomLenShift) & atomLenMask)

	if a.IsEmbedded() && n > 4 {
		return 4
	}

	return n
}

// Bytes returns the bytes for the atom
func (a Atom) Bytes() []byte {
	if a.IsEmpty() {
		return nil
	}

	n := a.Len()

	if a.IsEmbedded() {
		switch n {
		case 1:
			return []byte{byte(a)}
		case 2:
			return []byte{byte(uint32(a) >> 8), byte(a)}
		case 3:
			return []byte{byte(uint32(a) >> 16), byte(uint32(a) >> 8), byte(a)}
		default:
			return []byte{byte(uint32(a)>>24) & atomLenMask, byte(uint32(a) >> 16), byte(uint32(a) >> 8), byte(a)}
		}
	}

	off := int(a & 0xFFFFFF)

	return atomData.Bytes()[off : off+n]
}

func (a Atom) String() string {
	return string(a.Bytes())
}

const (
	embeddedTag  = 1 << 31
	atomLenMask  = 0x7F
	atomLenShift = 24
)

var atomData *bytes.Buffer
var atomCache Cache
var atomLock sync.RWMutex

func reset(data []byte, cache Cache) {
	atomLock.Lock()
	atomData = bytes.NewBuffer(data)
	if cache == nil {
		atomCache = make(Cache)
	} else {
		atomCache = cache
	}
	atomLock.Unlock()
}

func init() {
	reset(nil, nil)
}

func newAtom(off, size int) Atom {
	if size == 0 || size > MaxAtomLen || off > atomData.Len() {
		return Empty
	}

	return Atom((size << atomLenShift) + off)
}

func embedAtom(s []byte) Atom {
	n := len(s)

	if n == 4 && ((uint8(s[0]) & 0x80) == 0) {
		return Atom(embeddedTag | ((uint32(s[0]) << 24) + (uint32(s[1]) << 16) + (uint32(s[2]) << 8) + uint32(s[3])))
	}

	switch n {
	case 3:
		return Atom(embeddedTag | ((uint32(n) << 24) + (uint32(s[0]) << 16) + (uint32(s[1]) << 8) + uint32(s[2])))
	case 2:
		return Atom(embeddedTag | ((uint32(n) << 24) + (uint32(s[0]) << 8) + uint32(s[1])))
	case 1:
		return Atom(embeddedTag | ((uint32(n) << 24) + uint32(s[0])))
	}

	return Empty
}

func hashAtom(s []byte) uint64 {
	h := fnv.New64a()

	h.Write(s)

	return h.Sum64()
}

func cacheAtom(a Atom) Atom {
	h := hashAtom(a.Bytes())

	atomLock.Lock()
	atomCache[h] = a
	atomLock.Unlock()

	return a
}

func findAtomInData(s []byte) Atom {
	off := bytes.Index(atomData.Bytes(), s)

	if off != -1 {
		return cacheAtom(newAtom(off, len(s)))
	}

	return Empty
}

func findAtomInCache(s []byte) Atom {
	h := hashAtom([]byte(s))

	atomLock.RLock()
	a, exists := atomCache[h]
	atomLock.RUnlock()

	if exists {
		return a
	}

	return Empty
}

func addAtom(s []byte) Atom {
	n, _ := atomData.Write(s)

	return cacheAtom(newAtom(atomData.Len()-n, n))
}

// Lookup returns the atom whose name is s.
//
// It returns Empty if there is no such atom. The lookup is case sensitive.
func Lookup(s string) Atom {
	n := len(s)

	if n == 0 || n > MaxAtomLen {
		return Empty
	}

	if a := embedAtom([]byte(s)); !a.IsEmpty() {
		return a
	}

	if a := findAtomInCache([]byte(s)); !a.IsEmpty() {
		return a
	}

	return findAtomInData([]byte(s))
}

// New return an exists atom or create it whose name is s.
//
// It returns Empty if s is longer than MaxAtomLen
func New(s string) Atom {
	n := len(s)

	if n == 0 || n > MaxAtomLen {
		return Empty
	}

	if a := Lookup(s); !a.IsEmpty() {
		return a
	}

	return addAtom([]byte(s))
}

// Save save atoms data and cache
func Save() ([]byte, Cache) {
	cache := make(Cache)

	atomLock.RLock()

	for hash, atom := range atomCache {
		cache[hash] = atom
	}

	atomLock.RUnlock()

	return atomData.Bytes(), cache
}

// Load load atoms data and cache
func Load(data []byte, cache Cache) {
	reset(data, cache)
}
