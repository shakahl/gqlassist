package gqlassist

import (
	"bytes"
	"context"
	"encoding/json"
	"go/format"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/shakahl/gqlassist/internal/apiclient"
	"github.com/shakahl/gqlassist/internal/statikdata"
	"github.com/shakahl/gqlassist/internal/utils"
)

func init() {

}

type GraphQLTypeDefGeneratorConfig struct {
	Endpoint        string
	AuthHeader      string
	AuthToken       string
	OutputDirectory string
	OutputPackage   string
}

type GraphQLTypeDefGenerator struct {
	options   GraphQLTypeDefGeneratorConfig
	client    *apiclient.ApiClient
	logger    *log.Logger
	templates map[string]*template.Template
}

func New(options GraphQLTypeDefGeneratorConfig, logger *log.Logger) *GraphQLTypeDefGenerator {
	if logger == nil {
		logger = log.New(os.Stdout, "gqlassist", log.LstdFlags|log.LUTC)
	}
	g := &GraphQLTypeDefGenerator{
		options: options,
		logger:  logger,
	}
	return g
}

// getTemplates returns a template map (filename->template)
func (g *GraphQLTypeDefGenerator) GetTemplates(params map[string]interface{}) map[string]*template.Template {
	if g.templates != nil {
		return g.templates
	}

	statikTemplates := map[string]string{
		"gen_graphql_enums.go":         "/graphql_enums.gotmpl",
		"gen_graphql_scalars.go":       "/graphql_scalars.gotmpl",
		"gen_graphql_input_objects.go": "/graphql_input_objects.gotmpl",
		"gen_graphql_objects.go":       "/graphql_objects.gotmpl",
	}

	g.templates = make(map[string]*template.Template)

	for name, path := range statikTemplates {
		content := statikdata.MustReadFileString(path)
		g.templates[name] = renderGeneratorTemplate(name, content)
	}

	return g.templates
}

func (g *GraphQLTypeDefGenerator) getWorkingDirectory() string {
	pwd, err := os.Getwd()
	if err != nil {
		return ""
	}
	return filepath.FromSlash(pwd)
}

func (g *GraphQLTypeDefGenerator) GetOutputDirectory() string {
	return filepath.FromSlash(g.options.OutputDirectory)
}

func (g *GraphQLTypeDefGenerator) getOutputFilePath(filename string) string {
	return filepath.Join(g.GetOutputDirectory(), filename)
}

func (g *GraphQLTypeDefGenerator) getClient() *apiclient.ApiClient {
	if g.client != nil {
		return g.client
	}

	c := apiclient.NewApiClient(g.options.Endpoint, nil)
	if g.options.AuthToken != "" {
		c.Header.Set(g.options.AuthHeader, g.options.AuthToken)
	}
	g.client = c
	return c
}

func (g *GraphQLTypeDefGenerator) fetchGraphQLSchema() (string, error) {
	client := g.getClient()
	result, err := client.SendGraphQLQuery(context.Background(), GeneratorTemplateIntrospectionQuery, nil)
	if err != nil {
		return "", err
	}
	return string(result.GetBody()), nil
}

func (g *GraphQLTypeDefGenerator) decodeStringToInterface(schema string) (map[string]interface{}, error) {
	var target map[string]interface{}
	s := strings.NewReader(schema)
	if err := json.NewDecoder(s).Decode(&target); err != nil {
		return nil, err
	}
	return target, nil
}

func (g *GraphQLTypeDefGenerator) Generate() error {
	var err error

	schemaFile := g.getOutputFilePath("schema.json")

	if err = utils.EnsurePathDirectoriesExists(schemaFile); err != nil {
		return err
	}

	schemaJson, err := g.fetchGraphQLSchema()
	if err != nil {
		return err
	}

	if err := utils.WriteToFile(schemaFile, schemaJson); err != nil {
		return err
	}

	g.logger.Printf("GraphQL schema written into: %s", schemaFile)

	schema, err := g.decodeStringToInterface(schemaJson)
	if err != nil {
		return err
	}

	params := make(map[string]interface{})
	params["FileHeader"] = FileHeaderText
	params["Schema"] = schema
	params["PackageName"] = g.options.OutputPackage
	params["FeatureFlags"] = map[string]interface{}{
		"UseIntegerEnums": true,
	}
	// use_integer_enums

	templates := g.GetTemplates(params)

	for filename, t := range templates {
		outputFile := g.getOutputFilePath(filename)
		g.logger.Printf("Processing template: %s\n", filename)
		var buf bytes.Buffer
		err := t.Execute(&buf, params)
		if err != nil {
			g.logger.Printf("ERR: %+v\n", err)
			return err
		}
		out, err := format.Source(buf.Bytes())
		if err != nil {
			g.logger.Println(err)
			out = []byte("// gofmt error: " + err.Error() + "\n\n" + buf.String())
		}
		g.logger.Printf("writing %s\n", outputFile)
		err = ioutil.WriteFile(outputFile, out, 0644)
		if err != nil {
			return err
		}
	}

	return nil
}
