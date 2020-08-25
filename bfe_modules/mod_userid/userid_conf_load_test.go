// Copyright (c) 2019 The BFE Authors.
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

package mod_userid

import (
	"reflect"
	"testing"
)

func TestNewConfigFromFile(t *testing.T) {
	tests := []struct {
		name        string
		fileName    string
		want        *Config
		valdateFunc func(a *Config) bool
		wantErr     bool
	}{
		{
			name:     "case:succ",
			fileName: "./testdata/mod_userid/userid_rule.data",
			valdateFunc: func(a *Config) bool {
				return a != nil && len(a.Products) == 1
			},
			wantErr: false,
		},
		{
			name:     "case:fail bad file name",
			fileName: "./testdata/mod_userid/userid_rule.data_not existed",
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewConfigFromFile(tt.fileName)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewConfigFromFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.valdateFunc != nil {
				if !tt.valdateFunc(got) {
					t.Errorf("NewConfigFromFile() = %v, want %v", got, nil)
				}
			} else if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewConfigFromFile() = %v, want %v", got, tt.want)
			}
		})
	}
}
