package lng

var ru = map[string]string{
	"welcome":          "Добро пожаловать в ABT Miner!",
	"description":      "ABT Miner – это мини приложение Телеграм для фарминга ABT Coins.",
	"coins":            "ABT coins – это уникальная, интегральная монета всей экосистемы ABT, в которую входят игровые и бизнес проекты направленные на извлечение прибыли и капитализацию личных усилий.",
	"about":            "Подробнее можно посмотреть здесь: https://aibetrade.com/abt_coin/abtcoins_ru.html",
	"ecosystem":        "Экосистема ABT - Ваш новый, легкий способ накопить ABT coins и в последствии заработать в мире криптовалют.",
	"harvest":          "⛏️ Фармите монету – от 0.9 ABT Coin за 24 часа.",
	"invite":           "🤝 Собирайте свою команду — от +1 ABT Coin за каждого новичка, 10 уровней выплат реферального бонуса за фарминг команды.",
	"exchange":         "🚀 Обменивайте и продавайте монету — капитализируйте свои усилия. Количество опций будет увеличиваться по мере развития проекта.",
	"earn":             "💵 Начните зарабатывать на криптовалютах с нашими сервисами.",
	"post":             "Сделай важный шаг на пути к личному успеху с ABT сейчас! 🌟",
	"register":         "Привет!\nПоздравляем тебя с новым последователем и сообщаем, что по твоей ссылке зарегистрировался:",
	"inviteReward":     "Ты получаешь +%s ABT за приглашение.",
	"invited":          "Всего тобой уже приглашено друзей:",
	"link":             "https://t.me/+AqZL_Ou8PrM5Njk6",
	"openApp":          "Открыть приложение",
	"ourCommunity":     "Наше сообщество",
	"autoclaimSuccess": "Автофарминг успешно совершен! +%s ABT.",
	"autoclaimEnd":     "Автофарминг завершен.",
	"refBuy":           "За покупку друга вы получаете +%s %s.",
}

var en = map[string]string{
	"welcome":          "Welcome to ABT Miner!",
	"description":      "ABT Miner is a mini Telegram application for farming ABT Coins.",
	"coins":            "ABT coins are a unique, integral currency of the entire ABT ecosystem, which includes gaming and business projects aimed at generating profit and capitalizing on individual efforts.",
	"about":            "Learn more here: https://aibetrade.com/abt_coin/abtcoins.html",
	"ecosystem":        "The ABT ecosystem is your new, easy way to accumulate ABT coins and subsequently earn in the world of cryptocurrencies.",
	"harvest":          "⛏️ Farm the coin – earn from 0.9 ABT Coin in 24 hours.",
	"invite":           "🤝 Build your team – get from +1 ABT Coin for each new recruit, plus 10 levels of referral farming bonuses from your team.",
	"exchange":         "🚀 Exchange and sell the coin – capitalize on your efforts. The number of options will increase as the project develops.",
	"earn":             "💵 Start earning with cryptocurrencies through our services.",
	"post":             "Take an important step towards personal success with ABT now! 🌟",
	"register":         "Hi!\nCongratulations to you with your new follower! And we inform you that new user registered using your link:",
	"inviteReward":     "You are getting +%s ABT for inviting.",
	"invited":          "Total of your already invited friends:",
	"link":             "https://t.me/aibetradecom",
	"openApp":          "Open App",
	"ourCommunity":     "Our Community",
	"autoclaimSuccess": "Auto-farming successful! +%s ABT.",
	"autoclaimEnd":     "Auto-farming ended.",
	"refBuy":           "For friend's purchase you get +%s %s.",
}

func Get(keylang ...string) string {
	if len(keylang) == 0 {
		return ""
	}
	key := keylang[0]
	locales := en
	if len(keylang) > 1 && keylang[1] == "ru" {
		locales = ru
	}
	k, ok := locales[key]
	if !ok {
		return ""
	}
	return k
}
