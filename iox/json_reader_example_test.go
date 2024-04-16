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

package iox_test

import (
	"fmt"
	"net/http"

	"github.com/ecodeclub/ekit/iox"
)

func ExampleNewJSONReader() {
	val := iox.NewJSONReader(User{Name: "Tom"})
	_, err := http.NewRequest(http.MethodPost, "/hello", val)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("OK")
}

type User struct {
	Name string `json:"name"`
}
