// Copyright 2020 Delving B.V.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package search

// Filter is used to limit the results of a SearchRequest.
//
// It supports both first level objects such as 'Meta' and 'Tree' and nested
// items from resource.entries via a NestedFilter.
//
//
type Filter struct {
	// SearchLabel is a short namespaced version of a URI.
	Field string `json:"searchLabel,omitempty"`
	Value string `json:"value,omitempty"`
	//Exclude     bool          `json:"exclude,omitempty"`
	//Nested      *NestedFilter `json:"nested,omitempty"`
}

// NestedFilter is used to filter in the nested RDF structure of the RecordGraph.
type NestedFilter struct {
	//SearchLabel string              `json:"searchLabel,omitempty"`
	//Value       string              `json:"value,omitempty"`
	//Level1      *ContextQueryFilter `json:"level1,omitempty"`
	//Level2      *ContextQueryFilter `json:"level2,omitempty"`
	//TypeClass   string              `json:"typeClass,omitempty"`
	//ID          bool                `json:"id,omitempty"`
	//Type        QueryFilterType     `json:"type,omitempty"`
	//Lte         string              `json:"lte,omitempty"`
	//Gte         string              `json:"gte,omitempty"`

}

// ContextFilter is used to specify the path to filter the nested resources.
// TypeClass is optional and can be used to specify the RDF class of the resource.
type ContextFilter struct {
	SearchLabel string `json:"SearchLabel,omitempty"`
	TypeClass   string `json:"TypeClass,omitempty"`
}
