# SOAP Mock Server

AuthProbe ve SOAP authentication flow testleri için mock server.

## Quick Start

```bash
# Server'ı başlat
./soap-mock-server

# Browser'da aç
open http://localhost:8099
```

## Credentials

| Username | Password |
|----------|----------|
| testuser | testpass123 |

## Data Stores

Each endpoint persists its own JSON file under `service/`:

- `session.json`
- `wsse.json`
- `basic.json`
- `ntlm.json`
- `noauth.json`

Override the directory with `SOAP_MOCK_DATA_DIR`.

## Endpoints

### 1. SOAP Session Auth (DIVA-like)
```
WSDL: http://localhost:8099/session/service.wsdl
SOAP: http://localhost:8099/session/soap
```

**Operations:**
- `Login` - Session token al
- `Logout` - Session sonlandır
- `ERP operations` - Customers, Stocks, Cash Accounts, Orders, Invoices, Cash Transactions (create/update/get/list/approve/cancel)

**Session taşıma yöntemleri (hepsi desteklenir):**
- HTTP Cookie: `session`
- HTTP Header: `X-Session-Token`
- SOAP Header: `SessionHeader` (sessionToken)

### 2. WS-Security UsernameToken
```
WSDL: http://localhost:8099/wsse/service.wsdl
SOAP: http://localhost:8099/wsse/soap
```

**Operations:**
- `ERP operations` - Customers, Stocks, Cash Accounts, Orders, Invoices, Cash Transactions (create/update/get/list/approve/cancel)

### 3. HTTP Basic Auth
```
WSDL: http://localhost:8099/basic/service.wsdl
SOAP: http://localhost:8099/basic/soap
```

**Operations:**
- `ERP operations` - Customers, Stocks, Cash Accounts, Orders, Invoices, Cash Transactions (create/update/get/list/approve/cancel)

### 4. NTLM Auth
```
WSDL: http://localhost:8099/ntlm/service.wsdl
SOAP: http://localhost:8099/ntlm/soap
```

**Operations:**
- `ERP operations` - Customers, Stocks, Cash Accounts, Orders, Invoices, Cash Transactions (create/update/get/list/approve/cancel)

**Credentials:** `TESTDOMAIN\testuser` / `testpass123`

### 5. No Auth (Public)
```
WSDL: http://localhost:8099/noauth/service.wsdl
SOAP: http://localhost:8099/noauth/soap
```

**Operations:**
- `GetCountries` - Auth gerektirmez
- `GetCities` - Auth gerektirmez
- `ERP operations` - Customers, Stocks, Cash Accounts, Orders, Invoices, Cash Transactions (create/update/get/list/approve/cancel)

## ERP Operations

- Customers: `GetCustomer`, `GetCustomers`/`ListCustomers`, `CreateCustomer`, `UpdateCustomer`
- Stocks: `GetStock`, `GetStocks`/`ListStocks`, `CreateStock`, `UpdateStock`
- Cash Accounts: `GetCashAccount`, `GetCashAccounts`/`ListCashAccounts`, `CreateCashAccount`, `UpdateCashAccount`
- Orders: `GetOrder`, `GetOrders`/`ListOrders`, `CreateOrder`, `UpdateOrder`, `ApproveOrder`, `CancelOrder`
- Invoices: `GetInvoice`, `GetInvoices`/`ListInvoices`, `CreateInvoice`, `CreateInvoiceFromOrder`, `CancelInvoice`
- Cash Transactions: `GetCashTransaction`, `GetCashTransactions`/`ListCashTransactions`, `CreateCashTransaction`, `ReverseCashTransaction`

---

## Test Örnekleri

### Session Auth - Login
```bash
curl -X POST http://localhost:8099/session/soap \
  -H "Content-Type: text/xml" \
  -H "SOAPAction: http://mock.imaxis.com/session/Login" \
  -d '<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <Login xmlns="http://mock.imaxis.com/session">
      <username>testuser</username>
      <password>testpass123</password>
    </Login>
  </soap:Body>
</soap:Envelope>'
```

### Session Auth - GetCustomers (Cookie ile)
```bash
# Login response Set-Cookie'den gelen session değerini kullan
curl -X POST http://localhost:8099/session/soap \
  -b "session=YOUR_TOKEN_HERE" \
  -H "Content-Type: text/xml" \
  -H "SOAPAction: http://mock.imaxis.com/session/GetCustomers" \
  -d '<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <GetCustomers xmlns="http://mock.imaxis.com/session">
      <maxResults>10</maxResults>
    </GetCustomers>
  </soap:Body>
</soap:Envelope>'
```

### Session Auth - GetCustomers (SOAP Header ile)
```bash
curl -X POST http://localhost:8099/session/soap \
  -H "Content-Type: text/xml" \
  -H "SOAPAction: http://mock.imaxis.com/session/GetCustomers" \
  -d '<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Header>
    <SessionHeader xmlns="http://mock.imaxis.com/session">
      <sessionToken>YOUR_TOKEN_HERE</sessionToken>
    </SessionHeader>
  </soap:Header>
  <soap:Body>
    <GetCustomers xmlns="http://mock.imaxis.com/session">
      <maxResults>10</maxResults>
    </GetCustomers>
  </soap:Body>
</soap:Envelope>'
```

### WS-Security - GetCustomers
```bash
curl -X POST http://localhost:8099/wsse/soap \
  -H "Content-Type: text/xml" \
  -H "SOAPAction: http://mock.imaxis.com/wsse/GetCustomers" \
  -d '<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/"
               xmlns:wsse="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-secext-1.0.xsd">
  <soap:Header>
    <wsse:Security>
      <wsse:UsernameToken>
        <wsse:Username>testuser</wsse:Username>
        <wsse:Password>testpass123</wsse:Password>
      </wsse:UsernameToken>
    </wsse:Security>
  </soap:Header>
  <soap:Body>
    <GetCustomers xmlns="http://mock.imaxis.com/wsse">
      <maxResults>10</maxResults>
    </GetCustomers>
  </soap:Body>
</soap:Envelope>'
```

### Basic Auth - GetCustomers
```bash
curl -X POST http://localhost:8099/basic/soap \
  -u testuser:testpass123 \
  -H "Content-Type: text/xml" \
  -H "SOAPAction: http://mock.imaxis.com/basic/GetCustomers" \
  -d '<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <GetCustomers xmlns="http://mock.imaxis.com/basic">
      <maxResults>10</maxResults>
    </GetCustomers>
  </soap:Body>
</soap:Envelope>'
```

### ERP - CreateCustomer (Basic Auth)
```bash
curl -X POST http://localhost:8099/basic/soap \
  -u testuser:testpass123 \
  -H "Content-Type: text/xml" \
  -H "SOAPAction: http://mock.imaxis.com/basic/CreateCustomer" \
  -d '<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <CreateCustomer xmlns="http://mock.imaxis.com/basic">
      <customer>
        <code>CUST-0100</code>
        <name>Yeni Musteri AS</name>
        <taxNumber>1234567890</taxNumber>
        <email>info@yenimusteri.com</email>
        <phone>5551112233</phone>
        <address>Istanbul</address>
        <currency>TRY</currency>
        <riskLimit>50000</riskLimit>
        <status>ACTIVE</status>
      </customer>
    </CreateCustomer>
  </soap:Body>
</soap:Envelope>'
```

### ERP - CreateOrder (Session Auth)
```bash
TOKEN=$(curl -s -X POST http://localhost:8099/session/soap \
  -H "Content-Type: text/xml" \
  -H "SOAPAction: http://mock.imaxis.com/session/Login" \
  -d '<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <Login xmlns="http://mock.imaxis.com/session">
      <username>testuser</username>
      <password>testpass123</password>
    </Login>
  </soap:Body>
</soap:Envelope>' | tr -d '\n' | sed -n 's:.*<sessionToken>\([^<]*\)</sessionToken>.*:\1:p')

curl -X POST http://localhost:8099/session/soap \
  -H "X-Session-Token: $TOKEN" \
  -H "Content-Type: text/xml" \
  -H "SOAPAction: http://mock.imaxis.com/session/CreateOrder" \
  -d '<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <CreateOrder xmlns="http://mock.imaxis.com/session">
      <order>
        <orderNo>SO-2024-TEST-01</orderNo>
        <customerCode>CUST-0001</customerCode>
        <status>DRAFT</status>
        <orderDate>2024-12-01</orderDate>
        <currency>TRY</currency>
        <notes>Test order</notes>
        <lines>
          <line>
            <stockCode>STK-0002</stockCode>
            <quantity>3</quantity>
            <unitPrice>350</unitPrice>
            <discountRate>0</discountRate>
            <taxRate>18</taxRate>
          </line>
        </lines>
      </order>
    </CreateOrder>
  </soap:Body>
</soap:Envelope>'
```

### No Auth - GetCountries
```bash
curl -X POST http://localhost:8099/noauth/soap \
  -H "Content-Type: text/xml" \
  -H "SOAPAction: http://mock.imaxis.com/public/GetCountries" \
  -d '<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <GetCountries xmlns="http://mock.imaxis.com/public">
      <continent>Europe</continent>
    </GetCountries>
  </soap:Body>
</soap:Envelope>'
```

---

## WSDL Fetch

```bash
# Session WSDL
curl http://localhost:8099/session/service.wsdl

# WSSE WSDL
curl http://localhost:8099/wsse/service.wsdl

# Basic WSDL
curl http://localhost:8099/basic/service.wsdl

# NTLM WSDL
curl http://localhost:8099/ntlm/service.wsdl

# NoAuth WSDL
curl http://localhost:8099/noauth/service.wsdl
```

---

## AuthProbe Test Senaryoları

| Senaryo | WSDL | Beklenen Detection |
|---------|------|-------------------|
| SOAP Session | /session/service.wsdl | `soap_session` with Login/Logout operations |
| WS-Security | /wsse/service.wsdl | `soap_wsse` |
| Basic Auth | /basic/service.wsdl | `basic_auth` |
| NTLM | /ntlm/service.wsdl | `ntlm` |
| No Auth | /noauth/service.wsdl | `none` |

---

## Port Değiştirmek

`main.go` içinde:
```go
port := ":8099"  // Bu satırı değiştir
```

Rebuild:
```bash
go build -o soap-mock-server .
```
