package catalog

import (
	"fmt"
	"testing"
)

func TestTransform(t *testing.T) {

	// create test response
	var req = &CategoryRequest{
		ID:                    100,
		ParentID:              1,
		Name:                  "Test Category",
		Permalink:             "test-category",
		Position:              1,
		ShowroomPosition:      23,
		IncludeInNavigation:   false,
		IncludeInShowroom:     false,
		NavigationDisplayType: "",
		DisplayBanner:         false,
		HideProductRelations:  nil,
		Meta: struct {
		Title       string `json:"title"`
		Keywords    string `json:"keywords"`
		Description string `json:"description"`
	}{},
		Type:               "opa",
		VisibleForSegments: nil,
	}

	result := req.transform()

	fmt.Printf("transformat: %+v\n", result)
}