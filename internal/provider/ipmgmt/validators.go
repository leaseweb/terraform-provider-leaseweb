package ipmgmt

import "regexp"

// AutomaticUnnullingAtMessage is the validation message shown against a
// Terraform attribute. It is exported so that the acceptance tests assert
// against the same value the schema reports, rather than a copy that can
// drift from it.
const AutomaticUnnullingAtMessage = "must be specified using the `2019-09-08 00:00:00 +0000 UTC` format"

// automaticUnnullingAtRegex matches the datetime format the API accepts.
var automaticUnnullingAtRegex = regexp.MustCompile(
	`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2} \+\d{4} UTC$`,
)
