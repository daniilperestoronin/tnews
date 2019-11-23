package lang

import "github.com/abadojack/whatlanggo"

func DetectLanguage(text string) string {
	info := whatlanggo.Detect(text)
	return info.Lang.Iso6391()
}
