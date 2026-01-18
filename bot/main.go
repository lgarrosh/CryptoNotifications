package main

import (
	"fmt"
	"log"
	"os"

	telebot "gopkg.in/telebot.v3"
)

func main() {
	// Получаем токен бота из переменной окружения
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN не установлен. Установите переменную окружения.")
	}

	// Получаем API ключ CoinMarketCap
	cmcAPIKey := os.Getenv("COINMARKETCAP_API_KEY")
	if cmcAPIKey == "" {
		log.Fatal("COINMARKETCAP_API_KEY не установлен. Установите переменную окружения.")
	}

	// Создаем клиент для CoinMarketCap API
	cmcClient := NewCoinMarketCapClient(cmcAPIKey)

	// Настраиваем бота
	pref := telebot.Settings{
		Token:  botToken,
		Poller: &telebot.LongPoller{Timeout: 10},
	}

	bot, err := telebot.NewBot(pref)
	if err != nil {
		log.Fatal(err)
	}

	// Обработчик команды /start
	bot.Handle("/start", func(c telebot.Context) error {
		message := "👋 Привет! Я бот для получения котировок криптовалют.\n\n" +
			"Доступные команды:\n" +
			"/price <символ> - получить цену криптовалюты (например: /price BTC)\n" +
			"/price <символ1,символ2,...> - получить цены нескольких криптовалют (например: /price BTC,ETH,BNB)\n" +
			"/help - показать эту справку"
		return c.Send(message)
	})

	// Обработчик команды /help
	bot.Handle("/help", func(c telebot.Context) error {
		message := "📖 Справка по командам:\n\n" +
			"/price <символ> - получить цену одной криптовалюты\n" +
			"Пример: /price BTC\n\n" +
			"/price <символ1,символ2,...> - получить цены нескольких криптовалют\n" +
			"Пример: /price BTC,ETH,BNB\n\n" +
			"/help - показать эту справку"
		return c.Send(message)
	})

	// Обработчик команды /price для одной или нескольких криптовалют
	bot.Handle("/price", func(c telebot.Context) error {
		symbols := c.Message().Payload
		userID := c.Sender().ID
		username := c.Sender().Username

		if symbols == "" {
			log.Printf("[INFO] Пользователь %d (@%s) запросил /price без символов", userID, username)
			return c.Send("❌ Пожалуйста, укажите символ(ы) криптовалюты.\nПример: /price BTC или /price BTC,ETH,BNB")
		}

		log.Printf("[INFO] Пользователь %d (@%s) запросил котировки для: %s", userID, username, symbols)

		// Отправляем сообщение о загрузке
		msg, _ := c.Bot().Send(c.Chat(), "⏳ Загружаю данные...")

		// Получаем данные о криптовалюте(ах)
		cryptos, err := cmcClient.GetCryptocurrencyQuotes(symbols)
		if err != nil {
			log.Printf("[ERROR] Ошибка при получении котировок для пользователя %d (@%s), символы: %s, ошибка: %v", userID, username, symbols, err)
			c.Bot().Delete(msg)
			return c.Send("❌ Ошибка при получении данных: " + err.Error())
		}

		log.Printf("[INFO] Успешно получены котировки для пользователя %d (@%s), количество: %d", userID, username, len(cryptos))

		// Форматируем ответ в зависимости от количества криптовалют
		var response string
		if len(cryptos) == 1 {
			response = formatCryptoResponse(cryptos[0])
		} else {
			response = formatMultipleCryptoResponse(cryptos)
		}

		c.Bot().Delete(msg)
		return c.Send(response, telebot.ModeMarkdown)
	})

	log.Println("Бот запущен и готов к работе!")
	bot.Start()
}

// Форматирует ответ для одной криптовалюты
func formatCryptoResponse(crypto *Cryptocurrency) string {
	return "💰 *" + crypto.Name + " (" + crypto.Symbol + ")*\n\n" +
		"💵 Цена: $" + formatPrice(crypto.Price) + "\n" +
		"📊 Изменение за 24ч: " + formatPercentChange(crypto.PercentChange24h) + "\n" +
		"📈 Рыночная капитализация: $" + formatMarketCap(crypto.MarketCap) + "\n" +
		"💹 Объем за 24ч: $" + formatVolume(crypto.Volume24h)
}

// Форматирует ответ для нескольких криптовалют
func formatMultipleCryptoResponse(cryptos []*Cryptocurrency) string {
	response := "💰 *Котировки криптовалют:*\n\n"
	for _, crypto := range cryptos {
		response += "• *" + crypto.Symbol + "* - $" + formatPrice(crypto.Price) +
			" (" + formatPercentChange(crypto.PercentChange24h) + ")\n"
	}
	return response
}

// Форматирует цену
func formatPrice(price float64) string {
	if price >= 1 {
		return formatNumber(price, 2)
	}
	return formatNumber(price, 8)
}

// Форматирует процент изменения
func formatPercentChange(change float64) string {
	sign := ""
	if change > 0 {
		sign = "📈 +"
	} else if change < 0 {
		sign = "📉 "
	}
	return sign + formatNumber(change, 2) + "%"
}

// Форматирует рыночную капитализацию
func formatMarketCap(marketCap float64) string {
	if marketCap >= 1e12 {
		return formatNumber(marketCap/1e12, 2) + "T"
	} else if marketCap >= 1e9 {
		return formatNumber(marketCap/1e9, 2) + "B"
	} else if marketCap >= 1e6 {
		return formatNumber(marketCap/1e6, 2) + "M"
	}
	return formatNumber(marketCap, 2)
}

// Форматирует объем
func formatVolume(volume float64) string {
	if volume >= 1e9 {
		return formatNumber(volume/1e9, 2) + "B"
	} else if volume >= 1e6 {
		return formatNumber(volume/1e6, 2) + "M"
	}
	return formatNumber(volume, 2)
}

// Форматирует число с заданным количеством знаков после запятой
func formatNumber(num float64, decimals int) string {
	format := fmt.Sprintf("%%.%df", decimals)
	return fmt.Sprintf(format, num)
}
