package eclient

import "time"

// A CatalogYAML contains a single root node of the catalog.
type CatalogYAML struct {
	Endpoints []string     `yaml:"endpoints"`
	Category  CategoryYAML `yaml:"catalog"`
}

// A CategoryYAML is a single node in the categories tree in the YAML file.
type CategoryYAML struct {
	Segment    string          `yaml:"segment"`
	Name       string          `yaml:"name"`
	Categories []*CategoryYAML `yaml:"categories,omitempty"`
}

// ProductContainerYAML container for a YAML product.
type ProductContainerYAML struct {
	Product ProductApplyYAML `yaml:"product"`
}

// ProductImageApplyYAML contains the product image data.
type ProductImageApplyYAML struct {
	Path  string `yaml:"path"`
	Title string `yaml:"title"`
}

// PriceYAML YAML price
type PriceYAML struct {
	Break     int `yaml:"break"`
	UnitPrice int `yaml:"unit_price"`
}

// ProductApplyYAML contains fields used when applying a product.
type ProductApplyYAML struct {
	Path    string                   `yaml:"path"`
	SKU     string                   `yaml:"sku"`
	Name    string                   `yaml:"name"`
	Images  []*ProductImageApplyYAML `yaml:"images"`
	Prices  map[string][]PriceYAML   `yaml:"prices"`
	Content interface{}              `yaml:"content"`
}

// ProductYAML contains all the fields that comprise a product in the catalog.
type ProductYAML struct {
	SKU      string                     `yaml:"sku"`
	EAN      string                     `yaml:"ean"`
	Path     string                     `yaml:"path"`
	Name     string                     `yaml:"name"`
	Images   []*ProductImage            `yaml:"images"`
	Pricing  map[string]*ProductPricing `yaml:"pricing"`
	Content  interface{}                `yaml:"content"`
	Created  time.Time                  `json:"created"`
	Modified time.Time                  `json:"modified"`
}

// ProductCategoryRelationsYAML for the product to category relations.
type ProductCategoryRelationsYAML struct {
	Rels map[string]*ProductSetYAML `yaml:"product_category_relations"`
}

// An ProductSetYAML holds a single product to catalog association.
type ProductSetYAML struct {
	Products []string `yaml:"products"`
}
