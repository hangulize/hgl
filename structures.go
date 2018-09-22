package hgl

import (
	"fmt"
)

// HGL is a decoding result of an HGL code.
type HGL map[string]Section

// Liner holds the line number.
//
// Section and Pair is a liner.
//
type Liner interface {
	Line() int
}

// Section contains pairs.
type Section interface {
	Liner
	Pairs() []Pair
	addPair(string, []string) error
}

// -----------------------------------------------------------------------------
// Pair

// Pair is a left-right tuple:
//
//  aa -> "ㅏ", "ㅐ"
//  ^^^^^^^^^^^^^^^^
//
type Pair struct {
	l string
	r []string

	line int
}

func (p *Pair) String() string {
	return fmt.Sprintf("Pair{%#v, %#v}", p.l, p.r)
}

// Left is a string. It is used for as keys in dict:
//
//  english = "Italian"
//  ^^^^^^^
//
// Or as left of pair:
//
//  aa -> "ㅏ", "ㅐ"
//  ^^
//
func (p Pair) Left() string {
	return p.l
}

// Right is a string array. It is used for as values in dict:
//
//  english = "Italian"
//	          ^^^^^^^^^
//
// Or as right of pair:
//
//  aa -> "ㅏ", "ㅐ"
//        ^^^^^^^^^^
//
func (p Pair) Right() []string {
	return p.r
}

// Line returns the line number where the pair is defined.
func (p *Pair) Line() int {
	return p.line
}

// -----------------------------------------------------------------------------
// ListSection

// ListSection has an ordered list of pairs.
type ListSection struct {
	pairs []Pair
	line  int
}

// newListSection creates an empty list section.
func newListSection(line int) *ListSection {
	return &ListSection{make([]Pair, 0), line}
}

// Pairs returns underlying pairs as an array.
func (s *ListSection) Pairs() []Pair {
	return s.pairs
}

// addPair adds a pair into a list section. It never fails.
func (s *ListSection) addPair(l string, r []string, line int) error {
	s.pairs = append(s.pairs, Pair{l, r, line})
	return nil
}

// Line returns the line number where the list section is defined.
func (s *ListSection) Line() int {
	return s.line
}

// Array returns the underying pair array of a list section.
func (s *ListSection) Array() []Pair {
	return s.pairs
}

// -----------------------------------------------------------------------------
// DictSection

// DictSection has an unordered list of pairs.
// Each left of underlying pairs is unique.
type DictSection struct {
	dict map[string]Pair
	line int
}

// newDictSection creates an empty dict section.
func newDictSection(line int) *DictSection {
	return &DictSection{make(map[string]Pair), line}
}

// Pairs returns dict key-values as an array of pairs.
func (s *DictSection) Pairs() []Pair {
	pairs := make([]Pair, len(s.dict))

	i := 0
	for _, pair := range s.dict {
		pairs[i] = pair
		i++
	}

	return pairs
}

// addPair adds a pair into a dict section. If there's already a pair having
// same left, it will fails.
func (s *DictSection) addPair(l string, r []string, line int) error {
	_, ok := s.dict[l]
	if ok {
		return fmt.Errorf("left of pair duplicated: %#v", l)
	}

	s.dict[l] = Pair{l, r, line}
	return nil
}

// Line returns the line number where the dict section is defined.
func (s *DictSection) Line() int {
	return s.line
}

// Map returns the underying map of a dict section.
func (s *DictSection) Map() map[string][]string {
	m := make(map[string][]string, len(s.dict))

	for _, pair := range s.dict {
		m[pair.Left()] = pair.Right()
	}

	return m
}

// Injective returns the underying 1-to-1 map of a dict section.
// If some right (values) has multiple values, it returns an error.
func (s *DictSection) Injective() (map[string]string, error) {
	oneToOne := make(map[string]string, len(s.dict))

	for _, pair := range s.dict {
		right := pair.Right()

		if len(right) != 1 {
			err := fmt.Errorf("right %#v has multiple values", right)
			return nil, err
		}

		oneToOne[pair.Left()] = right[0]
	}

	return oneToOne, nil
}

// One assumes the given left (key) has only one right (values). Then returns
// the only right value.
func (s *DictSection) One(left string) string {
	pair, ok := s.dict[left]

	if !ok {
		return ""
	}

	right := pair.Right()

	if len(right) == 0 {
		return ""
	}

	return right[0]
}

// All returns the right values.
func (s *DictSection) All(left string) []string {
	pair, ok := s.dict[left]

	if !ok {
		return make([]string, 0)
	}

	return pair.Right()
}
