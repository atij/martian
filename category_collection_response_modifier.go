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
	parse.Register("catalog.CategoryCollectionModifier", categoryCollectionModifierFromJSON)
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

	res.Body.Close()

	err = json.Unmarshal(b, &list)
	//err := json.NewDecoder(res.Body).Decode(&r)
	if err != nil {
		return err
	}

	var result []*CategoryResponse
	for _, item := range list {
		//var i CategoryRequest
		var s = new(CategoryRequest)
		err = json.Unmarshal(item, &s)
		if err != nil {
			return err
		}

		result = append(result, s.transform())
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

type CategoryCollectionModifierJSON struct {
	ContentType string               `json:"contentType"`
	Scope       []parse.ModifierType `json:"scope"`
}

func categoryCollectionModifierFromJSON(b []byte) (*parse.Result, error) {
	msg := &CategoryCollectionModifierJSON{}

	if err := json.Unmarshal(b, msg); err != nil {
		return nil, err
	}

	mod := CategoryCollectionNewModifier(msg.ContentType)
	return parse.NewResult(mod, msg.Scope)
}
