package apiclient

// GQLApiErrors represents the "errors" array in a response from a GraphQL server.
// If returned via error interface, the slice is expected to contain at least 1 element.
//
// Specification: https://facebook.github.io/graphql/#sec-Errors.
type GQLApiErrors []struct {
	Message   string
	Locations []struct {
		Line   int
		Column int
	}
}

// Error implements error interface for GQLApiErrors.
func (e GQLApiErrors) Error() string {
	if len(e) > 0 {
		return e[0].Message
	}
	return ""
}
