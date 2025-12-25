package main

import (
    "bytes"
    "encoding/xml"
    "fmt"
    "net/http"
    "strconv"
    "strings"
    "time"
)

type CustomerRequest struct {
    Customer CustomerPayload `xml:"customer"`
}

type CustomerPayload struct {
    ID        int    `xml:"id"`
    Code      string `xml:"code"`
    Name      string `xml:"name"`
    TaxNumber string `xml:"taxNumber"`
    Email     string `xml:"email"`
    Phone     string `xml:"phone"`
    Address   string `xml:"address"`
    Currency  string `xml:"currency"`
    RiskLimit string `xml:"riskLimit"`
    Status    string `xml:"status"`
}

type CustomerLookupRequest struct {
    ID   int    `xml:"id"`
    Code string `xml:"code"`
}

type ListCustomersRequest struct {
    Code       string `xml:"code"`
    Status     string `xml:"status"`
    TaxNumber  string `xml:"taxNumber"`
    MaxResults int    `xml:"maxResults"`
    Offset     int    `xml:"offset"`
}

type StockRequest struct {
    Stock StockPayload `xml:"stock"`
}

type StockPayload struct {
    ID          int    `xml:"id"`
    Code        string `xml:"code"`
    Name        string `xml:"name"`
    Unit        string `xml:"unit"`
    VatRate     string `xml:"vatRate"`
    Price       string `xml:"price"`
    StockOnHand string `xml:"stockOnHand"`
    MinStock    string `xml:"minStock"`
    Status      string `xml:"status"`
}

type StockLookupRequest struct {
    ID   int    `xml:"id"`
    Code string `xml:"code"`
}

type ListStocksRequest struct {
    Code         string `xml:"code"`
    Status       string `xml:"status"`
    NameContains string `xml:"nameContains"`
    MaxResults   int    `xml:"maxResults"`
    Offset       int    `xml:"offset"`
}

type CashAccountRequest struct {
    CashAccount CashAccountPayload `xml:"cashAccount"`
}

type CashAccountPayload struct {
    ID       int    `xml:"id"`
    Code     string `xml:"code"`
    Name     string `xml:"name"`
    Currency string `xml:"currency"`
    Balance  string `xml:"balance"`
    Status   string `xml:"status"`
}

type CashAccountLookupRequest struct {
    ID   int    `xml:"id"`
    Code string `xml:"code"`
}

type ListCashAccountsRequest struct {
    Status     string `xml:"status"`
    Currency   string `xml:"currency"`
    MaxResults int    `xml:"maxResults"`
    Offset     int    `xml:"offset"`
}

type OrderRequest struct {
    Order OrderPayload `xml:"order"`
}

type OrderPayload struct {
    ID           int               `xml:"id"`
    OrderNo      string            `xml:"orderNo"`
    CustomerID   int               `xml:"customerId"`
    CustomerCode string            `xml:"customerCode"`
    Status       string            `xml:"status"`
    OrderDate    string            `xml:"orderDate"`
    Currency     string            `xml:"currency"`
    Notes        string            `xml:"notes"`
    Lines        []OrderLinePayload `xml:"lines>line"`
}

type OrderLinePayload struct {
    StockID      int    `xml:"stockId"`
    StockCode    string `xml:"stockCode"`
    Quantity     string `xml:"quantity"`
    UnitPrice    string `xml:"unitPrice"`
    DiscountRate string `xml:"discountRate"`
    TaxRate      string `xml:"taxRate"`
}

type OrderLookupRequest struct {
    ID      int    `xml:"id"`
    OrderNo string `xml:"orderNo"`
}

type ListOrdersRequest struct {
    CustomerID int    `xml:"customerId"`
    Status     string `xml:"status"`
    FromDate   string `xml:"fromDate"`
    ToDate     string `xml:"toDate"`
    MaxResults int    `xml:"maxResults"`
    Offset     int    `xml:"offset"`
}

type InvoiceRequest struct {
    Invoice InvoicePayload `xml:"invoice"`
}

type InvoicePayload struct {
    ID           int                `xml:"id"`
    InvoiceNo    string             `xml:"invoiceNo"`
    OrderID      int                `xml:"orderId"`
    OrderNo      string             `xml:"orderNo"`
    CustomerID   int                `xml:"customerId"`
    CustomerCode string             `xml:"customerCode"`
    InvoiceDate  string             `xml:"invoiceDate"`
    DueDate      string             `xml:"dueDate"`
    Currency     string             `xml:"currency"`
    Lines        []InvoiceLinePayload `xml:"lines>line"`
}

type InvoiceLinePayload struct {
    StockID      int    `xml:"stockId"`
    StockCode    string `xml:"stockCode"`
    Quantity     string `xml:"quantity"`
    UnitPrice    string `xml:"unitPrice"`
    DiscountRate string `xml:"discountRate"`
    TaxRate      string `xml:"taxRate"`
}

type InvoiceLookupRequest struct {
    ID        int    `xml:"id"`
    InvoiceNo string `xml:"invoiceNo"`
}

type CreateInvoiceFromOrderRequest struct {
    OrderID     int    `xml:"orderId"`
    OrderNo     string `xml:"orderNo"`
    InvoiceDate string `xml:"invoiceDate"`
    DueDate     string `xml:"dueDate"`
}

type ListInvoicesRequest struct {
    CustomerID int    `xml:"customerId"`
    Status     string `xml:"status"`
    FromDate   string `xml:"fromDate"`
    ToDate     string `xml:"toDate"`
    MaxResults int    `xml:"maxResults"`
    Offset     int    `xml:"offset"`
}

type CashTransactionRequest struct {
    CashTransaction CashTransactionPayload `xml:"cashTransaction"`
}

type CashTransactionPayload struct {
    ID              int    `xml:"id"`
    TxnNo           string `xml:"txnNo"`
    CashAccountID   int    `xml:"cashAccountId"`
    CashAccountCode string `xml:"cashAccountCode"`
    CustomerID      int    `xml:"customerId"`
    CustomerCode    string `xml:"customerCode"`
    InvoiceID       int    `xml:"invoiceId"`
    InvoiceNo       string `xml:"invoiceNo"`
    Type            string `xml:"type"`
    Amount          string `xml:"amount"`
    Currency        string `xml:"currency"`
    Method          string `xml:"method"`
    TransactionDate string `xml:"transactionDate"`
    Description     string `xml:"description"`
}

type CashTransactionLookupRequest struct {
    ID    int    `xml:"id"`
    TxnNo string `xml:"txnNo"`
}

type ListCashTransactionsRequest struct {
    CashAccountID int    `xml:"cashAccountId"`
    Type          string `xml:"type"`
    FromDate      string `xml:"fromDate"`
    ToDate        string `xml:"toDate"`
    MaxResults    int    `xml:"maxResults"`
    Offset        int    `xml:"offset"`
}

type CountriesRequest struct {
    Continent string `xml:"continent"`
}

type CitiesRequest struct {
    CountryCode string `xml:"countryCode"`
}

type customerResponsePayload struct {
    XMLName  xml.Name
    Xmlns    string   `xml:"xmlns,attr"`
    Customer Customer `xml:"customer"`
}

type stockResponsePayload struct {
    XMLName xml.Name
    Xmlns   string    `xml:"xmlns,attr"`
    Stock   StockItem `xml:"stock"`
}

type cashAccountResponsePayload struct {
    XMLName xml.Name
    Xmlns       string     `xml:"xmlns,attr"`
    CashAccount CashAccount `xml:"cashAccount"`
}

type orderResponsePayload struct {
    XMLName xml.Name
    Xmlns   string   `xml:"xmlns,attr"`
    Order   Order    `xml:"order"`
}

type invoiceResponsePayload struct {
    XMLName xml.Name
    Xmlns   string   `xml:"xmlns,attr"`
    Invoice Invoice  `xml:"invoice"`
}

type cashTransactionResponsePayload struct {
    XMLName xml.Name
    Xmlns            string           `xml:"xmlns,attr"`
    CashTransaction  CashTransaction  `xml:"cashTransaction"`
}

type customersResponsePayload struct {
    XMLName xml.Name
    Xmlns     string     `xml:"xmlns,attr"`
    Customers []Customer `xml:"customers>customer"`
}

type stocksResponsePayload struct {
    XMLName xml.Name
    Xmlns   string      `xml:"xmlns,attr"`
    Stocks  []StockItem `xml:"stocks>stock"`
}

type cashAccountsResponsePayload struct {
    XMLName xml.Name
    Xmlns        string       `xml:"xmlns,attr"`
    CashAccounts []CashAccount `xml:"cashAccounts>cashAccount"`
}

type ordersResponsePayload struct {
    XMLName xml.Name
    Xmlns   string   `xml:"xmlns,attr"`
    Orders  []Order  `xml:"orders>order"`
}

type invoicesResponsePayload struct {
    XMLName xml.Name
    Xmlns    string    `xml:"xmlns,attr"`
    Invoices []Invoice `xml:"invoices>invoice"`
}

type cashTransactionsResponsePayload struct {
    XMLName xml.Name
    Xmlns           string            `xml:"xmlns,attr"`
    CashTransactions []CashTransaction `xml:"cashTransactions>cashTransaction"`
}

type countriesResponsePayload struct {
    XMLName xml.Name
    Xmlns    string    `xml:"xmlns,attr"`
    Countries []Country `xml:"countries>country"`
}

type citiesResponsePayload struct {
    XMLName xml.Name
    Xmlns   string   `xml:"xmlns,attr"`
    Cities  []City   `xml:"cities>city"`
}

func decodeSOAPOperation(body []byte, operation string, out interface{}) error {
    decoder := xml.NewDecoder(bytes.NewReader(body))
    for {
        token, err := decoder.Token()
        if err != nil {
            return err
        }
        if start, ok := token.(xml.StartElement); ok {
            if start.Name.Local == operation {
                return decoder.DecodeElement(out, &start)
            }
        }
    }
}

func handleERPOperation(w http.ResponseWriter, soapVersion, operation string, store *Store, ns string, body []byte) bool {
    switch operation {
    case "GetCustomers", "ListCustomers":
        var req ListCustomersRequest
        if err := decodeSOAPOperation(body, operation, &req); err != nil {
            writeSOAPFault(w, soapVersion, http.StatusBadRequest, "soap:Client", "", "Invalid request")
            return true
        }
        customers := store.ListCustomers(CustomerFilter{
            Code:       req.Code,
            Status:     req.Status,
            TaxNumber:  req.TaxNumber,
            MaxResults: req.MaxResults,
            Offset:     req.Offset,
        })
        responseName := "GetCustomersResponse"
        if operation == "ListCustomers" {
            responseName = "ListCustomersResponse"
        }
        payload := customersResponsePayload{
            XMLName:   xml.Name{Local: responseName},
            Xmlns:     ns,
            Customers: customers,
        }
        writeSOAPResponse(w, soapVersion, http.StatusOK, payload)
        return true

    case "GetCustomer":
        var req CustomerLookupRequest
        if err := decodeSOAPOperation(body, operation, &req); err != nil {
            writeSOAPFault(w, soapVersion, http.StatusBadRequest, "soap:Client", "", "Invalid request")
            return true
        }
        if req.ID == 0 && strings.TrimSpace(req.Code) == "" {
            writeSOAPFault(w, soapVersion, http.StatusBadRequest, "soap:Client", "", "id or code is required")
            return true
        }
        customer, found := store.GetCustomer(req.ID, req.Code)
        if !found {
            writeSOAPFault(w, soapVersion, http.StatusNotFound, "soap:Client", "", "Customer not found")
            return true
        }
        payload := customerResponsePayload{
            XMLName: xml.Name{Local: "GetCustomerResponse"},
            Xmlns:   ns,
            Customer: customer,
        }
        writeSOAPResponse(w, soapVersion, http.StatusOK, payload)
        return true

    case "CreateCustomer":
        var req CustomerRequest
        if err := decodeSOAPOperation(body, operation, &req); err != nil {
            writeSOAPFault(w, soapVersion, http.StatusBadRequest, "soap:Client", "", "Invalid request")
            return true
        }
        draft, err := customerDraftFromPayload(req.Customer)
        if err != nil {
            writeSOAPFault(w, soapVersion, http.StatusBadRequest, "soap:Client", "", err.Error())
            return true
        }
        customer, err := store.CreateCustomer(draft)
        if err != nil {
            writeSOAPFault(w, soapVersion, http.StatusBadRequest, "soap:Client", "", err.Error())
            return true
        }
        payload := customerResponsePayload{
            XMLName: xml.Name{Local: "CreateCustomerResponse"},
            Xmlns:   ns,
            Customer: customer,
        }
        writeSOAPResponse(w, soapVersion, http.StatusOK, payload)
        return true

    case "UpdateCustomer":
        var req CustomerRequest
        if err := decodeSOAPOperation(body, operation, &req); err != nil {
            writeSOAPFault(w, soapVersion, http.StatusBadRequest, "soap:Client", "", "Invalid request")
            return true
        }
        if req.Customer.ID == 0 && strings.TrimSpace(req.Customer.Code) == "" {
            writeSOAPFault(w, soapVersion, http.StatusBadRequest, "soap:Client", "", "id or code is required")
            return true
        }
        draft, err := customerDraftFromPayload(req.Customer)
        if err != nil {
            writeSOAPFault(w, soapVersion, http.StatusBadRequest, "soap:Client", "", err.Error())
            return true
        }
        customer, err := store.UpdateCustomer(req.Customer.ID, req.Customer.Code, draft)
        if err != nil {
            writeSOAPFault(w, soapVersion, http.StatusBadRequest, "soap:Client", "", err.Error())
            return true
        }
        payload := customerResponsePayload{
            XMLName: xml.Name{Local: "UpdateCustomerResponse"},
            Xmlns:   ns,
            Customer: customer,
        }
        writeSOAPResponse(w, soapVersion, http.StatusOK, payload)
        return true

    case "GetStocks", "ListStocks":
        var req ListStocksRequest
        if err := decodeSOAPOperation(body, operation, &req); err != nil {
            writeSOAPFault(w, soapVersion, http.StatusBadRequest, "soap:Client", "", "Invalid request")
            return true
        }
        stocks := store.ListStocks(StockFilter{
            Code:         req.Code,
            Status:       req.Status,
            NameContains: req.NameContains,
            MaxResults:   req.MaxResults,
            Offset:       req.Offset,
        })
        responseName := "GetStocksResponse"
        if operation == "ListStocks" {
            responseName = "ListStocksResponse"
        }
        payload := stocksResponsePayload{
            XMLName: xml.Name{Local: responseName},
            Xmlns:   ns,
            Stocks:  stocks,
        }
        writeSOAPResponse(w, soapVersion, http.StatusOK, payload)
        return true

    case "GetStock":
        var req StockLookupRequest
        if err := decodeSOAPOperation(body, operation, &req); err != nil {
            writeSOAPFault(w, soapVersion, http.StatusBadRequest, "soap:Client", "", "Invalid request")
            return true
        }
        if req.ID == 0 && strings.TrimSpace(req.Code) == "" {
            writeSOAPFault(w, soapVersion, http.StatusBadRequest, "soap:Client", "", "id or code is required")
            return true
        }
        stock, found := store.GetStock(req.ID, req.Code)
        if !found {
            writeSOAPFault(w, soapVersion, http.StatusNotFound, "soap:Client", "", "Stock not found")
            return true
        }
        payload := stockResponsePayload{
            XMLName: xml.Name{Local: "GetStockResponse"},
            Xmlns:   ns,
            Stock: stock,
        }
        writeSOAPResponse(w, soapVersion, http.StatusOK, payload)
        return true

    case "CreateStock":
        var req StockRequest
        if err := decodeSOAPOperation(body, operation, &req); err != nil {
            writeSOAPFault(w, soapVersion, http.StatusBadRequest, "soap:Client", "", "Invalid request")
            return true
        }
        draft, err := stockDraftFromPayload(req.Stock)
        if err != nil {
            writeSOAPFault(w, soapVersion, http.StatusBadRequest, "soap:Client", "", err.Error())
            return true
        }
        stock, err := store.CreateStock(draft)
        if err != nil {
            writeSOAPFault(w, soapVersion, http.StatusBadRequest, "soap:Client", "", err.Error())
            return true
        }
        payload := stockResponsePayload{
            XMLName: xml.Name{Local: "CreateStockResponse"},
            Xmlns:   ns,
            Stock: stock,
        }
        writeSOAPResponse(w, soapVersion, http.StatusOK, payload)
        return true

    case "UpdateStock":
        var req StockRequest
        if err := decodeSOAPOperation(body, operation, &req); err != nil {
            writeSOAPFault(w, soapVersion, http.StatusBadRequest, "soap:Client", "", "Invalid request")
            return true
        }
        if req.Stock.ID == 0 && strings.TrimSpace(req.Stock.Code) == "" {
            writeSOAPFault(w, soapVersion, http.StatusBadRequest, "soap:Client", "", "id or code is required")
            return true
        }
        draft, err := stockDraftFromPayload(req.Stock)
        if err != nil {
            writeSOAPFault(w, soapVersion, http.StatusBadRequest, "soap:Client", "", err.Error())
            return true
        }
        stock, err := store.UpdateStock(req.Stock.ID, req.Stock.Code, draft)
        if err != nil {
            writeSOAPFault(w, soapVersion, http.StatusBadRequest, "soap:Client", "", err.Error())
            return true
        }
        payload := stockResponsePayload{
            XMLName: xml.Name{Local: "UpdateStockResponse"},
            Xmlns:   ns,
            Stock: stock,
        }
        writeSOAPResponse(w, soapVersion, http.StatusOK, payload)
        return true

    case "GetCashAccounts", "ListCashAccounts":
        var req ListCashAccountsRequest
        if err := decodeSOAPOperation(body, operation, &req); err != nil {
            writeSOAPFault(w, soapVersion, http.StatusBadRequest, "soap:Client", "", "Invalid request")
            return true
        }
        accounts := store.ListCashAccounts(CashAccountFilter{
            Status:     req.Status,
            Currency:   req.Currency,
            MaxResults: req.MaxResults,
            Offset:     req.Offset,
        })
        responseName := "GetCashAccountsResponse"
        if operation == "ListCashAccounts" {
            responseName = "ListCashAccountsResponse"
        }
        payload := cashAccountsResponsePayload{
            XMLName:      xml.Name{Local: responseName},
            Xmlns:        ns,
            CashAccounts: accounts,
        }
        writeSOAPResponse(w, soapVersion, http.StatusOK, payload)
        return true

    case "GetCashAccount":
        var req CashAccountLookupRequest
        if err := decodeSOAPOperation(body, operation, &req); err != nil {
            writeSOAPFault(w, soapVersion, http.StatusBadRequest, "soap:Client", "", "Invalid request")
            return true
        }
        if req.ID == 0 && strings.TrimSpace(req.Code) == "" {
            writeSOAPFault(w, soapVersion, http.StatusBadRequest, "soap:Client", "", "id or code is required")
            return true
        }
        account, found := store.GetCashAccount(req.ID, req.Code)
        if !found {
            writeSOAPFault(w, soapVersion, http.StatusNotFound, "soap:Client", "", "Cash account not found")
            return true
        }
        payload := cashAccountResponsePayload{
            XMLName: xml.Name{Local: "GetCashAccountResponse"},
            Xmlns:   ns,
            CashAccount: account,
        }
        writeSOAPResponse(w, soapVersion, http.StatusOK, payload)
        return true

    case "CreateCashAccount":
        var req CashAccountRequest
        if err := decodeSOAPOperation(body, operation, &req); err != nil {
            writeSOAPFault(w, soapVersion, http.StatusBadRequest, "soap:Client", "", "Invalid request")
            return true
        }
        draft, err := cashAccountDraftFromPayload(req.CashAccount)
        if err != nil {
            writeSOAPFault(w, soapVersion, http.StatusBadRequest, "soap:Client", "", err.Error())
            return true
        }
        account, err := store.CreateCashAccount(draft)
        if err != nil {
            writeSOAPFault(w, soapVersion, http.StatusBadRequest, "soap:Client", "", err.Error())
            return true
        }
        payload := cashAccountResponsePayload{
            XMLName: xml.Name{Local: "CreateCashAccountResponse"},
            Xmlns:   ns,
            CashAccount: account,
        }
        writeSOAPResponse(w, soapVersion, http.StatusOK, payload)
        return true

    case "UpdateCashAccount":
        var req CashAccountRequest
        if err := decodeSOAPOperation(body, operation, &req); err != nil {
            writeSOAPFault(w, soapVersion, http.StatusBadRequest, "soap:Client", "", "Invalid request")
            return true
        }
        if req.CashAccount.ID == 0 && strings.TrimSpace(req.CashAccount.Code) == "" {
            writeSOAPFault(w, soapVersion, http.StatusBadRequest, "soap:Client", "", "id or code is required")
            return true
        }
        draft, err := cashAccountDraftFromPayload(req.CashAccount)
        if err != nil {
            writeSOAPFault(w, soapVersion, http.StatusBadRequest, "soap:Client", "", err.Error())
            return true
        }
        account, err := store.UpdateCashAccount(req.CashAccount.ID, req.CashAccount.Code, draft)
        if err != nil {
            writeSOAPFault(w, soapVersion, http.StatusBadRequest, "soap:Client", "", err.Error())
            return true
        }
        payload := cashAccountResponsePayload{
            XMLName: xml.Name{Local: "UpdateCashAccountResponse"},
            Xmlns:   ns,
            CashAccount: account,
        }
        writeSOAPResponse(w, soapVersion, http.StatusOK, payload)
        return true

    case "GetOrders", "ListOrders":
        var req ListOrdersRequest
        if err := decodeSOAPOperation(body, operation, &req); err != nil {
            writeSOAPFault(w, soapVersion, http.StatusBadRequest, "soap:Client", "", "Invalid request")
            return true
        }
        orders, err := store.ListOrders(OrderFilter{
            CustomerID: req.CustomerID,
            Status:     req.Status,
            FromDate:   req.FromDate,
            ToDate:     req.ToDate,
            MaxResults: req.MaxResults,
            Offset:     req.Offset,
        })
        if err != nil {
            writeSOAPFault(w, soapVersion, http.StatusBadRequest, "soap:Client", "", err.Error())
            return true
        }
        responseName := "GetOrdersResponse"
        if operation == "ListOrders" {
            responseName = "ListOrdersResponse"
        }
        payload := ordersResponsePayload{
            XMLName: xml.Name{Local: responseName},
            Xmlns:   ns,
            Orders:  orders,
        }
        writeSOAPResponse(w, soapVersion, http.StatusOK, payload)
        return true

    case "GetOrder":
        var req OrderLookupRequest
        if err := decodeSOAPOperation(body, operation, &req); err != nil {
            writeSOAPFault(w, soapVersion, http.StatusBadRequest, "soap:Client", "", "Invalid request")
            return true
        }
        if req.ID == 0 && strings.TrimSpace(req.OrderNo) == "" {
            writeSOAPFault(w, soapVersion, http.StatusBadRequest, "soap:Client", "", "id or orderNo is required")
            return true
        }
        order, found := store.GetOrder(req.ID, req.OrderNo)
        if !found {
            writeSOAPFault(w, soapVersion, http.StatusNotFound, "soap:Client", "", "Order not found")
            return true
        }
        payload := orderResponsePayload{
            XMLName: xml.Name{Local: "GetOrderResponse"},
            Xmlns:   ns,
            Order: order,
        }
        writeSOAPResponse(w, soapVersion, http.StatusOK, payload)
        return true

    case "CreateOrder":
        var req OrderRequest
        if err := decodeSOAPOperation(body, operation, &req); err != nil {
            writeSOAPFault(w, soapVersion, http.StatusBadRequest, "soap:Client", "", "Invalid request")
            return true
        }
        draft, err := orderDraftFromPayload(req.Order)
        if err != nil {
            writeSOAPFault(w, soapVersion, http.StatusBadRequest, "soap:Client", "", err.Error())
            return true
        }
        order, err := store.CreateOrder(draft)
        if err != nil {
            writeSOAPFault(w, soapVersion, http.StatusBadRequest, "soap:Client", "", err.Error())
            return true
        }
        payload := orderResponsePayload{
            XMLName: xml.Name{Local: "CreateOrderResponse"},
            Xmlns:   ns,
            Order: order,
        }
        writeSOAPResponse(w, soapVersion, http.StatusOK, payload)
        return true

    case "UpdateOrder":
        var req OrderRequest
        if err := decodeSOAPOperation(body, operation, &req); err != nil {
            writeSOAPFault(w, soapVersion, http.StatusBadRequest, "soap:Client", "", "Invalid request")
            return true
        }
        if req.Order.ID == 0 && strings.TrimSpace(req.Order.OrderNo) == "" {
            writeSOAPFault(w, soapVersion, http.StatusBadRequest, "soap:Client", "", "id or orderNo is required")
            return true
        }
        draft, err := orderDraftFromPayload(req.Order)
        if err != nil {
            writeSOAPFault(w, soapVersion, http.StatusBadRequest, "soap:Client", "", err.Error())
            return true
        }
        order, err := store.UpdateOrder(req.Order.ID, req.Order.OrderNo, draft)
        if err != nil {
            writeSOAPFault(w, soapVersion, http.StatusBadRequest, "soap:Client", "", err.Error())
            return true
        }
        payload := orderResponsePayload{
            XMLName: xml.Name{Local: "UpdateOrderResponse"},
            Xmlns:   ns,
            Order: order,
        }
        writeSOAPResponse(w, soapVersion, http.StatusOK, payload)
        return true

    case "ApproveOrder":
        var req OrderLookupRequest
        if err := decodeSOAPOperation(body, operation, &req); err != nil {
            writeSOAPFault(w, soapVersion, http.StatusBadRequest, "soap:Client", "", "Invalid request")
            return true
        }
        if req.ID == 0 && strings.TrimSpace(req.OrderNo) == "" {
            writeSOAPFault(w, soapVersion, http.StatusBadRequest, "soap:Client", "", "id or orderNo is required")
            return true
        }
        order, err := store.ApproveOrder(req.ID, req.OrderNo)
        if err != nil {
            writeSOAPFault(w, soapVersion, http.StatusBadRequest, "soap:Client", "", err.Error())
            return true
        }
        payload := orderResponsePayload{
            XMLName: xml.Name{Local: "ApproveOrderResponse"},
            Xmlns:   ns,
            Order: order,
        }
        writeSOAPResponse(w, soapVersion, http.StatusOK, payload)
        return true

    case "CancelOrder":
        var req OrderLookupRequest
        if err := decodeSOAPOperation(body, operation, &req); err != nil {
            writeSOAPFault(w, soapVersion, http.StatusBadRequest, "soap:Client", "", "Invalid request")
            return true
        }
        if req.ID == 0 && strings.TrimSpace(req.OrderNo) == "" {
            writeSOAPFault(w, soapVersion, http.StatusBadRequest, "soap:Client", "", "id or orderNo is required")
            return true
        }
        order, err := store.CancelOrder(req.ID, req.OrderNo)
        if err != nil {
            writeSOAPFault(w, soapVersion, http.StatusBadRequest, "soap:Client", "", err.Error())
            return true
        }
        payload := orderResponsePayload{
            XMLName: xml.Name{Local: "CancelOrderResponse"},
            Xmlns:   ns,
            Order: order,
        }
        writeSOAPResponse(w, soapVersion, http.StatusOK, payload)
        return true

    case "GetInvoices", "ListInvoices":
        var req ListInvoicesRequest
        if err := decodeSOAPOperation(body, operation, &req); err != nil {
            writeSOAPFault(w, soapVersion, http.StatusBadRequest, "soap:Client", "", "Invalid request")
            return true
        }
        invoices, err := store.ListInvoices(InvoiceFilter{
            CustomerID: req.CustomerID,
            Status:     req.Status,
            FromDate:   req.FromDate,
            ToDate:     req.ToDate,
            MaxResults: req.MaxResults,
            Offset:     req.Offset,
        })
        if err != nil {
            writeSOAPFault(w, soapVersion, http.StatusBadRequest, "soap:Client", "", err.Error())
            return true
        }
        responseName := "GetInvoicesResponse"
        if operation == "ListInvoices" {
            responseName = "ListInvoicesResponse"
        }
        payload := invoicesResponsePayload{
            XMLName:  xml.Name{Local: responseName},
            Xmlns:    ns,
            Invoices: invoices,
        }
        writeSOAPResponse(w, soapVersion, http.StatusOK, payload)
        return true

    case "GetInvoice":
        var req InvoiceLookupRequest
        if err := decodeSOAPOperation(body, operation, &req); err != nil {
            writeSOAPFault(w, soapVersion, http.StatusBadRequest, "soap:Client", "", "Invalid request")
            return true
        }
        if req.ID == 0 && strings.TrimSpace(req.InvoiceNo) == "" {
            writeSOAPFault(w, soapVersion, http.StatusBadRequest, "soap:Client", "", "id or invoiceNo is required")
            return true
        }
        invoice, found := store.GetInvoice(req.ID, req.InvoiceNo)
        if !found {
            writeSOAPFault(w, soapVersion, http.StatusNotFound, "soap:Client", "", "Invoice not found")
            return true
        }
        payload := invoiceResponsePayload{
            XMLName: xml.Name{Local: "GetInvoiceResponse"},
            Xmlns:   ns,
            Invoice: invoice,
        }
        writeSOAPResponse(w, soapVersion, http.StatusOK, payload)
        return true

    case "CreateInvoice":
        var req InvoiceRequest
        if err := decodeSOAPOperation(body, operation, &req); err != nil {
            writeSOAPFault(w, soapVersion, http.StatusBadRequest, "soap:Client", "", "Invalid request")
            return true
        }
        draft, err := invoiceDraftFromPayload(req.Invoice)
        if err != nil {
            writeSOAPFault(w, soapVersion, http.StatusBadRequest, "soap:Client", "", err.Error())
            return true
        }
        if draft.OrderID != 0 || strings.TrimSpace(draft.OrderNo) != "" {
            writeSOAPFault(w, soapVersion, http.StatusBadRequest, "soap:Client", "", "Use CreateInvoiceFromOrder for order-linked invoices")
            return true
        }
        invoice, err := store.CreateInvoice(draft)
        if err != nil {
            writeSOAPFault(w, soapVersion, http.StatusBadRequest, "soap:Client", "", err.Error())
            return true
        }
        payload := invoiceResponsePayload{
            XMLName: xml.Name{Local: "CreateInvoiceResponse"},
            Xmlns:   ns,
            Invoice: invoice,
        }
        writeSOAPResponse(w, soapVersion, http.StatusOK, payload)
        return true

    case "CreateInvoiceFromOrder":
        var req CreateInvoiceFromOrderRequest
        if err := decodeSOAPOperation(body, operation, &req); err != nil {
            writeSOAPFault(w, soapVersion, http.StatusBadRequest, "soap:Client", "", "Invalid request")
            return true
        }
        if req.OrderID == 0 && strings.TrimSpace(req.OrderNo) == "" {
            writeSOAPFault(w, soapVersion, http.StatusBadRequest, "soap:Client", "", "orderId or orderNo is required")
            return true
        }
        invoiceDate := req.InvoiceDate
        if strings.TrimSpace(invoiceDate) == "" {
            invoiceDate = time.Now().UTC().Format("2006-01-02")
        }
        invoice, err := store.CreateInvoiceFromOrder(req.OrderID, req.OrderNo, invoiceDate, req.DueDate)
        if err != nil {
            writeSOAPFault(w, soapVersion, http.StatusBadRequest, "soap:Client", "", err.Error())
            return true
        }
        payload := invoiceResponsePayload{
            XMLName: xml.Name{Local: "CreateInvoiceFromOrderResponse"},
            Xmlns:   ns,
            Invoice: invoice,
        }
        writeSOAPResponse(w, soapVersion, http.StatusOK, payload)
        return true

    case "CancelInvoice":
        var req InvoiceLookupRequest
        if err := decodeSOAPOperation(body, operation, &req); err != nil {
            writeSOAPFault(w, soapVersion, http.StatusBadRequest, "soap:Client", "", "Invalid request")
            return true
        }
        if req.ID == 0 && strings.TrimSpace(req.InvoiceNo) == "" {
            writeSOAPFault(w, soapVersion, http.StatusBadRequest, "soap:Client", "", "id or invoiceNo is required")
            return true
        }
        invoice, err := store.CancelInvoice(req.ID, req.InvoiceNo)
        if err != nil {
            writeSOAPFault(w, soapVersion, http.StatusBadRequest, "soap:Client", "", err.Error())
            return true
        }
        payload := invoiceResponsePayload{
            XMLName: xml.Name{Local: "CancelInvoiceResponse"},
            Xmlns:   ns,
            Invoice: invoice,
        }
        writeSOAPResponse(w, soapVersion, http.StatusOK, payload)
        return true

    case "GetCashTransactions", "ListCashTransactions":
        var req ListCashTransactionsRequest
        if err := decodeSOAPOperation(body, operation, &req); err != nil {
            writeSOAPFault(w, soapVersion, http.StatusBadRequest, "soap:Client", "", "Invalid request")
            return true
        }
        txns, err := store.ListCashTransactions(CashTransactionFilter{
            CashAccountID: req.CashAccountID,
            Type:          req.Type,
            FromDate:      req.FromDate,
            ToDate:        req.ToDate,
            MaxResults:    req.MaxResults,
            Offset:        req.Offset,
        })
        if err != nil {
            writeSOAPFault(w, soapVersion, http.StatusBadRequest, "soap:Client", "", err.Error())
            return true
        }
        responseName := "GetCashTransactionsResponse"
        if operation == "ListCashTransactions" {
            responseName = "ListCashTransactionsResponse"
        }
        payload := cashTransactionsResponsePayload{
            XMLName:         xml.Name{Local: responseName},
            Xmlns:           ns,
            CashTransactions: txns,
        }
        writeSOAPResponse(w, soapVersion, http.StatusOK, payload)
        return true

    case "GetCashTransaction":
        var req CashTransactionLookupRequest
        if err := decodeSOAPOperation(body, operation, &req); err != nil {
            writeSOAPFault(w, soapVersion, http.StatusBadRequest, "soap:Client", "", "Invalid request")
            return true
        }
        if req.ID == 0 && strings.TrimSpace(req.TxnNo) == "" {
            writeSOAPFault(w, soapVersion, http.StatusBadRequest, "soap:Client", "", "id or txnNo is required")
            return true
        }
        txn, found := store.GetCashTransaction(req.ID, req.TxnNo)
        if !found {
            writeSOAPFault(w, soapVersion, http.StatusNotFound, "soap:Client", "", "Cash transaction not found")
            return true
        }
        payload := cashTransactionResponsePayload{
            XMLName: xml.Name{Local: "GetCashTransactionResponse"},
            Xmlns:   ns,
            CashTransaction: txn,
        }
        writeSOAPResponse(w, soapVersion, http.StatusOK, payload)
        return true

    case "CreateCashTransaction":
        var req CashTransactionRequest
        if err := decodeSOAPOperation(body, operation, &req); err != nil {
            writeSOAPFault(w, soapVersion, http.StatusBadRequest, "soap:Client", "", "Invalid request")
            return true
        }
        draft, err := cashTransactionDraftFromPayload(req.CashTransaction)
        if err != nil {
            writeSOAPFault(w, soapVersion, http.StatusBadRequest, "soap:Client", "", err.Error())
            return true
        }
        txn, err := store.CreateCashTransaction(draft)
        if err != nil {
            writeSOAPFault(w, soapVersion, http.StatusBadRequest, "soap:Client", "", err.Error())
            return true
        }
        payload := cashTransactionResponsePayload{
            XMLName: xml.Name{Local: "CreateCashTransactionResponse"},
            Xmlns:   ns,
            CashTransaction: txn,
        }
        writeSOAPResponse(w, soapVersion, http.StatusOK, payload)
        return true

    case "ReverseCashTransaction":
        var req CashTransactionLookupRequest
        if err := decodeSOAPOperation(body, operation, &req); err != nil {
            writeSOAPFault(w, soapVersion, http.StatusBadRequest, "soap:Client", "", "Invalid request")
            return true
        }
        if req.ID == 0 && strings.TrimSpace(req.TxnNo) == "" {
            writeSOAPFault(w, soapVersion, http.StatusBadRequest, "soap:Client", "", "id or txnNo is required")
            return true
        }
        txn, err := store.ReverseCashTransaction(req.ID, req.TxnNo)
        if err != nil {
            writeSOAPFault(w, soapVersion, http.StatusBadRequest, "soap:Client", "", err.Error())
            return true
        }
        payload := cashTransactionResponsePayload{
            XMLName: xml.Name{Local: "ReverseCashTransactionResponse"},
            Xmlns:   ns,
            CashTransaction: txn,
        }
        writeSOAPResponse(w, soapVersion, http.StatusOK, payload)
        return true
    }

    return false
}

func handlePublicOperation(w http.ResponseWriter, soapVersion, operation string, store *Store, ns string, body []byte) bool {
    switch operation {
    case "GetCountries":
        var req CountriesRequest
        if err := decodeSOAPOperation(body, operation, &req); err != nil {
            writeSOAPFault(w, soapVersion, http.StatusBadRequest, "soap:Client", "", "Invalid request")
            return true
        }
        countries := store.ListCountries(req.Continent)
        payload := countriesResponsePayload{
            XMLName:  xml.Name{Local: "GetCountriesResponse"},
            Xmlns:    ns,
            Countries: countries,
        }
        writeSOAPResponse(w, soapVersion, http.StatusOK, payload)
        return true

    case "GetCities":
        var req CitiesRequest
        if err := decodeSOAPOperation(body, operation, &req); err != nil {
            writeSOAPFault(w, soapVersion, http.StatusBadRequest, "soap:Client", "", "Invalid request")
            return true
        }
        cities, err := store.ListCities(req.CountryCode)
        if err != nil {
            writeSOAPFault(w, soapVersion, http.StatusBadRequest, "soap:Client", "", err.Error())
            return true
        }
        payload := citiesResponsePayload{
            XMLName: xml.Name{Local: "GetCitiesResponse"},
            Xmlns:   ns,
            Cities:  cities,
        }
        writeSOAPResponse(w, soapVersion, http.StatusOK, payload)
        return true
    }
    return false
}

func writeSOAPResponse(w http.ResponseWriter, soapVersion string, status int, payload interface{}) {
    body, err := soapEnvelopeFromPayload(soapVersion, payload)
    if err != nil {
        writeSOAPFault(w, soapVersion, http.StatusInternalServerError, "soap:Server", "", "Failed to build response")
        return
    }
    writeSOAPXML(w, soapVersion, status, body)
}

func soapEnvelopeFromPayload(version string, payload interface{}) (string, error) {
    output, err := xml.Marshal(payload)
    if err != nil {
        return "", err
    }
    body := "        " + string(output)
    return soapEnvelope(version, body), nil
}

func customerDraftFromPayload(payload CustomerPayload) (CustomerDraft, error) {
    riskLimit, err := parseFloat(payload.RiskLimit)
    if err != nil {
        return CustomerDraft{}, fmt.Errorf("invalid riskLimit")
    }
    return CustomerDraft{
        Code:      payload.Code,
        Name:      payload.Name,
        TaxNumber: payload.TaxNumber,
        Email:     payload.Email,
        Phone:     payload.Phone,
        Address:   payload.Address,
        Currency:  payload.Currency,
        RiskLimit: riskLimit,
        Status:    payload.Status,
    }, nil
}

func stockDraftFromPayload(payload StockPayload) (StockDraft, error) {
    vatRate, err := parseFloat(payload.VatRate)
    if err != nil {
        return StockDraft{}, fmt.Errorf("invalid vatRate")
    }
    price, err := parseFloat(payload.Price)
    if err != nil {
        return StockDraft{}, fmt.Errorf("invalid price")
    }
    stockOnHand, err := parseFloat(payload.StockOnHand)
    if err != nil {
        return StockDraft{}, fmt.Errorf("invalid stockOnHand")
    }
    minStock, err := parseFloat(payload.MinStock)
    if err != nil {
        return StockDraft{}, fmt.Errorf("invalid minStock")
    }
    return StockDraft{
        Code:        payload.Code,
        Name:        payload.Name,
        Unit:        payload.Unit,
        VatRate:     vatRate,
        Price:       price,
        StockOnHand: stockOnHand,
        MinStock:    minStock,
        Status:      payload.Status,
    }, nil
}

func cashAccountDraftFromPayload(payload CashAccountPayload) (CashAccountDraft, error) {
    balance, err := parseFloat(payload.Balance)
    if err != nil {
        return CashAccountDraft{}, fmt.Errorf("invalid balance")
    }
    return CashAccountDraft{
        Code:     payload.Code,
        Name:     payload.Name,
        Currency: payload.Currency,
        Balance:  balance,
        Status:   payload.Status,
    }, nil
}

func orderDraftFromPayload(payload OrderPayload) (OrderDraft, error) {
    lines, err := orderLinesFromPayload(payload.Lines)
    if err != nil {
        return OrderDraft{}, err
    }
    return OrderDraft{
        OrderNo:      payload.OrderNo,
        CustomerID:   payload.CustomerID,
        CustomerCode: payload.CustomerCode,
        Status:       payload.Status,
        OrderDate:    payload.OrderDate,
        Currency:     payload.Currency,
        Notes:        payload.Notes,
        Lines:        lines,
    }, nil
}

func orderLinesFromPayload(lines []OrderLinePayload) ([]OrderLineDraft, error) {
    if len(lines) == 0 {
        return nil, fmt.Errorf("order lines are required")
    }
    var result []OrderLineDraft
    for idx, line := range lines {
        quantity, err := parseFloat(line.Quantity)
        if err != nil {
            return nil, fmt.Errorf("line %d: invalid quantity", idx+1)
        }
        unitPrice, err := parseFloat(line.UnitPrice)
        if err != nil {
            return nil, fmt.Errorf("line %d: invalid unitPrice", idx+1)
        }
        discountRate, err := parseFloat(line.DiscountRate)
        if err != nil {
            return nil, fmt.Errorf("line %d: invalid discountRate", idx+1)
        }
        taxRate, err := parseFloat(line.TaxRate)
        if err != nil {
            return nil, fmt.Errorf("line %d: invalid taxRate", idx+1)
        }
        result = append(result, OrderLineDraft{
            StockID:      line.StockID,
            StockCode:    line.StockCode,
            Quantity:     quantity,
            UnitPrice:    unitPrice,
            DiscountRate: discountRate,
            TaxRate:      taxRate,
        })
    }
    return result, nil
}

func invoiceDraftFromPayload(payload InvoicePayload) (InvoiceDraft, error) {
    lines, err := invoiceLinesFromPayload(payload.Lines)
    if err != nil {
        return InvoiceDraft{}, err
    }
    return InvoiceDraft{
        InvoiceNo:    payload.InvoiceNo,
        OrderID:      payload.OrderID,
        OrderNo:      payload.OrderNo,
        CustomerID:   payload.CustomerID,
        CustomerCode: payload.CustomerCode,
        InvoiceDate:  payload.InvoiceDate,
        DueDate:      payload.DueDate,
        Currency:     payload.Currency,
        Lines:        lines,
    }, nil
}

func invoiceLinesFromPayload(lines []InvoiceLinePayload) ([]InvoiceLineDraft, error) {
    if len(lines) == 0 {
        return nil, fmt.Errorf("invoice lines are required")
    }
    var result []InvoiceLineDraft
    for idx, line := range lines {
        quantity, err := parseFloat(line.Quantity)
        if err != nil {
            return nil, fmt.Errorf("line %d: invalid quantity", idx+1)
        }
        unitPrice, err := parseFloat(line.UnitPrice)
        if err != nil {
            return nil, fmt.Errorf("line %d: invalid unitPrice", idx+1)
        }
        discountRate, err := parseFloat(line.DiscountRate)
        if err != nil {
            return nil, fmt.Errorf("line %d: invalid discountRate", idx+1)
        }
        taxRate, err := parseFloat(line.TaxRate)
        if err != nil {
            return nil, fmt.Errorf("line %d: invalid taxRate", idx+1)
        }
        result = append(result, InvoiceLineDraft{
            StockID:      line.StockID,
            StockCode:    line.StockCode,
            Quantity:     quantity,
            UnitPrice:    unitPrice,
            DiscountRate: discountRate,
            TaxRate:      taxRate,
        })
    }
    return result, nil
}

func cashTransactionDraftFromPayload(payload CashTransactionPayload) (CashTransactionDraft, error) {
    amount, err := parseFloat(payload.Amount)
    if err != nil {
        return CashTransactionDraft{}, fmt.Errorf("invalid amount")
    }
    return CashTransactionDraft{
        TxnNo:           payload.TxnNo,
        CashAccountID:   payload.CashAccountID,
        CashAccountCode: payload.CashAccountCode,
        CustomerID:      payload.CustomerID,
        CustomerCode:    payload.CustomerCode,
        InvoiceID:       payload.InvoiceID,
        InvoiceNo:       payload.InvoiceNo,
        Type:            payload.Type,
        Amount:          amount,
        Currency:        payload.Currency,
        Method:          payload.Method,
        TransactionDate: payload.TransactionDate,
        Description:     payload.Description,
    }, nil
}

func parseFloat(value string) (float64, error) {
    trimmed := strings.TrimSpace(value)
    if trimmed == "" {
        return 0, nil
    }
    parsed, err := strconv.ParseFloat(trimmed, 64)
    if err != nil {
        return 0, err
    }
    return parsed, nil
}

