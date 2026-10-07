package objectstorage

import "regexp"

var (
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
