package utils

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// RateLimiter 全局限流器，控制请求频率
type RateLimiter struct {
	ticker *time.Ticker
	tokens chan struct{}
}

// 全局限流器实例（每秒最多 2 个 Codeforces 请求）
var codeforcesLimiter = NewRateLimiter(2, time.Second)

// NewRateLimiter 创建新的限流器
// rate: 每个周期允许的请求数
// per: 时间周期
func NewRateLimiter(rate int, per time.Duration) *RateLimiter {
	limiter := &RateLimiter{
		ticker: time.NewTicker(per / time.Duration(rate)),
		tokens: make(chan struct{}, rate),
	}

	// 初始化 token bucket
	for i := 0; i < rate; i++ {
		limiter.tokens <- struct{}{}
	}

	// 定期补充 tokens
	go func() {
		for range limiter.ticker.C {
			select {
			case limiter.tokens <- struct{}{}:
			default:
			}
		}
	}()

	return limiter
}

// Wait 等待获取一个 token
func (rl *RateLimiter) Wait() {
	<-rl.tokens
}

// getRandomUserAgent 返回随机 User-Agent
func getRandomUserAgent() string {
	userAgents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:124.0) Gecko/20100101 Firefox/124.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.4 Safari/605.1.15",
	}
	rand.Seed(time.Now().UnixNano())
	return userAgents[rand.Intn(len(userAgents))]
}

// createHTTPClient 创建带超时的 HTTP 客户端
func createHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 15 * time.Second, // 15 秒超时
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 10,
			IdleConnTimeout:     90 * time.Second,
		},
	}
}

// FetchCodeforcesRatingWithRetry 带重试的 Codeforces 分数获取
func FetchCodeforcesRatingWithRetry(username string, maxRetries int) (string, error) {
	var lastErr error

	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			// 重试前等待，使用指数退避
			waitTime := time.Duration(attempt) * 2 * time.Second
			time.Sleep(waitTime)
		}

		rating, err := FetchCodeforcesRatingOnce(username)
		if err == nil {
			return rating, nil
		}

		lastErr = err
	}

	return "", fmt.Errorf("尝试 %d 次后失败: %v", maxRetries, lastErr)
}

// FetchCodeforcesRatingOnce 单次 Codeforces 分数获取（带限流）
func FetchCodeforcesRatingOnce(username string) (string, error) {
	// 等待限流器允许请求
	codeforcesLimiter.Wait()

	url := fmt.Sprintf("https://codeforces.com/profile/%s", username)

	client := createHTTPClient()
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %v", err)
	}

	// 使用随机 User-Agent
	req.Header.Set("User-Agent", getRandomUserAgent())
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	req.Header.Set("Connection", "keep-alive")

	res, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("请求失败: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return "", fmt.Errorf("请求失败，状态码: %d", res.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return "", fmt.Errorf("解析失败: %v", err)
	}

	var rating string
	// 查找包含 Contest rating 的 <li>，然后找它的第一个带 class user-* 的 <span>
	doc.Find("li").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if strings.Contains(s.Text(), "Contest rating:") {
			s.Find("span").EachWithBreak(func(j int, span *goquery.Selection) bool {
				class, _ := span.Attr("class")
				if strings.HasPrefix(class, "user-") {
					rating = strings.TrimSpace(span.Text())
					return false // 找到当前分数就跳出
				}
				return true
			})
			return false // 找到包含 Contest rating 的 li 就跳出
		}
		return true
	})

	if rating == "" {
		return "", fmt.Errorf("未能找到当前 rating")
	}
	return rating, nil
}

// FetchCodeforcesRating 保持向后兼容的接口（内部使用重试版本）
func FetchCodeforcesRating(username string) (string, error) {
	return FetchCodeforcesRatingWithRetry(username, 3) // 默认重试 3 次
}

// FetchAtCoderRatingWithRetry 带重试的 AtCoder 分数获取
func FetchAtCoderRatingWithRetry(username string, maxRetries int) (string, error) {
	var lastErr error

	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			waitTime := time.Duration(attempt) * 2 * time.Second
			time.Sleep(waitTime)
		}

		rating, err := FetchAtCoderRatingOnce(username)
		if err == nil {
			return rating, nil
		}

		lastErr = err
	}

	return "", fmt.Errorf("尝试 %d 次后失败: %v", maxRetries, lastErr)
}

// FetchAtCoderRatingOnce 单次 AtCoder 分数获取
func FetchAtCoderRatingOnce(username string) (string, error) {
	url := fmt.Sprintf("https://atcoder.jp/users/%s", username)

	client := createHTTPClient()
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %v", err)
	}

	req.Header.Set("User-Agent", getRandomUserAgent())
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")

	res, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("请求失败: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return "", fmt.Errorf("请求失败，状态码: %d", res.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return "", fmt.Errorf("解析失败: %v", err)
	}

	var rating string
	doc.Find("table.dl-table tr").EachWithBreak(func(i int, s *goquery.Selection) bool {
		thText := strings.TrimSpace(s.Find("th").Text())
		if thText == "Rating" {
			s.Find("td span").EachWithBreak(func(j int, span *goquery.Selection) bool {
				class, _ := span.Attr("class")
				if strings.HasPrefix(class, "user-") {
					rating = strings.TrimSpace(span.Text())
					return false
				}
				return true
			})
			return false
		}
		return true
	})

	if rating == "" {
		return "", fmt.Errorf("未能找到当前 AtCoder rating")
	}
	return rating, nil
}

// FetchAtCoderRating 保持向后兼容的接口
func FetchAtCoderRating(username string) (string, error) {
	return FetchAtCoderRatingWithRetry(username, 3)
}

// FetchNowcoderRatingWithRetry 带重试的 Nowcoder 分数获取
func FetchNowcoderRatingWithRetry(userID string, maxRetries int) (string, error) {
	var lastErr error

	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			waitTime := time.Duration(attempt) * 2 * time.Second
			time.Sleep(waitTime)
		}

		rating, err := FetchNowcoderRatingOnce(userID)
		if err == nil {
			return rating, nil
		}

		lastErr = err
	}

	return "", fmt.Errorf("尝试 %d 次后失败: %v", maxRetries, lastErr)
}

// FetchNowcoderRatingOnce 单次 Nowcoder 分数获取
func FetchNowcoderRatingOnce(userID string) (string, error) {
	url := fmt.Sprintf("https://ac.nowcoder.com/acm/contest/profile/%s", userID)

	client := createHTTPClient()
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %v", err)
	}

	req.Header.Set("User-Agent", getRandomUserAgent())
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")

	res, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("请求失败: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return "", fmt.Errorf("请求失败，状态码: %d", res.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return "", fmt.Errorf("解析失败: %v", err)
	}

	var rating string
	found := false
	for i := 1; i <= 7; i++ {
		selector := fmt.Sprintf("a.state-num.rate-score%d", i)
		doc.Find(selector).EachWithBreak(func(_ int, s *goquery.Selection) bool {
			text := strings.TrimSpace(s.Text())
			if text != "" {
				rating = text
				found = true
				return false
			}
			return true
		})
		if found {
			break
		}
	}

	if rating == "" {
		return "", fmt.Errorf("未能找到 Nowcoder Rating")
	}
	return rating, nil
}

// FetchNowcoderRating 保持向后兼容的接口
func FetchNowcoderRating(userID string) (string, error) {
	return FetchNowcoderRatingWithRetry(userID, 3)
}
