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

package ekit

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIfThenElse(t *testing.T) {
	i := 7
	i = IfThenElse(false, i, 0)
	assert.Equal(t, i, 0)
}

func ExampleIfThenElse() {
	result := IfThenElse(true, "yes", "no")
	fmt.Println(result)

	// Output:
	// yes
}

func TestIfThenElseFunc(t *testing.T) {
	err := IfThenElseFunc(true, func() error {
		return nil
	}, func() error {
		return errors.New("some error")
	})
	assert.NoError(t, err)
}

func ExampleIfThenElseFunc() {
	err := IfThenElseFunc(false, func() error {
		// do something when condition is true
		// ...
		return nil
	}, func() error {
		// do something when condition is false
		// ...
		return errors.New("some error when execute func2")
	})
	fmt.Println(err)

	// Output:
	// some error when execute func2
}
