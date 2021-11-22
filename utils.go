package jd_cookie

import (
    "strings"
)

func translate(str string, isWechat bool) string {

	if !isWechat {
		return str
	}

	tempMsg := str
	tempMsg = strings.Replace(tempMsg, "⭕", "[emoji=\\u2b55]", -1)
	tempMsg = strings.Replace(tempMsg, "🧧", "[emoji=\\uD83E\\uDDE7]", -1)
	tempMsg = strings.Replace(tempMsg, "🥚", "[emoji=\\ud83e\\udd5a]", -1)
	tempMsg = strings.Replace(tempMsg, "💰", "[emoji=\\ud83d\\udcb0]", -1)
	tempMsg = strings.Replace(tempMsg, "⏰", "[emoji=\\u23f0]", -1)
	tempMsg = strings.Replace(tempMsg, "🍒", "[emoji=\\ud83c\\udf52]", -1)
	tempMsg = strings.Replace(tempMsg, "🐶", "[emoji=\\ud83d\\udc36]", -1)
	tempMsg = strings.Replace(tempMsg, "🎰", "[emoji=\\ud83c\\udfb0]", -1)
	tempMsg = strings.Replace(tempMsg, "🌂", "[emoji=\\ud83c\\udf02]", -1)
	return tempMsg
}