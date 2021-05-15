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
	parse.Register("catalog.CategoryCollectionModifier", categoryModifierFromJSON)
}

type CategoryCollectionModifier struct {
	ContentType string
}

func (m *CategoryCollectionModifier) ModifyResponse(res *http.Response) error {
	log.Debugf("catalog.CategoryModifier.ModifyResponse: request: %s", res.Request.URL)

	var r CategoryCollectionRequest

	err := json.NewDecoder(res.Body).Decode(&r)
	if err != nil {
		return err
	}

	result := r.transform()

	var buffer bytes.Buffer
	json.NewEncoder(&buffer).Encode(&result)

	res.Body = ioutil.NopCloser(&buffer)
	res.Header.Set("Content-Type", "application/json")

	return nil
}

func CategoryCollectionNewModifier(contentType string) martian.ResponseModifier {
	return &CategoryModifier{
		ContentType: contentType,
	}
}

type CategoryCollectionRequest struct {
	Items []CategoryRequest
}

type CategoryCollectionResponse struct {
	Items []*CategoryResponse
}

func (r *CategoryCollectionRequest) transform() *CategoryCollectionResponse {
	log.Debugf("transform")
	var result CategoryCollectionResponse
	for _, item := range r.Items {
		result.Items = append(result.Items, item.transform())
	}
	return &result
}
