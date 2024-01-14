package config

type InputFlags = map[string]string

var ValidFlags = InputFlags{
	"v":   "v",   // print version
	"y":   "y",   // do not ask user confirmation
	"q":   "q",   // do not print stats
	"nfg": "nfg", // do not use folder groups
}

func IsValidFlag(flag string) bool {
	if _, ok := ValidFlags[flag]; ok {
		return true
	}
	return false
}

func IsFlagExists(flag string, flags *InputFlags) bool {
	if _, ok := ValidFlags[flag]; ok {
		if _, _ok := (*flags)[flag]; _ok {
			return true
		}
	}
	return false
}
