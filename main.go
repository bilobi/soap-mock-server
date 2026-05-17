package main

import (
	"bytes"
	"crypto/rand"
	"crypto/sha1"
	"crypto/subtle"
	"encoding/base64"
	"encoding/binary"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

// ============================================================================
// SOAP Mock Server - Multiple Auth Patterns for Testing
// ============================================================================
// Endpoints:
//   /session/service.wsdl    - SOAP Session Auth (DIVA-like)
//   /wsse/service.wsdl       - WS-Security UsernameToken
//   /basic/service.wsdl      - Basic Auth over SOAP
//   /ntlm/service.wsdl       - NTLM Authentication (Windows Auth simulation)
//   /noauth/service.wsdl     - No Authentication
// ============================================================================

// Session store for SOAP Session auth
var (
	sessions     = make(map[string]*Session)
	sessionMutex sync.RWMutex
)

// SAP session store — keyed by SAP_SESSIONID cookie value
var (
	sapSessions     = make(map[string]*SAPSession)
	sapSessionMutex sync.RWMutex
)

type SAPSession struct {
	SessionID string
	CSRFToken string
	SAPClient string
	Username  string
	CreatedAt time.Time
	ExpiresAt time.Time
}

type NTLMChallenge struct {
	Challenge string
	CreatedAt time.Time
	ClientIP  string
}

type Session struct {
	Token     string
	Username  string
	CreatedAt time.Time
	ExpiresAt time.Time
}

// Test credentials
const (
	TestUsername      = "testuser"
	TestPassword      = "testpass123"
	SessionCookieName = "session"
	BasicAuthRealm    = "IMAXIS Mock Server"
	Soap11Namespace   = "http://schemas.xmlsoap.org/soap/envelope/"
	Soap12Namespace   = "http://www.w3.org/2003/05/soap-envelope"
)

type SOAPMessageInfo struct {
	Version         string
	Operation       string
	OperationNS     string
	BodyValues      map[string]string
	SessionToken    string
	WSSECredentials *WSSECredentials
}

type WSSECredentials struct {
	Username     string
	Password     string
	PasswordType string
	Nonce        string
	Created      string
}

func main() {
	stores, err := initStores()
	if err != nil {
		log.Fatalf("failed to load JSON stores: %v", err)
	}
	endpointStores = stores

	mux := http.NewServeMux()

	// WSDL endpoints
	mux.HandleFunc("/session/service.wsdl", handleSessionWSDL)
	mux.HandleFunc("/wsse/service.wsdl", handleWSSEWSDL)
	mux.HandleFunc("/basic/service.wsdl", handleBasicWSDL)
	mux.HandleFunc("/ntlm/service.wsdl", handleNTLMWSDL)
	mux.HandleFunc("/noauth/service.wsdl", handleNoAuthWSDL)
	mux.HandleFunc("/sap/service.wsdl", handleSAPWSDL)

	// SOAP endpoints
	mux.HandleFunc("/session/soap", handleSessionSOAP)
	mux.HandleFunc("/wsse/soap", handleWSSESOAP)
	mux.HandleFunc("/basic/soap", handleBasicSOAP)
	mux.HandleFunc("/ntlm/soap", handleNTLMSOAP)
	mux.HandleFunc("/noauth/soap", handleNoAuthSOAP)
	mux.HandleFunc("/sap/soap", handleSAPSOAP)

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// Info page
	mux.HandleFunc("/", handleIndex)

	port := os.Getenv("SOAP_MOCK_PORT")
	if port == "" {
		port = "8099"
	}
	addr = ":" + port

	log.Printf("🚀 SOAP Mock Server starting on http://localhost%s", addr)
	log.Printf("📋 Available WSDLs:")
	log.Printf("   - http://localhost%s/session/service.wsdl (SOAP Session Auth)", addr)
	log.Printf("   - http://localhost%s/wsse/service.wsdl (WS-Security UsernameToken + WS-Policy)", addr)
	log.Printf("   - http://localhost%s/basic/service.wsdl (Basic Auth + WS-Policy)", addr)
	log.Printf("   - http://localhost%s/ntlm/service.wsdl (NTLM/Windows Auth + WS-Policy)", addr)
	log.Printf("   - http://localhost%s/noauth/service.wsdl (No Auth)", addr)
	log.Printf("   - http://localhost%s/sap/service.wsdl (SAP Session Auth)", addr)
	log.Printf("🔑 Test credentials: %s / %s", TestUsername, TestPassword)
	log.Printf("🔑 NTLM credentials: TESTDOMAIN\\%s / %s", TestUsername, TestPassword)
	log.Printf("🔑 SAP credentials: %s / %s (sap-client: 100)", TestUsername, TestPassword)
	log.Printf("────────────────────────────────────────────────────────────────────")

	// Wrap with logging middleware
	handler := loggingMiddleware(mux)

	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatal(err)
	}
}

// ============================================================================
// Logging Middleware
// ============================================================================

// responseWriter wrapper to capture status code and response body
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	body       bytes.Buffer
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
	lrw.body.Write(b) // Capture response body
	return lrw.ResponseWriter.Write(b)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Skip logging for health checks and index page
		if r.URL.Path == "/health" || r.URL.Path == "/" {
			next.ServeHTTP(w, r)
			return
		}

		// Read request body
		var requestBody []byte
		if r.Body != nil {
			requestBody, _ = io.ReadAll(r.Body)
			r.Body = io.NopCloser(bytes.NewBuffer(requestBody)) // Restore body for handler
		}

		// Log request
		log.Printf("════════════════════════════════════════════════════════════════════")
		log.Printf("📥 REQUEST: %s %s", r.Method, r.URL.Path)
		log.Printf("   From: %s", r.RemoteAddr)

		// Log important headers
		logHeaders := []string{"Content-Type", "SOAPAction", "Authorization", "Cookie", "X-Session-Token"}
		for _, h := range logHeaders {
			if v := r.Header.Get(h); v != "" {
				// Mask sensitive data
				if h == "Authorization" && len(v) > 20 {
					v = v[:20] + "..."
				}
				log.Printf("   %s: %s", h, v)
			}
		}

		// Log request body (truncated for readability)
		if len(requestBody) > 0 {
			bodyStr := string(requestBody)
			// Pretty print or truncate
			if len(bodyStr) > 2000 {
				log.Printf("   Body (%d bytes): %s...[truncated]", len(bodyStr), bodyStr[:2000])
			} else {
				log.Printf("   Body (%d bytes):\n%s", len(bodyStr), indentXML(bodyStr))
			}
		}

		// Wrap response writer to capture response
		lrw := &loggingResponseWriter{
			ResponseWriter: w,
			statusCode:     200, // Default status
		}

		// Call the actual handler
		next.ServeHTTP(lrw, r)

		duration := time.Since(start)

		// Log response
		log.Printf("📤 RESPONSE: %d %s (took %v)", lrw.statusCode, http.StatusText(lrw.statusCode), duration)

		// Log response body (truncated)
		responseBody := lrw.body.String()
		if len(responseBody) > 0 {
			if len(responseBody) > 2000 {
				log.Printf("   Body (%d bytes): %s...[truncated]", len(responseBody), responseBody[:2000])
			} else {
				log.Printf("   Body (%d bytes):\n%s", len(responseBody), indentXML(responseBody))
			}
		}
		log.Printf("────────────────────────────────────────────────────────────────────")
	})
}

// indentXML adds simple indentation for readability
func indentXML(xmlStr string) string {
	// Simple indent for log readability
	lines := strings.Split(xmlStr, "\n")
	var result []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			result = append(result, "      "+trimmed)
		}
	}
	return strings.Join(result, "\n")
}

// ============================================================================
// Index Page
// ============================================================================

func handleIndex(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>SOAP Mock Server</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; max-width: 1120px; margin: 40px auto; padding: 20px; color: #202124; line-height: 1.5; }
        h1 { color: #333; margin-bottom: 6px; }
        h2 { margin-top: 34px; padding-bottom: 6px; border-bottom: 1px solid #ddd; }
        h3 { margin-bottom: 8px; }
        p { margin: 8px 0; }
        .endpoint { background: #f5f5f5; padding: 15px; margin: 10px 0; border-radius: 8px; }
        .endpoint h3 { margin-top: 0; color: #0066cc; }
        code { background: #e0e0e0; padding: 2px 6px; border-radius: 4px; }
        pre { background: #1f2933; color: #f8fafc; padding: 14px; border-radius: 8px; overflow-x: auto; line-height: 1.4; }
        pre code { background: transparent; color: inherit; padding: 0; }
        .auth-type { display: inline-block; padding: 4px 8px; border-radius: 4px; font-size: 12px; font-weight: bold; }
        .session { background: #4CAF50; color: white; }
        .wsse { background: #2196F3; color: white; }
        .basic { background: #FF9800; color: white; }
        .noauth { background: #9E9E9E; color: white; }
        .note { background: #fff8e1; border-left: 4px solid #f9ab00; padding: 10px 12px; border-radius: 6px; }
        .grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(260px, 1fr)); gap: 12px; }
        .section { background: #fafafa; border: 1px solid #e0e0e0; border-radius: 8px; padding: 14px; margin: 12px 0; }
        .section h3 { margin-top: 0; color: #174ea6; }
        .muted { color: #5f6368; }
        .pill { display: inline-block; background: #e8f0fe; color: #174ea6; padding: 2px 8px; border-radius: 999px; margin: 2px; font-size: 12px; }
        .method { font-weight: 700; color: #137333; }
        .optional { color: #5f6368; font-size: 12px; }
        details { border: 1px solid #ddd; border-radius: 8px; padding: 10px 12px; margin: 10px 0; background: #fff; }
        summary { cursor: pointer; font-weight: 700; color: #174ea6; }
        table { width: 100%; border-collapse: collapse; margin: 10px 0; }
        td, th { padding: 8px; text-align: left; border-bottom: 1px solid #ddd; }
        th { background: #f1f3f4; }
        ul { margin-top: 8px; }
    </style>
</head>
<body>
    <h1>🧪 SOAP Mock Server</h1>
    <p>Test different authentication patterns for AuthProbe validation.</p>
    
    <h2>Credentials</h2>
    <table>
        <tr><th>Username</th><td><code>testuser</code></td></tr>
        <tr><th>Password</th><td><code>testpass123</code></td></tr>
    </table>

    <h2>Data Stores</h2>
    <p>Each endpoint persists its own JSON in <code>service/</code> (session.json, wsse.json, basic.json, ntlm.json, noauth.json, sap.json). Override with <code>SOAP_MOCK_DATA_DIR</code>.</p>

    <h2>Endpoints</h2>
    
    <div class="endpoint">
        <h3><span class="auth-type session">SOAP SESSION</span> DIVA-like Session Authentication</h3>
        <p>WSDL: <code>/session/service.wsdl</code></p>
        <p>Operations: <code>Login</code>, <code>Logout</code>, ERP CRUD (Customers, Stocks, Cash Accounts, Orders, Invoices, Cash Transactions)</p>
        <p>Flow: Login sets <code>session</code> cookie and returns token; send via Cookie, X-Session-Token, or SOAP SessionHeader.</p>
        <p>SOAP: 1.1 and 1.2 supported</p>
    </div>

    <div class="endpoint">
        <h3><span class="auth-type wsse">WS-SECURITY</span> UsernameToken Authentication</h3>
        <p>WSDL: <code>/wsse/service.wsdl</code></p>
        <p>Operations: ERP CRUD (Customers, Stocks, Cash Accounts, Orders, Invoices, Cash Transactions)</p>
        <p>Auth: WS-Security UsernameToken in SOAP Header (PasswordText/Digest)</p>
    </div>

    <div class="endpoint">
        <h3><span class="auth-type basic">BASIC AUTH</span> HTTP Basic Authentication</h3>
        <p>WSDL: <code>/basic/service.wsdl</code></p>
        <p>Operations: ERP CRUD (Customers, Stocks, Cash Accounts, Orders, Invoices, Cash Transactions)</p>
        <p>Auth: HTTP Authorization header</p>
    </div>

    <div class="endpoint">
        <h3><span class="auth-type ntlm" style="background: #9C27B0;">NTLM</span> Windows Authentication</h3>
        <p>WSDL: <code>/ntlm/service.wsdl</code></p>
        <p>Operations: ERP CRUD (Customers, Stocks, Cash Accounts, Orders, Invoices, Cash Transactions)</p>
        <p>Auth: NTLM challenge-response (3-way handshake)</p>
        <p>Credentials: <code>TESTDOMAIN\testuser</code> / <code>testpass123</code></p>
        <p><small>Supports Type1 (Negotiate) → Type2 (Challenge) → Type3 (Authenticate) flow</small></p>
    </div>

    <div class="endpoint">
        <h3><span class="auth-type noauth">NO AUTH</span> Public Service</h3>
        <p>WSDL: <code>/noauth/service.wsdl</code></p>
        <p>Operations: <code>GetCountries</code>, <code>GetCities</code>, ERP CRUD (no auth)</p>
        <p>Auth: None required</p>
    </div>

    <div class="endpoint">
        <h3><span class="auth-type session" style="background: #E91E63;">SAP SESSION</span> SAP Session Authentication</h3>
        <p>WSDL: <code>/sap/service.wsdl</code></p>
        <p>Operations: ERP CRUD (Customers, Stocks, Cash Accounts, Orders, Invoices, Cash Transactions)</p>
        <p>Auth: HTTP Basic Auth for first request + <code>sap-client</code> header</p>
        <p>Flow: First request creates session → <code>SAP_SESSIONID_MKS_100</code> cookie + <code>x-csrf-token</code> returned. Subsequent requests use cookie.</p>
        <p>Credentials: <code>testuser</code> / <code>testpass123</code> (sap-client: <code>100</code>)</p>
    </div>

    <h2>Service Usage</h2>
    <div class="section">
        <h3>SOAP Request Rules</h3>
        <table>
            <tr><th>Item</th><th>Value</th></tr>
            <tr><td>SOAP method</td><td><span class="method">POST</span> to <code>/{auth}/soap</code></td></tr>
            <tr><td>WSDL method</td><td><span class="method">GET</span> <code>/{auth}/service.wsdl</code></td></tr>
            <tr><td>Content-Type</td><td><code>text/xml</code> for SOAP 1.1 or <code>application/soap+xml</code> for SOAP 1.2</td></tr>
            <tr><td>SOAPAction</td><td><code>http://mock.imaxis.com/{endpoint}/{Operation}</code></td></tr>
            <tr><td>Operation detection</td><td>Operation can be read from the SOAP body or from <code>SOAPAction</code>.</td></tr>
        </table>
        <p class="note">Use namespace <code>http://mock.imaxis.com/session</code>, <code>/wsse</code>, <code>/basic</code>, <code>/ntlm</code>, <code>/sap</code>, or <code>/public</code> according to the endpoint. For no-auth calls, the namespace is <code>http://mock.imaxis.com/public</code>.</p>
    </div>

    <div class="section">
        <h3>Endpoint Map</h3>
        <table>
            <tr><th>Auth mode</th><th>WSDL</th><th>SOAP</th><th>Namespace</th></tr>
            <tr><td>Session</td><td><code>/session/service.wsdl</code></td><td><code>/session/soap</code></td><td><code>http://mock.imaxis.com/session</code></td></tr>
            <tr><td>WS-Security</td><td><code>/wsse/service.wsdl</code></td><td><code>/wsse/soap</code></td><td><code>http://mock.imaxis.com/wsse</code></td></tr>
            <tr><td>Basic</td><td><code>/basic/service.wsdl</code></td><td><code>/basic/soap</code></td><td><code>http://mock.imaxis.com/basic</code></td></tr>
            <tr><td>NTLM</td><td><code>/ntlm/service.wsdl</code></td><td><code>/ntlm/soap</code></td><td><code>http://mock.imaxis.com/ntlm</code></td></tr>
            <tr><td>SAP</td><td><code>/sap/service.wsdl</code></td><td><code>/sap/soap</code></td><td><code>http://mock.imaxis.com/sap</code></td></tr>
            <tr><td>No Auth</td><td><code>/noauth/service.wsdl</code></td><td><code>/noauth/soap</code></td><td><code>http://mock.imaxis.com/public</code></td></tr>
        </table>
    </div>

    <h2>ERP Operation Catalog</h2>
    <div class="grid">
        <div class="section">
            <h3>Customers</h3>
            <p><span class="pill">GetCustomers</span><span class="pill">ListCustomers</span><span class="pill">GetCustomer</span><span class="pill">CreateCustomer</span><span class="pill">UpdateCustomer</span></p>
            <p><strong>List filters:</strong> <code>code</code>, <code>status</code>, <code>taxNumber</code>, <code>maxResults</code>, <code>offset</code></p>
            <p><strong>Lookup:</strong> <code>id</code> or <code>code</code></p>
            <p><strong>Payload:</strong> <code>customer</code> with code, name, taxNumber, email, phone, address, currency, riskLimit, status</p>
        </div>
        <div class="section">
            <h3>Stocks</h3>
            <p><span class="pill">GetStocks</span><span class="pill">ListStocks</span><span class="pill">GetStock</span><span class="pill">CreateStock</span><span class="pill">UpdateStock</span></p>
            <p><strong>List filters:</strong> <code>code</code>, <code>status</code>, <code>nameContains</code>, <code>maxResults</code>, <code>offset</code></p>
            <p><strong>Lookup:</strong> <code>id</code> or <code>code</code></p>
            <p><strong>Payload:</strong> <code>stock</code> with code, name, unit, vatRate, price, stockOnHand, minStock, status</p>
        </div>
        <div class="section">
            <h3>Cash Accounts</h3>
            <p><span class="pill">GetCashAccounts</span><span class="pill">ListCashAccounts</span><span class="pill">GetCashAccount</span><span class="pill">CreateCashAccount</span><span class="pill">UpdateCashAccount</span></p>
            <p><strong>List filters:</strong> <code>status</code>, <code>currency</code>, <code>maxResults</code>, <code>offset</code></p>
            <p><strong>Lookup:</strong> <code>id</code> or <code>code</code></p>
            <p><strong>Payload:</strong> <code>cashAccount</code> with code, name, currency, balance, status</p>
        </div>
        <div class="section">
            <h3>Orders</h3>
            <p><span class="pill">GetOrders</span><span class="pill">ListOrders</span><span class="pill">GetOrder</span><span class="pill">CreateOrder</span><span class="pill">UpdateOrder</span><span class="pill">ApproveOrder</span><span class="pill">CancelOrder</span></p>
            <p><strong>List filters:</strong> <code>customerId</code>, <code>status</code>, <code>fromDate</code>, <code>toDate</code>, <code>maxResults</code>, <code>offset</code></p>
            <p><strong>Lookup/action:</strong> <code>id</code> or <code>orderNo</code></p>
            <p><strong>Payload:</strong> <code>order</code> with orderNo, customerId/customerCode, status, orderDate, currency, notes, lines</p>
        </div>
        <div class="section">
            <h3>Invoices</h3>
            <p><span class="pill">GetInvoices</span><span class="pill">ListInvoices</span><span class="pill">GetInvoice</span><span class="pill">CreateInvoice</span><span class="pill">CreateInvoiceFromOrder</span><span class="pill">CancelInvoice</span></p>
            <p><strong>List filters:</strong> <code>customerId</code>, <code>status</code>, <code>fromDate</code>, <code>toDate</code>, <code>maxResults</code>, <code>offset</code></p>
            <p><strong>Lookup/action:</strong> <code>id</code> or <code>invoiceNo</code></p>
            <p><strong>Payload:</strong> <code>invoice</code> with invoiceNo, orderId/orderNo, customerId/customerCode, invoiceDate, dueDate, currency, lines</p>
        </div>
        <div class="section">
            <h3>Cash Transactions</h3>
            <p><span class="pill">GetCashTransactions</span><span class="pill">ListCashTransactions</span><span class="pill">GetCashTransaction</span><span class="pill">CreateCashTransaction</span><span class="pill">ReverseCashTransaction</span></p>
            <p><strong>List filters:</strong> <code>cashAccountId</code>, <code>type</code>, <code>fromDate</code>, <code>toDate</code>, <code>maxResults</code>, <code>offset</code></p>
            <p><strong>Lookup/action:</strong> <code>id</code> or <code>txnNo</code></p>
            <p><strong>Payload:</strong> <code>cashTransaction</code> with txnNo, cashAccountId/cashAccountCode, customerId/customerCode, invoiceId/invoiceNo, type, amount, currency, method, transactionDate, description</p>
        </div>
    </div>

    <div class="section">
        <h3>Public Operations</h3>
        <p><span class="pill">GetCountries</span><span class="pill">GetCities</span></p>
        <p><strong>GetCountries:</strong> optional <code>continent</code> filter.</p>
        <p><strong>GetCities:</strong> optional <code>countryCode</code> filter.</p>
    </div>

    <details>
        <summary>Common payload shapes</summary>
        <table>
            <tr><th>Payload</th><th>Fields</th></tr>
            <tr><td><code>customer</code></td><td><code>id</code> <span class="optional">(update)</span>, <code>code</code>, <code>name</code>, <code>taxNumber</code>, <code>email</code>, <code>phone</code>, <code>address</code>, <code>currency</code>, <code>riskLimit</code>, <code>status</code></td></tr>
            <tr><td><code>stock</code></td><td><code>id</code> <span class="optional">(update)</span>, <code>code</code>, <code>name</code>, <code>unit</code>, <code>vatRate</code>, <code>price</code>, <code>stockOnHand</code>, <code>minStock</code>, <code>status</code></td></tr>
            <tr><td><code>cashAccount</code></td><td><code>id</code> <span class="optional">(update)</span>, <code>code</code>, <code>name</code>, <code>currency</code>, <code>balance</code>, <code>status</code></td></tr>
            <tr><td><code>order</code></td><td><code>id</code> <span class="optional">(update)</span>, <code>orderNo</code>, <code>customerId</code> or <code>customerCode</code>, <code>status</code>, <code>orderDate</code>, <code>currency</code>, <code>notes</code>, <code>lines&gt;line</code></td></tr>
            <tr><td><code>invoice</code></td><td><code>id</code> <span class="optional">(update)</span>, <code>invoiceNo</code>, <code>orderId</code> or <code>orderNo</code>, <code>customerId</code> or <code>customerCode</code>, <code>invoiceDate</code>, <code>dueDate</code>, <code>currency</code>, <code>lines&gt;line</code></td></tr>
            <tr><td><code>line</code></td><td><code>stockId</code> or <code>stockCode</code>, <code>quantity</code>, <code>unitPrice</code>, <code>discountRate</code>, <code>taxRate</code></td></tr>
            <tr><td><code>cashTransaction</code></td><td><code>id</code> <span class="optional">(reverse)</span>, <code>txnNo</code>, <code>cashAccountId</code> or <code>cashAccountCode</code>, <code>customerId</code> or <code>customerCode</code>, <code>invoiceId</code> or <code>invoiceNo</code>, <code>type</code>, <code>amount</code>, <code>currency</code>, <code>method</code>, <code>transactionDate</code>, <code>description</code></td></tr>
        </table>
        <p class="muted">Create operations can omit IDs. Update/action operations accept ID fields or the related business number/code where available.</p>
    </details>

    <div class="section">
        <h3>Common Values</h3>
        <table>
            <tr><th>Field</th><th>Accepted values / format</th></tr>
            <tr><td><code>currency</code></td><td><code>TRY</code>, <code>USD</code>, <code>EUR</code></td></tr>
            <tr><td><code>vatRate</code>, <code>taxRate</code></td><td><code>0</code>, <code>1</code>, <code>8</code>, <code>18</code></td></tr>
            <tr><td>Cash transaction type</td><td><code>COLLECTION</code>, <code>PAYMENT</code></td></tr>
            <tr><td>Cash transaction method</td><td><code>CASH</code>, <code>BANK</code>, <code>POS</code>, <code>TRANSFER</code></td></tr>
            <tr><td><code>orderDate</code>, <code>invoiceDate</code>, <code>dueDate</code>, <code>transactionDate</code></td><td><code>YYYY-MM-DD</code></td></tr>
            <tr><td>Customer/stock status</td><td><code>ACTIVE</code>, <code>INACTIVE</code></td></tr>
            <tr><td>Cash account status</td><td><code>OPEN</code>, <code>CLOSED</code></td></tr>
            <tr><td>Order status</td><td><code>DRAFT</code>, <code>APPROVED</code>, <code>CANCELLED</code>, <code>INVOICED</code></td></tr>
            <tr><td>Invoice status</td><td><code>ISSUED</code>, <code>PAID</code>, <code>CANCELLED</code></td></tr>
            <tr><td>Cash transaction status</td><td><code>ISSUED</code>, <code>REVERSED</code></td></tr>
        </table>
    </div>

    <h2>Example Requests</h2>
    <details open>
        <summary>No Auth - List Customers</summary>
        <pre><code>curl -X POST http://localhost:8099/noauth/soap \
  -H "Content-Type: text/xml" \
  -H "SOAPAction: http://mock.imaxis.com/public/ListCustomers" \
  -d '&lt;?xml version="1.0" encoding="UTF-8"?&gt;
&lt;soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/"&gt;
  &lt;soap:Body&gt;
    &lt;ListCustomers xmlns="http://mock.imaxis.com/public"&gt;
      &lt;status&gt;ACTIVE&lt;/status&gt;
      &lt;maxResults&gt;10&lt;/maxResults&gt;
    &lt;/ListCustomers&gt;
  &lt;/soap:Body&gt;
&lt;/soap:Envelope&gt;'</code></pre>
    </details>

    <details>
        <summary>Basic Auth - Create Order</summary>
        <pre><code>curl -X POST http://localhost:8099/basic/soap \
  -u testuser:testpass123 \
  -H "Content-Type: text/xml" \
  -H "SOAPAction: http://mock.imaxis.com/basic/CreateOrder" \
  -d '&lt;?xml version="1.0" encoding="UTF-8"?&gt;
&lt;soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/"&gt;
  &lt;soap:Body&gt;
    &lt;CreateOrder xmlns="http://mock.imaxis.com/basic"&gt;
      &lt;order&gt;
        &lt;orderNo&gt;SO-2024-TEST-01&lt;/orderNo&gt;
        &lt;customerCode&gt;CUST-0001&lt;/customerCode&gt;
        &lt;status&gt;DRAFT&lt;/status&gt;
        &lt;orderDate&gt;2024-12-01&lt;/orderDate&gt;
        &lt;currency&gt;TRY&lt;/currency&gt;
        &lt;lines&gt;
          &lt;line&gt;
            &lt;stockCode&gt;STK-0002&lt;/stockCode&gt;
            &lt;quantity&gt;3&lt;/quantity&gt;
            &lt;unitPrice&gt;350&lt;/unitPrice&gt;
            &lt;taxRate&gt;18&lt;/taxRate&gt;
          &lt;/line&gt;
        &lt;/lines&gt;
      &lt;/order&gt;
    &lt;/CreateOrder&gt;
  &lt;/soap:Body&gt;
&lt;/soap:Envelope&gt;'</code></pre>
    </details>

    <details>
        <summary>Session Auth - Login then Get Stock</summary>
        <pre><code>TOKEN=$(curl -s -X POST http://localhost:8099/session/soap \
  -H "Content-Type: text/xml" \
  -H "SOAPAction: http://mock.imaxis.com/session/Login" \
  -d '&lt;soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/"&gt;
  &lt;soap:Body&gt;
    &lt;Login xmlns="http://mock.imaxis.com/session"&gt;
      &lt;username&gt;testuser&lt;/username&gt;
      &lt;password&gt;testpass123&lt;/password&gt;
    &lt;/Login&gt;
  &lt;/soap:Body&gt;
&lt;/soap:Envelope&gt;' | tr -d '\n' | sed -n 's:.*&lt;sessionToken&gt;\([^&lt;]*\)&lt;/sessionToken&gt;.*:\1:p')

curl -X POST http://localhost:8099/session/soap \
  -H "X-Session-Token: $TOKEN" \
  -H "Content-Type: text/xml" \
  -H "SOAPAction: http://mock.imaxis.com/session/GetStock" \
  -d '&lt;soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/"&gt;
  &lt;soap:Body&gt;
    &lt;GetStock xmlns="http://mock.imaxis.com/session"&gt;
      &lt;code&gt;STK-0002&lt;/code&gt;
    &lt;/GetStock&gt;
  &lt;/soap:Body&gt;
&lt;/soap:Envelope&gt;'</code></pre>
    </details>

    <details>
        <summary>SAP Session - First request creates cookie and CSRF token</summary>
        <pre><code>curl -i -c sap-cookie.txt -X POST http://localhost:8099/sap/soap \
  -u testuser:testpass123 \
  -H "sap-client: 100" \
  -H "x-csrf-token: Fetch" \
  -H "Content-Type: text/xml" \
  -H "SOAPAction: http://mock.imaxis.com/sap/ListInvoices" \
  -d '&lt;soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/"&gt;
  &lt;soap:Body&gt;
    &lt;ListInvoices xmlns="http://mock.imaxis.com/sap"&gt;
      &lt;maxResults&gt;5&lt;/maxResults&gt;
    &lt;/ListInvoices&gt;
  &lt;/soap:Body&gt;
&lt;/soap:Envelope&gt;'</code></pre>
        <p class="muted">The response includes <code>SAP_SESSIONID_MKS_100</code> and <code>x-csrf-token</code>. Send the cookie on later SAP requests.</p>
    </details>

    <h2>Responses and State</h2>
    <div class="section">
        <ul>
            <li>Responses are SOAP envelopes with <code>{Operation}Response</code> in the body.</li>
            <li>Validation and business errors are returned as SOAP faults with HTTP status codes such as <code>400</code> or <code>401</code>.</li>
            <li>List operations return arrays under names such as <code>customers</code>, <code>stocks</code>, <code>orders</code>, <code>invoices</code>, and <code>cashTransactions</code>.</li>
            <li>Create operations generate missing IDs and default timestamps, then persist into that endpoint's JSON store.</li>
            <li>Approve, cancel, reverse, and invoice-from-order operations update related entity state in the same endpoint store.</li>
        </ul>
    </div>

    <h2>Standards Compliance</h2>
    <table>
        <tr><th>Feature</th><th>Status</th></tr>
        <tr><td>SOAP 1.1 &amp; 1.2</td><td>All endpoints support both</td></tr>
        <tr><td>WS-Policy</td><td>WSSE, Basic, NTLM WSDLs include policy assertions</td></tr>
        <tr><td>WS-Security</td><td>UsernameToken (PasswordText + PasswordDigest)</td></tr>
        <tr><td>NTLM 3-way handshake</td><td>Type1→Type2→Type3 flow</td></tr>
        <tr><td>SAP session</td><td>Cookie + CSRF token flow</td></tr>
        <tr><td>Dynamic WSDL addresses</td><td>Derived from request Host header</td></tr>
    </table>
</body>
</html>`
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}

// ============================================================================
// WSDL Handlers
// ============================================================================

// baseURL builds the service base URL from the incoming request Host header.
// Falls back to localhost with the configured port if Host is empty.
func baseURL(r *http.Request) string {
	host := r.Host
	if host == "" {
		host = "localhost" + addr
	}
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	return scheme + "://" + host
}

// addr is set in main() and used by baseURL as fallback
var addr string

func handleSessionWSDL(w http.ResponseWriter, r *http.Request) {
	base := baseURL(r)
	wsdl := buildWSDL(wsdlConfig{
		Namespace:            "http://mock.imaxis.com/session",
		ServiceName:          "SessionAuthService",
		PortTypeName:         "SessionAuthPortType",
		BindingName:          "SessionAuthBinding",
		Address:              base + "/session/soap",
		Endpoint:             "session",
		IncludeSessionOps:    true,
		IncludeSessionHeader: true,
		IncludePublic:        false,
		Policy:               PolicyNone,
	})
	w.Header().Set("Content-Type", "text/xml")
	w.Write([]byte(wsdl))
}

func handleWSSEWSDL(w http.ResponseWriter, r *http.Request) {
	base := baseURL(r)
	wsdl := buildWSDL(wsdlConfig{
		Namespace:            "http://mock.imaxis.com/wsse",
		ServiceName:          "WSSEAuthService",
		PortTypeName:         "WSSEAuthPortType",
		BindingName:          "WSSEAuthBinding",
		Address:              base + "/wsse/soap",
		Endpoint:             "wsse",
		IncludeSessionOps:    false,
		IncludeSessionHeader: false,
		IncludePublic:        false,
		Policy:               PolicyWSSEUsernameToken,
	})
	w.Header().Set("Content-Type", "text/xml")
	w.Write([]byte(wsdl))
}

func handleBasicWSDL(w http.ResponseWriter, r *http.Request) {
	base := baseURL(r)
	wsdl := buildWSDL(wsdlConfig{
		Namespace:            "http://mock.imaxis.com/basic",
		ServiceName:          "BasicAuthService",
		PortTypeName:         "BasicAuthPortType",
		BindingName:          "BasicAuthBinding",
		Address:              base + "/basic/soap",
		Endpoint:             "basic",
		IncludeSessionOps:    false,
		IncludeSessionHeader: false,
		IncludePublic:        false,
		Policy:               PolicyBasicAuth,
	})
	w.Header().Set("Content-Type", "text/xml")
	w.Write([]byte(wsdl))
}

func handleNTLMWSDL(w http.ResponseWriter, r *http.Request) {
	base := baseURL(r)
	wsdl := buildWSDL(wsdlConfig{
		Namespace:            "http://mock.imaxis.com/ntlm",
		ServiceName:          "NTLMAuthService",
		PortTypeName:         "NTLMAuthPortType",
		BindingName:          "NTLMAuthBinding",
		Address:              base + "/ntlm/soap",
		Endpoint:             "ntlm",
		IncludeSessionOps:    false,
		IncludeSessionHeader: false,
		IncludePublic:        false,
		Policy:               PolicyNTLMNegotiate,
	})
	w.Header().Set("Content-Type", "text/xml")
	w.Write([]byte(wsdl))
}

func handleNoAuthWSDL(w http.ResponseWriter, r *http.Request) {
	base := baseURL(r)
	wsdl := buildWSDL(wsdlConfig{
		Namespace:            "http://mock.imaxis.com/public",
		ServiceName:          "PublicService",
		PortTypeName:         "PublicPortType",
		BindingName:          "PublicBinding",
		Address:              base + "/noauth/soap",
		Endpoint:             "public",
		IncludeSessionOps:    false,
		IncludeSessionHeader: false,
		IncludePublic:        true,
		Policy:               PolicyNone,
	})
	w.Header().Set("Content-Type", "text/xml")
	w.Write([]byte(wsdl))
}

func handleSAPWSDL(w http.ResponseWriter, r *http.Request) {
	base := baseURL(r)
	wsdl := buildWSDL(wsdlConfig{
		Namespace:            "http://mock.imaxis.com/sap",
		ServiceName:          "SAPService",
		PortTypeName:         "SAPPortType",
		BindingName:          "SAPBinding",
		Address:              base + "/sap/soap",
		Endpoint:             "sap",
		IncludeSessionOps:    false,
		IncludeSessionHeader: false,
		IncludePublic:        false,
		Policy:               PolicyBasicAuth,
	})
	w.Header().Set("Content-Type", "text/xml")
	w.Write([]byte(wsdl))
}

// ============================================================================
// SOAP Handlers
// ============================================================================

func handleSessionSOAP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeSOAPFault(w, soapVersionFromContentType(r.Header.Get("Content-Type")), http.StatusBadRequest, "soap:Client", "", "Failed to read SOAP body")
		return
	}

	info, parseErr := parseSOAPMessage(body)
	soapVersion := soapVersionFromContentType(r.Header.Get("Content-Type"))
	if parseErr == nil && info.Version != "" {
		soapVersion = info.Version
	}

	soapAction := getSOAPAction(r)
	operation := ""
	if parseErr == nil {
		operation = info.Operation
	}
	if operation == "" && soapAction != "" {
		operation = operationFromSOAPAction(soapAction)
	}

	log.Printf("[SESSION] SOAPAction: %s Operation: %s SOAP %s", soapAction, operation, soapVersion)

	store, storeErr := storeForEndpoint("session")
	if storeErr != nil {
		writeSOAPFault(w, soapVersion, http.StatusInternalServerError, "soap:Server", "", "Store not available")
		return
	}

	switch operation {
	case "Login":
		username := lookupBodyValue(info, []string{"username", "user", "userid", "Username", "User", "UserId"})
		password := lookupBodyValue(info, []string{"password", "pass", "Password", "Pass"})
		log.Printf("[SESSION] Login attempt - username: '%s', password length: %d", username, len(password))
		if !validateCredentials(username, password) {
			log.Printf("[SESSION] ❌ Login FAILED - Invalid credentials for user: '%s'", username)
			writeSOAPXML(w, soapVersion, http.StatusOK, loginResponse(soapVersion, "", 0, false))
			return
		}

		token := generateToken()
		session := &Session{
			Token:     token,
			Username:  username,
			CreatedAt: time.Now(),
			ExpiresAt: time.Now().Add(30 * time.Minute),
		}

		sessionMutex.Lock()
		sessions[token] = session
		sessionMutex.Unlock()

		http.SetCookie(w, &http.Cookie{
			Name:     SessionCookieName,
			Value:    token,
			Path:     "/",
			HttpOnly: true,
			Expires:  session.ExpiresAt,
		})
		w.Header().Set("X-Session-Token", token)

		log.Printf("[SESSION] Login successful for user: %s, token: %s...", username, token[:min(16, len(token))])
		writeSOAPXML(w, soapVersion, http.StatusOK, loginResponse(soapVersion, token, int(time.Until(session.ExpiresAt).Seconds()), true))
		return

	case "Logout":
		token := extractSessionToken(info, r)
		if token == "" || !isValidSession(token) {
			writeSOAPFault(w, soapVersion, http.StatusUnauthorized, "soap:Client", "", "Invalid or expired session token")
			return
		}

		sessionMutex.Lock()
		delete(sessions, token)
		sessionMutex.Unlock()

		http.SetCookie(w, &http.Cookie{
			Name:     SessionCookieName,
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			HttpOnly: true,
		})

		log.Printf("[SESSION] Logout for token: %s...", token[:min(16, len(token))])
		writeSOAPXML(w, soapVersion, http.StatusOK, logoutResponse(soapVersion, true))
		return

	default:
		token := extractSessionToken(info, r)
		if token == "" || !isValidSession(token) {
			writeSOAPFault(w, soapVersion, http.StatusUnauthorized, "soap:Client", "", "Invalid or expired session token")
			return
		}

		sessionMutex.RLock()
		session := sessions[token]
		sessionMutex.RUnlock()
		log.Printf("[SESSION] Auth OK - user: %s", session.Username)

		if handleERPOperation(w, soapVersion, operation, store, "http://mock.imaxis.com/session", body) {
			return
		}
		writeSOAPFault(w, soapVersion, http.StatusBadRequest, "soap:Client", "", "Unknown operation")
		return
	}
}

func handleWSSESOAP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeSOAPFault(w, soapVersionFromContentType(r.Header.Get("Content-Type")), http.StatusBadRequest, "soap:Client", "", "Failed to read SOAP body")
		return
	}

	info, parseErr := parseSOAPMessage(body)
	soapVersion := soapVersionFromContentType(r.Header.Get("Content-Type"))
	if parseErr == nil && info.Version != "" {
		soapVersion = info.Version
	}

	if parseErr != nil || info == nil {
		writeSOAPFault(w, soapVersion, http.StatusBadRequest, "soap:Client", "", "Invalid SOAP message")
		return
	}

	if info.WSSECredentials == nil {
		writeSOAPFault(w, soapVersion, http.StatusInternalServerError, "soap:Sender", "wsse:MissingSecurity", "WS-Security header required")
		return
	}

	if !validateWSSE(info.WSSECredentials) {
		writeSOAPFault(w, soapVersion, http.StatusInternalServerError, "soap:Sender", "wsse:FailedAuthentication", "Invalid WS-Security credentials")
		return
	}

	soapAction := getSOAPAction(r)
	operation := info.Operation
	if operation == "" && soapAction != "" {
		operation = operationFromSOAPAction(soapAction)
	}

	log.Printf("[WSSE] SOAPAction: %s Operation: %s - Auth OK", soapAction, operation)

	store, storeErr := storeForEndpoint("wsse")
	if storeErr != nil {
		writeSOAPFault(w, soapVersion, http.StatusInternalServerError, "soap:Server", "", "Store not available")
		return
	}

	if handleERPOperation(w, soapVersion, operation, store, "http://mock.imaxis.com/wsse", body) {
		return
	}
	writeSOAPFault(w, soapVersion, http.StatusBadRequest, "soap:Client", "", "Unknown operation")
}

func handleBasicSOAP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check Basic Auth
	username, password, ok := r.BasicAuth()
	if !ok || !validateCredentials(username, password) {
		soapVersion := soapVersionFromContentType(r.Header.Get("Content-Type"))
		w.Header().Set("WWW-Authenticate", fmt.Sprintf(`Basic realm="%s"`, BasicAuthRealm))
		writeSOAPFault(w, soapVersion, http.StatusUnauthorized, "soap:Client", "", "Unauthorized - Basic authentication required")
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeSOAPFault(w, soapVersionFromContentType(r.Header.Get("Content-Type")), http.StatusBadRequest, "soap:Client", "", "Failed to read SOAP body")
		return
	}

	info, parseErr := parseSOAPMessage(body)
	soapVersion := soapVersionFromContentType(r.Header.Get("Content-Type"))
	if parseErr == nil && info.Version != "" {
		soapVersion = info.Version
	}

	soapAction := getSOAPAction(r)
	operation := ""
	if parseErr == nil {
		operation = info.Operation
	}
	if operation == "" && soapAction != "" {
		operation = operationFromSOAPAction(soapAction)
	}

	log.Printf("[BASIC] SOAPAction: %s Operation: %s - Auth OK for user: %s", soapAction, operation, username)

	store, storeErr := storeForEndpoint("basic")
	if storeErr != nil {
		writeSOAPFault(w, soapVersion, http.StatusInternalServerError, "soap:Server", "", "Store not available")
		return
	}

	if handleERPOperation(w, soapVersion, operation, store, "http://mock.imaxis.com/basic", body) {
		return
	}
	writeSOAPFault(w, soapVersion, http.StatusBadRequest, "soap:Client", "", "Unknown operation")
}

func handleNTLMSOAP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	soapVersion := soapVersionFromContentType(r.Header.Get("Content-Type"))
	authHeader := r.Header.Get("Authorization")

	// No Authorization header - send initial 401 challenge
	if authHeader == "" {
		log.Printf("[NTLM] No auth header - sending initial 401 challenge")
		w.Header().Set("WWW-Authenticate", "NTLM")
		w.Header().Set("Connection", "keep-alive")
		writeSOAPFault(w, soapVersion, http.StatusUnauthorized, "soap:Client", "", "NTLM authentication required")
		return
	}

	// Check if it's an NTLM auth header
	if !strings.HasPrefix(authHeader, "NTLM ") && !strings.HasPrefix(authHeader, "Negotiate ") {
		log.Printf("[NTLM] Invalid auth scheme: %s", authHeader[:min(20, len(authHeader))])
		w.Header().Set("WWW-Authenticate", "NTLM")
		writeSOAPFault(w, soapVersion, http.StatusUnauthorized, "soap:Client", "", "NTLM authentication required")
		return
	}

	// Extract the base64-encoded NTLM message
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 {
		w.Header().Set("WWW-Authenticate", "NTLM")
		writeSOAPFault(w, soapVersion, http.StatusUnauthorized, "soap:Client", "", "Invalid NTLM header")
		return
	}

	ntlmData, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		log.Printf("[NTLM] Failed to decode NTLM data: %v", err)
		w.Header().Set("WWW-Authenticate", "NTLM")
		writeSOAPFault(w, soapVersion, http.StatusUnauthorized, "soap:Client", "", "Invalid NTLM encoding")
		return
	}

	// Check NTLM message type (byte 8 in the message)
	if len(ntlmData) < 12 {
		log.Printf("[NTLM] NTLM message too short: %d bytes", len(ntlmData))
		w.Header().Set("WWW-Authenticate", "NTLM")
		writeSOAPFault(w, soapVersion, http.StatusUnauthorized, "soap:Client", "", "Invalid NTLM message")
		return
	}

	// Verify NTLM signature "NTLMSSP\0"
	if string(ntlmData[0:7]) != "NTLMSSP" || ntlmData[7] != 0 {
		log.Printf("[NTLM] Invalid NTLM signature")
		w.Header().Set("WWW-Authenticate", "NTLM")
		writeSOAPFault(w, soapVersion, http.StatusUnauthorized, "soap:Client", "", "Invalid NTLM signature")
		return
	}

	messageType := ntlmData[8]

	switch messageType {
	case 1: // Type 1 - Negotiate message
		log.Printf("[NTLM] Received Type1 (Negotiate) - sending Type2 (Challenge)")
		challenge := generateNTLMType2Challenge()
		w.Header().Set("WWW-Authenticate", "NTLM "+challenge)
		w.Header().Set("Connection", "keep-alive")
		writeSOAPFault(w, soapVersion, http.StatusUnauthorized, "soap:Client", "", "NTLM challenge")
		return

	case 3: // Type 3 - Authenticate message
		log.Printf("[NTLM] Received Type3 (Authenticate)")

		// Extract username and domain from Type3 message
		username, domain := parseNTLMType3(ntlmData)
		log.Printf("[NTLM] Type3 - Domain: %s, User: %s", domain, username)

		// For mock purposes, we accept any well-formed Type3 message
		// with our test credentials (simplified validation)
		if !validateNTLMCredentials(username, domain) {
			log.Printf("[NTLM] Authentication failed for user: %s\\%s", domain, username)
			w.Header().Set("WWW-Authenticate", "NTLM")
			writeSOAPFault(w, soapVersion, http.StatusUnauthorized, "soap:Client", "", "NTLM authentication failed")
			return
		}

		log.Printf("[NTLM] Authentication successful for user: %s\\%s", domain, username)

		// Authentication successful - process the SOAP request
		body, err := io.ReadAll(r.Body)
		if err != nil {
			writeSOAPFault(w, soapVersion, http.StatusBadRequest, "soap:Client", "", "Failed to read SOAP body")
			return
		}

		info, parseErr := parseSOAPMessage(body)
		if parseErr == nil && info.Version != "" {
			soapVersion = info.Version
		}

		soapAction := getSOAPAction(r)
		operation := ""
		if parseErr == nil {
			operation = info.Operation
		}
		if operation == "" && soapAction != "" {
			operation = operationFromSOAPAction(soapAction)
		}

		log.Printf("[NTLM] SOAPAction: %s Operation: %s - Auth OK for user: %s\\%s", soapAction, operation, domain, username)

		store, storeErr := storeForEndpoint("ntlm")
		if storeErr != nil {
			writeSOAPFault(w, soapVersion, http.StatusInternalServerError, "soap:Server", "", "Store not available")
			return
		}

		if handleERPOperation(w, soapVersion, operation, store, "http://mock.imaxis.com/ntlm", body) {
			return
		}
		writeSOAPFault(w, soapVersion, http.StatusBadRequest, "soap:Client", "", "Unknown operation")
		return

	default:
		log.Printf("[NTLM] Unknown NTLM message type: %d", messageType)
		w.Header().Set("WWW-Authenticate", "NTLM")
		writeSOAPFault(w, soapVersion, http.StatusUnauthorized, "soap:Client", "", "Invalid NTLM message type")
		return
	}
}

// generateNTLMType2Challenge creates a mock NTLM Type 2 (Challenge) message
func generateNTLMType2Challenge() string {
	// NTLM Type 2 message structure (simplified):
	// Bytes 0-7:   "NTLMSSP\0" signature
	// Bytes 8-11:  Message type (2)
	// Bytes 12-19: Target name (security buffer)
	// Bytes 20-23: Flags
	// Bytes 24-31: Server challenge (8 bytes)
	// ... (rest is optional target info)

	targetName := "TESTDOMAIN"
	targetNameBytes := encodeUTF16LE(targetName)

	// Calculate offsets
	targetNameLen := len(targetNameBytes)
	targetNameOffset := 56 // Start of payload

	// Build the message
	msg := make([]byte, targetNameOffset+targetNameLen)

	// Signature
	copy(msg[0:8], []byte("NTLMSSP\x00"))

	// Message type (2)
	msg[8] = 2
	msg[9] = 0
	msg[10] = 0
	msg[11] = 0

	// Target name security buffer (len, maxlen, offset)
	binary.LittleEndian.PutUint16(msg[12:14], uint16(targetNameLen))
	binary.LittleEndian.PutUint16(msg[14:16], uint16(targetNameLen))
	binary.LittleEndian.PutUint32(msg[16:20], uint32(targetNameOffset))

	// Flags (NTLMSSP_NEGOTIATE_UNICODE | NTLMSSP_NEGOTIATE_NTLM | NTLMSSP_TARGET_TYPE_DOMAIN)
	flags := uint32(0x00008201)
	binary.LittleEndian.PutUint32(msg[20:24], flags)

	// Server challenge (8 random bytes)
	challenge := make([]byte, 8)
	rand.Read(challenge)
	copy(msg[24:32], challenge)

	// Reserved (8 bytes of zeros)
	// Already zero from make()

	// Target name payload
	copy(msg[targetNameOffset:], targetNameBytes)

	return base64.StdEncoding.EncodeToString(msg)
}

// parseNTLMType3 extracts username and domain from Type 3 message
func parseNTLMType3(data []byte) (username, domain string) {
	if len(data) < 72 {
		return "", ""
	}

	// Domain name security buffer at offset 28
	domainLen := binary.LittleEndian.Uint16(data[28:30])
	domainOffset := binary.LittleEndian.Uint32(data[32:36])

	// User name security buffer at offset 36
	userLen := binary.LittleEndian.Uint16(data[36:38])
	userOffset := binary.LittleEndian.Uint32(data[40:44])

	// Extract domain
	if domainOffset+uint32(domainLen) <= uint32(len(data)) {
		domain = decodeUTF16LE(data[domainOffset : domainOffset+uint32(domainLen)])
	}

	// Extract username
	if userOffset+uint32(userLen) <= uint32(len(data)) {
		username = decodeUTF16LE(data[userOffset : userOffset+uint32(userLen)])
	}

	return username, domain
}

// validateNTLMCredentials checks if the username/domain match our test credentials
func validateNTLMCredentials(username, domain string) bool {
	// Accept our test user with any domain, or TESTDOMAIN specifically
	if strings.EqualFold(username, TestUsername) {
		// Accept any domain or TESTDOMAIN
		return domain == "" || strings.EqualFold(domain, "TESTDOMAIN")
	}
	return false
}

// encodeUTF16LE encodes a string to UTF-16LE bytes
func encodeUTF16LE(s string) []byte {
	runes := []rune(s)
	result := make([]byte, len(runes)*2)
	for i, r := range runes {
		result[i*2] = byte(r)
		result[i*2+1] = byte(r >> 8)
	}
	return result
}

// decodeUTF16LE decodes UTF-16LE bytes to a string
func decodeUTF16LE(b []byte) string {
	if len(b)%2 != 0 {
		return ""
	}
	runes := make([]rune, len(b)/2)
	for i := 0; i < len(b); i += 2 {
		runes[i/2] = rune(b[i]) | rune(b[i+1])<<8
	}
	return string(runes)
}

func handleNoAuthSOAP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeSOAPFault(w, soapVersionFromContentType(r.Header.Get("Content-Type")), http.StatusBadRequest, "soap:Client", "", "Failed to read SOAP body")
		return
	}

	info, parseErr := parseSOAPMessage(body)
	soapVersion := soapVersionFromContentType(r.Header.Get("Content-Type"))
	if parseErr == nil && info.Version != "" {
		soapVersion = info.Version
	}

	soapAction := getSOAPAction(r)
	operation := ""
	if parseErr == nil {
		operation = info.Operation
	}
	if operation == "" && soapAction != "" {
		operation = operationFromSOAPAction(soapAction)
	}

	log.Printf("[NOAUTH] SOAPAction: %s Operation: %s", soapAction, operation)

	store, storeErr := storeForEndpoint("noauth")
	if storeErr != nil {
		writeSOAPFault(w, soapVersion, http.StatusInternalServerError, "soap:Server", "", "Store not available")
		return
	}

	if handlePublicOperation(w, soapVersion, operation, store, "http://mock.imaxis.com/public", body) {
		return
	}

	if handleERPOperation(w, soapVersion, operation, store, "http://mock.imaxis.com/public", body) {
		return
	}

	writeSOAPFault(w, soapVersion, http.StatusBadRequest, "soap:Client", "", "Unknown operation")
}

// ============================================================================
// SAP Session Handler
// ============================================================================
// SAP session flow:
//   1. Client sends first request with HTTP Basic Auth + sap-client header
//      → Server validates credentials, creates session
//      → Response includes Set-Cookie: SAP_SESSIONID_xxx=<token> and x-csrf-token
//   2. Subsequent requests send the cookie + x-csrf-token header (for mutating ops)
//      → No Basic Auth needed, session cookie is sufficient

const (
	SAPClientHeader  = "sap-client"
	SAPCSRFHeader    = "x-csrf-token"
	SAPSessionCookie = "SAP_SESSIONID_MKS_100"
	SAPDefaultClient = "100"
)

func handleSAPSOAP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	soapVersion := soapVersionFromContentType(r.Header.Get("Content-Type"))

	// Check for existing SAP session via cookie
	var sapSession *SAPSession
	if cookie, err := r.Cookie(SAPSessionCookie); err == nil {
		sapSessionMutex.RLock()
		s, exists := sapSessions[cookie.Value]
		sapSessionMutex.RUnlock()
		if exists && time.Now().Before(s.ExpiresAt) {
			sapSession = s
		}
	}

	// If no valid session, require Basic Auth to create one (SAP login flow)
	if sapSession == nil {
		username, password, ok := r.BasicAuth()
		if !ok || !validateCredentials(username, password) {
			w.Header().Set("WWW-Authenticate", `Basic realm="SAP Mock Server"`)
			writeSOAPFault(w, soapVersion, http.StatusUnauthorized, "soap:Client", "", "SAP authentication required - provide Basic Auth credentials with sap-client header")
			return
		}

		sapClient := r.Header.Get(SAPClientHeader)
		if sapClient == "" {
			sapClient = SAPDefaultClient
		}

		// Create SAP session
		sessionID := generateToken()
		csrfToken := generateToken()
		sapSession = &SAPSession{
			SessionID: sessionID,
			CSRFToken: csrfToken,
			SAPClient: sapClient,
			Username:  username,
			CreatedAt: time.Now(),
			ExpiresAt: time.Now().Add(30 * time.Minute),
		}

		sapSessionMutex.Lock()
		sapSessions[sessionID] = sapSession
		sapSessionMutex.Unlock()

		// Set SAP session cookie and CSRF token in response
		http.SetCookie(w, &http.Cookie{
			Name:     SAPSessionCookie,
			Value:    sessionID,
			Path:     "/sap/",
			HttpOnly: true,
			Expires:  sapSession.ExpiresAt,
		})
		w.Header().Set(SAPCSRFHeader, csrfToken)
		w.Header().Set(SAPClientHeader, sapClient)

		log.Printf("[SAP] New session created for user: %s, sap-client: %s, session: %s...", username, sapClient, sessionID[:min(16, len(sessionID))])
	} else {
		log.Printf("[SAP] Session auth OK - user: %s, sap-client: %s", sapSession.Username, sapSession.SAPClient)
	}

	// Handle CSRF token fetch request (SAP convention: GET/HEAD with x-csrf-token: Fetch)
	// Some SAP clients also do this via POST, so check the header
	if strings.EqualFold(r.Header.Get(SAPCSRFHeader), "Fetch") {
		w.Header().Set(SAPCSRFHeader, sapSession.CSRFToken)
		w.Header().Set(SAPClientHeader, sapSession.SAPClient)
		w.WriteHeader(http.StatusOK)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeSOAPFault(w, soapVersion, http.StatusBadRequest, "soap:Client", "", "Failed to read SOAP body")
		return
	}

	info, parseErr := parseSOAPMessage(body)
	if parseErr == nil && info.Version != "" {
		soapVersion = info.Version
	}

	soapAction := getSOAPAction(r)
	operation := ""
	if parseErr == nil {
		operation = info.Operation
	}
	if operation == "" && soapAction != "" {
		operation = operationFromSOAPAction(soapAction)
	}

	log.Printf("[SAP] SOAPAction: %s Operation: %s", soapAction, operation)

	store, storeErr := storeForEndpoint("sap")
	if storeErr != nil {
		writeSOAPFault(w, soapVersion, http.StatusInternalServerError, "soap:Server", "", "Store not available")
		return
	}

	if handleERPOperation(w, soapVersion, operation, store, "http://mock.imaxis.com/sap", body) {
		return
	}
	writeSOAPFault(w, soapVersion, http.StatusBadRequest, "soap:Client", "", "Unknown operation")
}

// ============================================================================
// Helper Functions
// ============================================================================

func generateToken() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("fallback-%d", time.Now().UnixNano())))
	}
	return base64.StdEncoding.EncodeToString(b)
}

func validateCredentials(username, password string) bool {
	return constantTimeEquals(username, TestUsername) && constantTimeEquals(password, TestPassword)
}

func constantTimeEquals(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}

func validateWSSE(creds *WSSECredentials) bool {
	if creds == nil {
		return false
	}
	if !constantTimeEquals(creds.Username, TestUsername) {
		return false
	}

	passwordType := normalizePasswordType(creds.PasswordType)
	if passwordType == "PasswordDigest" {
		return validateWSSEDigest(creds)
	}
	return constantTimeEquals(creds.Password, TestPassword)
}

func normalizePasswordType(passwordType string) string {
	if passwordType == "" {
		return "PasswordText"
	}
	lower := strings.ToLower(passwordType)
	if strings.Contains(lower, "passworddigest") {
		return "PasswordDigest"
	}
	if strings.Contains(lower, "passwordtext") {
		return "PasswordText"
	}
	return passwordType
}

func validateWSSEDigest(creds *WSSECredentials) bool {
	if creds == nil || creds.Nonce == "" || creds.Created == "" || creds.Password == "" {
		return false
	}
	nonceBytes, err := base64.StdEncoding.DecodeString(creds.Nonce)
	if err != nil {
		return false
	}

	h := sha1.New()
	h.Write(nonceBytes)
	h.Write([]byte(creds.Created))
	h.Write([]byte(TestPassword))
	digest := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return constantTimeEquals(creds.Password, digest)
}

func parseSOAPMessage(body []byte) (*SOAPMessageInfo, error) {
	info := &SOAPMessageInfo{
		Version:    "1.1",
		BodyValues: make(map[string]string),
	}

	decoder := xml.NewDecoder(bytes.NewReader(body))
	inHeader := false
	inBody := false
	inUsernameToken := false
	wsse := WSSECredentials{}

	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return info, err
		}

		switch t := token.(type) {
		case xml.StartElement:
			local := t.Name.Local
			switch local {
			case "Envelope":
				if t.Name.Space == Soap12Namespace {
					info.Version = "1.2"
				} else if t.Name.Space == Soap11Namespace {
					info.Version = "1.1"
				}
			case "Header":
				inHeader = true
			case "Body":
				inBody = true
			default:
				if inBody && info.Operation == "" && local != "Fault" {
					info.Operation = local
					info.OperationNS = t.Name.Space
				}

				if inHeader {
					if local == "UsernameToken" {
						inUsernameToken = true
						continue
					}

					if inUsernameToken {
						switch local {
						case "Username":
							var val string
							if err := decoder.DecodeElement(&val, &t); err == nil {
								wsse.Username = strings.TrimSpace(val)
							}
							continue
						case "Password":
							wsse.PasswordType = attrValue(t.Attr, "Type")
							var val string
							if err := decoder.DecodeElement(&val, &t); err == nil {
								wsse.Password = strings.TrimSpace(val)
							}
							continue
						case "Nonce":
							var val string
							if err := decoder.DecodeElement(&val, &t); err == nil {
								wsse.Nonce = strings.TrimSpace(val)
							}
							continue
						case "Created":
							var val string
							if err := decoder.DecodeElement(&val, &t); err == nil {
								wsse.Created = strings.TrimSpace(val)
							}
							continue
						}
					}

					if info.SessionToken == "" && isSessionTokenName(local) {
						var val string
						if err := decoder.DecodeElement(&val, &t); err == nil {
							info.SessionToken = strings.TrimSpace(val)
						}
						continue
					}
				}

				if inBody && isTrackedBodyValue(local) {
					if _, exists := info.BodyValues[strings.ToLower(local)]; !exists {
						var val string
						if err := decoder.DecodeElement(&val, &t); err == nil {
							info.BodyValues[strings.ToLower(local)] = strings.TrimSpace(val)
							if info.SessionToken == "" && isSessionTokenName(local) {
								info.SessionToken = strings.TrimSpace(val)
							}
						}
					}
					continue
				}
			}

		case xml.EndElement:
			switch t.Name.Local {
			case "Header":
				inHeader = false
			case "Body":
				inBody = false
			case "UsernameToken":
				inUsernameToken = false
			}
		}
	}

	if wsse.Username != "" || wsse.Password != "" {
		info.WSSECredentials = &wsse
	}

	return info, nil
}

func attrValue(attrs []xml.Attr, name string) string {
	for _, attr := range attrs {
		if strings.EqualFold(attr.Name.Local, name) {
			return attr.Value
		}
	}
	return ""
}

func isTrackedBodyValue(name string) bool {
	switch strings.ToLower(name) {
	case "username", "user", "userid", "password", "pass", "sessiontoken", "sessionid", "token":
		return true
	default:
		return false
	}
}

func isSessionTokenName(name string) bool {
	switch strings.ToLower(name) {
	case "sessiontoken", "sessionid", "token":
		return true
	default:
		return false
	}
}

func soapVersionFromContentType(contentType string) string {
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err == nil && strings.EqualFold(mediaType, "application/soap+xml") {
		return "1.2"
	}
	if strings.Contains(strings.ToLower(contentType), "application/soap+xml") {
		return "1.2"
	}
	return "1.1"
}

func normalizeSOAPVersion(version string) string {
	if version == "1.2" {
		return "1.2"
	}
	return "1.1"
}

func getSOAPAction(r *http.Request) string {
	soapAction := strings.Trim(r.Header.Get("SOAPAction"), "\"")
	if soapAction != "" {
		return soapAction
	}

	mediaType, params, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err == nil && strings.EqualFold(mediaType, "application/soap+xml") {
		if action, ok := params["action"]; ok {
			return strings.Trim(action, "\"")
		}
	}
	return ""
}

func operationFromSOAPAction(action string) string {
	action = strings.TrimSpace(strings.Trim(action, "\""))
	if action == "" {
		return ""
	}
	if idx := strings.LastIndexAny(action, "/#"); idx != -1 && idx+1 < len(action) {
		return action[idx+1:]
	}
	return action
}

func lookupBodyValue(info *SOAPMessageInfo, keys []string) string {
	if info == nil || info.BodyValues == nil {
		return ""
	}
	for _, key := range keys {
		if val := info.BodyValues[strings.ToLower(key)]; val != "" {
			return val
		}
	}
	return ""
}

func extractSessionToken(info *SOAPMessageInfo, r *http.Request) string {
	if info != nil {
		if info.SessionToken != "" {
			return normalizeSessionToken(info.SessionToken)
		}
		if val := lookupBodyValue(info, []string{"sessionToken", "sessionId", "token"}); val != "" {
			return normalizeSessionToken(val)
		}
	}

	if headerToken := r.Header.Get("X-Session-Token"); headerToken != "" {
		return normalizeSessionToken(headerToken)
	}

	for _, name := range []string{SessionCookieName, "JSESSIONID", "ASP.NET_SessionId"} {
		if cookie, err := r.Cookie(name); err == nil {
			return normalizeSessionToken(cookie.Value)
		}
	}
	return ""
}

func normalizeSessionToken(value string) string {
	token := strings.TrimSpace(value)
	if idx := strings.Index(token, ";"); idx != -1 {
		token = token[:idx]
	}

	lower := strings.ToLower(token)
	for _, name := range []string{SessionCookieName, "jsessionid", "asp.net_sessionid"} {
		prefix := strings.ToLower(name) + "="
		if strings.HasPrefix(lower, prefix) {
			return token[len(prefix):]
		}
	}
	return token
}

func isValidSession(token string) bool {
	sessionMutex.RLock()
	session, exists := sessions[token]
	sessionMutex.RUnlock()
	return exists && time.Now().Before(session.ExpiresAt)
}

func writeSOAPXML(w http.ResponseWriter, version string, status int, body string) {
	contentType := "text/xml; charset=utf-8"
	if version == "1.2" {
		contentType = "application/soap+xml; charset=utf-8"
	}
	w.Header().Set("Content-Type", contentType)
	if status > 0 {
		w.WriteHeader(status)
	}
	w.Write([]byte(body))
}

func writeSOAPFault(w http.ResponseWriter, version string, status int, code, subcode, message string) {
	writeSOAPXML(w, version, status, soapFault(version, code, subcode, message))
}

func normalizeFaultCode(version, code string) string {
	if version == "1.2" {
		switch code {
		case "soap:Client":
			return "soap:Sender"
		case "soap:Server":
			return "soap:Receiver"
		default:
			return code
		}
	}

	switch code {
	case "soap:Sender":
		return "soap:Client"
	case "soap:Receiver":
		return "soap:Server"
	default:
		return code
	}
}

func soapFault(version, code, subcode, message string) string {
	version = normalizeSOAPVersion(version)
	faultCode := normalizeFaultCode(version, code)
	escapedMessage := html.EscapeString(message)

	if version == "1.2" {
		subcodeBlock := ""
		wsseNS := ""
		if subcode != "" {
			subcodeBlock = fmt.Sprintf("<soap:Subcode><soap:Value>%s</soap:Value></soap:Subcode>", html.EscapeString(subcode))
			if strings.HasPrefix(subcode, "wsse:") {
				wsseNS = ` xmlns:wsse="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-secext-1.0.xsd"`
			}
		}
		return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="%s"%s>
    <soap:Body>
        <soap:Fault>
            <soap:Code>
                <soap:Value>%s</soap:Value>
                %s
            </soap:Code>
            <soap:Reason>
                <soap:Text xml:lang="en">%s</soap:Text>
            </soap:Reason>
        </soap:Fault>
    </soap:Body>
</soap:Envelope>`, Soap12Namespace, wsseNS, faultCode, subcodeBlock, escapedMessage)
	}

	detail := ""
	if subcode != "" {
		detail = fmt.Sprintf("<detail>%s</detail>", html.EscapeString(subcode))
	}
	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="%s">
    <soap:Body>
        <soap:Fault>
            <faultcode>%s</faultcode>
            <faultstring>%s</faultstring>
            %s
        </soap:Fault>
    </soap:Body>
</soap:Envelope>`, Soap11Namespace, faultCode, escapedMessage, detail)
}

func soapEnvelope(version, body string) string {
	soapNS := Soap11Namespace
	if version == "1.2" {
		soapNS = Soap12Namespace
	}
	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="%s">
    <soap:Body>
%s
    </soap:Body>
</soap:Envelope>`, soapNS, body)
}

func loginResponse(version, token string, expiresIn int, success bool) string {
	body := fmt.Sprintf(`        <LoginResponse xmlns="http://mock.imaxis.com/session">
            <sessionToken>%s</sessionToken>
            <expiresIn>%d</expiresIn>
            <success>%t</success>
        </LoginResponse>`, html.EscapeString(token), expiresIn, success)
	return soapEnvelope(version, body)
}

func logoutResponse(version string, success bool) string {
	body := fmt.Sprintf(`        <LogoutResponse xmlns="http://mock.imaxis.com/session">
            <success>%t</success>
        </LogoutResponse>`, success)
	return soapEnvelope(version, body)
}

func customersResponse(ns, version string) string {
	body := fmt.Sprintf(`        <GetCustomersResponse xmlns="%s">
            <customers>
                <customer>
                    <id>1</id>
                    <name>Ahmet Yılmaz</name>
                    <email>ahmet@example.com</email>
                </customer>
                <customer>
                    <id>2</id>
                    <name>Mehmet Demir</name>
                    <email>mehmet@example.com</email>
                </customer>
                <customer>
                    <id>3</id>
                    <name>Ayşe Kaya</name>
                    <email>ayse@example.com</email>
                </customer>
            </customers>
        </GetCustomersResponse>`, ns)
	return soapEnvelope(version, body)
}

func ordersResponse(ns, version string) string {
	body := fmt.Sprintf(`        <GetOrdersResponse xmlns="%s">
            <orders>
                <order>
                    <id>1001</id>
                    <date>2024-01-15</date>
                    <amount>1500.00</amount>
                </order>
                <order>
                    <id>1002</id>
                    <date>2024-01-20</date>
                    <amount>2300.50</amount>
                </order>
            </orders>
        </GetOrdersResponse>`, ns)
	return soapEnvelope(version, body)
}

func countriesResponse(version string) string {
	body := `        <GetCountriesResponse xmlns="http://mock.imaxis.com/public">
            <countries>
                <country>
                    <code>TR</code>
                    <name>Turkey</name>
                    <continent>Europe</continent>
                </country>
                <country>
                    <code>DE</code>
                    <name>Germany</name>
                    <continent>Europe</continent>
                </country>
                <country>
                    <code>US</code>
                    <name>United States</name>
                    <continent>North America</continent>
                </country>
            </countries>
        </GetCountriesResponse>`
	return soapEnvelope(version, body)
}

func citiesResponse(version string) string {
	body := `        <GetCitiesResponse xmlns="http://mock.imaxis.com/public">
            <cities>
                <city>
                    <name>Istanbul</name>
                    <population>15000000</population>
                </city>
                <city>
                    <name>Ankara</name>
                    <population>5500000</population>
                </city>
                <city>
                    <name>Izmir</name>
                    <population>4300000</population>
                </city>
            </cities>
        </GetCitiesResponse>`
	return soapEnvelope(version, body)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
