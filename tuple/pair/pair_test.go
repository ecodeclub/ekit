// Copyright 2021 ecodeclub
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package pair_test

import (
	"fmt"
	"sort"
	"testing"

	"github.com/ecodeclub/ekit/mapx"
	"github.com/ecodeclub/ekit/tuple/pair"
	"github.com/stretchr/testify/suite"
)

type testPairSuite struct{ suite.Suite }

func (s *testPairSuite) TestString() {
	{
		p := pair.NewPair(100, "23333")
		s.Assert().Equal("<100, \"23333\">", p.String())
	}
	{
		p := pair.NewPair("testStruct", map[int]int{
			11: 1,
			22: 2,
			33: 3,
		})
		s.Assert().Equal("<\"testStruct\", map[int]int{11:1, 22:2, 33:3}>", p.String())
	}
}

func (s *testPairSuite) TestNewPairs() {
	type caseType struct {
		// input
		keys   []int
		values []string
		// expected
		pairs []pair.Pair[int, string]
		err   error
	}
	for _, c := range []caseType{
		{
			keys:   []int{1, 2, 3, 4, 5},
			values: []string{"1", "2", "3", "4", "5"},
			pairs: []pair.Pair[int, string]{
				pair.NewPair(1, "1"),
				pair.NewPair(2, "2"),
				pair.NewPair(3, "3"),
				pair.NewPair(4, "4"),
				pair.NewPair(5, "5"),
			},
			err: nil,
		},
		{
			keys:   nil,
			values: []string{"1"},
			pairs:  nil,
			err:    fmt.Errorf("keys与values均不可为nil"),
		},
		{
			keys:   []int{1},
			values: nil,
			pairs:  nil,
			err:    fmt.Errorf("keys与values均不可为nil"),
		},
		{
			keys:   nil,
			values: nil,
			pairs:  nil,
			err:    fmt.Errorf("keys与values均不可为nil"),
		},
		{
			keys:   []int{1, 2},
			values: []string{"1"},
			pairs:  nil,
			err:    fmt.Errorf("keys与values的长度不同, len(keys)=2, len(values)=1"),
		},
	} {
		pairs, err := pair.NewPairs(c.keys, c.values)
		s.Assert().Equal(c.err, err)
		s.Assert().EqualValues(c.pairs, pairs)
	}
}

func (s *testPairSuite) TestSplitPairs() {
	type caseType struct {
		// input
		pairs []pair.Pair[int, string]
		// expected
		keys   []int
		values []string
	}
	for _, c := range []caseType{
		{
			pairs: []pair.Pair[int, string]{
				pair.NewPair(1, "1"),
				pair.NewPair(2, "2"),
				pair.NewPair(3, "3"),
				pair.NewPair(4, "4"),
				pair.NewPair(5, "5"),
			},
			keys:   []int{1, 2, 3, 4, 5},
			values: []string{"1", "2", "3", "4", "5"},
		},
		{
			pairs: nil,

			keys:   nil,
			values: nil,
		},
		{
			pairs:  []pair.Pair[int, string]{},
			keys:   []int{},
			values: []string{},
		},
	} {
		keys, values := pair.SplitPairs(c.pairs)
		if c.pairs == nil {
			s.Assert().Nil(keys)
			s.Assert().Nil(values)
		} else {
			s.Assert().Len(keys, len(c.pairs))
			s.Assert().Len(values, len(c.pairs))
			for i, pair := range c.pairs {
				s.Assert().Equal(pair.Key, keys[i])
				s.Assert().Equal(pair.Value, values[i])
			}
		}
	}
}

func (s *testPairSuite) TestFlattenPairs() {
	type caseType struct {
		pairs      []pair.Pair[int, string]
		flattPairs []any
	}

	for _, c := range []caseType{
		{
			pairs: []pair.Pair[int, string]{
				pair.NewPair(1, "1"),
				pair.NewPair(2, "2"),
				pair.NewPair(3, "3"),
				pair.NewPair(4, "4"),
				pair.NewPair(5, "5"),
			},
			flattPairs: []any{1, "1", 2, "2", 3, "3", 4, "4", 5, "5"},
		},
		{
			pairs:      nil,
			flattPairs: nil,
		},
		{
			pairs:      []pair.Pair[int, string]{},
			flattPairs: []any{},
		},
	} {
		flatPairs := pair.FlattenPairs(c.pairs)
		s.Assert().EqualValues(c.flattPairs, flatPairs)
	}
}

func (s *testPairSuite) TestPackPairs() {
	type caseType struct {
		flattPairs []any
		pairs      []pair.Pair[int, string]
	}

	for _, c := range []caseType{
		{
			flattPairs: []any{1, "1", 2, "2", 3, "3", 4, "4", 5, "5"},
			pairs: []pair.Pair[int, string]{
				pair.NewPair(1, "1"),
				pair.NewPair(2, "2"),
				pair.NewPair(3, "3"),
				pair.NewPair(4, "4"),
				pair.NewPair(5, "5"),
			},
		},
		{
			flattPairs: nil,
			pairs:      nil,
		},
		{
			flattPairs: []any{},
			pairs:      []pair.Pair[int, string]{},
		},
	} {
		pairs := pair.PackPairs[int, string](c.flattPairs)
		s.Assert().EqualValues(c.pairs, pairs)
	}
}

func (s *testPairSuite) TestMapPairMapping() {
	// map to pairs
	expectedMap := map[int]string{
		1: "1",
		2: "2",
		3: "3",
	}
	expectedPairs := []pair.Pair[int, string]{
		pair.NewPair(1, "1"),
		pair.NewPair(2, "2"),
		pair.NewPair(3, "3"),
	}

	// 可以用这种方式实现map到[]Pair的映射
	pairs, err := pair.NewPairs(mapx.KeysValues(expectedMap))
	s.Assert().Nil(err)
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].Key < pairs[j].Key
	})
	s.Assert().EqualValues(expectedPairs, pairs)

	// 可以用这种方式实现[]Pair到map的映射
	mp, err := mapx.ToMap(pair.SplitPairs(expectedPairs))
	s.Assert().Nil(err)
	for k, v := range mp {
		s.Assert().Equal(expectedMap[k], v)
	}
}

func TestPair(t *testing.T) {
	suite.Run(t, new(testPairSuite))
}
