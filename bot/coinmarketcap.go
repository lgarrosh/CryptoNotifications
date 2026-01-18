package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

const baseURL = "https://pro-api.coinmarketcap.com/v2"

// CoinMarketCapClient представляет клиент для работы с CoinMarketCap API
type CoinMarketCapClient struct {
	apiKey     string
	httpClient *http.Client
	baseURL    string
}

// Cryptocurrency представляет данные о криптовалюте
type Cryptocurrency struct {
	ID               int     `json:"id"`
	Name             string  `json:"name"`
	Symbol           string  `json:"symbol"`
	Price            float64 `json:"price"`
	PercentChange24h float64 `json:"percent_change_24h"`
	MarketCap        float64 `json:"market_cap"`
	Volume24h        float64 `json:"volume_24h"`
	LastUpdated      string  `json:"last_updated"`
}

// apiTag представляет тег криптовалюты
type apiTag struct {
	Slug     string `json:"slug"`
	Name     string `json:"name"`
	Category string `json:"category"`
}

// apiStatus представляет структуру статуса ответа API v2
type apiStatus struct {
	Timestamp    time.Time   `json:"timestamp"`
	ErrorCode    int         `json:"error_code"`
	ErrorMessage interface{} `json:"error_message"`
	Elapsed      int         `json:"elapsed"`
	CreditCount  int         `json:"credit_count"`
	Notice       interface{} `json:"notice"`
}

// apiQuoteResponse представляет структуру ответа API v2 для котировок
type apiQuoteResponse struct {
	Status apiStatus                          `json:"status"`
	Data   map[string][]apiCryptocurrencyData `json:"data"`
}

// apiCryptocurrencyData представляет данные криптовалюты из API v2
type apiCryptocurrencyData struct {
	ID                            int                     `json:"id"`
	Name                          string                  `json:"name"`
	Symbol                        string                  `json:"symbol"`
	Slug                          string                  `json:"slug"`
	NumMarketPairs                int                     `json:"num_market_pairs"`
	DateAdded                     time.Time               `json:"date_added"`
	Tags                          []apiTag                `json:"tags"`
	MaxSupply                     *float64                `json:"max_supply"`
	CirculatingSupply             float64                 `json:"circulating_supply"`
	TotalSupply                   float64                 `json:"total_supply"`
	IsActive                      int                     `json:"is_active"`
	InfiniteSupply                bool                    `json:"infinite_supply"`
	MintedMarketCap               float64                 `json:"minted_market_cap"`
	Platform                      interface{}             `json:"platform"`
	CmcRank                       int                     `json:"cmc_rank"`
	IsFiat                        int                     `json:"is_fiat"`
	SelfReportedCirculatingSupply interface{}             `json:"self_reported_circulating_supply"`
	SelfReportedMarketCap         interface{}             `json:"self_reported_market_cap"`
	TvlRatio                      interface{}             `json:"tvl_ratio"`
	LastUpdated                   time.Time               `json:"last_updated"`
	Quote                         map[string]apiQuoteData `json:"quote"`
}

// apiQuoteData представляет данные котировки из API v2
type apiQuoteData struct {
	Price                 float64     `json:"price"`
	Volume24h             float64     `json:"volume_24h"`
	VolumeChange24h       float64     `json:"volume_change_24h"`
	PercentChange1h       float64     `json:"percent_change_1h"`
	PercentChange24h      float64     `json:"percent_change_24h"`
	PercentChange7d       float64     `json:"percent_change_7d"`
	PercentChange30d      float64     `json:"percent_change_30d"`
	PercentChange60d      float64     `json:"percent_change_60d"`
	PercentChange90d      float64     `json:"percent_change_90d"`
	MarketCap             float64     `json:"market_cap"`
	MarketCapDominance    float64     `json:"market_cap_dominance"`
	FullyDilutedMarketCap float64     `json:"fully_diluted_market_cap"`
	Tvl                   interface{} `json:"tvl"`
	LastUpdated           time.Time   `json:"last_updated"`
}

// NewCoinMarketCapClient создает новый клиент для CoinMarketCap API
func NewCoinMarketCapClient(apiKey string) *CoinMarketCapClient {
	return &CoinMarketCapClient{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL: baseURL,
	}
}

// makeRequest выполняет HTTP запрос к CoinMarketCap API
func (c *CoinMarketCapClient) makeRequest(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("[ERROR] Ошибка создания HTTP запроса: %v, URL: %s", err, url)
		return nil, fmt.Errorf("ошибка создания запроса: %v", err)
	}

	req.Header.Set("X-CMC_PRO_API_KEY", c.apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Printf("[ERROR] Ошибка выполнения HTTP запроса: %v, URL: %s", err, url)
		return nil, fmt.Errorf("ошибка выполнения запроса: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)

		// Обработка ошибок 4xx (клиентские ошибки)
		if resp.StatusCode >= 400 && resp.StatusCode < 500 {
			log.Printf("[ERROR] Ошибка клиента (4xx): статус %d, URL: %s, Response: %s", resp.StatusCode, url, string(body))
			return nil, fmt.Errorf("некорректный запрос или не найден символ")
		}

		// Обработка ошибок 5xx (ошибки сервера)
		if resp.StatusCode >= 500 && resp.StatusCode < 600 {
			log.Printf("[ERROR] Ошибка сервера (5xx): статус %d, URL: %s, Response: %s", resp.StatusCode, url, string(body))
			return nil, fmt.Errorf("временная ошибка сервера, попробуйте позже")
		}

		// Обработка других статусов
		log.Printf("[ERROR] API вернул неожиданный статус: %d, URL: %s, Response: %s", resp.StatusCode, url, string(body))
		return nil, fmt.Errorf("API вернул статус %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[ERROR] Ошибка чтения тела ответа: %v, URL: %s", err, url)
		return nil, fmt.Errorf("ошибка чтения ответа: %v", err)
	}

	return body, nil
}

// parseQuoteResponse парсит JSON ответ и преобразует его в массив объектов Cryptocurrency
// Для каждого символа выбирается элемент с минимальным ID
func (c *CoinMarketCapClient) parseQuoteResponse(body []byte) ([]*Cryptocurrency, error) {
	var apiResp apiQuoteResponse

	if err := json.Unmarshal(body, &apiResp); err != nil {
		log.Printf("[ERROR] Ошибка парсинга JSON ответа: %v, Body length: %d", err, len(body))
		return nil, fmt.Errorf("ошибка парсинга JSON: %v", err)
	}

	if apiResp.Status.ErrorCode != 0 {
		errorMsg := "неизвестная ошибка"
		if msg, ok := apiResp.Status.ErrorMessage.(string); ok {
			errorMsg = msg
		}
		log.Printf("[ERROR] API вернул ошибку: ErrorCode=%d, ErrorMessage=%v", apiResp.Status.ErrorCode, apiResp.Status.ErrorMessage)
		return nil, fmt.Errorf("API ошибка: %s", errorMsg)
	}

	var cryptos []*Cryptocurrency

	// В v2 API data - это map[string][]apiCryptocurrencyData, где ключ - символ криптовалюты
	// Для каждого символа выбираем элемент с минимальным ID
	for symbol, cryptoArray := range apiResp.Data {
		if len(cryptoArray) == 0 {
			log.Printf("[WARN] Пустой массив для символа: %s", symbol)
			continue
		}

		// Находим элемент с минимальным ID
		minIDCrypto := cryptoArray[0]
		for _, item := range cryptoArray {
			if item.ID < minIDCrypto.ID {
				minIDCrypto = item
			}
		}

		// Получаем данные из USD котировки
		usdQuote, ok := minIDCrypto.Quote["USD"]
		if !ok {
			log.Printf("[WARN] USD котировка не найдена для символа: %s (ID: %d)", symbol, minIDCrypto.ID)
			continue
		}

		cryptos = append(cryptos, &Cryptocurrency{
			ID:               minIDCrypto.ID,
			Name:             minIDCrypto.Name,
			Symbol:           minIDCrypto.Symbol,
			Price:            usdQuote.Price,
			PercentChange24h: usdQuote.PercentChange24h,
			MarketCap:        usdQuote.MarketCap,
			Volume24h:        usdQuote.Volume24h,
			LastUpdated:      usdQuote.LastUpdated.Format(time.RFC3339),
		})
	}

	if len(cryptos) == 0 {
		log.Printf("[ERROR] Криптовалюты не найдены в ответе API")
		return nil, fmt.Errorf("криптовалюты не найдены")
	}

	log.Printf("[INFO] Успешно распарсено %d криптовалют", len(cryptos))
	return cryptos, nil
}

// GetCryptocurrencyQuotes получает котировки одной или нескольких криптовалют
// symbols может быть одним символом или несколькими через запятую
func (c *CoinMarketCapClient) GetCryptocurrencyQuotes(symbols string) ([]*Cryptocurrency, error) {
	// Разделяем символы по запятой и очищаем
	symbolList := strings.Split(symbols, ",")
	var cleanSymbols []string
	for _, s := range symbolList {
		cleaned := strings.ToUpper(strings.TrimSpace(s))
		if cleaned != "" {
			cleanSymbols = append(cleanSymbols, cleaned)
		}
	}

	if len(cleanSymbols) == 0 {
		log.Printf("[ERROR] Не указаны символы криптовалют в запросе: %s", symbols)
		return nil, fmt.Errorf("не указаны символы криптовалют")
	}

	symbolsParam := strings.Join(cleanSymbols, ",")
	url := fmt.Sprintf("%s/cryptocurrency/quotes/latest?symbol=%s", c.baseURL, symbolsParam)

	log.Printf("[INFO] Запрос котировок для символов: %s", symbolsParam)
	body, err := c.makeRequest(url)
	if err != nil {
		log.Printf("[ERROR] Ошибка при запросе котировок для символов %s: %v", symbolsParam, err)
		return nil, err
	}

	cryptos, err := c.parseQuoteResponse(body)
	if err != nil {
		log.Printf("[ERROR] Ошибка при парсинге ответа для символов %s: %v", symbolsParam, err)
		return nil, err
	}

	return cryptos, nil
}
