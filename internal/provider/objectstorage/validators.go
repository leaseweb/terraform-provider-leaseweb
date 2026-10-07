package objectstorage

import "regexp"

// Validation messages shown against a Terraform attribute. They are declared
// here so that each wording lives in one place, rather than being repeated
// wherever it is reported.
const (
	bucketNameMessage = "must start and end with an alphanumeric character and may only contain alphanumeric characters and hyphens"
	groupIDMessage    = "must be a group ID"
	rfc3339Message    = "must be specified using the RFC3339 format (`yyyy-mm-ddThh:mm:ssZ`)"
)

var (
	// bucketNameRegex matches the bucket names the API accepts.
	bucketNameRegex = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9]$`)

	// uuidRegex matches the group identifiers the API accepts, which it
	// documents as `format: uuid`.
	uuidRegex = regexp.MustCompile(
		`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`,
	)

	// rfc3339Regex matches the datetime format the API accepts, so a malformed
	// value is rejected during the plan rather than part way through an apply.
	rfc3339Regex = regexp.MustCompile(
		`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(\.\d+)?(Z|[+-]\d{2}:\d{2})$`,
	)
)
