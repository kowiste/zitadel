package domain

// BrandingConfig holds the branding configuration for an organization
type BrandingConfig struct {
	PrimaryColor         string `json:"primaryColor"`
	BackgroundColor      string `json:"backgroundColor"`
	WarnColor            string `json:"warnColor"`
	FontColor            string `json:"fontColor"`
	PrimaryColorDark     string `json:"primaryColorDark"`
	BackgroundColorDark  string `json:"backgroundColorDark"`
	WarnColorDark        string `json:"warnColorDark"`
	FontColorDark        string `json:"fontColorDark"`
	HideLoginNameSuffix  bool   `json:"hideLoginNameSuffix"`
	DisableWatermark     bool   `json:"disableWatermark"`
}

// DefaultBrandingConfig returns the default branding configuration
func DefaultBrandingConfig() *BrandingConfig {
	return &BrandingConfig{
		PrimaryColor:         "#5469d4",
		BackgroundColor:      "#ffffff",
		WarnColor:            "#ff3b5b",
		FontColor:            "#000000",
		PrimaryColorDark:     "#5469d4",
		BackgroundColorDark:  "#1a1a1a",
		WarnColorDark:        "#ff3b5b",
		FontColorDark:        "#ffffff",
		HideLoginNameSuffix:  true,
		DisableWatermark:     true,
	}
}
