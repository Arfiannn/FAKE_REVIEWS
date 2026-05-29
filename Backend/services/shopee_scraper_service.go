package services

import (
	"BE_FAKE_REVIEW/config"
	"BE_FAKE_REVIEW/models"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/playwright-community/playwright-go"
)

type ShopeeScraperService struct {
	Config config.Config
}

func NewShopeeScraperService(cfg config.Config) ShopeeScraperService {
	return ShopeeScraperService{
		Config: cfg,
	}
}

type shopeeRatingResponse struct {
	Data struct {
		Ratings []struct {
			CmtID          int64  `json:"cmtid"`
			AuthorUsername string `json:"author_username"`
			RatingStar     int    `json:"rating_star"`
			Comment        string `json:"comment"`
			CTime          int64  `json:"ctime"`
			ShopID         int64  `json:"shopid"`
			ProductItems   []struct {
				Name string `json:"name"`
			} `json:"product_items"`
		} `json:"ratings"`
	} `json:"data"`
}

type shopeeShopResponse struct {
	Data struct {
		Name string `json:"name"`
	} `json:"data"`
}

func (s ShopeeScraperService) ScrapeReviews(productURL string, limit int) (models.ShopeeScraperResult, error) {
	if limit <= 0 {
		limit = 1
	}

	pw, err := playwright.Run()
	if err != nil {
		return models.ShopeeScraperResult{}, err
	}
	defer pw.Stop()

	headless := strings.ToLower(s.Config.ShopeeHeadless) == "true"

	fmt.Println("Chrome User Data Dir:", s.Config.ShopeeUserDataDir)
	fmt.Println("Shopee Headless:", s.Config.ShopeeHeadless)

	context, err := pw.Chromium.LaunchPersistentContext(
		s.Config.ShopeeUserDataDir,
		playwright.BrowserTypeLaunchPersistentContextOptions{
			Headless: playwright.Bool(headless),
			Channel:  playwright.String("chrome"),
			Args: []string{
				"--disable-blink-features=AutomationControlled",
				"--start-maximized",
			},
			Viewport: nil,
		},
	)
	if err != nil {
		return models.ShopeeScraperResult{}, err
	}
	defer context.Close()

	page, err := context.NewPage()
	if err != nil {
		return models.ShopeeScraperResult{}, err
	}

	results := make([]models.ShopeeReview, 0)
	seenIDs := map[int64]bool{}
	productName := ""
	shopName := ""
	var mu sync.Mutex

	page.OnResponse(func(response playwright.Response) {
		if !strings.Contains(response.URL(), "get_ratings") {
			return
		}

		go func() {
			body, err := response.Body()
			if err != nil {
				fmt.Println("Gagal membaca response get_ratings:", err)
				return
			}

			var ratingResponse shopeeRatingResponse
			if err := json.Unmarshal(body, &ratingResponse); err != nil {
				fmt.Println("Gagal parsing response get_ratings:", err)
				return
			}

			mu.Lock()
			defer mu.Unlock()

			for _, r := range ratingResponse.Data.Ratings {
				if len(results) >= limit {
					break
				}

				if r.CmtID == 0 || seenIDs[r.CmtID] {
					continue
				}

				if strings.TrimSpace(r.Comment) == "" {
					continue
				}

				seenIDs[r.CmtID] = true

				itemName := productName
				if itemName == "" && len(r.ProductItems) > 0 {
					itemName = r.ProductItems[0].Name
				}

				if shopName == "" && r.ShopID != 0 {
					name, err := s.getShopName(page, r.ShopID)
					if err == nil {
						shopName = name
					}
				}

				reviewDate := ""
				if r.CTime > 0 {
					reviewDate = time.Unix(r.CTime, 0).Format("2006-01-02 15:04:05")
				}

				results = append(results, models.ShopeeReview{
					ProductName: itemName,
					ShopName:    shopName,
					Username:    r.AuthorUsername,
					Rating:      r.RatingStar,
					Review:      r.Comment,
					Date:        reviewDate,
				})
			}

			fmt.Println("Total review terkumpul:", len(results))
		}()
	})

	fmt.Println("Membuka Shopee homepage...")

	_, _ = page.Goto("https://shopee.co.id", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateDomcontentloaded,
		Timeout:   playwright.Float(60000),
	})

	page.WaitForTimeout(5000)

	fmt.Println("Membuka halaman produk:", productURL)

	_, err = page.Goto(productURL, playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateCommit,
		Timeout:   playwright.Float(60000),
	})
	if err != nil {
		page.WaitForTimeout(5000)

		_, err = page.Goto(productURL, playwright.PageGotoOptions{
			WaitUntil: playwright.WaitUntilStateCommit,
			Timeout:   playwright.Float(60000),
		})

		if err != nil {
			return models.ShopeeScraperResult{}, err
		}
	}

	page.WaitForTimeout(10000)

	title, err := page.Title()
	if err == nil {
		productName = title
	}

	if h1, err := page.Locator("h1").First().InnerText(playwright.LocatorInnerTextOptions{
		Timeout: playwright.Float(7000),
	}); err == nil {
		productName = h1
	}

	for i := 0; i < 25; i++ {
		mu.Lock()
		currentTotal := len(results)
		mu.Unlock()

		if currentTotal >= limit {
			break
		}

		_, _ = page.Evaluate(`window.scrollBy(0, 400)`)
		page.WaitForTimeout(800)

		count, _ := page.Locator("text=Penilaian").Count()
		if count > 0 {
			break
		}
	}

	penilaian := page.Locator("text=Penilaian").First()
	if count, _ := penilaian.Count(); count > 0 {
		_ = penilaian.Click(playwright.LocatorClickOptions{
			Timeout: playwright.Float(7000),
		})
		page.WaitForTimeout(5000)
	}

	pageNumber := 1
	maxPage := 50
	lastTotal := 0
	sameCount := 0

	for pageNumber <= maxPage {
		mu.Lock()
		currentTotal := len(results)
		mu.Unlock()

		if currentTotal >= limit {
			break
		}

		fmt.Println("Halaman review:", pageNumber, "| total:", currentTotal)

		if currentTotal == lastTotal {
			sameCount++
			if sameCount >= 3 {
				break
			}
		} else {
			sameCount = 0
		}

		lastTotal = currentTotal

		nextSuccess := goNextReviewPage(page)
		if !nextSuccess {
			break
		}

		pageNumber++
	}

	page.WaitForTimeout(3000)

	mu.Lock()
	defer mu.Unlock()

	if len(results) == 0 {
		return models.ShopeeScraperResult{}, fmt.Errorf("review tidak ditemukan atau response get_ratings tidak muncul")
	}

	if len(results) > limit {
		results = results[:limit]
	}

	return models.ShopeeScraperResult{
		Success:      true,
		TotalReviews: len(results),
		Reviews:      results,
	}, nil
}

func goNextReviewPage(page playwright.Page) bool {
	nextBtn := page.Locator(".shopee-icon-button--right").First()

	count, err := nextBtn.Count()
	if err != nil || count == 0 {
		return false
	}

	className, _ := nextBtn.GetAttribute("class")
	if strings.Contains(className, "disabled") {
		return false
	}

	_ = nextBtn.ScrollIntoViewIfNeeded()
	page.WaitForTimeout(500)

	err = nextBtn.Click(playwright.LocatorClickOptions{
		Timeout: playwright.Float(7000),
	})
	if err != nil {
		return false
	}

	page.WaitForTimeout(6000)
	return true
}

func (s ShopeeScraperService) getShopName(page playwright.Page, shopID int64) (string, error) {
	url := "https://shopee.co.id/api/v4/shop/get_shop_base?shopid=" + strconv.FormatInt(shopID, 10)

	response, err := page.Request().Get(url)
	if err != nil {
		return "", err
	}

	body, err := response.Body()
	if err != nil {
		return "", err
	}

	var shopResponse shopeeShopResponse
	if err := json.Unmarshal(body, &shopResponse); err != nil {
		return "", err
	}

	return shopResponse.Data.Name, nil
}
