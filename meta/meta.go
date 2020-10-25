/*
Copyright © 2020 Soma Szelpal <szelpalsoma@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package meta

const (
	// Version is the application version in !Semver! format
	Version = "0.0.0"

	// BuildID is the build number
	BuildID = "-"

	// BuildDate holds the timestamp of the build
	BuildDate = "-"

	// BuildPlatform hold the platform where the binary was built
	BuildPlatform = "-"

	// RepositoryUrl is an url to a GitHub repository
	RepositoryUrl = "https://github.com/shakahl/gqlassist"

	// ModuleName is the name if the main Go module
	ModuleName = "github.com/shakahl/gqlassist"

	// ProjectID is the project identifier
	ProjectID = "gqlassist"

	// ProjectName is the project repository name
	ProjectName = "gqlassist"

	// ProjectTitle is the project's human-readable name
	ProjectTitle = "Golang GraphQL type definitions gqlassist for consuming GraphQL server resources"

	// BinaryName is the binary name
	BinaryName = "gqlassist"

	// ShortDescription is a short description
	ShortDescription = ProjectTitle

	// LongDescription is a long description
	LongDescription = `
GQLAssist is a CLI tool for Go that helps you working with GraphQL servers by
generating struct definitions for remote GraphQL schema.
`

	// ConfigName is the name of config file
	ConfigName = ProjectID

	// ConfigFileNameBase is the name of config file
	ConfigFileNameBase = "." + ConfigName

	// ConfigFileName is the name of config file
	ConfigFileName = ConfigFileNameBase + ".yaml"
)
