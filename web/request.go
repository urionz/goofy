package web

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

type Request struct {
	*http.Request
	json *ParameterBag
	form *ParameterBag
}

func NewRequest(r *http.Request) *Request {
	formParameters := Parameters{}
	r.ParseForm()
	for k, v := range r.Form {
		for _, vv := range v {
			formParameters[k] = append(formParameters[k], vv)
		}
	}
	request := &Request{
		Request: r,
		form:    NewParameterBag(formParameters),
	}
	request.setJson()
	return request
}

// Determine if the request is sending JSON.
func (r *Request) IsJson() bool {
	contentType := r.Header.Get("Content-Type")
	accepts := []string{"/json", "+json"}
	for _, accept := range accepts {
		if strings.Contains(contentType, accept) {
			return true
		}
	}
	return false
}

// Retrieve an input item from the request.
func (r *Request) Input() *ParameterBag {
	inputParameters := r.getInputSource().All()
	queryParameters := r.form.All()
	for k, v := range inputParameters {
		queryParameters[k] = v
	}
	return NewParameterBag(queryParameters)
}

// Get the input source for the request.
func (r *Request) getInputSource() *ParameterBag {
	if r.IsJson() {
		return r.Json()
	}
	return r.form
}

func (r *Request) setJson() {
	jsonParameters := Parameters{}
	var data map[string]interface{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		r.json = NewEmptyParameterBag()
		return
	}
	r.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	if err := json.Unmarshal(body, &data); err != nil {
		r.json = NewEmptyParameterBag()
		return
	}
	for k, v := range data {
		jsonParameters[k] = append(jsonParameters[k], v)
	}
	r.json = NewParameterBag(jsonParameters)
}

func (r *Request) Json() *ParameterBag {
	if r.json == nil {
		r.setJson()
	}
	return r.json
}
