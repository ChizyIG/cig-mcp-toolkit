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
	"os"

	"github.com/ChizyIG/cig-mcp-toolkit/go/cigmcp"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// LookupCompanyInput is the JSON input schema for the lookup_company tool.
type LookupCompanyInput struct {
	Ticker string `json:"ticker" jsonschema:"stock ticker symbol such as AAPL or MSFT"`
}

// ListFilingsInput is the JSON input schema for the list_filings tool.
type ListFilingsInput struct {
	Ticker string `json:"ticker"          jsonschema:"stock ticker symbol"`
	Form   string `json:"form,omitempty"  jsonschema:"optional form-type filter such as 10-K, 10-Q, 8-K"`
	Limit  int    `json:"limit,omitempty" jsonschema:"max number of filings to return (default 10, max 100)"`
}

// ListFilingsOutput wraps a slice for cleaner JSON serialization on the wire.
type ListFilingsOutput struct {
	Filings []Filing `json:"filings"`
}

// Run starts the EDGAR MCP server over stdio. It blocks until ctx is canceled
// or the underlying transport returns.
//
// The User-Agent header sent to SEC can be overridden with the EDGAR_USER_AGENT
// environment variable, which SEC requires to identify automated traffic.
func Run(ctx context.Context) error {
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "cig-mcp-edgar",
		Version: cigmcp.Version,
	}, nil)

	client := New()
	if ua := os.Getenv("EDGAR_USER_AGENT"); ua != "" {
		client.UserAgent = ua
	}

	mcp.AddTool(server, &mcp.Tool{
		Name: "lookup_company",
		Description: "Look up an SEC-registered company by ticker symbol. " +
			"Returns CIK, legal name, SIC industry code, and listed exchanges.",
	}, func(
		ctx context.Context, _ *mcp.CallToolRequest, in LookupCompanyInput,
	) (*mcp.CallToolResult, Company, error) {
		co, err := client.LookupCompany(ctx, in.Ticker)
		if err != nil {
			return nil, Company{}, err
		}
		return nil, *co, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name: "list_filings",
		Description: "List recent SEC filings for a ticker. Optionally filter by " +
			"form type (10-K, 10-Q, 8-K, etc.). Returns form, filing date, " +
			"accession number, primary document, and a URL to the filing.",
	}, func(
		ctx context.Context, _ *mcp.CallToolRequest, in ListFilingsInput,
	) (*mcp.CallToolResult, ListFilingsOutput, error) {
		filings, err := client.ListFilings(ctx, in.Ticker, in.Form, in.Limit)
		if err != nil {
			return nil, ListFilingsOutput{}, err
		}
		return nil, ListFilingsOutput{Filings: filings}, nil
	})

	return server.Run(ctx, &mcp.StdioTransport{})
}
