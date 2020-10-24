package config

import (
	"errors"

	"github.com/shakahl/gqlassist/internal/utils"
)

const (
	DefaultOutputDirectory   = "internal/graphqltypes"
	DefaultOutputPackage     = "graphqltypes"
	DefaultGraphQLEndpoint   = ""
	DefaultGraphQLAuthHeader = "Authorization"
	DefaultGraphQLAuthToken  = ""
)

type (
	ConfigSchema struct {
		Config            string `json:"config" ignored:"true"`
		Debug             bool   `json:"debug" envconfig:"debug" default:"false"`
		OutputDirectory   string `json:"output_directory" envconfig:"output_directory" default:"internal/graphqltypes" split_words:"true"`
		OutputPackage     string `json:"output_package" envconfig:"output_package" required:"true" split_words:"true" default:"graphqltypes"`
		GraphQLEndpoint   string `json:"graphql_endpoint" envconfig:"graphql_endpoint" required:"true" split_words:"true"`
		GraphQLAuthHeader string `json:"graphql_auth_header" envconfig:"graphql_auth_header" required:"true" split_words:"true" default:"Authorization"`
		GraphQLAuthToken  string `json:"graphql_auth_token" envconfig:"graphql_auth_token" split_words:"true"`
	}
)

func (c *ConfigSchema) Validate() error {
	if c.OutputDirectory == "" {
		return errors.New("invalid output directory")
	}
	if !utils.IsPathExists(c.OutputDirectory) {
		return errors.New("output directory does not exists")
	}
	if c.OutputPackage == "" {
		return errors.New("invalid output package")
	}
	if c.GraphQLEndpoint == "" || !utils.IsValidURL(c.GraphQLEndpoint) {
		return errors.New("invalid graphql endpoint")
	}
	if c.GraphQLAuthHeader == "" {
		return errors.New("invalid graphql auth header")
	}
	if c.GraphQLAuthToken == "" {
		return errors.New("invalid graphql auth token")
	}
	return nil
}

func NewSchema() *ConfigSchema {
	return &ConfigSchema{
		OutputDirectory:   DefaultOutputDirectory,
		OutputPackage:     DefaultOutputPackage,
		GraphQLEndpoint:   DefaultGraphQLEndpoint,
		GraphQLAuthHeader: DefaultGraphQLAuthHeader,
		GraphQLAuthToken:  DefaultGraphQLAuthToken,
	}
}
