package gqlassist

const GeneratorTemplateScalar = `
{{- .FileHeader }}

package {{ .PackageName }}

import (
	"time"

	"github.com/google/uuid"
)

{{ if false }}
type (
	// Boolean represents true or false values.
	Boolean bool

	// Float represents signed double-precision fractional values as
	// specified by IEEE 754.
	Float float64

	// ID represents a unique identifier that is Base64 obfuscated. It
	// is often used to refetch an object or as key for a cache. The ID
	// type appears in a JSON response as a String; however, it is not
	// intended to be human-readable. When expected as an input type,
	// any string (such as "VXNlci0xMA==") or integer (such as 4) input
	// value will be accepted as an ID.
	// ID interface{}
	ID string

	// Int represents non-fractional signed whole numeric values.
	// Int can represent values between -(2^31) and 2^31 - 1.
	Int int32

	// String represents textual data as UTF-8 character sequences.
	// This type is most often used by GraphQL to represent free-form
	// human-readable text.
	String string
)

// NewBoolean is a helper to make a new *Boolean.
func NewBoolean(v Boolean) *Boolean { return &v }

// NewFloat is a helper to make a new *Float.
func NewFloat(v Float) *Float { return &v }

// NewID is a helper to make a new *ID.
func NewID(v ID) *ID { return &v }

// NewInt is a helper to make a new *Int.
func NewInt(v Int) *Int { return &v }

// NewString is a helper to make a new *String.
func NewString(v String) *String { return &v }
{{ end -}}

{{ range .Schema.data.__schema.types | filterBy "kind" "SCALAR" | sortByName -}}
{{- if and (eq .kind "SCALAR") (not (isExcluded .name)) -}}
{{- template "scalar" . -}}
{{- end -}}{{- end }}

{{- define "scalar" -}}
// {{ .name | scalarIdentifier }} @graphql="{{ .name }}" {{if .description}}{{.description | clean | formatDescription}}{{end}}
type {{ .name | scalarIdentifier }} {{ .name | scalarType }}

// New{{ .name | scalarIdentifier }} is a helper to make a new *{{ .name | scalarIdentifier }}.
func New{{ .name | scalarIdentifier }}(v {{ .name | scalarIdentifier }}) *{{ .name | scalarIdentifier }} { return &v }

{{ end -}}

`

const GeneratorTemplateEnum = `
{{- .FileHeader }}

package {{ .PackageName }}

import (
	"bytes"
	"encoding/json"
	"strings"
	// "fmt"
)

{{range .Schema.data.__schema.types | filterBy "kind" "ENUM" | sortByName}}{{if and (eq .kind "ENUM") (not (isExcluded .name))}}
{{template "enum" .}}
{{end}}{{end}}


{{- define "enum" -}}
// {{ .name | enumType }} @graphql="{{ .name }}" {{if .description}}{{.description | clean | formatDescription}}{{end}}
{{ if feature "use_integer_enums" -}}
type {{ .name | enumType }} int
{{- else -}}
type {{ .name | enumType }} string
{{- end }}

// {{ .name | enumType }} @graphql="{{ .name }}"  enum type values. {{if .description}}{{.description | clean | formatDescription}}{{else}}{{end}}
const (
	{{ if feature "use_integer_enums" -}}
	{{- range $i, $e := .enumValues | sortByNameRev }}
		{{- $.name | enumType}}{{$e.name | enumIdentifierValueSuffix}} {{if first $i $}}{{$.name | enumType}} = iota + 1 {{end}} // {{if .description}}{{.description | clean | formatDescription}}{{end}}
	{{ end -}}
	{{- else -}}
	{{- range .enumValues sortByNameRev}}
	{{$.name | enumType}}{{.name | enumIdentifierValueSuffix}} {{$.name | identifier}} = {{.name | quote}} // @graphql="{{$.name}}#{{.name}}" {{if .description}}{{.description | clean | formatDescription}}{{end}}
	{{ end -}}
	{{- end -}}
)
var (
	{{ if feature "use_integer_enums" -}}
	{{- .name | enumType}}__toString = map[{{- .name | enumType}}]string{
		{{ range .enumValues | sortByNameRev -}}
		{{ $.name | enumType }}{{ .name | enumIdentifierValueSuffix }}: {{.name | quote -}},
		{{ end }}
	}

	{{ .name | enumType}}__toID = map[string]{{- .name | enumType}}{
		{{ range .enumValues | sortByNameRev -}}
		{{ .name | quote }}: {{- $.name | enumType }}{{.name | enumIdentifierValueSuffix }},
		{{ end }}
	}
	/*
	{{.name | enumAllValuesIdentifier}} = []string{
		{{range .enumValues | sortByNameRev}}{{.name | quote}},{{ end }}
	}
	*/
	{{- else -}}
	{{ .name | enumAllValuesIdentifier }} = map[string]string{
		{{ range .enumValues | sortByNameRev -}}
			{{- .name | quote -}}: {{ .name | quote -}},
		{{ end }}
	}

	{{- end }}
)

func (e {{.name | enumType}}) String() string {
	{{ if feature "use_integer_enums" -}}
	return {{ .name | enumType }}__toString[e]
	/*
	return [...]string{
		{{- range .enumValues | sortByNameRev -}}
			{{- .name | quote -}},
		{{- end -}}
	}[e]
	*/
	{{- else -}}
	s := string(e)
	if v, ok := {{.name | enumAllValuesIdentifier}}[s]; ok {
		return v
	}
	panic(fmt.Errorf("Invalid ENUM value detected: %+v", e))
	{{- end }}
}

// MarshalJSON marshals the enum as a quoted json string
func (e {{.name | enumType}}) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("\"")
	buffer.WriteString({{- .name | enumType -}}__toString[e])
	buffer.WriteString("\"")
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (e *{{.name | enumType}}) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	j = strings.ToUpper(j)
	// Note that if the string cannot be found then it will be set to the zero value.
	*e = {{ .name | enumType }}__toID[j]
	return nil
}

func New{{.name | enumType}}(s string) *{{.name | enumType}} {
	if v, ok := {{ .name | enumType }}__toID[s]; ok {
		return &v
	}
	return nil
}

{{- end -}}

`

const GeneratorTemplateInputObjects = `
{{- .FileHeader }}

package {{ .PackageName }}

import (
	"time"

	"github.com/google/uuid"
)

// Input represents one of the Input structs:
//
// {{join (extractField "name" (filterBy "kind" "INPUT_OBJECT" .Schema.data.__schema.types)) ", "}}.

type Input interface{}

var (
	TypeMapInputObjects = map[string]string{
		{{ range .Schema.data.__schema.types | filterBy "kind" "INPUT_OBJECT" | sortByName -}}{{if eq .kind "INPUT_OBJECT" -}}
		{{ .name | identifier | quote }}: {{ .name | quote }},
		{{ end }}{{- end }}
	}
	ReverseTypeMapInputObjects = map[string]string{
		{{ range .Schema.data.__schema.types | filterBy "kind" "INPUT_OBJECT" | sortByName -}}{{ if eq .kind "INPUT_OBJECT" -}}
		{{ .name | quote }}: {{ .name | identifier| quote }},
		{{ end }}{{- end }}
	}
)

{{range .Schema.data.__schema.types | filterBy "kind" "INPUT_OBJECT" | sortByName}}{{if eq .kind "INPUT_OBJECT"}}
{{template "inputObject" .}}
{{end}}{{end}}


{{- define "inputObject" -}}
// {{ .name | identifier }} @graphql="{{ .name }}" {{ if .description }}{{ .description | clean | formatDescription }}{{ end }}
type {{.name | identifier }} struct {{- "{" }}
	{{- range .inputFields -}}
	{{- if eq .type.kind "NON_NULL" }}
	// {{ .name | identifier }} @graphql="{{ .name }}" {{ if .description }}{{ .description | clean | formatDescription }}{{ end }} (Required.)
	{{ .name | identifier }} {{ .type | type }} ` + "`" + `json:"{{ .name }}" graphql:"!{{ .name }}"` + "`" + `
	{{- end -}}
	{{- if ne .type.kind "NON_NULL" }}
	// {{ .name | identifier }} @graphql="{{ .name }}" {{ if .description }}{{ .description | clean | formatDescription }}{{ end }} (Optional.)
	{{ .name | identifier }} {{ .type | type }} ` + "`" + `json:"{{ .name }},omitempty" graphql:"{{ .name }}"` + "`" + `
	{{- end -}}
	{{- end -}}
}
{{- end -}}
`

const GeneratorTemplateObjects = `
{{- .FileHeader }}

package {{ .PackageName }}

import (
	"time"

	"github.com/google/uuid"
)

// Object represents one of the Input structs:
//
// {{join (extractField "name" (filterBy "kind" "OBJECT" .Schema.data.__schema.types)) ", "}}.

type Object interface{}

type GQLMeta struct {
	// GQLTypeName @graphql="__typename" Is a meta field containing the GraphQL type name
	GQLTypeName *string ` + "`" + `json:"__typename,omitempty"` + "`" + `
}

var (
	TypeMapObjects = map[string]string{
		{{ range .Schema.data.__schema.types | filterBy "kind" "OBJECT" | sortByName -}}{{if eq .kind "OBJECT" -}}
		{{ .name | identifier | quote }}: {{ .name | quote }},
		{{ end }}{{- end }}
	}
	ReverseTypeMapObjects = map[string]string{
		{{ range .Schema.data.__schema.types | filterBy "kind" "OBJECT" | sortByName -}}{{ if eq .kind "OBJECT" -}}
		{{ .name | quote }}: {{ .name | identifier| quote }},
		{{ end }}{{- end }}
	}
)

{{range .Schema.data.__schema.types | filterBy "kind" "OBJECT" | sortByName}}{{if eq .kind "OBJECT"}}
{{template "object" .}}
{{end}}{{end}}


{{- define "object" -}}
// {{ .name | identifier }} @graphql="{{ .name }}" {{ if .description }}{{ .description | clean | formatDescription }}{{ end }}
type {{.name | identifier }} struct {{- "{" }}
	GQLMeta
	{{- range .fields -}}
	{{- if eq .type.kind "NON_NULL" }}
	// {{ .name | identifier }} @graphql="{{ .name }}" {{ if .description }}{{ .description | clean | formatDescription }}{{ end }} (Required.)
	{{ .name | identifier }} {{ .type | type }} ` + "`" + `json:"{{ .name }}" graphql:"!{{ .name }}"` + "`" + `
	{{- end -}}
	{{- if ne .type.kind "NON_NULL" }}
	// {{ .name | identifier }} @graphql="{{ .name }}" {{ if .description }}{{ .description | clean | formatDescription }}{{ end }} (Optional.)
	{{ .name | identifier }} {{ .type | type }} ` + "`" + `json:"{{ .name }},omitempty" graphql:"{{ .name }}"` + "`" + `
	{{- end -}}
	{{- end -}}
}
{{- end -}}
`

const GeneratorTemplateIntrospectionQuery = `
	query IntrospectionQuery {
      __schema {
        queryType { name }
        mutationType { name }
        subscriptionType { name }
        types {
          ...FullType
        }
        directives {
          name
          description
          locations
          args {
            ...InputValue
          }
        }
      }
    }

    fragment FullType on __Type {
      kind
      name
      description
      fields(includeDeprecated: true) {
        name
        description
        args {
          ...InputValue
        }
        type {
          ...TypeRef
        }
        isDeprecated
        deprecationReason
      }
      inputFields {
        ...InputValue
      }
      interfaces {
        ...TypeRef
      }
      enumValues(includeDeprecated: true) {
        name
        description
        isDeprecated
        deprecationReason
      }
      possibleTypes {
        ...TypeRef
      }
    }

    fragment InputValue on __InputValue {
      name
      description
      type { ...TypeRef }
      defaultValue
    }

    fragment TypeRef on __Type {
      kind
      name
      ofType {
        kind
        name
        ofType {
          kind
          name
          ofType {
            kind
            name
            ofType {
              kind
              name
              ofType {
                kind
                name
                ofType {
                  kind
                  name
                  ofType {
                    kind
                    name
                  }
                }
              }
            }
          }
        }
      }
    }
`
