#!/bin/bash

# SOAP Mock Server Test Script
# Tests authentication patterns and ERP flows

BASE_URL="http://localhost:8099"
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

STEP=1
TOTAL=6

step_title() {
  echo -e "${YELLOW}[$STEP/$TOTAL] $1${NC}"
  STEP=$((STEP+1))
}

echo "========================================"
echo "🧪 SOAP Mock Server Test Suite"
echo "========================================"
echo ""

# Health Check
step_title "Health Check"
HEALTH=$(curl -s $BASE_URL/health)
if [ "$HEALTH" = "OK" ]; then
  echo -e "${GREEN}✓ Server is running${NC}"
else
  echo -e "${RED}✗ Server is not responding${NC}"
  exit 1
fi
echo ""

# Test WSDL Endpoints
step_title "WSDL Endpoints"
for endpoint in session wsse basic ntlm noauth; do
  WSDL=$(curl -s -o /dev/null -w "%{http_code}" $BASE_URL/$endpoint/service.wsdl)
  if [ "$WSDL" = "200" ]; then
    echo -e "${GREEN}✓ $endpoint/service.wsdl${NC}"
  else
    echo -e "${RED}✗ $endpoint/service.wsdl (HTTP $WSDL)${NC}"
  fi
done
echo ""

# Test No Auth
step_title "No Auth - GetCountries"
RESPONSE=$(curl -s -X POST $BASE_URL/noauth/soap \
  -H "Content-Type: text/xml" \
  -H "SOAPAction: http://mock.imaxis.com/public/GetCountries" \
  -d '<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <GetCountries xmlns="http://mock.imaxis.com/public">
      <continent>Europe</continent>
    </GetCountries>
  </soap:Body>
</soap:Envelope>')

if echo "$RESPONSE" | grep -q "Turkey"; then
  echo -e "${GREEN}✓ No Auth works - Got countries${NC}"
else
  echo -e "${RED}✗ No Auth failed${NC}"
  echo "$RESPONSE"
fi
echo ""

# Test Basic Auth
step_title "Basic Auth - ERP"
RESPONSE=$(curl -s -X POST $BASE_URL/basic/soap \
  -u testuser:testpass123 \
  -H "Content-Type: text/xml" \
  -H "SOAPAction: http://mock.imaxis.com/basic/GetCustomers" \
  -d '<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <GetCustomers xmlns="http://mock.imaxis.com/basic">
      <maxResults>5</maxResults>
    </GetCustomers>
  </soap:Body>
</soap:Envelope>')

if echo "$RESPONSE" | grep -q "CUST-0001"; then
  echo -e "${GREEN}✓ Basic Auth works - Got customers${NC}"
else
  echo -e "${RED}✗ Basic Auth failed${NC}"
  echo "$RESPONSE"
fi

NEW_CODE="CUST-TEST-$(date +%s)"
CREATE_RESPONSE=$(curl -s -X POST $BASE_URL/basic/soap \
  -u testuser:testpass123 \
  -H "Content-Type: text/xml" \
  -H "SOAPAction: http://mock.imaxis.com/basic/CreateCustomer" \
  -d "<?xml version=\"1.0\" encoding=\"UTF-8\"?>
<soap:Envelope xmlns:soap=\"http://schemas.xmlsoap.org/soap/envelope/\">
  <soap:Body>
    <CreateCustomer xmlns=\"http://mock.imaxis.com/basic\">
      <customer>
        <code>$NEW_CODE</code>
        <name>Test Musteri AS</name>
        <taxNumber>1234567890</taxNumber>
        <email>test@firma.com</email>
        <phone>5551112233</phone>
        <address>Istanbul</address>
        <currency>TRY</currency>
        <riskLimit>50000</riskLimit>
        <status>ACTIVE</status>
      </customer>
    </CreateCustomer>
  </soap:Body>
</soap:Envelope>")

if echo "$CREATE_RESPONSE" | grep -q "$NEW_CODE"; then
  echo -e "${GREEN}✓ CreateCustomer succeeded${NC}"
else
  echo -e "${RED}✗ CreateCustomer failed${NC}"
  echo "$CREATE_RESPONSE"
fi

# Wrong credentials
RESPONSE=$(curl -s -X POST $BASE_URL/basic/soap \
  -u wronguser:wrongpass \
  -H "Content-Type: text/xml" \
  -H "SOAPAction: http://mock.imaxis.com/basic/GetCustomers" \
  -d '<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <GetCustomers xmlns="http://mock.imaxis.com/basic"/>
  </soap:Body>
</soap:Envelope>')

if echo "$RESPONSE" | grep -q "Unauthorized"; then
  echo -e "${GREEN}✓ Basic Auth rejects wrong credentials${NC}"
else
  echo -e "${RED}✗ Basic Auth should reject wrong credentials${NC}"
fi
echo ""

# Test WS-Security
step_title "WS-Security - GetCustomers"
RESPONSE=$(curl -s -X POST $BASE_URL/wsse/soap \
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
      <maxResults>3</maxResults>
    </GetCustomers>
  </soap:Body>
</soap:Envelope>')

if echo "$RESPONSE" | grep -q "CUST-0001"; then
  echo -e "${GREEN}✓ WS-Security works - Got customers${NC}"
else
  echo -e "${RED}✗ WS-Security failed${NC}"
  echo "$RESPONSE"
fi

echo ""

# Test Session Auth
step_title "Session Auth - ERP Flow"
LOGIN_RESPONSE=$(curl -s -X POST $BASE_URL/session/soap \
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
</soap:Envelope>')

TOKEN=$(echo "$LOGIN_RESPONSE" | tr -d '\n' | sed -n 's:.*<sessionToken>\([^<]*\)</sessionToken>.*:\1:p')
if [ -z "$TOKEN" ]; then
  echo -e "${RED}✗ Login failed${NC}"
  echo "$LOGIN_RESPONSE"
  exit 1
fi

echo -e "${GREEN}✓ Login successful - Token: ${TOKEN:0:20}...${NC}"

ORDER_NO="SO-TEST-$(date +%s)"
ORDER_RESPONSE=$(curl -s -X POST $BASE_URL/session/soap \
  -H "X-Session-Token: $TOKEN" \
  -H "Content-Type: text/xml" \
  -H "SOAPAction: http://mock.imaxis.com/session/CreateOrder" \
  -d "<?xml version=\"1.0\" encoding=\"UTF-8\"?>
<soap:Envelope xmlns:soap=\"http://schemas.xmlsoap.org/soap/envelope/\">
  <soap:Body>
    <CreateOrder xmlns=\"http://mock.imaxis.com/session\">
      <order>
        <orderNo>$ORDER_NO</orderNo>
        <customerCode>CUST-0001</customerCode>
        <status>DRAFT</status>
        <orderDate>2024-12-01</orderDate>
        <currency>TRY</currency>
        <notes>Test order</notes>
        <lines>
          <line>
            <stockCode>STK-0002</stockCode>
            <quantity>1</quantity>
            <unitPrice>350</unitPrice>
            <discountRate>0</discountRate>
            <taxRate>18</taxRate>
          </line>
        </lines>
      </order>
    </CreateOrder>
  </soap:Body>
</soap:Envelope>")

if echo "$ORDER_RESPONSE" | grep -q "$ORDER_NO"; then
  echo -e "${GREEN}✓ CreateOrder succeeded${NC}"
else
  echo -e "${RED}✗ CreateOrder failed${NC}"
  echo "$ORDER_RESPONSE"
fi

APPROVE_RESPONSE=$(curl -s -X POST $BASE_URL/session/soap \
  -H "X-Session-Token: $TOKEN" \
  -H "Content-Type: text/xml" \
  -H "SOAPAction: http://mock.imaxis.com/session/ApproveOrder" \
  -d "<?xml version=\"1.0\" encoding=\"UTF-8\"?>
<soap:Envelope xmlns:soap=\"http://schemas.xmlsoap.org/soap/envelope/\">
  <soap:Body>
    <ApproveOrder xmlns=\"http://mock.imaxis.com/session\">
      <orderNo>$ORDER_NO</orderNo>
    </ApproveOrder>
  </soap:Body>
</soap:Envelope>")

if echo "$APPROVE_RESPONSE" | grep -q "APPROVED"; then
  echo -e "${GREEN}✓ ApproveOrder succeeded${NC}"
else
  echo -e "${RED}✗ ApproveOrder failed${NC}"
  echo "$APPROVE_RESPONSE"
fi

INVOICE_RESPONSE=$(curl -s -X POST $BASE_URL/session/soap \
  -H "X-Session-Token: $TOKEN" \
  -H "Content-Type: text/xml" \
  -H "SOAPAction: http://mock.imaxis.com/session/CreateInvoiceFromOrder" \
  -d "<?xml version=\"1.0\" encoding=\"UTF-8\"?>
<soap:Envelope xmlns:soap=\"http://schemas.xmlsoap.org/soap/envelope/\">
  <soap:Body>
    <CreateInvoiceFromOrder xmlns=\"http://mock.imaxis.com/session\">
      <orderNo>$ORDER_NO</orderNo>
      <invoiceDate>2024-12-02</invoiceDate>
      <dueDate>2024-12-30</dueDate>
    </CreateInvoiceFromOrder>
  </soap:Body>
</soap:Envelope>")

INVOICE_NO=$(echo "$INVOICE_RESPONSE" | tr -d '\n' | sed -n 's:.*<invoiceNo>\([^<]*\)</invoiceNo>.*:\1:p')
if [ -n "$INVOICE_NO" ]; then
  echo -e "${GREEN}✓ CreateInvoiceFromOrder succeeded - $INVOICE_NO${NC}"
else
  echo -e "${RED}✗ CreateInvoiceFromOrder failed${NC}"
  echo "$INVOICE_RESPONSE"
fi

if [ -n "$INVOICE_NO" ]; then
  TXN_RESPONSE=$(curl -s -X POST $BASE_URL/session/soap \
    -H "X-Session-Token: $TOKEN" \
    -H "Content-Type: text/xml" \
    -H "SOAPAction: http://mock.imaxis.com/session/CreateCashTransaction" \
    -d "<?xml version=\"1.0\" encoding=\"UTF-8\"?>
<soap:Envelope xmlns:soap=\"http://schemas.xmlsoap.org/soap/envelope/\">
  <soap:Body>
    <CreateCashTransaction xmlns=\"http://mock.imaxis.com/session\">
      <cashTransaction>
        <cashAccountId>1</cashAccountId>
        <customerCode>CUST-0001</customerCode>
        <invoiceNo>$INVOICE_NO</invoiceNo>
        <type>COLLECTION</type>
        <amount>1000</amount>
        <currency>TRY</currency>
        <method>CASH</method>
        <transactionDate>2024-12-03</transactionDate>
        <description>Test collection</description>
      </cashTransaction>
    </CreateCashTransaction>
  </soap:Body>
</soap:Envelope>")

  if echo "$TXN_RESPONSE" | grep -q "COLLECTION"; then
    echo -e "${GREEN}✓ CreateCashTransaction succeeded${NC}"
  else
    echo -e "${RED}✗ CreateCashTransaction failed${NC}"
    echo "$TXN_RESPONSE"
  fi
fi

LOGOUT_RESPONSE=$(curl -s -X POST $BASE_URL/session/soap \
  -H "Content-Type: text/xml" \
  -H "SOAPAction: http://mock.imaxis.com/session/Logout" \
  -d "<?xml version=\"1.0\" encoding=\"UTF-8\"?>
<soap:Envelope xmlns:soap=\"http://schemas.xmlsoap.org/soap/envelope/\">
  <soap:Body>
    <Logout xmlns=\"http://mock.imaxis.com/session\">
      <sessionToken>$TOKEN</sessionToken>
    </Logout>
  </soap:Body>
</soap:Envelope>")

if echo "$LOGOUT_RESPONSE" | grep -q "true"; then
  echo -e "${GREEN}✓ Logout successful${NC}"
fi

# Try to use expired token
RESPONSE=$(curl -s -X POST $BASE_URL/session/soap \
  -H "Content-Type: text/xml" \
  -H "SOAPAction: http://mock.imaxis.com/session/GetCustomers" \
  -d "<?xml version=\"1.0\" encoding=\"UTF-8\"?>
<soap:Envelope xmlns:soap=\"http://schemas.xmlsoap.org/soap/envelope/\">
  <soap:Body>
    <GetCustomers xmlns=\"http://mock.imaxis.com/session\">
      <sessionToken>$TOKEN</sessionToken>
    </GetCustomers>
  </soap:Body>
</soap:Envelope>")

if echo "$RESPONSE" | grep -q "Invalid or expired"; then
  echo -e "${GREEN}✓ Session invalidated after logout${NC}"
else
  echo -e "${RED}✗ Session should be invalid after logout${NC}"
fi

echo ""
echo "========================================"
echo "🎉 Test Suite Complete"
echo "========================================"
