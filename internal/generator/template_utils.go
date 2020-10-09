package generator

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"text/template"

	"github.com/shurcooL/graphql/ident"
)

func renderGeneratorTemplate(name, text string, params map[string]string) *template.Template {
	funcMap := makeTemplateFuncMap(params)
	return template.Must(template.New(name).Funcs(funcMap).Parse(text))
}

func makeTemplateFuncMap(params map[string]string) template.FuncMap {
	isExcluded := func(s string) bool {
		return false
	}
	ucFirst := func(s string) string {
		s = strings.ToUpper(s[0:1]) + s[1:]
		return s
	}
	lcFirst := func(s string) string {
		s = strings.ToLower(s[0:1]) + s[1:]
		return s
	}
	identifier := func(s string) string {
		s = strings.TrimLeft(s, "_")
		// return "T_" + s
		return ident.ParseScreamingSnakeCase(s).ToMixedCaps()
		// return ident.ParseLowerCamelCase(name).ToMixedCaps()
	}
	scalarIdentifier := func(s string) string {
		return ident.ParseScreamingSnakeCase(s).ToMixedCaps()
		// s = strings.TrimLeft(s, "_")
		// s = strings.ToUpper(s[0:1]) + s[1:]
		// return s
		// return "T_" + s
		// return ident.ParseScreamingSnakeCase(s).ToMixedCaps()
		// return ident.ParseLowerCamelCase(name).ToMixedCaps()
	}
	enumIdentifierValueSuffix := func(s string) string {
		// return s
		// return ident.ParseScreamingSnakeCase(name).ToMixedCaps()
		return "_" + strings.ToUpper(s)
	}
	enumTypeString := func(s string) string {
		return identifier(s)
	}
	enumAllValuesIdentifier := func(s string) string {
		return enumTypeString(s) + "__LIST"
	}
	scalarTypeString := func(gqltype string) string {
		tnorm := strings.ToLower(gqltype)
		switch tnorm {
		case "order_by":
			return "string"
		case "time":
			return "time.Time"
		case "timestamp":
			return "time.Time"
		case "timestamptz":
			return "time.Time"
		case "date":
			return "time.Time"
		case "datetime":
			return "time.Time"
		case "uuid":
			return "uuid.UUID"
		case "id":
			return "string"
		case "string":
			return "string"
		case "boolean":
			return "bool"
		case "float":
			return "float64"
		case "integer":
			return "time.Time"
		case "int":
			return "int32"
		case "json":
			return "map[string]interface{}"
		case "jsonb":
			return "map[string]interface{}"
		default:
			return scalarIdentifier(gqltype)
		}
	}

	// typeString returns a string representation of GraphQL type t.
	var typeString func(t map[string]interface{}) string
	typeString = func(t map[string]interface{}) string {
		switch t["kind"] {
		case "SCALAR":
			return "*" + scalarTypeString(t["name"].(string))
		case "NON_NULL":
			s := typeString(t["ofType"].(map[string]interface{}))
			if !strings.HasPrefix(s, "*") {
				panic(fmt.Errorf("nullable type %q doesn't begin with '*'", s))
			}
			return s[1:] // Strip star from nullable type to make it non-null.
		case "LIST":
			return "*[]" + typeString(t["ofType"].(map[string]interface{}))
		case "ENUM":
			return "*" + enumTypeString(t["name"].(string))
		case "INPUT_OBJECT":
		case "OBJECT":
		default:
			break
			// 	return "*" + identifier(t["name"].(string))
			// 	// return "*" + scalarType(t["name"].(string), t["kind"].(string))
		}
		return "*" + identifier(t["name"].(string))
	}

	// settings["integer_enums"] = true

	return template.FuncMap{
		"ucFirst": func(s string) string {
			return ucFirst(s)
		},
		"lcFirst": func(s string) string {
			return lcFirst(s)
		},
		"pkgname": func() string {
			pkg, ok := params["pkg"]
			if ok {
				return pkg
			}
			return ""
		},
		"feature": func(field string) bool {
			switch field {
			case "use_integer_enums":
				return true
			default:
				return false
			}
		},
		"first": func(x int, a interface{}) bool {
			return x == 0
		},
		"last": func(x int, a interface{}) bool {
			return x == reflect.ValueOf(a).Len()-1
		},
		"toUpper": func(s string) string {
			return strings.ToUpper(s)
		},
		"toLower": func(s string) string {
			return strings.ToLower(s)
		},
		"isGraphQLMeta": func(s string) bool {
			return false
			// return strings.HasPrefix(s, "__")
		},
		"isExcluded": isExcluded,
		"quote": func(s string) string {
			return strconv.Quote(s)
		},
		"join": func(elems []string, sep string) string {
			return strings.Join(elems, sep)
		},
		"sortByName": func(types []interface{}) []interface{} {
			sort.Slice(types, func(i, j int) bool {
				ni := types[i].(map[string]interface{})["name"].(string)
				nj := types[j].(map[string]interface{})["name"].(string)
				return ni < nj
			})
			return types
		},
		"sortByNameRev": func(types []interface{}) []interface{} {
			sort.Slice(types, func(i, j int) bool {
				ni := types[i].(map[string]interface{})["name"].(string)
				nj := types[j].(map[string]interface{})["name"].(string)
				return ni > nj
			})
			return types
		},
		"filterBy": func(field string, kind string, types []interface{}) []interface{} {
			var filtered = []interface{}{}
			for _, t := range types {
				if val, ok := t.(map[string]interface{})[field]; ok && val.(string) == kind && !isExcluded(val.(string)) {
					filtered = append(filtered, t)
				}
				continue
			}
			return filtered
		},
		"extractField": func(field string, types []interface{}) []string {
			var values = []string{}
			for _, t := range types {
				if val, ok := t.(map[string]interface{})[field]; ok {
					values = append(values, val.(string))
				}
				continue
			}
			return values
		},
		"inputObjects": func(types []interface{}) []string {
			var names []string
			for _, t := range types {
				t := t.(map[string]interface{})
				if t["kind"].(string) != "INPUT_OBJECT" {
					continue
				}
				names = append(names, t["name"].(string))
			}
			sort.Strings(names)
			return names
		},
		"objects": func(types []interface{}) []string {
			var names []string
			for _, t := range types {
				t := t.(map[string]interface{})
				if t["kind"].(string) != "OBJECT" {
					continue
				}
				names = append(names, t["name"].(string))
			}
			sort.Strings(names)
			return names
		},
		"identifier":                identifier,
		"type":                      typeString,
		"enumType":                  enumTypeString,
		"enumIdentifierValueSuffix": enumIdentifierValueSuffix,
		"enumAllValuesIdentifier":   enumAllValuesIdentifier,
		"scalarIdentifier":          scalarIdentifier,
		"scalarType":                scalarTypeString,
		"clean": func(s string) string {
			return strings.Join(strings.Fields(s), " ")
		},
		"formatDescription": func(s string) string {
			s = strings.ToLower(s[0:1]) + s[1:]
			if !strings.HasSuffix(s, ".") {
				s += "."
			}
			return s
		},
	}
}
