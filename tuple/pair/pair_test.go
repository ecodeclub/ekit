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
	"testing"

	"github.com/ecodeclub/ekit/tuple/pair"
	"github.com/stretchr/testify/suite"
)

type testPairSuite struct{ suite.Suite }

func (s *testPairSuite) TestGetterAndSetter() {
	p := pair.NewPair(100, "23333")
	s.Assert().Equal(100, p.First())
	s.Assert().Equal("23333", p.Second())
	p.SetSecond("10000")
	s.Assert().Equal("10000", p.Second())
	p.SetFirst(-1000)
	s.Assert().Equal(-1000, p.First())
}

func (s *testPairSuite) TestToString() {
	{
		p := pair.NewPair(100, "23333")
		s.Assert().Equal("<100, \"23333\">", p.ToString())
	}
	{
		p := pair.NewPair("testStruct", map[int]int{
			11: 1,
			22: 2,
			33: 3,
		})
		s.Assert().Equal("<\"testStruct\", map[int]int{11:1, 22:2, 33:3}>", p.ToString())
	}
}

func (s *testPairSuite) TestToArray() {
	p := pair.NewPair(100, "23333")
	arr := p.ToArray()
	s.Assert().Len(arr, 2)
	s.Assert().Equal(100, arr[0])
	s.Assert().Equal("23333", arr[1])
}

func (s *testPairSuite) TestCopy() {
	p := pair.NewPair(100, "23333")
	pcopy := p.Copy()
	s.Assert().Equal(p.First(), pcopy.First())
	s.Assert().Equal(p.Second(), pcopy.Second())
	pcopy.SetFirst(200)
	s.Assert().NotEqual(p.First(), pcopy.First())
}

func (s *testPairSuite) TestJson() {
	{
		// parse success.
		p := pair.NewPair(100, "23333")
		jsonByte, err := p.ToJson()
		s.Assert().Nil(err)
		pFromJson := pair.NewEmptyPair[int, string]()
		s.Assert().Nil(pFromJson.FromJson(jsonByte))
		s.Assert().Equal(p.First(), pFromJson.First())
		s.Assert().Equal(p.Second(), pFromJson.Second())
	}
	{
		// parse failed.
		p := pair.NewPair(100, "23333")
		s.Assert().NotNil(p.FromJson([]byte("errro string")))
		// FromJson失败不会影响原来的值
		s.Assert().Equal(100, p.First())
		s.Assert().Equal("23333", p.Second())
	}
}

func (s *testPairSuite) TestMergeFrom() {
	{
		p := pair.NewPair(map[int]int{
			11: 1,
			22: 2,
			33: 3,
		}, "23333")
		p.MergeFrom(pair.NewPair[map[int]int, string](nil, "23333"), false)
		s.Assert().Nil(p.First())
		s.Assert().Equal("23333", p.Second())
	}
	{
		p := pair.NewPair(map[int]int{
			11: 1,
			22: 2,
			33: 3,
		}, "23333")
		p.MergeFrom(pair.NewPair[map[int]int, string](nil, "23333"), true)
		s.Assert().Equal(map[int]int{
			11: 1,
			22: 2,
			33: 3,
		}, p.First())
		s.Assert().Equal("23333", p.Second())
	}
	{
		p := pair.NewPair("23333", map[int]int{
			11: 1,
			22: 2,
			33: 3,
		})
		p.MergeFrom(pair.NewPair[string, map[int]int]("23333", nil), false)
		s.Assert().Equal("23333", p.First())
		s.Assert().Nil(p.Second())
	}
	{
		p := pair.NewPair("23333", map[int]int{
			11: 1,
			22: 2,
			33: 3,
		})
		p.MergeFrom(pair.NewPair[string, map[int]int]("23333", nil), true)
		s.Assert().Equal("23333", p.First())
		s.Assert().Equal(map[int]int{
			11: 1,
			22: 2,
			33: 3,
		}, p.Second())
	}
}

func TestPair(t *testing.T) {
	suite.Run(t, new(testPairSuite))
}
