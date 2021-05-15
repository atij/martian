package catalog

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/google/martian"
	"github.com/google/martian/parse"
	"github.com/google/martian/v3/log"
)

func init() {
	parse.Register("catalog.CategoryModifier", categoryModifierFromJSON)
}

type CategoryModifier struct {
	ContentType string
}

func (m *CategoryModifier) ModifyResponse(res *http.Response) error {
	log.Debugf("catalog.CategoryModifier.ModifyResponse: request: %s", res.Request.URL)

	var r CategoryRequest

	err := json.NewDecoder(res.Body).Decode(&r)
	if err != nil {
		return err
	}

	res.Body.Close()

	result := r.transform()

	b, err := json.Marshal(&result)
	if err != nil {
		return err
	}

	//res.ContentLength = int64(len(b))
	res.Body = ioutil.NopCloser(bytes.NewReader(b))
	res.Header.Set("Content-Type", "application/json")

	return nil
}

func CategoryNewModifier(contentType string) martian.ResponseModifier {
	log.Debugf("catalog.CategoryNewModifier: contentType %s", contentType)
	return &CategoryModifier{
		ContentType: contentType,
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

	mod := CategoryNewModifier(msg.ContentType)
	return parse.NewResult(mod, msg.Scope)
}

type CategoryRequest struct {
	ID                    int                  `json:"id"`
	Name                  string               `json:"name"`
	Permalink             string               `json:"permalink"`
	ParentID              interface{}          `json:"parent_id"`
	Enabled               bool                 `json:"enabled"`
	Position              int                  `json:"position"`
	Anchor                bool                 `json:"anchor"`
	IncludeInNavigation   bool                 `json:"include_in_navigation"`
	IncludeInShowroom     bool                 `json:"include_in_showroom"`
	DisplayBanner         bool                 `json:"display_banner"`
	ShowroomPosition      interface{}          `json:"showroom_position"`
	NavigationDisplayType string               `json:"navigation_display_type"`
	HideProductRelations  []interface{}        `json:"hide_product_relations"`
	VisibleFor            []string             `json:"visible_for"`
	MetaTitle             string               `json:"meta_title"`
	MetaKeywords          string               `json:"meta_keywords"`
	MetaDescription       string               `json:"meta_description"`
	VisibleForSegments    []VisibleForSegments `json:"visible_for_segments"`
	Type                  string               `json:"type"`
	ProductsCount         int                  `json:"products_count"`
	LastBqViewName        interface{}          `json:"last_bq_view_name"`
}

type VisibleForSegments struct {
	Type   string   `json:"type"`
	Values []string `json:"values"`
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
	Meta       Meta          `json:"meta"`
}

type Meta struct {
	Title       string `json:"title"`
	Keywords    string `json:"keywords"`
	Description string `json:"description"`
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
		s := Segment(c)
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
		Meta: Meta{
			Title:       r.MetaTitle,
			Keywords:    r.MetaKeywords,
			Description: r.MetaDescription,
		},
	}
}
