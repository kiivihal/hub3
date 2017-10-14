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

package hub3

import (
	"bitbucket.org/delving/rapid/config"
	"github.com/renevanderark/goharvest/oai"
)

// ProcessVerb processes different OAI-PMH verbs
func ProcessVerb(r *oai.Request) interface{} {
	switch r.Verb {
	case "Identify":
		return renderIdentify(r)
	case "ListMetadataFormats":
		formats := []oai.MetadataFormat{
			oai.MetadataFormat{
				MetadataPrefix:    "edm",
				Schema:            "",
				MetadataNamespace: "http://www.europeana.eu/schemas/edm/",
			},
		}

		return oai.ListMetadataFormats{
			MetadataFormat: formats,
		}
	case "ListSets":
		return "sets"
	case "ListIdentifiers":
		return "identifiers"
	case "ListRecords":
		return "records"
	case "GetRecord":
		return "record"
	default:
		return "badVerb"
	}
}

// renderIdentify returns the identify response of the repository
func renderIdentify(r *oai.Request) interface{} {
	return oai.Identify{
		RepositoryName:    config.Config.OAIPMH.RepositoryName,
		BaseURL:           r.BaseURL,
		ProtocolVersion:   "2.0",
		AdminEmail:        config.Config.OAIPMH.AdminEmails,
		DeletedRecord:     "persistent",
		EarliestDatestamp: "1970-01-01T00:00:00Z",
		Granularity:       "YYYY-MM-DDThh:mm:ssZ",
	}
}
