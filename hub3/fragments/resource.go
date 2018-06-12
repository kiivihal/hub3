// Copyright © 2017 Delving B.V. <info@delving.eu>
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

package fragments

import (
	fmt "fmt"
	"log"
	"net/url"
	"strings"

	c "github.com/delving/rapid-saas/config"
	r "github.com/kiivihal/rdf2go"
	elastic "gopkg.in/olivere/elastic.v5"
)

// FragmentReferrerContext holds the referrer in formation for creating new fragments
type FragmentReferrerContext struct {
	Subject      string   `json:"subject"`
	SubjectClass []string `json:"subjectClass"`
	Predicate    string   `json:"predicate"`
	SearchLabel  string   `json:"searchLabel"`
	Level        int      `json:"level"`
	ObjectID     string   `json:"objectID"`
	// todo: decide if the sortKey belongs here
	//SortKey         int      `json:"sortKey"`
}

// NewContext returns the context for the current fragmentresource
func (fr *FragmentResource) NewContext(predicate, objectID string) *FragmentReferrerContext {
	label, err := c.Config.NameSpaceMap.GetSearchLabel(predicate)
	if err != nil {
		log.Printf("Unable to create search label for %s  due to %s\n", predicate, err)
		label = ""
	}

	return &FragmentReferrerContext{
		Subject:      fr.ID,
		SubjectClass: fr.Types,
		Predicate:    predicate,
		Level:        fr.GetLevel(),
		ObjectID:     objectID,
		SearchLabel:  label,
	}
}

// ResourceMap is a convenience structure to hold the resourceMap data and functions
type ResourceMap struct {
	resources map[string]*FragmentResource `json:"resources"`
}

// FragmentResource holds all the conttext information for a resource
// It works together with the FragmentBuilder to create the linked fragments
type FragmentResource struct {
	ID                   string                      `json:"id"`
	Types                []string                    `json:"types"`
	GraphExternalContext []*FragmentReferrerContext  `json:"graphExternalContext"`
	Context              []*FragmentReferrerContext  `json:"context"`
	Predicates           map[string][]*FragmentEntry `json:""`
	ObjectIDs            []*FragmentReferrerContext  `json:"objectIDs"`
}

// GetLabel returns the label and language for a resource
// This is used to present a label for a link in the interface
func (fr *FragmentResource) GetLabel() (label, language string) {
	if fr.ID == "" {
		return "", ""
	}
	labels := []string{
		"http://www.w3.org/2004/02/skos/core#prefLabel",
		"http://xmlns.com/foaf/0.1/name",
	}
	for _, labelPredicate := range labels {
		o, ok := fr.Predicates[labelPredicate]
		if ok && len(o) != 0 {
			return o[0].Value, o[0].Language
		}
	}
	return "", ""
}

// GetSubject returns the root FragmentResource based on its subject URI
// todo: remove this function later
func (rm *ResourceMap) GetSubject(uri string) (*FragmentResource, bool) {
	subject, ok := rm.GetResource(uri)
	return subject, ok
}

// SetContextLevels sets FragmentReferrerContext to each level from the root
func (rm *ResourceMap) SetContextLevels(subjectURI string) error {
	subject, ok := rm.GetResource(subjectURI)
	if !ok {
		return fmt.Errorf("Subject %s is not part of the graph", subjectURI)
	}
	for _, level1 := range subject.ObjectIDs {
		level2Resource, ok := rm.GetResource(level1.ObjectID)
		if !ok {
			log.Printf("unknown target URI: %s", level1.ObjectID)
			continue
		}
		level1.Level = 2
		level2Resource.AppendContext(level1)

		// loop into the next level, i.e. level 3
		for _, level2 := range level2Resource.ObjectIDs {
			level2.Level = 3
			level3Resource, ok := rm.GetResource(level2.ObjectID)
			if !ok {
				log.Printf("unknown target URI: %s", level2.ObjectID)
				continue
			}
			level3Resource.AppendContext(level1, level2)
		}
	}

	return nil
}

// AppendContext adds the referrerContext to the FragmentResource
// This action increments the level count
func (fr *FragmentResource) AppendContext(ctxs ...*FragmentReferrerContext) {
	for _, ctx := range ctxs {
		//ctx.Level = fr.GetLevel()
		if !containsContext(fr.Context, ctx) {
			fr.Context = append(fr.Context, ctx)
		}
	}
}

/*


 TODO: restructure fragments into blocks with header, geoblock, context block (maybe nested)

*/

// FragmentEntry holds all the information for the object of a rdf2go.Triple
type FragmentEntry struct {
	ID        string            `json:"@id,omitempty"`
	Value     string            `json:"@value,omitempty"`
	Language  string            `json:"@language,omitempty"`
	Datatype  string            `json:"@type,omitempty"`
	Entrytype string            `json:"entrytype"`
	Triple    string            `json:"triple"`
	Inline    *FragmentResource `json:"inline"`
}

// NewResourceMap creates a map for all the resources in the rdf2go.Graph
func NewResourceMap(g *r.Graph) (*ResourceMap, error) {
	rm := &ResourceMap{make(map[string]*FragmentResource)}

	if g.Len() == 0 {
		return rm, fmt.Errorf("The graph cannot be empty")
	}

	for t := range g.IterTriples() {
		err := AppendTriple(rm.resources, t)
		if err != nil {
			return rm, err
		}
	}
	return rm, nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// containsContext determines if a FragmentReferrerContext is already part of list
// this deduplication is important to not provide false counts for context levels
func containsContext(s []*FragmentReferrerContext, e *FragmentReferrerContext) bool {
	for _, a := range s {
		if a.ObjectID == e.ObjectID && a.Predicate == e.Predicate {
			return true
		}
	}
	return false
}

// containtsEntry determines if a FragmentEntry is already part of a predicate list
func containsEntry(s []*FragmentEntry, e *FragmentEntry) bool {
	for _, a := range s {
		if a.ID == e.ID && a.Value == e.Value {
			return true
		}
	}
	return false
}

// debrack removes the brackets around a string representation of a triple
func debrack(s string) string {
	if len(s) < 2 {
		return s
	}
	if s[0] != '<' {
		return s
	}
	if s[len(s)-1] != '>' {
		return s
	}
	return s[1 : len(s)-1]
}

// CreateFragmentEntry creates a FragmentEntry from a triple
func CreateFragmentEntry(t *r.Triple) (*FragmentEntry, string) {
	entry := &FragmentEntry{Triple: t.String()}
	switch o := t.Object.(type) {
	case *r.Resource:
		id := r.GetResourceID(o)
		entry.ID = r.GetResourceID(o)
		entry.Entrytype = "Resource"
		return entry, id
	case *r.BlankNode:
		id := r.GetResourceID(o)
		entry.ID = r.GetResourceID(o)
		entry.Entrytype = "Bnode"
		return entry, id
	case *r.Literal:
		entry.Value = o.Value
		entry.Entrytype = "Literal"
		if o.Datatype != nil && len(o.Datatype.String()) > 0 {
			if o.Datatype.String() != "<http://www.w3.org/2001/XMLSchema#string>" {
				entry.Datatype = debrack(o.Datatype.String())
			}
		}
		if len(o.Language) > 0 {
			entry.Language = o.Language
		}
	}
	return entry, ""
}

// AppendTriple appends a triple to a subject map
func AppendTriple(resources map[string]*FragmentResource, t *r.Triple) error {
	id := t.GetSubjectID()
	fr, ok := resources[id]
	if !ok {
		fr = &FragmentResource{}
		fr.ID = id
		resources[id] = fr
		fr.Predicates = make(map[string][]*FragmentEntry)
	}

	ttype, ok := t.GetRDFType()
	if ok {
		if !contains(fr.Types, ttype) {
			fr.Types = append(fr.Types, ttype)
		}
		return nil
	}

	p := r.GetResourceID(t.Predicate)
	predicates, ok := fr.Predicates[p]
	if !ok {
		predicates = []*FragmentEntry{}
	}
	entry, fragID := CreateFragmentEntry(t)
	if fragID != "" {
		if fragID != id {
			ctx := fr.NewContext(p, fragID)
			if !containsContext(fr.ObjectIDs, ctx) {
				fr.ObjectIDs = append(fr.ObjectIDs, ctx)
			}
		}
	}
	if !containsEntry(predicates, entry) {
		fr.Predicates[p] = append(predicates, entry)
	}

	return nil
}

// Resources returns the map
func (rm *ResourceMap) Resources() map[string]*FragmentResource {
	return rm.resources
}

// GetResource returns a Fragment resource from the ResourceMap
func (rm *ResourceMap) GetResource(subject string) (*FragmentResource, bool) {
	fr, ok := rm.resources[subject]
	return fr, ok
}

// GetLevel returns the relative level that this resource has from the root
// or parent resource
func (fr *FragmentResource) GetLevel() int {
	return len(fr.Context) + 1
}

// CreateHeader Linked Data Fragment entry for ElasticSearch
// as described here: http://linkeddatafragments.org/.
//
// The goal of this document is to support Linked Data Fragments based resolving
// for all stored RDF triples in the Hub3 system.
func (fg *FragmentGraph) CreateHeader(docType string) *Header {
	h := &Header{
		OrgID:    fg.Meta.OrgID,
		Spec:     fg.Meta.Spec,
		Revision: fg.Meta.Revision,
		HubID:    fg.Meta.HubID,
		DocType:  docType,
	}
	return h
}

// AddTags adds a tag string to the tags array of the Header
func (h *Header) AddTags(tags ...string) {
	for _, tag := range tags {
		if !contains(h.Tags, tag) {
			h.Tags = append(h.Tags, tag)
		}
	}
}

// CreateLodKey returns the path including the # fragments from the subject URL
// This is used for the Linked Open Data resolving
func (fr *FragmentResource) CreateLodKey() (string, error) {
	u, err := url.Parse(fr.ID)
	if err != nil {
		return "", err
	}
	lodKey := u.Path
	if c.Config.LOD.SingleEndpoint == "" {
		lodResourcePrefix := fmt.Sprintf("/%s", c.Config.LOD.Resource)
		if !strings.HasPrefix(u.Path, lodResourcePrefix) {
			return "", nil
		}
		lodKey = strings.TrimPrefix(u.Path, lodResourcePrefix)
	}
	if u.Fragment != "" {
		lodKey = fmt.Sprintf("%s#%s", lodKey, u.Fragment)
	}
	return lodKey, nil
}

// NormalisedResource creates a unique BlankNode key
// Normal resources are returned as is.
//
// This function is used so that you can query via the Fragment API for
// unique BlankNodes
func (fg *FragmentGraph) NormalisedResource(uri string) string {
	if !strings.HasPrefix(uri, "_:") {
		return uri
	}
	return fmt.Sprintf("%s-%s", uri, CreateHash(fg.GetNamedGraphURI()))
}

// CreateFragments creates ElasticSearch documents for each
// RDF triple in the FragmentResource
func (fr *FragmentResource) CreateFragments(fg *FragmentGraph) ([]*Fragment, error) {
	fragments := []*Fragment{}

	lodKey, _ := fr.CreateLodKey()

	// add type links
	for _, ttype := range fr.Types {
		frag := &Fragment{
			Meta:          fg.CreateHeader(FragmentDocType),
			Subject:       fg.NormalisedResource(fr.ID),
			Predicate:     RDFType,
			Object:        ttype,
			NamedGraphURI: fg.GetNamedGraphURI(),
		}
		if strings.HasPrefix(fr.ID, "_:") {
			frag.Triple = fmt.Sprintf("%s <%s> <%s> .", frag.Subject, RDFType, ttype)
		} else {
			frag.Triple = fmt.Sprintf("<%s> <%s> <%s> .", fr.ID, RDFType, ttype)
		}
		frag.Meta.AddTags("typelink", "Resource")
		if lodKey != "" {
			frag.LodKey = lodKey
		}
		fragments = append(fragments, frag)
	}

	// add entries
	for predicate, entries := range fr.Predicates {
		for _, entry := range entries {
			frag := &Fragment{
				Meta:          fg.CreateHeader(FragmentDocType),
				Subject:       fg.NormalisedResource(fr.ID),
				Predicate:     predicate,
				DataType:      entry.Datatype,
				Language:      entry.Language,
				NamedGraphURI: fg.GetNamedGraphURI(),
			}
			if entry.ID != "" {
				frag.Object = fg.NormalisedResource(entry.ID)
			} else {
				frag.Object = entry.Value
			}
			frag.Triple = strings.Replace(entry.Triple, entry.ID, fg.NormalisedResource(entry.ID), -1)
			frag.Triple = strings.Replace(frag.Triple, fr.ID, frag.Subject, -1)
			frag.Meta.AddTags(entry.Entrytype)
			if lodKey != "" {
				frag.LodKey = lodKey
			}

			fragments = append(fragments, frag)
		}
	}
	return fragments, nil
}

// GetXSDLabel returns a namespaced label for the RDF datatype
func (fe *FragmentEntry) GetXSDLabel() string {
	return strings.Replace(fe.Datatype, "http://www.w3.org/2001/XMLSchema#", "xsd:", 1)
}

// IndexFragments updates the Fragments for standalone indexing and adds them to the Elasti BulkProcessorService
func (fb *FragmentBuilder) IndexFragments(p *elastic.BulkProcessor) error {
	rm, err := fb.ResourceMap()
	if err != nil {
		return err
	}

	for _, fr := range rm.Resources() {
		fragments, err := fr.CreateFragments(fb.FragmentGraph())
		if err != nil {
			return err
		}
		for _, frag := range fragments {
			err := frag.AddTo(p)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// AddGraphExternalContext

// InlineFragmentResources
