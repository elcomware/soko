package commerce_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"elcomware/app/facades"
	"elcomware/app/models"
	"github.com/stretchr/testify/require"
)

func TestConcurrentCartAddsPreserveEveryUpdate(t *testing.T) {
	p := models.Product{ID: "serialized-cart-product", Title: "Concurrent cart", Type: models.TypePhysical, Format: models.FormatRetail, Price: 20, MOQ: 1, Stock: 100}
	require.NoError(t, facades.Orm().Query().Create(&p))
	server := httptest.NewServer(facades.Route())
	defer server.Close()
	jar, err := cookiejar.New(nil)
	require.NoError(t, err)
	client := &http.Client{Jar: jar, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
	response, err := client.Get(server.URL + "/cart")
	require.NoError(t, err)
	io.Copy(io.Discard, response.Body)
	response.Body.Close()
	const count = 24
	start := make(chan struct{})
	failures := make(chan error, count)
	var wg sync.WaitGroup
	for i := 0; i < count; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			req, _ := http.NewRequest("POST", server.URL+"/cart/items", strings.NewReader(`{"product_id":"serialized-cart-product","quantity":1}`))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Sec-Fetch-Site", "same-origin")
			resp, e := client.Do(req)
			if e != nil {
				failures <- e
				return
			}
			defer resp.Body.Close()
			io.Copy(io.Discard, resp.Body)
			if resp.StatusCode != 303 {
				failures <- fmt.Errorf("unexpected status %d", resp.StatusCode)
			}
		}()
	}
	close(start)
	wg.Wait()
	close(failures)
	for err := range failures {
		require.NoError(t, err)
	}
	req, _ := http.NewRequest("GET", server.URL+"/cart", nil)
	req.Header.Set("X-Inertia", "true")
	response, err = client.Do(req)
	require.NoError(t, err)
	defer response.Body.Close()
	var page struct {
		Props struct {
			Cart struct {
				Items []struct {
					Quantity int `json:"quantity"`
				} `json:"items"`
			} `json:"cart"`
		} `json:"props"`
	}
	require.NoError(t, json.NewDecoder(response.Body).Decode(&page))
	require.Len(t, page.Props.Cart.Items, 1)
	require.Equal(t, count, page.Props.Cart.Items[0].Quantity, "every concurrent add must survive the session read/write lifecycle")
}
