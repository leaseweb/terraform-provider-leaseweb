package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const defaultErrMsg = "An error has occurred in the program. Please consider opening an issue."
const errTitle = "Unexpected Error"

// GeneralError should be called when general errors need to be handled.
func GeneralError(diags *diag.Diagnostics, ctx context.Context, err error) {
	if err != nil {
		logDebug(err.Error(), ctx)
	}
	ReportError(defaultErrMsg, diags)
}

// ImportOnlyError should be used in resource Read() functions for resources that can only be imported.
func ImportOnlyError(diags *diag.Diagnostics) {
	diags.AddError("Resource can only be imported, not created.", "")
}

// UnexpectedImportIdentifierError should be used in Import() functions where the identifier is incorrect.
func UnexpectedImportIdentifierError(diags *diag.Diagnostics, format string, got string) {
	diags.AddError(
		"Unexpected Import Identifier",
		fmt.Sprintf(
			"Expected import identifier with format: %q. Got: %q",
			format,
			got,
		),
	)
}

// SdkError should be used to handle errors returned by the SDK.
func SdkError(
	ctx context.Context,
	diags *diag.Diagnostics,
	err error,
	resp *http.Response,
) {
	// At a minimum diagnostics & error need to be set.
	if diags == nil {
		logDebug("unable to record error details.", ctx)
		log.Fatal(errTitle, defaultErrMsg)
	}

	if err == nil {
		logDebug("No error detail found.", ctx)
		ReportError(defaultErrMsg, diags)
		return
	}

	// Without a response we only need to handle the error.
	if resp == nil {
		ReportError(err.Error(), diags)
		return
	}

	// Handle the error response.
	// Close response body with direct defer reference for clarity
	defer func() {
		if err := resp.Body.Close(); err != nil {
			logDebug(
				fmt.Sprintf("error closing response body: %v", err),
				ctx,
			)
		}
	}()

	// For certain http responses we don't need to analyze the response body.
	if resp.StatusCode == 504 {
		logDebug(fmt.Sprintf("server response: %v", resp.Body), ctx)
		ReportError("The server took too long to respond.", diags)
		return
	}
	if resp.StatusCode == 404 {
		logDebug(fmt.Sprintf("server response: %v", resp.Body), ctx)
		ReportError("Resource not found.", diags)
		return
	}

	// Always log the response body for debugging purposes.
	logDebug(fmt.Sprintf("response body: %v", resp.Body), ctx)

	// Parse the response body. If it can't be parsed throw a general error.
	var errorResponse struct {
		ErrorDetails map[string]any `json:"errorDetails,omitempty"`
		ErrorMessage string         `json:"errorMessage,omitempty"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
		logDebug(
			fmt.Sprintf("error decoding HTTP response body: %v", err),
			ctx,
		)
		ReportError(defaultErrMsg, diags)
		return
	}

	// Show returned errors to the end user.
	if len(errorResponse.ErrorDetails) > 0 {
		handleValidationError(errorResponse.ErrorDetails, diags)
		if diags.HasError() {
			return
		}
	}

	if len(errorResponse.ErrorMessage) > 0 {
		ReportError(errorResponse.ErrorMessage, diags)
		return
	}

	// Show general error if any part of the error response cannot be parsed.
	jsonOutput, err := json.MarshalIndent(errorResponse, "", "  ")
	if err != nil {
		SdkError(
			ctx,
			diags,
			fmt.Errorf("failed to format error response as JSON: %v", err),
			nil,
		)
		return
	}
	GeneralError(diags, ctx, fmt.Errorf(string(jsonOutput), err))
}

func handleValidationError(errorDetails map[string]any, diags *diag.Diagnostics) {
	for errorKey, errorCollections := range errorDetails {
		// Generated a normalized error key that we can work with
		var normalizedErrorKey string
		// If the key contains dots, replace the dots with "_"
		if strings.Contains(errorKey, ".") {
			normalizedErrorKey = strings.ToLower(
				strings.ReplaceAll(errorKey, ".", "_"),
			)
		}
		// If no dots are found, assume camel case and separate the words with "_".
		if normalizedErrorKey == "" {
			m := regexp.MustCompile("[A-Z]")
			res := m.ReplaceAllStringFunc(
				errorKey,
				func(s string) string {
					return "_" + s
				},
			)
			normalizedErrorKey = strings.ToLower(res)
		}

		// Get attribute path from normalized key.
		mapKeys := strings.Split(normalizedErrorKey, "_")
		attributePath := path.Root(mapKeys[0])
		// Every element in the map goes one level deeper.
		for _, mapKey := range mapKeys[1:] {
			attributePath = attributePath.AtMapKey(mapKey)
		}

		// Handle string array errorCollections
		stringErrorCollection, ok := errorCollections.([]interface{})
		if ok {
			handleStringErrorCollection(diags, attributePath, stringErrorCollection)
			continue
		}

		// Handle errorCollections that are a map of string arrays
		errorMapCollection, ok := errorCollections.(map[string]interface{})
		if ok {
			for _, errorMap := range errorMapCollection {
				stringErrorCollection, ok := errorMap.([]interface{})
				if ok {
					handleStringErrorCollection(diags, attributePath, stringErrorCollection)
					continue
				}
			}
			continue
		}
	}
}

func logDebug(details string, ctx context.Context) {
	tflog.Debug(ctx, fmt.Sprintf("Details: %v", details))
}

func ReportError(details string, diags *diag.Diagnostics) {
	diags.AddError(errTitle, details)
}

func handleStringErrorCollection(
	diags *diag.Diagnostics,
	attributePath path.Path,
	stringErrorCollection []interface{},
) bool {
	containsUnhandledErrors := false

	for _, stringError := range stringErrorCollection {
		parsedStringError, ok := stringError.(string)
		if ok {
			diags.AddAttributeError(
				attributePath,
				errTitle,
				parsedStringError,
			)
			continue
		}

		containsUnhandledErrors = true
	}

	return containsUnhandledErrors
}
