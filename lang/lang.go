package lang

import "github.com/abadojack/whatlanggo"

func detectLanguage(text string) string {
	info := whatlanggo.Detect(text)
	return info.Lang.String()
}

