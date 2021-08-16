package cj

type CollectionJsonType struct {
	Collection CollectionType `json:"collection"` // REQUIRED
}

type CollectionType struct {
	Version  string       `json:"version"`
	Href     URIType      `json:"href"`
	Links    []LinkType   `json:"links"`
	Items    []ItemType   `json:"items"`
	Queries  []QueryType  `json:"queries"`
	Template TemplateType `json:"template"`
	Error    ErrorType    `json:"error"`
}

type ItemType struct {
	Href  URIType    `json:"href"`
	Data  []DataType `json:"data"`
	Links []LinkType `json:"links"` // OPTIONAL
}

type URIType string

type LinkType struct {
	Href   URIType `json:"href"`   // REQUIRED
	Rel    string  `json:"rel"`    // REQUIRED
	Name   string  `json:"name"`   // OPTIONAL
	Render string  `json:"render"` // OPTIONAL MUST be "image" or "link"
	Prompt string  `json:"prompt"` // OPTIONAL
}

type QueryType struct {
	Href   URIType    `json:"href"`   // REQUIRED
	Rel    string     `json:"rel"`    // REQUIRED
	Name   string     `json:"name"`   // OPTIONAL
	Prompt string     `json:"prompt"` // OPTIONAL
	Data   []DataType `json:"data"`   // OPTIONAL
}

type TemplateType interface{}

type DataType struct {
	Name   string    `json:"name"`   // REQUIRED
	Value  ValueType `json:"value"`  // OPTIONAL
	Prompt string    `json:"prompt"` // OPTIONAL
}

type ValueType interface{}

type ErrorType struct {
	Title   string `json:"title"`
	Code    string `json:"code"`
	Message string `json:"message"`
}
