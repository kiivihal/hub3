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
	"context"
	fmt "fmt"
	"log"
	"net/url"
	"sort"
	"strings"

	c "github.com/delving/rapid-saas/config"
	"github.com/delving/rapid-saas/hub3/index"
	r "github.com/kiivihal/rdf2go"
	elastic "gopkg.in/olivere/elastic.v5"
)

const (
	literal  = "Literal"
	resource = "Resource"
	bnode    = "Bnode"
)

var ctx context.Context

func init() {
	ctx = context.Background()
}

// FragmentReferrerContext holds the referrer in formation for creating new fragments
//type FragmentReferrerContext struct {
//Subject      string   `json:"subject"`
//SubjectClass []string `json:"subjectClass"`
//Predicate    string   `json:"predicate"`
//SearchLabel  string   `json:"searchLabel"`
//Level        int      `json:"level"`
//ObjectID     string   `json:"objectID"`
//// todo: decide if the sortKey belongs here
////SortKey         int      `json:"sortKey"`
//}

// NewContext returns the context for the current fragmentresource
func (fr *FragmentResource) NewContext(predicate, objectID string) *FragmentReferrerContext {
	searchLabel, err := c.Config.NameSpaceMap.GetSearchLabel(predicate)
	if err != nil {
		log.Printf("Unable to create search label for %s  due to %s\n", predicate, err)
		searchLabel = ""
	}

	label, _ := fr.GetLabel()

	return &FragmentReferrerContext{
		Subject:      fr.ID,
		SubjectClass: fr.Types,
		Predicate:    predicate,
		Level:        fr.GetLevel(),
		ObjectID:     objectID,
		SearchLabel:  searchLabel,
		Label:        label,
	}
}

// ResourceMap is a convenience structure to hold the resourceMap data and functions
type ResourceMap struct {
	resources map[string]*FragmentResource `json:"resources"`
}

// FragmentGraph is a container for all entries of an RDF Named Graph
type FragmentGraph struct {
	Meta       *Header                   `json:"meta"`
	Resources  []*FragmentResource       `json:"resources,omitempty"`
	Summary    *ResultSummary            `json:"summary,omitempty"`
	JSONLD     []map[string]interface{}  `json:"jsonld,omitempty"`
	Fields     map[string][]string       `json:"fields,omitempty"`
	Highlights []*ResourceEntryHighlight `json:"highlights,omitempty"`
}

// ResourceEntryHighlight holds the values of the ElasticSearch highlight fiel
type ResourceEntryHighlight struct {
	SearchLabel string `json:"searchLabel"`
	MarkDown    string `json:"markdown"`
}

// GenerateJSONLD converts a FragmenResource into a JSON-LD entry
func (fr *FragmentResource) GenerateJSONLD() map[string]interface{} {
	m := map[string]interface{}{}
	m["@id"] = fr.ID
	if len(fr.Types) > 0 {
		m["@type"] = fr.Types
	}
	entries := map[string][]*ResourceEntry{}
	for _, p := range fr.Entries {
		entries[p.Predicate] = append(entries[p.Predicate], p)
	}
	for k, v := range entries {
		for _, p := range v {
			m[k] = p.AsLdObject()
		}
	}
	return m
}

// Collapsed holds each entry of a FieldCollapse elasticsearch result
type Collapsed struct {
	Field    string           `json:"field"`
	Title    string           `json:"title"`
	HitCount int64            `json:"hitCount"`
	Items    []*FragmentGraph `json:"items"`
}

// ScrollResultV4 intermediate non-protobuf search results
type ScrollResultV4 struct {
	Pager     *ScrollPager     `json:"pager"`
	Query     *Query           `json:"query"`
	Items     []*FragmentGraph `json:"items,omitempty"`
	Collapsed []*Collapsed     `json:"collapse,omitempty"`
	Peek      map[string]int64 `json:"peek,omitempty"`
	Facets    []*QueryFacet    `json:"facets,omitempty"`
}

// QueryFacet contains all the information for an ElasticSearch Aggregation
type QueryFacet struct {
	Name        string       `json:"name"`
	Field       string       `json:"field"`
	IsSelected  bool         `json:"isSelected"`
	I18n        string       `json:"i18N,omitempty"`
	Total       int64        `json:"total"`
	MissingDocs int64        `json:"missingDocs"`
	OtherDocs   int64        `json:"otherDocs"`
	Links       []*FacetLink `json:"links"`
}

// FacetLink contains all the information for creating a filter for this facet
type FacetLink struct {
	URL           string `json:"url"`
	IsSelected    bool   `json:"isSelected"`
	Value         string `json:"value"`
	DisplayString string `json:"displayString"`
	Count         int64  `json:"count"`
}

// FragmentResource holds all the conttext information for a resource
// It works together with the FragmentBuilder to create the linked fragments
type FragmentResource struct {
	ID                   string                      `json:"id"`
	Types                []string                    `json:"types"`
	GraphExternalContext []*FragmentReferrerContext  `json:"graphExternalContext"`
	Context              []*FragmentReferrerContext  `json:"context"`
	predicates           map[string][]*FragmentEntry `json:"predicates"`
	objectIDs            []*FragmentReferrerContext  `json:"objectIDs"`
	Entries              []*ResourceEntry            `json:"entries"`
	Tags                 []string                    `json:"tags,omitempty"`
}

// ObjectIDs returns an array of FragmentReferrerContext
func (fr *FragmentResource) ObjectIDs() []*FragmentReferrerContext {
	return fr.objectIDs
}

// Predicates returns a map of FragmentEntry
func (fr *FragmentResource) Predicates() map[string][]*FragmentEntry {
	return fr.predicates
}

// SetEntries sets the ResourceEntries for indexing
func (fr *FragmentResource) SetEntries(rm *ResourceMap) error {
	fr.Entries = []*ResourceEntry{}
	for predicate, entries := range fr.predicates {
		for _, entry := range entries {
			re, err := entry.NewResourceEntry(predicate, fr.GetLevel(), rm)
			if err != nil {
				return err
			}
			fr.Entries = append(fr.Entries, re)
		}
	}
	// sort entries by order
	sort.Slice(fr.Entries[:], func(i, j int) bool {
		return fr.Entries[i].Order < fr.Entries[j].Order
	})
	return nil
}

// AsLdObject generates an rdf2go.LdObject for JSON-LD generation
func (fe *FragmentEntry) AsLdObject() *r.LdObject {
	return &r.LdObject{
		ID:       fe.ID,
		Value:    fe.Value,
		Language: fe.Language,
		Datatype: fe.DataType,
	}
}

// NewResourceEntry creates a resource entry for indexing
func (fe *FragmentEntry) NewResourceEntry(predicate string, level int32, rm *ResourceMap) (*ResourceEntry, error) {
	label, err := c.Config.NameSpaceMap.GetSearchLabel(predicate)
	if err != nil {
		log.Printf("Unable to create search label for %s  due to %s\n", predicate, err)
		label = ""
	}
	re := &ResourceEntry{
		ID:          fe.ID,
		Value:       fe.Value,
		Language:    fe.Language,
		DataType:    fe.DataType,
		EntryType:   fe.EntryType,
		Predicate:   predicate,
		Level:       level,
		SearchLabel: label,
		Order:       fe.Order,
	}

	if re.ID != "" {
		r, ok := rm.GetResource(re.ID)
		if ok {
			re.Value, _ = r.GetLabel()
		}
	}

	// add label for resolved
	if fe.Resolved {
		re.AddTags("resolved")
	}

	labels, ok := c.Config.RDFTagMap.Get(predicate)
	if ok {
		re.AddTags(labels...)
		if re.Value != "" {
			// TODO add validation for the values here
			for _, label := range labels {
				switch label {
				case "isoDate":
					re.Date = re.Value
					//log.Printf("Date value: %s", re.Date)
				case "dateRange":
					re.DateRange = re.Value
				case "latLong":
					re.LatLong = re.Value
				}
			}
		}
	}
	return re, nil
}

// GetLabel returns the label and language for a resource
// This is used to present a label for a link in the interface
func (fr *FragmentResource) GetLabel() (label, language string) {
	if fr.ID == "" {
		return "", ""
	}
	for _, labelPredicate := range c.Config.RDFTag.Label {
		o, ok := fr.predicates[labelPredicate]
		if ok && len(o) != 0 {
			return o[0].Value, o[0].Language
		}
	}
	return "", ""
}

// SetContextLevels sets FragmentReferrerContext to each level from the root
func (rm *ResourceMap) SetContextLevels(subjectURI string) error {
	subject, ok := rm.GetResource(subjectURI)
	if !ok {
		return fmt.Errorf("Subject %s is not part of the graph", subjectURI)
	}

	for _, level1 := range subject.objectIDs {
		level2Resource, ok := rm.GetResource(level1.ObjectID)
		if !ok {
			log.Printf("unknown target URI: %s", level1.ObjectID)
			continue
		}
		level1.Level = 1
		if len(level1.GetSubjectClass()) == 0 {
			level1.SubjectClass = subject.Types
		}
		// validate context
		level2Resource.AppendContext(level1)

		// loop into the next level, i.e. level 3
		for _, level2 := range level2Resource.objectIDs {
			level2.Level = 2
			level3Resource, ok := rm.GetResource(level2.ObjectID)
			if !ok {
				log.Printf("unknown target URI: %s", level2.ObjectID)
				continue
			}
			if len(level2.GetSubjectClass()) == 0 {
				level2.SubjectClass = level2Resource.Types
			}
			level3Resource.AppendContext(level1, level2)
		}
	}

	return nil
}

// AppendContext adds the referrerContext to the FragmentResource
// This action increments nilthe level count
func (fr *FragmentResource) AppendContext(ctxs ...*FragmentReferrerContext) {
	for _, ctx := range ctxs {
		if !containsContext(fr.Context, ctx) {
			fr.Context = append(fr.Context, ctx)
		}
	}
}

// FragmentEntry holds all the information for the object of a rdf2go.Triple
type FragmentEntry struct {
	ID        string `json:"@id,omitempty"`
	Value     string `json:"@value,omitempty"`
	Language  string `json:"@language,omitempty"`
	DataType  string `json:"@type,omitempty"`
	EntryType string `json:"entrytype"`
	Triple    string `json:"triple"`
	Resolved  bool   `json:"resolved"`
	Order     int    `json:"order"`
}

// ResourceEntry contains all the indexed entries for FragmentResources
type ResourceEntry struct {
	ID          string            `json:"@id,omitempty"`
	Value       string            `json:"@value,omitempty"`
	Language    string            `json:"@language,omitempty"`
	DataType    string            `json:"@type,omitempty"`
	EntryType   string            `json:"entrytype,omitempty"`
	Predicate   string            `json:"predicate,omitempty"`
	SearchLabel string            `json:"searchLabel,omitempty"`
	Level       int32             `json:"level"`
	Tags        []string          `json:"tags,omitempty"`
	Date        string            `json:"date,omitempty"`
	DateRange   string            `json:"dateRange,omitempty"`
	LatLong     string            `json:"latLong,omitempty"`
	Inline      *FragmentResource `json:"inline,omitempty"`
	Order       int               `json:"order"`
}

// AsLdObject generates an rdf2go.LdObject for JSON-LD generation
func (re *ResourceEntry) AsLdObject() *r.LdObject {
	o := &r.LdObject{
		ID:       re.ID,
		Language: re.Language,
		Datatype: re.DataType,
	}
	if re.ID == "" {
		o.Value = re.Value
	}
	return o
}

// NewResourceMap creates a map for all the resources in the rdf2go.Graph
func NewResourceMap(g *r.Graph) (*ResourceMap, error) {
	rm := &ResourceMap{make(map[string]*FragmentResource)}

	if g.Len() == 0 {
		return rm, fmt.Errorf("The graph cannot be empty")
	}

	for t := range g.IterTriples() {
		err := rm.AppendTriple(t, false)
		if err != nil {
			return rm, err
		}
	}
	return rm, nil
}

// NewEmptyResourceMap returns an initialised ResourceMap
func NewEmptyResourceMap() *ResourceMap {
	return &ResourceMap{make(map[string]*FragmentResource)}
}

// ResolveObjectIDs queries the fragmentstore for additional context
func (rm *ResourceMap) ResolveObjectIDs(excludeHubID string) error {
	objectIDs := []string{}
	for _, fr := range rm.Resources() {
		if contains(fr.Types, "http://www.europeana.eu/schemas/edm/WebResource") {
			objectIDs = append(objectIDs, fr.ID)

		}
	}
	req := NewFragmentRequest()
	req.Subject = objectIDs
	req.ExcludeHubID = excludeHubID
	frags, _, err := req.Find(ctx, index.ESClient())
	if err != nil {
		log.Printf("unable to find fragments: %s", err.Error())
		return err
	}
	for _, f := range frags {
		t := f.CreateTriple()
		err = rm.AppendTriple(t, true)
		if err != nil {
			return err
		}
	}
	return nil
}

// CreateTriple creates a *rdf2go.Triple from a Fragment
func (f *Fragment) CreateTriple() *r.Triple {
	s := r.NewResource(f.Subject)
	p := r.NewResource(f.Predicate)
	var o r.Term

	switch f.ObjectType {
	case resource:
		o = r.NewResource(f.Object)
	case bnode:
		o = r.NewBlankNode(f.Object)
	default:
		if f.DataType != "" {
			o = r.NewLiteralWithDatatype(
				f.Object,
				r.NewResource(f.DataType),
			)
			t := r.NewTriple(s, p, o)
			return t
		}
		o = r.NewLiteralWithLanguage(f.Object, f.Language)
	}

	t := r.NewTriple(s, p, o)
	return t
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
func CreateFragmentEntry(t *r.Triple, resolved bool, order int) (*FragmentEntry, string) {
	entry := &FragmentEntry{Triple: t.String()}
	entry.Order = order
	entry.Resolved = resolved

	switch o := t.Object.(type) {
	case *r.Resource:
		id := r.GetResourceID(o)
		entry.ID = r.GetResourceID(o)
		entry.EntryType = resource
		return entry, id
	case *r.BlankNode:
		id := r.GetResourceID(o)
		entry.ID = r.GetResourceID(o)
		entry.EntryType = bnode
		return entry, id
	case *r.Literal:
		entry.Value = o.Value
		entry.EntryType = literal
		if o.Datatype != nil && len(o.Datatype.String()) > 0 {
			if o.Datatype.String() != "<http://www.w3.org/2001/XMLSchema#string>" {
				entry.DataType = debrack(o.Datatype.String())
			}
		}
		if len(o.Language) > 0 {
			entry.Language = o.Language
		}
	}
	return entry, ""
}

// AppendTriple appends a triple to a subject map
func (rm *ResourceMap) AppendTriple(t *r.Triple, resolved bool) error {
	return rm.AppendOrderedTriple(t, resolved, 0)
}

// AppendOrderedTriple appends a triple to a subject map
func (rm *ResourceMap) AppendOrderedTriple(t *r.Triple, resolved bool, order int) error {
	id := t.GetSubjectID()
	fr, ok := rm.resources[id]
	if !ok {
		fr = &FragmentResource{}
		fr.ID = id
		rm.resources[id] = fr
		fr.predicates = make(map[string][]*FragmentEntry)
	}

	ttype, ok := t.GetRDFType()
	if ok {
		if !contains(fr.Types, ttype) {
			fr.Types = append(fr.Types, ttype)
		}
		return nil
	}

	p := r.GetResourceID(t.Predicate)
	predicates, ok := fr.predicates[p]
	if !ok {
		predicates = []*FragmentEntry{}
	}

	entry, fragID := CreateFragmentEntry(t, resolved, order)
	if fragID != "" {
		if fragID != id {
			ctx := fr.NewContext(p, fragID)
			if !containsContext(fr.objectIDs, ctx) {
				fr.objectIDs = append(fr.objectIDs, ctx)
			}
		}
	}
	if !containsEntry(predicates, entry) {
		fr.predicates[p] = append(predicates, entry)
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
func (fr *FragmentResource) GetLevel() int32 {
	highestLevel := int32(0)
	for _, ctx := range fr.Context {
		if ctx.GetLevel() > highestLevel {
			highestLevel = ctx.GetLevel()
		}
	}
	return int32(highestLevel + 1)
}

// NewResultSummary creates a Summary from the FragmentGraph based on the
// RDFTag configuration.
func (fg *FragmentGraph) NewResultSummary() *ResultSummary {
	fg.Summary = &ResultSummary{}
	for _, rsc := range fg.Resources {
		for _, entry := range rsc.Entries {
			fg.Summary.AddEntry(entry)
		}

	}
	return fg.Summary
}

// NewFields returns a map of the triples sorted by their searchLabel
func (fg *FragmentGraph) NewFields() map[string][]string {
	fg.Fields = make(map[string][]string)
	for _, rsc := range fg.Resources {
		for _, entry := range rsc.Entries {
			fg.Fields[entry.SearchLabel] = append(fg.Fields[entry.SearchLabel], entry.Value)
		}
	}
	return fg.Fields
}

// NewJSONLD creates a JSON-LD version of the FragmentGraph
func (fg *FragmentGraph) NewJSONLD() []map[string]interface{} {
	fg.JSONLD = []map[string]interface{}{}
	for _, rsc := range fg.Resources {
		fg.JSONLD = append(fg.JSONLD, rsc.GenerateJSONLD())
	}
	return fg.JSONLD
}

// NewGrouped returns an inlined version of the FragmentResources in the FragmentGraph
func (fg *FragmentGraph) NewGrouped() (*FragmentResource, error) {
	rm := &ResourceMap{make(map[string]*FragmentResource)}

	// create the resource map
	for _, fr := range fg.Resources {
		rm.resources[fr.ID] = fr
	}

	// set the inlines
	for _, fr := range fg.Resources {
		for _, entry := range fr.Entries {
			if entry.ID != "" && fr.ID != entry.ID {
				target, ok := rm.GetResource(entry.ID)
				if ok {
					entry.Inline = target
				}
			}
		}
	}

	// only return the subject
	subject, ok := rm.GetResource(fg.GetAboutURI())

	if !ok {
		return nil, fmt.Errorf("unable to find root of the graph for %s", fg.GetAboutURI())
	}

	fg.Resources = []*FragmentResource{subject}
	return subject, nil
}

// AddEntry adds Summary fields based on the ResourceEntry tags
func (sum *ResultSummary) AddEntry(entry *ResourceEntry) {

	for _, tag := range entry.Tags {
		switch tag {
		case "title":
			if sum.Title == "" {
				sum.Title = entry.Value
			}
		case "thumbnail":
			if sum.Thumbnail == "" {
				sum.Thumbnail = entry.Value
			}
		case "subject":
			if sum.Subject == "" {
				sum.Subject = entry.Value
			}
		case "creator":
			if sum.Creator == "" {
				sum.Creator = entry.Value
			}
		case "description":
			if sum.Description == "" {
				sum.Description = entry.Value
			}
		case "landingPage":
			if sum.LandingPage == "" {
				sum.LandingPage = entry.Value
			}
		case "collection":
			if sum.Collection == "" {
				sum.Collection = entry.Value
			}
		case "subCollection":
			if sum.SubCollection == "" {
				sum.SubCollection = entry.Value
			}
		case "objectType":
			if sum.ObjectType == "" {
				sum.ObjectType = entry.Value
			}
		case "objectID":
			if sum.ObjectID == "" {
				sum.ObjectID = entry.Value
			}
		case "owner":
			if sum.Owner == "" {
				sum.Owner = entry.Value
			}
		case "date":
			if sum.Date == "" {
				sum.Date = entry.Value
			}
		}
	}
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
		Modified: NowInMillis(),
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

// AddTags adds a tag string to the tags array of the Header
func (re *ResourceEntry) AddTags(tags ...string) {
	for _, tag := range tags {
		if !contains(re.Tags, tag) {
			re.Tags = append(re.Tags, tag)
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
// unique BlankNodesThe named graph that this triple is part of
func (fg *FragmentGraph) NormalisedResource(uri string) string {
	if !strings.HasPrefix(uri, "_:") {
		return uri
	}
	return fmt.Sprintf("%s-%s", uri, CreateHash(fg.Meta.NamedGraphURI))
}

// CreateFragments creates ElasticSearch documents for each
// RDF triple in the FragmentResource
func (fr *FragmentResource) CreateFragments(fg *FragmentGraph) ([]*Fragment, error) {
	fragments := []*Fragment{}

	lodKey, _ := fr.CreateLodKey()

	// TODO add statistics path
	// type is searchLabel
	// @about is extra entry
	// add type links
	for _, ttype := range fr.Types {
		frag := &Fragment{
			Meta:       fg.CreateHeader(FragmentDocType),
			Subject:    fg.NormalisedResource(fr.ID),
			Predicate:  RDFType,
			Object:     ttype,
			ObjectType: resource,
		}
		frag.Meta.NamedGraphURI = fg.Meta.NamedGraphURI
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
	for predicate, entries := range fr.predicates {
		for _, entry := range entries {
			frag := &Fragment{
				Meta:       fg.CreateHeader(FragmentDocType),
				Subject:    fg.NormalisedResource(fr.ID),
				Predicate:  predicate,
				DataType:   entry.DataType,
				Language:   entry.Language,
				ObjectType: entry.EntryType,
				Order:      int32(entry.Order),
			}
			frag.Meta.NamedGraphURI = fg.Meta.NamedGraphURI
			if entry.ID != "" {
				frag.Object = fg.NormalisedResource(entry.ID)
			} else {
				frag.Object = entry.Value
			}
			frag.Triple = strings.Replace(entry.Triple, entry.ID, fg.NormalisedResource(entry.ID), -1)
			frag.Triple = strings.Replace(frag.Triple, fr.ID, frag.Subject, -1)
			frag.Meta.AddTags(entry.EntryType)
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
	return strings.Replace(fe.DataType, "http://www.w3.org/2001/XMLSchema#", "xsd:", 1)
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