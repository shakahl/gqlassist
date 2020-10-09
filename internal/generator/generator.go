package generator

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

	"github.com/shakahl/graphql-typedef-go/internal/apiclient"
	"github.com/shakahl/graphql-typedef-go/internal/utils"
)

func init() {

}

type GraphQLTypeDefGeneratorOptions struct {
	Endpoint        string
	AuthHeader      string
	AuthToken       string
	OutputDirectory string
	OutputPackage   string
}

type GraphQLTypeDefGenerator struct {
	options GraphQLTypeDefGeneratorOptions
	client  *apiclient.ApiClient
	logger  *log.Logger
}

func New(options GraphQLTypeDefGeneratorOptions, logger *log.Logger) *GraphQLTypeDefGenerator {
	if logger == nil {
		logger = log.New(os.Stdout, "graphql_type_def_generator", log.LstdFlags|log.LUTC)
	}
	g := &GraphQLTypeDefGenerator{
		options: options,
		logger:  logger,
	}
	return g
}

// getTemplates returns a template map (filename->template)
func (g *GraphQLTypeDefGenerator) getTemplates() map[string]*template.Template {
	// Filename -> Template.
	var params = map[string]string{
		"pkg": g.options.OutputPackage,
	}
	var templates = map[string]*template.Template{
		"gen_graphql_scalars.go":       renderGeneratorTemplate("gen_graphql_scalars.gotmpl", GeneratorTemplateScalar, params),
		"gen_graphql_enums.go":         renderGeneratorTemplate("gen_graphql_enums.gotmpl", GeneratorTemplateEnum, params),
		"gen_graphql_input_objects.go": renderGeneratorTemplate("gen_graphql_input_objects.gotmpl", GeneratorTemplateInputObjects, params),
		"gen_graphql_objects.go":       renderGeneratorTemplate("gen_graphql_objects.gotmpl", GeneratorTemplateObjects, params),
	}
	return templates
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

func (g *GraphQLTypeDefGenerator) decodeStringToInterface(schema string) (interface{}, error) {
	var target interface{}
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

	for filename, t := range g.getTemplates() {
		outputFile := g.getOutputFilePath(filename)
		g.logger.Printf("Processing template: %s\n", filename)
		var buf bytes.Buffer
		err := t.Execute(&buf, schema)
		if err != nil {
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
