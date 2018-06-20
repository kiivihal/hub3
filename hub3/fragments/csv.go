package fragments

import (
	"encoding/csv"
	fmt "fmt"
	"io"
	"strings"

	"github.com/delving/rapid-saas/config"
	r "github.com/kiivihal/rdf2go"
	elastic "gopkg.in/olivere/elastic.v5"
)

// CSVConvertor holds all values to convert a CSV to RDF
type CSVConvertor struct {
	SubjectColumn         string    `json:"subjectColumn"`
	Separator             string    `json:"separator"`
	PredicateURIBase      string    `json:"predicateURIBase"`
	SubjectClass          string    `json:"subjectClass"`
	SubjectURIBase        string    `json:"subjectURIBase"`
	ObjectURIFormat       string    `json:"objectURIFormat"`
	ObjectResourceColumns []string  `json:"objectResourceColumns"`
	DefaultSpec           string    `json:"defaultSpec"`
	InputFile             io.Reader `json:"inputFile"`
}

// NewCSVConvertor creates a CSV convertor from an net/http Form
func NewCSVConvertor() *CSVConvertor {
	return &CSVConvertor{}
}

// IndexFragments stores the fragments generated from the CSV into ElasticSearch
func (con *CSVConvertor) IndexFragments(p *elastic.BulkProcessor, revision int) (int, error) {

	fg := NewFragmentGraph()
	fg.Meta = &Header{
		OrgID:    config.Config.OrgID,
		Revision: int32(revision),
		DocType:  "csvUpload",
		Spec:     con.DefaultSpec,
		Tags:     []string{"csvUpload"},
	}

	rm, err := con.Convert()

	if err != nil {
		return 0, err
	}

	seen := 0
	for k, fr := range rm.Resources() {
		fg.Meta.EntryURI = k
		fg.Meta.NamedGraphURI = fmt.Sprintf("%s/graph", k)
		frags, err := fr.CreateFragments(fg)
		if err != nil {
			return 0, err
		}

		for _, frag := range frags {
			frag.Meta.AddTags("csvUpload")
			err := frag.AddTo(p)
			if err != nil {
				return 0, err
			}
			seen = seen + 1
		}
	}
	return seen, nil
}

//Convert converts the CSV InputFile to an RDF ResourceMap
func (con *CSVConvertor) Convert() (*ResourceMap, error) {
	rm := &ResourceMap{make(map[string]*FragmentResource)}

	triples, err := con.CreateTriples()
	if err != nil {
		return rm, err
	}

	if len(triples) == 0 {
		return rm, fmt.Errorf("the list of triples cannot be empty")
	}

	for _, t := range triples {
		err := AppendTriple(rm.resources, t)
		if err != nil {
			return rm, err
		}
	}

	return rm, nil

}

// CreateTriples converts a csv file to a list of Triples
func (con *CSVConvertor) CreateTriples() ([]*r.Triple, error) {

	records, err := con.GetReader()
	if err != nil {
		return nil, err
	}

	var header []string
	var headerMap map[int]r.Term
	var subjectColumnIdx int

	triples := []*r.Triple{}

	for idx, row := range records {
		if idx == 0 {
			header = row
			headerMap = con.CreateHeader(header)
			subjectColumnIdx, err = con.GetSubjectColumn(header)
			if err != nil {
				return nil, err
			}
			continue
		}

		s, sType := con.CreateSubjectResource(row[subjectColumnIdx])
		triples = append(triples, sType)

		for idx, column := range row {

			if idx == subjectColumnIdx {
				continue
			}
			p := headerMap[idx]
			triples = append(triples, con.CreateTriple(s, p, column))
			if err != nil {
				return nil, err
			}
		}

	}

	return triples, nil
}

// CreateHeader creates a map based on column id for the predicates
func (con *CSVConvertor) CreateHeader(row []string) map[int]r.Term {
	m := make(map[int]r.Term)
	for idx, column := range row {
		m[idx] = r.NewResource(
			fmt.Sprintf("%s/%s", strings.TrimSuffix(con.PredicateURIBase, "/"), strings.ToLower(column)),
		)
	}
	return m
}

// CreateTriple creates a rdf2go.Triple from the CSV column
func (con *CSVConvertor) CreateTriple(subject r.Term, predicate r.Term, column string) *r.Triple {
	return r.NewTriple(
		subject,
		predicate,
		r.NewLiteral(column),
	)
}

// CreateSubjectResource creates the Subject  URI and type triple for the subject column
func (con *CSVConvertor) CreateSubjectResource(subjectID string) (r.Term, *r.Triple) {
	s := r.NewResource(fmt.Sprintf("%s/%s", strings.TrimSuffix(con.SubjectURIBase, "/"), subjectID))
	t := r.NewTriple(
		s,
		r.NewResource("http://www.w3.org/1999/02/22-rdf-syntax-ns#type"),
		r.NewResource(con.SubjectClass),
	)
	return s, t
}

// GetReader returns a nested array of strings
func (con *CSVConvertor) GetReader() ([][]string, error) {
	r := csv.NewReader(con.InputFile)
	r.Comma = []rune(con.Separator)[0]
	r.Comment = '#'

	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	return records, nil
}

// GetSubjectColumn returns the index of the subject column
func (con *CSVConvertor) GetSubjectColumn(headers []string) (int, error) {
	for idx, column := range headers {
		if column == con.SubjectColumn {
			return idx, nil
		}
	}

	return 0, fmt.Errorf("subjectColumn %s not found in header", con.SubjectColumn)
}

// func Valid bool
// todo add curl example
