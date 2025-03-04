// SPDX-License-Identifier: Apache-2.0

package complytime

import (
	"os"
	"path/filepath"

	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/oscal-compass/oscal-sdk-go/generators"
)

// The assessment results md file needs the catalog title information
// LoadCatalogSource returns an OSCAL catalogs from a given application directory.

func LoadCatalogSource(appDir ApplicationDirectory) (*oscalTypes.Catalog, error) {
	var path string
	path = "controls/catalog.json"
	appDirPath := appDir.AppDir()

	// A path relative to the root of the complytime application directory
	path = filepath.Join(appDirPath, path)
	cleanedPath := filepath.Clean(path)

	sourceFile, err := os.Open(cleanedPath)
	if err != nil {
		return nil, err
	}
	defer sourceFile.Close()

	return generators.NewCatalog(sourceFile)
}
