package catalog

import (
	"encoding/json"
	"net/http"

	"github.com/google/martian"
	"github.com/google/martian/parse"
	"github.com/google/martian/v3/log"
)

func init() {
	log.Debugf("catalog.CategoryModifier register")
	parse.Register("catalog.CategoryModifier", categoryModifierFromJSON)
}

type CategoryModifier struct {
	contentType string
}

func (m *CategoryModifier) ModifyResponse(res *http.Response) error {
	log.Debugf("catalog.CategoryModifier.ModifyResponse: request: %s", res.Request.URL)

	/*
		var r CategoryRequest

		err := json.NewDecoder(res.Body).Decode(&r)
		if err != nil {
			return err
		}

		result := r.transform()

		b, err := json.Marshal(&result)
		if err != nil {
			return err
		}

		res.Body = ioutil.NopCloser(bytes.NewReader(b))
	*/
	return nil
}

func CategoryNewModifier(contentType string) martian.ResponseModifier {
	log.Debugf("catalog.CategoryNewModifier: contentType %s", contentType)
	return &CategoryModifier{
		contentType: contentType,
	}
}

type CategoryModifierJSON struct {
	ContentType string               `json:"contentType"`
	Scope       []parse.ModifierType `json:"scope"`
}

func categoryModifierFromJSON(b []byte) (*parse.Result, error) {
	msg := &CategoryModifierJSON{}

	if err := json.Unmarshal(b, msg); err != nil {
		return nil, err
	}

	return parse.NewResult(CategoryNewModifier(msg.ContentType), msg.Scope)
}

type CategoryRequest struct {
	ID                    int           `json:"id"`
	ParentID              interface{}   `json:"parent_id"`
	Name                  string        `json:"name"`
	Permalink             string        `json:"permalink"`
	Position              int           `json:"position"`
	ShowroomPosition      interface{}   `json:"showroom_position"`
	IncludeInNavigation   bool          `json:"include_in_navigation"`
	IncludeInShowroom     bool          `json:"include_in_showroom"`
	NavigationDisplayType string        `json:"navigation_display_type"`
	DisplayBanner         bool          `json:"display_banner"`
	HideProductRelations  []interface{} `json:"hide_product_relations"`
	Meta                  struct {
		Title       string `json:"title"`
		Keywords    string `json:"keywords"`
		Description string `json:"description"`
	} `json:"meta"`
	Type               string `json:"type"`
	VisibleForSegments []struct {
		Type   string   `json:"type"`
		Values []string `json:"values"`
	} `json:"visible_for_segments"`
}

type CategoryResponse struct {
	ID         int           `json:"id"`
	Type       string        `json:"type"`
	ParentID   interface{}   `json:"parent_id"`
	Name       string        `json:"name"`
	Position   int           `json:"position"`
	Permalink  string        `json:"permalink"`
	Hide       []interface{} `json:"hide"`
	Conditions []Conditions  `json:"conditions"`
	Meta       struct {
		Title       string `json:"title"`
		Keywords    string `json:"keywords"`
		Description string `json:"description"`
	} `json:"meta"`
}

type Conditions struct {
	Type  string    `json:"type"`
	Value []Segment `json:"value"`
}

type Segment struct {
	Type   string   `json:"type"`
	Values []string `json:"values"`
}

func (r *CategoryRequest) transform() *CategoryResponse {

	var conditions []Conditions
	var segments []Segment
	for _, c := range r.VisibleForSegments {
		s := Segment{
			Type:   c.Type,
			Values: c.Values,
		}
		segments = append(segments, s)
	}

	conditions = append(conditions, Conditions{
		Type:  "segment",
		Value: segments,
	})

	return &CategoryResponse{
		ID:         r.ID,
		Type:       r.Type,
		ParentID:   r.ParentID,
		Name:       r.Name,
		Position:   r.Position,
		Permalink:  r.Permalink,
		Hide:       r.HideProductRelations,
		Conditions: conditions,
		Meta:       r.Meta,
	}
}
