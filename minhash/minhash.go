package minhash

import (
	"io"
	"math"
	"sync"

	mapset "github.com/deckarep/golang-set"
	"github.com/suntong/deduper/text"
)

const (
	// The large prime used to hash shingles
	p1 = uint64(4294967311)

	// The large prime used to hash bands
	p2 = uint64(7562380294967317)
)

type hasher func(...uint32) uint32

// Match represents a matching document.
type Match struct {
	// ID is the unique ID of the document that was
	// given when the document was added.
	ID string `json:"id"`

	// Similarity is the Jaccard similarity from 0 to 1 of this document
	// to the document it was compared against.
	Similarity float64 `json:"similarity"`
}

// New creates a new MinHasher with the given band size, number of rows, and shingle size.
func New(B int, R int, shingleSize int) *MinHasher {
	return &MinHasher{
		hashers:       generateHahsers(B*R, p1),
		bandHashers:   generateHahsers(B, p2),
		Matrix:        make(Matrix, 0),
		R:             R,
		B:             B,
		N:             shingleSize,
		ColumnMapping: make(map[int]string),
		Ids:           mapset.NewSet(),
	}
}

// MinHasher provides near-similar matching capabilities on large
// strings of text.
type MinHasher struct {
	// The mapping of column indexes in the matrix to document Ids.
	ColumnMapping map[int]string

	// The hash functions used to hash the document's shingles.
	hashers []hasher

	// The hash functions used to hash the hash function results
	// into bands.
	bandHashers []hasher

	// The matrix of documents and hash values. Each vector
	// is a list of hash values for a document's shingles, eg element m[i,j] is the
	// value of h[i](document[j]).
	Matrix Matrix

	// The unique list of document Ids being stored.
	Ids mapset.Set

	// The band matrix generated with LSH.
	Bands Matrix

	// Locks the bands matrix.
	bandMutex sync.RWMutex

	// Locks the matrix.
	matrixMutex sync.RWMutex

	// Number of bands.
	B int

	// Number of rows.
	R int

	// N-shingles being used.
	N int
}

// Add adds a new document with the given ID to the collection of
// documents.
func (m *MinHasher) Add(id string, R io.Reader) {
	column := m.hashColumn(R)

	m.matrixMutex.Lock()
	m.Matrix = append(m.Matrix, column)
	m.ColumnMapping[len(m.Matrix)-1] = id
	m.matrixMutex.Unlock()

	m.Ids.Add(id)

	m.bandMutex.Lock()
	m.Bands = nil
	m.bandMutex.Unlock()
}

// FindSimilar returns a list of documents whose similarity to the given document
// is greater than or equal to the threshold provided.
func (m *MinHasher) FindSimilar(R io.Reader, threshold float64) []Match {
	col := m.hashColumn(R)
	col = m.bandColumn(col)

	similar := make([]Match, 0)

	m.bandMutex.RLock()
	if m.Bands == nil {
		m.bandMutex.RUnlock()

		// lock as writer
		m.bandMutex.Lock()
		m.Bands = m.bandMatrix()
		m.bandMutex.Unlock()

		// relock as reader
		m.bandMutex.RLock()
	}

	// for each document in the band Matrix
	for i, c := range m.Bands {
		// see if they share any common Bands with input
		for j := 0; j < len(col); j++ {
			if col[j] == c[j] {
				// needs deeper inspection ie jaccard similarity
				sim := jaccard(c, col)

				if sim >= threshold {
					similar = append(similar, Match{
						ID:         m.ColumnMapping[i],
						Similarity: sim,
					})
				}

				break
			}
		}
	}

	m.bandMutex.RUnlock()

	return similar
}

// Contains returns true if the MinHasher contains
// the document with the given id.
func (m *MinHasher) Contains(id string) bool {
	return m.Ids.Contains(id)
}

func (m *MinHasher) hashColumn(R io.Reader) vector {
	// the result which holds each minimum hash
	// value of h_i at the ith index of each N-gram
	column := make(vector, len(m.hashers))

	shingler := text.NewShingler(R, m.N)

	// initialize to max value to find the min
	for i, _ := range m.hashers {
		column[i] = uint32(math.MaxUint32)
	}

	for shingler.Scan() {
		sh := shingler.Text()

		// convert the string to a number by
		// hashing it... similar to GetHashCode
		// in C#
		v := hashCode(sh)

		for i, h := range m.hashers {
			hash := h(v)
			if hash < column[i] {
				column[i] = hash
			}
		}
	}

	return column
}

func (m *MinHasher) bandColumn(col vector) vector {
	bcol := make(vector, m.B)

	for i, hash := range m.bandHashers {
		for j := 0; j < len(col); j += m.R {
			rows := col[j : j+m.R]
			h := hash(rows...)

			bcol[i] = h
		}
	}

	return bcol
}

func (m *MinHasher) bandMatrix() Matrix {
	m.matrixMutex.RLock()
	defer m.matrixMutex.RUnlock()

	B := make(Matrix, len(m.Matrix))

	for i, col := range m.Matrix {
		bcol := m.bandColumn(col)
		B[i] = bcol
	}

	return B
}
