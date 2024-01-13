package config

type InputFlags = map[string]string

var ValidFlags = InputFlags{
	"v": "v", // print version
	"y": "y", // do not ask user confirmation
	"q": "q", // do not print stats
}
