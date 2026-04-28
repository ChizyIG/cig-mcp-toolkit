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

package edgar

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const tickersFixture = `{
  "0": {"cik_str": 320193, "ticker": "AAPL", "title": "Apple Inc."},
  "1": {"cik_str": 789019, "ticker": "MSFT", "title": "Microsoft Corp."}
}`

const aaplSubmissionsFixture = `{
  "name": "Apple Inc.",
  "sic": "3571",
  "sicDescription": "Electronic Computers",
  "exchanges": ["Nasdaq"],
  "tickers": ["AAPL"],
  "filings": {
    "recent": {
      "accessionNumber": ["0000320193-24-000123", "0000320193-24-000110", "0000320193-24-000099"],
      "filingDate":      ["2024-11-01",            "2024-08-02",            "2024-05-03"],
      "form":            ["10-K",                  "10-Q",                  "10-Q"],
      "primaryDocument": ["aapl-20240928.htm",     "aapl-q3.htm",           "aapl-q2.htm"]
    }
  }
}`

func newTestClient(t *testing.T) *Client {
	t.Helper()
	mux := http.NewServeMux()
	mux.HandleFunc("/tickers.json", func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("User-Agent"); got == "" {
			t.Errorf("expected User-Agent header, got empty")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(tickersFixture))
	})
	mux.HandleFunc("/submissions/CIK0000320193.json", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(aaplSubmissionsFixture))
	})
	mux.HandleFunc("/submissions/CIK0000999999.json", func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "not found", http.StatusNotFound)
	})

	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	c := New()
	c.TickersURL = srv.URL + "/tickers.json"
	c.SubmissionsURL = srv.URL + "/submissions"
	return c
}

func TestLookupCompany(t *testing.T) {
	c := newTestClient(t)
	co, err := c.LookupCompany(context.Background(), "aapl")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if co.CIK != "0000320193" {
		t.Errorf("CIK = %q, want 0000320193", co.CIK)
	}
	if co.Name != "Apple Inc." {
		t.Errorf("Name = %q, want Apple Inc.", co.Name)
	}
	if co.SIC != "3571" {
		t.Errorf("SIC = %q, want 3571", co.SIC)
	}
	if len(co.Exchanges) != 1 || co.Exchanges[0] != "Nasdaq" {
		t.Errorf("Exchanges = %v, want [Nasdaq]", co.Exchanges)
	}
}

func TestLookupCompanyUnknownTicker(t *testing.T) {
	c := newTestClient(t)
	_, err := c.LookupCompany(context.Background(), "ZZZZ")
	if err == nil {
		t.Fatal("expected error for unknown ticker")
	}
	if !strings.Contains(err.Error(), "unknown ticker") {
		t.Errorf("err = %v, want 'unknown ticker' in message", err)
	}
}

func TestLookupCompanyEmptyTicker(t *testing.T) {
	c := newTestClient(t)
	for _, in := range []string{"", "   ", "\t"} {
		if _, err := c.LookupCompany(context.Background(), in); err == nil {
			t.Errorf("LookupCompany(%q) expected error", in)
		}
	}
}

func TestListFilings(t *testing.T) {
	c := newTestClient(t)
	filings, err := c.ListFilings(context.Background(), "AAPL", "", 0)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(filings) != 3 {
		t.Fatalf("got %d filings, want 3", len(filings))
	}

	first := filings[0]
	if first.Form != "10-K" || first.FilingDate != "2024-11-01" {
		t.Errorf("first filing = %+v, unexpected", first)
	}
	wantURL := "https://www.sec.gov/Archives/edgar/data/320193/" +
		"000032019324000123/aapl-20240928.htm"
	if first.URL != wantURL {
		t.Errorf("URL = %q, want %q", first.URL, wantURL)
	}
}

func TestListFilingsFiltersByForm(t *testing.T) {
	c := newTestClient(t)
	filings, err := c.ListFilings(context.Background(), "AAPL", "10-Q", 0)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(filings) != 2 {
		t.Fatalf("got %d filings, want 2", len(filings))
	}
	for _, f := range filings {
		if f.Form != "10-Q" {
			t.Errorf("filing form = %q, want 10-Q", f.Form)
		}
	}
}

func TestListFilingsRespectsLimit(t *testing.T) {
	c := newTestClient(t)
	filings, err := c.ListFilings(context.Background(), "AAPL", "", 1)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(filings) != 1 {
		t.Fatalf("got %d filings, want 1", len(filings))
	}
}

func TestListFilingsFormFilterIsCaseInsensitive(t *testing.T) {
	c := newTestClient(t)
	filings, err := c.ListFilings(context.Background(), "AAPL", "10-q", 0)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(filings) != 2 {
		t.Errorf("got %d filings, want 2", len(filings))
	}
}

func TestSubmissionsHTTPErrorIsSurfaced(t *testing.T) {
	c := newTestClient(t)
	// Inject a ticker that maps to CIK 999999 which our test server 404s.
	c.tickers = map[string]int{"FAKE": 999999}
	c.once.Do(func() {})

	_, err := c.LookupCompany(context.Background(), "FAKE")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "status 404") {
		t.Errorf("err = %v, want 'status 404' in message", err)
	}
}
