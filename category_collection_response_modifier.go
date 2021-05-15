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

	var list []json.RawMessage
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, &list)
	//err := json.NewDecoder(res.Body).Decode(&r)
	if err != nil {
		return err
	}

	var result []CategoryResponse
	for _, item := range list {
		//var i CategoryRequest
		var s = new(CategoryRequest)
		err = json.Unmarshal(item, &s)
		if err != nil {
			return err
		}

		r := s.transform()

		result = append(result, *r)
	}

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
