package objectstorage

import "regexp"

// Validation messages shown against a Terraform attribute. They are exported
// so that the acceptance tests assert against the same value the schema
// reports, rather than a copy that can drift from it.
const (
	BucketNameMessage = "must start and end with an alphanumeric character and may only contain alphanumeric characters and hyphens"
	GroupIDMessage    = "must be a group ID"
	RFC3339Message    = "must be specified using the RFC3339 format to the second (`yyyy-mm-ddThh:mm:ssZ`)"
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
	//
	// Fractional seconds are deliberately excluded. The API returns them, but
	// a timestamp is stored back as RFC3339 to the second, so a configured
	// value carrying them would never match what is read back.
	rfc3339Regex = regexp.MustCompile(
		`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(Z|[+-]\d{2}:\d{2})$`,
	)
)
