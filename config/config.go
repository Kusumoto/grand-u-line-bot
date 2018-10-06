package config

// Config  model configuration for viper
type Config struct {
	CheckRegisterMailAPIUrl string `json:"checkRegisterMailAPIUrl"`
	CheckCheckBalanceAPIUrl string `json:"checkCheckBalanceAPIUrl"`
	PhoneNumber             string `json:"phoneNumber"`
	UnitID                  string `json:"unitID"`
	ProjectCode             string `json:"projectCode"`
	LineAccessToken         string `json:"lineAccessToken"`
	LineChannelID           string `json:"lineChannelID"`
	LineChannelSecret       string `json:"lineChannelSecret"`
}
