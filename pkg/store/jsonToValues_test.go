// Copyright 2019-present Open Networking Foundation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package store

import (
	"fmt"
	"gotest.tools/assert"
	"io/ioutil"
	"regexp"
	"strings"
	"testing"
)

func setUpJSONToValues() ([]byte, error) {
	sampleTree, err := ioutil.ReadFile("./testdata/sample-tree.json")
	if err != nil {
		return nil, err
	}
	return sampleTree, nil
}

func Test_DecomposeTree(t *testing.T) {
	const matchAllChars = `(/[a-zA-Z0-9\-:\[\]]*)*`
	const matchOnIndex = `(\[[0-9]]).*?`
	sampleTree, err := setUpJSONToValues()
	assert.NilError(t, err)

	assert.Assert(t, len(sampleTree) > 0, "Empty sample tree", (len(sampleTree)))

	values, err := DecomposeTree(sampleTree)
	assert.NilError(t, err)

	for _, v := range values {
		fmt.Printf("%s %s\n", (*v).Path, (*v).String())
	}
	assert.Equal(t, len(values), 22)

	rAllChars := regexp.MustCompile(matchAllChars)
	rOnIndex := regexp.MustCompile(matchOnIndex)
	for _, v := range values {
		match := rAllChars.FindString(v.Path)
		assert.Equal(t, match, v.Path)

		matches := rOnIndex.FindAllStringSubmatch(v.Path, -1)
		newPath := v.Path
		for _, m := range matches {
			newPath = strings.Replace(newPath, m[1], "[*]", -1)
		}

		switch newPath {
		case
			"/openconfig-interfaces:interfaces/default-type",
			"/openconfig-interfaces:interfaces/interface[*]/type",
			"/openconfig-interfaces:interfaces/interface[*]/name",
			"/openconfig-system:system/openconfig-openflow:openflow/controllers/controller[*]/connections/connection[*]/aux-id",
			"/openconfig-system:system/openconfig-openflow:openflow/controllers/controller[*]/connections/connection[*]/discombobulator",
			"/openconfig-system:system/openconfig-openflow:openflow/controllers/controller[*]/name",
			"/openconfig-system:system/openconfig-openflow:openflow/controllers/controller[*]/type",
			"/openconfig-system:system/openconfig-openflow:openflow/controllers/controller[*]/connections/connections-type",
			"/openconfig-system:system/openconfig-openflow:openflow/controllers/controller[*]/connections/connections-freq",
			"/openconfig-system:system/openconfig-openflow:openflow/controllers/controller[*]/connections/example-ll[*]",
			"/openconfig-system:system/openconfig-openflow:openflow/controllers/controller[*]/connections/connection[*]/conn-type":
			//NOOP
		default:
			t.Fatal("Unexpected jsonPath", newPath)
		}
	}
}