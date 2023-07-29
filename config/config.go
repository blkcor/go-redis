package config

type ServerProperties struct {
	Bind           string `cfg:"bind"`
	Port           int    `cfg:"port"`
	AppendOnly     bool   `cfg:"appendOnly"`
	AppendFilename string `cfg:"appendFilename"`
	MaxClients     int    `cfg:"maxClients"`
	RequirePass    string `cfg:"requirePass"`
	Databases      int    `cfg:"databases"`

	Peers []string `cfg:"peers"`
	Self  string   `cfg:"self"`
}

// Properties holds global config properties
var Properties *ServerProperties

func init() {
	//default config
	Properties = &ServerProperties{
		Bind:       "127.0.0.1",
		Port:       6379,
		AppendOnly: false,
	}
}
