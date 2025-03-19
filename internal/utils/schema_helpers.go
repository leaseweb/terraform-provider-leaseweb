package utils

import (
	"fmt"
	"log"
	"slices"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

type Action int

const (
	CreateAction Action = iota
	ReadAction
	UpdateAction
	DeleteAction
)

func (a Action) string(unsupportedActions []Action) string {
	var secondAction string
	var format = "Once %s, this resource cannot be %s"

	switch a {
	case CreateAction:
		secondAction = "created"
	case ReadAction:
		secondAction = "read"
	case UpdateAction:
		secondAction = "updated"
	case DeleteAction:
		secondAction = "deleted"
	default:
		log.Fatal(fmt.Printf("do not know how to handle action: %q", a))
	}

	firstAction := a.firstAction(unsupportedActions)

	if slices.Contains(unsupportedActions, CreateAction) && a == CreateAction {
		return "This resource cannot be created, only imported"
	}

	if firstAction == secondAction {
		log.Fatal(fmt.Printf(
			"firstAction %q and secondAction %q cannot be equal",
			firstAction,
			secondAction,
		))
	}

	return fmt.Sprintf(
		format,
		firstAction,
		secondAction,
	)
}

func (a Action) firstAction(unsupportedActions []Action) string {
	if slices.Contains(unsupportedActions, CreateAction) {
		return "imported"
	}

	return "created"
}

// AddUnsupportedActionsNotation lets the end user know which actions aren't supported in the markdown description.
func AddUnsupportedActionsNotation(
	response *resource.SchemaResponse,
	unsupportedActions []Action,
) {
	response.Schema.MarkdownDescription += "**Note:**"
	for _, action := range unsupportedActions {
		response.Schema.MarkdownDescription += fmt.Sprintf(
			"\n- %s.",
			action.string(unsupportedActions),
		)
	}
}
