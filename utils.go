package jd_cookie

import (
	"fmt"
    "strings"
)

func appendActivityPath(str *string) {
	
	wrap := "\n\n"
  	path := `⏰提醒:⏰
【极速金币】京东极速版->我的->金币(极速版使用)
【京东赚赚】微信->京东赚赚小程序->底部赚好礼->提现无门槛红包(京东使用)
【京东秒杀】京东->中间频道往右划找到京东秒杀->中间点立即签到->兑换无门槛红包(京东使用)
【东东萌宠】京东->我的->东东萌宠,完成是京东红包,可以用于京东app的任意商品
【领现金】京东->我的->东东萌宠->领现金(微信提现+京东红包)
【东东农场】京东->我的->东东农场,完成是京东红包,可以用于京东app的任意商品
【京喜工厂】京喜->我的->京喜工厂,完成是商品红包,用于购买指定商品(不兑换会过期)
【其他】京喜红包只能在京喜使用,其他同理
	`
 	*str = fmt.Sprintln(*str, wrap, path)
}

func translateEmoji(str *string, isWechat bool) {

	if !isWechat {
		return
	}

	*str = strings.Replace(*str, "⭕", "[emoji=\\u2b55]", -1)
	*str = strings.Replace(*str, "🧧", "[emoji=\\uD83E\\uDDE7]", -1)
	*str = strings.Replace(*str, "🥚", "[emoji=\\ud83e\\udd5a]", -1)
	*str = strings.Replace(*str, "💰", "[emoji=\\ud83d\\udcb0]", -1)
	*str = strings.Replace(*str, "⏰", "[emoji=\\u23f0]", -1)
	*str = strings.Replace(*str, "🍒", "[emoji=\\ud83c\\udf52]", -1)
	*str = strings.Replace(*str, "🐶", "[emoji=\\ud83d\\udc36]", -1)
	*str = strings.Replace(*str, "🎰", "[emoji=\\ud83c\\udfb0]", -1)
	*str = strings.Replace(*str, "🌂", "[emoji=\\ud83c\\udf02]", -1)
}