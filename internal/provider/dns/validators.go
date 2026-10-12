package dns

import "regexp"

// RecordNameMessage is the validation message shown against a Terraform
// attribute. It is exported so that the acceptance tests assert against the
// same value the schema reports, rather than a copy that can drift from it.
const RecordNameMessage = "must end in ."

// recordNameRegex matches the fully qualified names the API accepts.
var recordNameRegex = regexp.MustCompile(`^.*\.$`)
