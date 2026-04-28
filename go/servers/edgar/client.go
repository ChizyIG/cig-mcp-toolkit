// Copyright 2026 CIG Engineering
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package edgar wraps the public SEC EDGAR endpoints used by the
// cig-mcp-edgar MCP server.
package edgar

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

// DefaultUserAgent is sent on every request unless overridden. SEC requires a
// real contact in the User-Agent header — see https://www.sec.gov/os/accessing-edgar-data.
const DefaultUserAgent = "cig-mcp-toolkit/0.1.0 (legal@chizyig.com)"

const (
	defaultTickersURL     = "https://www.sec.gov/files/company_tickers.json"
	defaultSubmissionsURL = "https://data.sec.gov/submissions"
	defaultArchiveURL     = "https://www.sec.gov/Archives/edgar/data"
)

// Client talks to EDGAR's public JSON endpoints.
//
// The zero value is not usable — call New(). Override TickersURL,
// SubmissionsURL, and ArchiveBaseURL in tests to point at an httptest.Server.
type Client struct {
	HTTP           *http.Client
	UserAgent      string
	TickersURL     string
	SubmissionsURL string
	ArchiveBaseURL string

	mu      sync.RWMutex
	tickers map[string]int
}

// Company is the metadata returned by LookupCompany.
type Company struct {
	CIK            string   `json:"cik"`
	Name           string   `json:"name"`
	SIC            string   `json:"sic"`
	SICDescription string   `json:"sic_description"`
	Exchanges      []string `json:"exchanges"`
	Tickers        []string `json:"tickers"`
}

// Filing is one row from the recent-filings array on the submissions endpoint.
//
// URL is empty when the filing has no primary document (e.g., header-only
// submissions and some 8-K amendments). Callers that need to browse such
// filings can build the directory URL themselves from CIK + AccessionNumber.
type Filing struct {
	Form            string `json:"form"`
	FilingDate      string `json:"filing_date"`
	AccessionNumber string `json:"accession_number"`
	PrimaryDocument string `json:"primary_document"`
	URL             string `json:"url"`
}

// New returns a Client with sensible defaults for production use.
func New() *Client {
	return &Client{
		HTTP:           &http.Client{Timeout: 15 * time.Second},
		UserAgent:      DefaultUserAgent,
		TickersURL:     defaultTickersURL,
		SubmissionsURL: defaultSubmissionsURL,
		ArchiveBaseURL: defaultArchiveURL,
	}
}

func (c *Client) loadTickers(ctx context.Context) error {
	c.mu.RLock()
	cached := c.tickers
	c.mu.RUnlock()
	if cached != nil {
		return nil
	}

	// Fetch outside the lock so a slow upstream doesn't block readers, and so
	// a transient failure isn't permanently cached (the next call retries).
	m, err := c.fetchTickers(ctx)
	if err != nil {
		return err
	}

	c.mu.Lock()
	if c.tickers == nil {
		c.tickers = m
	}
	c.mu.Unlock()
	return nil
}

func (c *Client) fetchTickers(ctx context.Context) (map[string]int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.TickersURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", c.UserAgent)
	req.Header.Set("Accept", "application/json")

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch tickers: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch tickers: status %d", resp.StatusCode)
	}

	var raw map[string]struct {
		CIKStr int    `json:"cik_str"`
		Ticker string `json:"ticker"`
		Title  string `json:"title"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("decode tickers: %w", err)
	}

	m := make(map[string]int, len(raw))
	for _, entry := range raw {
		m[strings.ToUpper(entry.Ticker)] = entry.CIKStr
	}
	return m, nil
}

func (c *Client) cikForTicker(ctx context.Context, ticker string) (int, error) {
	cleaned := strings.ToUpper(strings.TrimSpace(ticker))
	if cleaned == "" {
		return 0, errors.New("ticker must be non-empty")
	}
	if err := c.loadTickers(ctx); err != nil {
		return 0, err
	}
	c.mu.RLock()
	cik, ok := c.tickers[cleaned]
	c.mu.RUnlock()
	if !ok {
		return 0, fmt.Errorf("unknown ticker %q", cleaned)
	}
	return cik, nil
}

type submissionsResponse struct {
	Name           string   `json:"name"`
	SIC            string   `json:"sic"`
	SICDescription string   `json:"sicDescription"`
	Exchanges      []string `json:"exchanges"`
	Tickers        []string `json:"tickers"`
	Filings        struct {
		Recent struct {
			AccessionNumber []string `json:"accessionNumber"`
			FilingDate      []string `json:"filingDate"`
			Form            []string `json:"form"`
			PrimaryDocument []string `json:"primaryDocument"`
		} `json:"recent"`
	} `json:"filings"`
}

func (c *Client) fetchSubmissions(ctx context.Context, cik int) (*submissionsResponse, error) {
	url := fmt.Sprintf("%s/CIK%010d.json", strings.TrimRight(c.SubmissionsURL, "/"), cik)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", c.UserAgent)
	req.Header.Set("Accept", "application/json")

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch submissions: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch submissions: status %d", resp.StatusCode)
	}

	var data submissionsResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("decode submissions: %w", err)
	}
	return &data, nil
}

// LookupCompany resolves a ticker to its EDGAR metadata.
func (c *Client) LookupCompany(ctx context.Context, ticker string) (*Company, error) {
	cik, err := c.cikForTicker(ctx, ticker)
	if err != nil {
		return nil, err
	}
	sub, err := c.fetchSubmissions(ctx, cik)
	if err != nil {
		return nil, err
	}
	return &Company{
		CIK:            fmt.Sprintf("%010d", cik),
		Name:           sub.Name,
		SIC:            sub.SIC,
		SICDescription: sub.SICDescription,
		Exchanges:      sub.Exchanges,
		Tickers:        sub.Tickers,
	}, nil
}

// ListFilings returns the most recent filings for a ticker, optionally
// filtered by form type. limit defaults to 10 and is capped at 100.
func (c *Client) ListFilings(
	ctx context.Context, ticker, formFilter string, limit int,
) ([]Filing, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	cik, err := c.cikForTicker(ctx, ticker)
	if err != nil {
		return nil, err
	}
	sub, err := c.fetchSubmissions(ctx, cik)
	if err != nil {
		return nil, err
	}

	rec := sub.Filings.Recent
	n := len(rec.AccessionNumber)
	if n != len(rec.FilingDate) || n != len(rec.Form) || n != len(rec.PrimaryDocument) {
		return nil, errors.New("submissions response: filings arrays differ in length")
	}

	formFilter = strings.ToUpper(strings.TrimSpace(formFilter))
	out := make([]Filing, 0, limit)
	for i := 0; i < n && len(out) < limit; i++ {
		if formFilter != "" && strings.ToUpper(rec.Form[i]) != formFilter {
			continue
		}
		accNoDash := strings.ReplaceAll(rec.AccessionNumber[i], "-", "")
		var url string
		if doc := rec.PrimaryDocument[i]; doc != "" {
			url = fmt.Sprintf(
				"%s/%d/%s/%s",
				strings.TrimRight(c.ArchiveBaseURL, "/"),
				cik, accNoDash, doc,
			)
		}
		out = append(out, Filing{
			Form:            rec.Form[i],
			FilingDate:      rec.FilingDate[i],
			AccessionNumber: rec.AccessionNumber[i],
			PrimaryDocument: rec.PrimaryDocument[i],
			URL:             url,
		})
	}
	return out, nil
}
