package utils

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func FetchCodeforcesRating(username string) (string, error) {
	url := fmt.Sprintf("https://codeforces.com/profile/%s", username)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %v", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")

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

func FetchAtCoderRating(username string) (string, error) {
	url := fmt.Sprintf("https://atcoder.jp/users/%s", username)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %v", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")

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
	// 找到 table 中包含 “Rating” 的 <th>，然后去找其下一个 <td> 里的 span.user-*
	doc.Find("table.dl-table tr").EachWithBreak(func(i int, s *goquery.Selection) bool {
		thText := strings.TrimSpace(s.Find("th").Text())
		if thText == "Rating" {
			// 找 td 下第一个 class="user-*"
			s.Find("td span").EachWithBreak(func(j int, span *goquery.Selection) bool {
				class, _ := span.Attr("class")
				if strings.HasPrefix(class, "user-") {
					rating = strings.TrimSpace(span.Text())
					return false // 找到了
				}
				return true
			})
			return false // 找到 Rating 行就退出
		}
		return true
	})

	if rating == "" {
		return "", fmt.Errorf("未能找到当前 AtCoder rating")
	}
	return rating, nil
}

func FetchNowcoderRating(userID string) (string, error) {
	url := fmt.Sprintf("https://ac.nowcoder.com/acm/contest/profile/%s", userID)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %v", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")

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
	doc.Find("a.rate-score4").Each(func(i int, s *goquery.Selection) {
		rating = strings.TrimSpace(s.Text())
	})

	if rating == "" {
		return "", fmt.Errorf("未能找到 Nowcoder Rating")
	}
	return rating, nil
}
