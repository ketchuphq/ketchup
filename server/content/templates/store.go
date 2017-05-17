package templates

import (
	"github.com/ketchuphq/ketchup/proto/ketchup/models"
	"github.com/ketchuphq/ketchup/proto/ketchup/packages"
)

// ThemeStore is a interface to support multiple theme backends.
type ThemeStore interface {
	// List all themes in the store
	List() ([]*models.Theme, error)

	// Add a theme from a theme file
	Add(*models.Theme) error

	AddPackage(*packages.Package) error

	// Get a theme by name from the store
	Get(string) (*models.Theme, error)

	// Get a template for the given theme from the store
	GetTemplate(t *models.Theme, template string) (*models.ThemeTemplate, error)

	// Get an asset for the given theme from the store
	GetAsset(t *models.Theme, asset string) (*models.ThemeAsset, error)
}
