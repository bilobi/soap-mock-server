package main

import (
    "encoding/json"
    "errors"
    "fmt"
    "math"
    "os"
    "path/filepath"
    "sort"
    "strings"
    "sync"
    "time"
)

type Store struct {
    mu   sync.RWMutex
    path string
    data StoreData
}

type StoreData struct {
    Meta             StoreMeta          `json:"meta"`
    Customers        []Customer         `json:"customers"`
    Stocks           []StockItem        `json:"stocks"`
    CashAccounts     []CashAccount      `json:"cashAccounts"`
    Orders           []Order            `json:"orders"`
    Invoices         []Invoice          `json:"invoices"`
    CashTransactions []CashTransaction  `json:"cashTransactions"`
    Countries        []Country          `json:"countries"`
    Cities           []City             `json:"cities"`
}

type StoreMeta struct {
    NextCustomerID        int `json:"nextCustomerId"`
    NextStockID           int `json:"nextStockId"`
    NextCashAccountID     int `json:"nextCashAccountId"`
    NextOrderID           int `json:"nextOrderId"`
    NextInvoiceID         int `json:"nextInvoiceId"`
    NextCashTransactionID int `json:"nextCashTransactionId"`
}

type Customer struct {
    ID        int     `json:"id" xml:"id"`
    Code      string  `json:"code" xml:"code"`
    Name      string  `json:"name" xml:"name"`
    TaxNumber string  `json:"taxNumber" xml:"taxNumber"`
    Email     string  `json:"email" xml:"email"`
    Phone     string  `json:"phone" xml:"phone"`
    Address   string  `json:"address" xml:"address"`
    Currency  string  `json:"currency" xml:"currency"`
    RiskLimit float64 `json:"riskLimit" xml:"riskLimit"`
    Status    string  `json:"status" xml:"status"`
    CreatedAt string  `json:"createdAt" xml:"createdAt"`
    UpdatedAt string  `json:"updatedAt" xml:"updatedAt"`
}

type StockItem struct {
    ID          int     `json:"id" xml:"id"`
    Code        string  `json:"code" xml:"code"`
    Name        string  `json:"name" xml:"name"`
    Unit        string  `json:"unit" xml:"unit"`
    VatRate     float64 `json:"vatRate" xml:"vatRate"`
    Price       float64 `json:"price" xml:"price"`
    StockOnHand float64 `json:"stockOnHand" xml:"stockOnHand"`
    MinStock    float64 `json:"minStock" xml:"minStock"`
    Status      string  `json:"status" xml:"status"`
    CreatedAt   string  `json:"createdAt" xml:"createdAt"`
    UpdatedAt   string  `json:"updatedAt" xml:"updatedAt"`
}

type CashAccount struct {
    ID        int     `json:"id" xml:"id"`
    Code      string  `json:"code" xml:"code"`
    Name      string  `json:"name" xml:"name"`
    Currency  string  `json:"currency" xml:"currency"`
    Balance   float64 `json:"balance" xml:"balance"`
    Status    string  `json:"status" xml:"status"`
    CreatedAt string  `json:"createdAt" xml:"createdAt"`
    UpdatedAt string  `json:"updatedAt" xml:"updatedAt"`
}

type Order struct {
    ID            int         `json:"id" xml:"id"`
    OrderNo       string      `json:"orderNo" xml:"orderNo"`
    CustomerID    int         `json:"customerId" xml:"customerId"`
    Status        string      `json:"status" xml:"status"`
    OrderDate     string      `json:"orderDate" xml:"orderDate"`
    Currency      string      `json:"currency" xml:"currency"`
    Notes         string      `json:"notes" xml:"notes"`
    Lines         []OrderLine `json:"lines" xml:"lines>line"`
    Subtotal      float64     `json:"subtotal" xml:"subtotal"`
    DiscountTotal float64     `json:"discountTotal" xml:"discountTotal"`
    TaxTotal      float64     `json:"taxTotal" xml:"taxTotal"`
    GrandTotal    float64     `json:"grandTotal" xml:"grandTotal"`
    CreatedAt     string      `json:"createdAt" xml:"createdAt"`
    UpdatedAt     string      `json:"updatedAt" xml:"updatedAt"`
}

type OrderLine struct {
    LineNo       int     `json:"lineNo" xml:"lineNo"`
    StockID      int     `json:"stockId" xml:"stockId"`
    StockCode    string  `json:"stockCode" xml:"stockCode"`
    Description  string  `json:"description" xml:"description"`
    Unit         string  `json:"unit" xml:"unit"`
    Quantity     float64 `json:"quantity" xml:"quantity"`
    UnitPrice    float64 `json:"unitPrice" xml:"unitPrice"`
    DiscountRate float64 `json:"discountRate" xml:"discountRate"`
    TaxRate      float64 `json:"taxRate" xml:"taxRate"`
    LineTotal    float64 `json:"lineTotal" xml:"lineTotal"`
}

type Invoice struct {
    ID            int           `json:"id" xml:"id"`
    InvoiceNo     string        `json:"invoiceNo" xml:"invoiceNo"`
    OrderID       int           `json:"orderId" xml:"orderId"`
    CustomerID    int           `json:"customerId" xml:"customerId"`
    Status        string        `json:"status" xml:"status"`
    InvoiceDate   string        `json:"invoiceDate" xml:"invoiceDate"`
    DueDate       string        `json:"dueDate" xml:"dueDate"`
    Currency      string        `json:"currency" xml:"currency"`
    Lines         []InvoiceLine `json:"lines" xml:"lines>line"`
    Subtotal      float64       `json:"subtotal" xml:"subtotal"`
    DiscountTotal float64       `json:"discountTotal" xml:"discountTotal"`
    TaxTotal      float64       `json:"taxTotal" xml:"taxTotal"`
    GrandTotal    float64       `json:"grandTotal" xml:"grandTotal"`
    PaidAmount    float64       `json:"paidAmount" xml:"paidAmount"`
    CreatedAt     string        `json:"createdAt" xml:"createdAt"`
    UpdatedAt     string        `json:"updatedAt" xml:"updatedAt"`
}

type InvoiceLine struct {
    LineNo       int     `json:"lineNo" xml:"lineNo"`
    StockID      int     `json:"stockId" xml:"stockId"`
    StockCode    string  `json:"stockCode" xml:"stockCode"`
    Description  string  `json:"description" xml:"description"`
    Unit         string  `json:"unit" xml:"unit"`
    Quantity     float64 `json:"quantity" xml:"quantity"`
    UnitPrice    float64 `json:"unitPrice" xml:"unitPrice"`
    DiscountRate float64 `json:"discountRate" xml:"discountRate"`
    TaxRate      float64 `json:"taxRate" xml:"taxRate"`
    LineTotal    float64 `json:"lineTotal" xml:"lineTotal"`
}

type CashTransaction struct {
    ID              int     `json:"id" xml:"id"`
    TxnNo           string  `json:"txnNo" xml:"txnNo"`
    CashAccountID   int     `json:"cashAccountId" xml:"cashAccountId"`
    CustomerID      int     `json:"customerId" xml:"customerId"`
    InvoiceID       int     `json:"invoiceId" xml:"invoiceId"`
    Type            string  `json:"type" xml:"type"`
    Amount          float64 `json:"amount" xml:"amount"`
    Currency        string  `json:"currency" xml:"currency"`
    Method          string  `json:"method" xml:"method"`
    TransactionDate string  `json:"transactionDate" xml:"transactionDate"`
    Description     string  `json:"description" xml:"description"`
    Status          string  `json:"status" xml:"status"`
    ReversalOf      int     `json:"reversalOf" xml:"reversalOf"`
    CreatedAt       string  `json:"createdAt" xml:"createdAt"`
}

type Country struct {
    Code      string `json:"code" xml:"code"`
    Name      string `json:"name" xml:"name"`
    Continent string `json:"continent" xml:"continent"`
}

type City struct {
    CountryCode string `json:"countryCode" xml:"countryCode"`
    Name        string `json:"name" xml:"name"`
    Population  int    `json:"population" xml:"population"`
}

type CustomerFilter struct {
    Code       string
    Status     string
    TaxNumber  string
    MaxResults int
    Offset     int
}

type StockFilter struct {
    Code         string
    Status       string
    NameContains string
    MaxResults   int
    Offset       int
}

type CashAccountFilter struct {
    Status     string
    Currency   string
    MaxResults int
    Offset     int
}

type OrderFilter struct {
    CustomerID int
    Status     string
    FromDate   string
    ToDate     string
    MaxResults int
    Offset     int
}

type InvoiceFilter struct {
    CustomerID int
    Status     string
    FromDate   string
    ToDate     string
    MaxResults int
    Offset     int
}

type CashTransactionFilter struct {
    CashAccountID int
    Type          string
    FromDate      string
    ToDate        string
    MaxResults    int
    Offset        int
}

type CustomerDraft struct {
    Code      string
    Name      string
    TaxNumber string
    Email     string
    Phone     string
    Address   string
    Currency  string
    RiskLimit float64
    Status    string
}

type StockDraft struct {
    Code        string
    Name        string
    Unit        string
    VatRate     float64
    Price       float64
    StockOnHand float64
    MinStock    float64
    Status      string
}

type CashAccountDraft struct {
    Code     string
    Name     string
    Currency string
    Balance  float64
    Status   string
}

type OrderLineDraft struct {
    StockID      int
    StockCode    string
    Quantity     float64
    UnitPrice    float64
    DiscountRate float64
    TaxRate      float64
}

type OrderDraft struct {
    OrderNo     string
    CustomerID  int
    CustomerCode string
    Status      string
    OrderDate   string
    Currency    string
    Notes       string
    Lines       []OrderLineDraft
}

type InvoiceLineDraft struct {
    StockID      int
    StockCode    string
    Quantity     float64
    UnitPrice    float64
    DiscountRate float64
    TaxRate      float64
}

type InvoiceDraft struct {
    InvoiceNo   string
    OrderID     int
    OrderNo     string
    CustomerID  int
    CustomerCode string
    InvoiceDate string
    DueDate     string
    Currency    string
    Lines       []InvoiceLineDraft
}

type CashTransactionDraft struct {
    TxnNo            string
    CashAccountID    int
    CashAccountCode  string
    CustomerID       int
    CustomerCode     string
    InvoiceID        int
    InvoiceNo        string
    Type             string
    Amount           float64
    Currency         string
    Method           string
    TransactionDate  string
    Description      string
}

func loadStore(path string) (*Store, error) {
    store := &Store{path: path}
    if err := store.load(); err != nil {
        return nil, err
    }
    return store, nil
}

func (s *Store) load() error {
    payload, err := os.ReadFile(s.path)
    if err != nil {
        if errors.Is(err, os.ErrNotExist) {
            s.data = defaultSeedData()
            return s.saveUnlocked()
        }
        return err
    }
    if len(strings.TrimSpace(string(payload))) == 0 {
        s.data = defaultSeedData()
        return s.saveUnlocked()
    }
    if err := json.Unmarshal(payload, &s.data); err != nil {
        return fmt.Errorf("failed to parse %s: %w", s.path, err)
    }
    s.normalize()
    return nil
}

func (s *Store) saveUnlocked() error {
    dir := filepath.Dir(s.path)
    if err := os.MkdirAll(dir, 0o755); err != nil {
        return err
    }
    tempFile, err := os.CreateTemp(dir, "store-*.json")
    if err != nil {
        return err
    }
    defer os.Remove(tempFile.Name())

    encoder := json.NewEncoder(tempFile)
    encoder.SetIndent("", "  ")
    if err := encoder.Encode(s.data); err != nil {
        _ = tempFile.Close()
        return err
    }
    if err := tempFile.Close(); err != nil {
        return err
    }
    return os.Rename(tempFile.Name(), s.path)
}

func (s *Store) normalize() {
    for i := range s.data.Customers {
        s.data.Customers[i].Code = strings.ToUpper(strings.TrimSpace(s.data.Customers[i].Code))
        s.data.Customers[i].Status = normalizeStatus(s.data.Customers[i].Status, "ACTIVE")
        s.data.Customers[i].Currency = normalizeCurrency(s.data.Customers[i].Currency)
    }
    for i := range s.data.Stocks {
        s.data.Stocks[i].Code = strings.ToUpper(strings.TrimSpace(s.data.Stocks[i].Code))
        s.data.Stocks[i].Status = normalizeStatus(s.data.Stocks[i].Status, "ACTIVE")
        s.data.Stocks[i].VatRate = roundMoney(s.data.Stocks[i].VatRate)
    }
    for i := range s.data.CashAccounts {
        s.data.CashAccounts[i].Code = strings.ToUpper(strings.TrimSpace(s.data.CashAccounts[i].Code))
        s.data.CashAccounts[i].Status = normalizeStatus(s.data.CashAccounts[i].Status, "OPEN")
        s.data.CashAccounts[i].Currency = normalizeCurrency(s.data.CashAccounts[i].Currency)
        s.data.CashAccounts[i].Balance = roundMoney(s.data.CashAccounts[i].Balance)
    }
    for i := range s.data.Orders {
        s.data.Orders[i].Status = normalizeStatus(s.data.Orders[i].Status, "DRAFT")
        recalcOrderTotals(&s.data.Orders[i])
    }
    for i := range s.data.Invoices {
        s.data.Invoices[i].Status = normalizeStatus(s.data.Invoices[i].Status, "ISSUED")
        recalcInvoiceTotals(&s.data.Invoices[i])
        s.data.Invoices[i].PaidAmount = roundMoney(s.data.Invoices[i].PaidAmount)
        if s.data.Invoices[i].PaidAmount >= s.data.Invoices[i].GrandTotal && s.data.Invoices[i].Status != "CANCELLED" {
            s.data.Invoices[i].Status = "PAID"
        }
    }
    for i := range s.data.CashTransactions {
        if s.data.CashTransactions[i].Status == "" {
            s.data.CashTransactions[i].Status = "POSTED"
        }
        s.data.CashTransactions[i].Type = normalizeStatus(s.data.CashTransactions[i].Type, "")
    }
    s.normalizeMeta()
}

func (s *Store) normalizeMeta() {
    maxCustomer := maxCustomerID(s.data.Customers)
    if s.data.Meta.NextCustomerID <= maxCustomer {
        s.data.Meta.NextCustomerID = maxCustomer + 1
    }
    maxStock := maxStockID(s.data.Stocks)
    if s.data.Meta.NextStockID <= maxStock {
        s.data.Meta.NextStockID = maxStock + 1
    }
    maxCash := maxCashAccountID(s.data.CashAccounts)
    if s.data.Meta.NextCashAccountID <= maxCash {
        s.data.Meta.NextCashAccountID = maxCash + 1
    }
    maxOrder := maxOrderID(s.data.Orders)
    if s.data.Meta.NextOrderID <= maxOrder {
        s.data.Meta.NextOrderID = maxOrder + 1
    }
    maxInvoice := maxInvoiceID(s.data.Invoices)
    if s.data.Meta.NextInvoiceID <= maxInvoice {
        s.data.Meta.NextInvoiceID = maxInvoice + 1
    }
    maxTxn := maxCashTransactionID(s.data.CashTransactions)
    if s.data.Meta.NextCashTransactionID <= maxTxn {
        s.data.Meta.NextCashTransactionID = maxTxn + 1
    }
}

func (s *Store) ListCustomers(filter CustomerFilter) []Customer {
    s.mu.RLock()
    defer s.mu.RUnlock()
    var result []Customer
    for _, customer := range s.data.Customers {
        if filter.Code != "" && !strings.EqualFold(customer.Code, filter.Code) {
            continue
        }
        if filter.Status != "" && !strings.EqualFold(customer.Status, filter.Status) {
            continue
        }
        if filter.TaxNumber != "" && customer.TaxNumber != filter.TaxNumber {
            continue
        }
        result = append(result, customer)
    }
    return paginateCustomers(result, filter.Offset, filter.MaxResults)
}

func (s *Store) GetCustomer(id int, code string) (Customer, bool) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    if id != 0 {
        for _, customer := range s.data.Customers {
            if customer.ID == id {
                return customer, true
            }
        }
    }
    if code != "" {
        for _, customer := range s.data.Customers {
            if strings.EqualFold(customer.Code, code) {
                return customer, true
            }
        }
    }
    return Customer{}, false
}

func (s *Store) CreateCustomer(draft CustomerDraft) (Customer, error) {
    s.mu.Lock()
    defer s.mu.Unlock()

    if err := validateCustomerDraft(draft); err != nil {
        return Customer{}, err
    }

    code := strings.ToUpper(strings.TrimSpace(draft.Code))
    if s.customerCodeExistsLocked(code, 0) {
        return Customer{}, fmt.Errorf("customer code already exists")
    }

    now := time.Now().UTC().Format(time.RFC3339)
    customer := Customer{
        ID:        s.data.Meta.NextCustomerID,
        Code:      code,
        Name:      strings.TrimSpace(draft.Name),
        TaxNumber: strings.TrimSpace(draft.TaxNumber),
        Email:     strings.TrimSpace(draft.Email),
        Phone:     strings.TrimSpace(draft.Phone),
        Address:   strings.TrimSpace(draft.Address),
        Currency:  normalizeCurrency(draft.Currency),
        RiskLimit: roundMoney(draft.RiskLimit),
        Status:    normalizeStatus(draft.Status, "ACTIVE"),
        CreatedAt: now,
        UpdatedAt: now,
    }

    s.data.Customers = append(s.data.Customers, customer)
    s.data.Meta.NextCustomerID++
    if err := s.saveUnlocked(); err != nil {
        return Customer{}, err
    }
    return customer, nil
}

func (s *Store) UpdateCustomer(id int, code string, draft CustomerDraft) (Customer, error) {
    s.mu.Lock()
    defer s.mu.Unlock()

    idx := s.findCustomerIndexLocked(id, code)
    if idx == -1 {
        return Customer{}, fmt.Errorf("customer not found")
    }

    if err := validateCustomerDraft(draft); err != nil {
        return Customer{}, err
    }

    normalizedCode := strings.ToUpper(strings.TrimSpace(draft.Code))
    if normalizedCode == "" {
        normalizedCode = s.data.Customers[idx].Code
    }
    if s.customerCodeExistsLocked(normalizedCode, s.data.Customers[idx].ID) {
        return Customer{}, fmt.Errorf("customer code already exists")
    }

    customer := s.data.Customers[idx]
    customer.Code = normalizedCode
    customer.Name = strings.TrimSpace(draft.Name)
    customer.TaxNumber = strings.TrimSpace(draft.TaxNumber)
    customer.Email = strings.TrimSpace(draft.Email)
    customer.Phone = strings.TrimSpace(draft.Phone)
    customer.Address = strings.TrimSpace(draft.Address)
    customer.Currency = normalizeCurrency(draft.Currency)
    customer.RiskLimit = roundMoney(draft.RiskLimit)
    customer.Status = normalizeStatus(draft.Status, customer.Status)
    customer.UpdatedAt = time.Now().UTC().Format(time.RFC3339)

    s.data.Customers[idx] = customer
    if err := s.saveUnlocked(); err != nil {
        return Customer{}, err
    }
    return customer, nil
}

func (s *Store) ListStocks(filter StockFilter) []StockItem {
    s.mu.RLock()
    defer s.mu.RUnlock()
    var result []StockItem
    for _, stock := range s.data.Stocks {
        if filter.Code != "" && !strings.EqualFold(stock.Code, filter.Code) {
            continue
        }
        if filter.Status != "" && !strings.EqualFold(stock.Status, filter.Status) {
            continue
        }
        if filter.NameContains != "" && !strings.Contains(strings.ToLower(stock.Name), strings.ToLower(filter.NameContains)) {
            continue
        }
        result = append(result, stock)
    }
    return paginateStocks(result, filter.Offset, filter.MaxResults)
}

func (s *Store) GetStock(id int, code string) (StockItem, bool) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    if id != 0 {
        for _, stock := range s.data.Stocks {
            if stock.ID == id {
                return stock, true
            }
        }
    }
    if code != "" {
        for _, stock := range s.data.Stocks {
            if strings.EqualFold(stock.Code, code) {
                return stock, true
            }
        }
    }
    return StockItem{}, false
}

func (s *Store) CreateStock(draft StockDraft) (StockItem, error) {
    s.mu.Lock()
    defer s.mu.Unlock()

    if err := validateStockDraft(draft); err != nil {
        return StockItem{}, err
    }

    code := strings.ToUpper(strings.TrimSpace(draft.Code))
    if s.stockCodeExistsLocked(code, 0) {
        return StockItem{}, fmt.Errorf("stock code already exists")
    }

    now := time.Now().UTC().Format(time.RFC3339)
    stock := StockItem{
        ID:          s.data.Meta.NextStockID,
        Code:        code,
        Name:        strings.TrimSpace(draft.Name),
        Unit:        strings.TrimSpace(draft.Unit),
        VatRate:     roundMoney(draft.VatRate),
        Price:       roundMoney(draft.Price),
        StockOnHand: roundQuantity(draft.StockOnHand),
        MinStock:    roundQuantity(draft.MinStock),
        Status:      normalizeStatus(draft.Status, "ACTIVE"),
        CreatedAt:   now,
        UpdatedAt:   now,
    }

    s.data.Stocks = append(s.data.Stocks, stock)
    s.data.Meta.NextStockID++
    if err := s.saveUnlocked(); err != nil {
        return StockItem{}, err
    }
    return stock, nil
}

func (s *Store) UpdateStock(id int, code string, draft StockDraft) (StockItem, error) {
    s.mu.Lock()
    defer s.mu.Unlock()

    idx := s.findStockIndexLocked(id, code)
    if idx == -1 {
        return StockItem{}, fmt.Errorf("stock not found")
    }

    if err := validateStockDraft(draft); err != nil {
        return StockItem{}, err
    }

    normalizedCode := strings.ToUpper(strings.TrimSpace(draft.Code))
    if normalizedCode == "" {
        normalizedCode = s.data.Stocks[idx].Code
    }
    if s.stockCodeExistsLocked(normalizedCode, s.data.Stocks[idx].ID) {
        return StockItem{}, fmt.Errorf("stock code already exists")
    }

    stock := s.data.Stocks[idx]
    stock.Code = normalizedCode
    stock.Name = strings.TrimSpace(draft.Name)
    stock.Unit = strings.TrimSpace(draft.Unit)
    stock.VatRate = roundMoney(draft.VatRate)
    stock.Price = roundMoney(draft.Price)
    stock.StockOnHand = roundQuantity(draft.StockOnHand)
    stock.MinStock = roundQuantity(draft.MinStock)
    stock.Status = normalizeStatus(draft.Status, stock.Status)
    stock.UpdatedAt = time.Now().UTC().Format(time.RFC3339)

    s.data.Stocks[idx] = stock
    if err := s.saveUnlocked(); err != nil {
        return StockItem{}, err
    }
    return stock, nil
}

func (s *Store) ListCashAccounts(filter CashAccountFilter) []CashAccount {
    s.mu.RLock()
    defer s.mu.RUnlock()
    var result []CashAccount
    for _, account := range s.data.CashAccounts {
        if filter.Status != "" && !strings.EqualFold(account.Status, filter.Status) {
            continue
        }
        if filter.Currency != "" && !strings.EqualFold(account.Currency, filter.Currency) {
            continue
        }
        result = append(result, account)
    }
    return paginateCashAccounts(result, filter.Offset, filter.MaxResults)
}

func (s *Store) GetCashAccount(id int, code string) (CashAccount, bool) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    if id != 0 {
        for _, account := range s.data.CashAccounts {
            if account.ID == id {
                return account, true
            }
        }
    }
    if code != "" {
        for _, account := range s.data.CashAccounts {
            if strings.EqualFold(account.Code, code) {
                return account, true
            }
        }
    }
    return CashAccount{}, false
}

func (s *Store) CreateCashAccount(draft CashAccountDraft) (CashAccount, error) {
    s.mu.Lock()
    defer s.mu.Unlock()

    if err := validateCashAccountDraft(draft); err != nil {
        return CashAccount{}, err
    }

    code := strings.ToUpper(strings.TrimSpace(draft.Code))
    if s.cashAccountCodeExistsLocked(code, 0) {
        return CashAccount{}, fmt.Errorf("cash account code already exists")
    }

    now := time.Now().UTC().Format(time.RFC3339)
    account := CashAccount{
        ID:        s.data.Meta.NextCashAccountID,
        Code:      code,
        Name:      strings.TrimSpace(draft.Name),
        Currency:  normalizeCurrency(draft.Currency),
        Balance:   roundMoney(draft.Balance),
        Status:    normalizeStatus(draft.Status, "OPEN"),
        CreatedAt: now,
        UpdatedAt: now,
    }

    s.data.CashAccounts = append(s.data.CashAccounts, account)
    s.data.Meta.NextCashAccountID++
    if err := s.saveUnlocked(); err != nil {
        return CashAccount{}, err
    }
    return account, nil
}

func (s *Store) UpdateCashAccount(id int, code string, draft CashAccountDraft) (CashAccount, error) {
    s.mu.Lock()
    defer s.mu.Unlock()

    idx := s.findCashAccountIndexLocked(id, code)
    if idx == -1 {
        return CashAccount{}, fmt.Errorf("cash account not found")
    }

    if err := validateCashAccountDraft(draft); err != nil {
        return CashAccount{}, err
    }

    normalizedCode := strings.ToUpper(strings.TrimSpace(draft.Code))
    if normalizedCode == "" {
        normalizedCode = s.data.CashAccounts[idx].Code
    }
    if s.cashAccountCodeExistsLocked(normalizedCode, s.data.CashAccounts[idx].ID) {
        return CashAccount{}, fmt.Errorf("cash account code already exists")
    }

    account := s.data.CashAccounts[idx]
    account.Code = normalizedCode
    account.Name = strings.TrimSpace(draft.Name)
    account.Currency = normalizeCurrency(draft.Currency)
    account.Balance = roundMoney(draft.Balance)
    account.Status = normalizeStatus(draft.Status, account.Status)
    account.UpdatedAt = time.Now().UTC().Format(time.RFC3339)

    s.data.CashAccounts[idx] = account
    if err := s.saveUnlocked(); err != nil {
        return CashAccount{}, err
    }
    return account, nil
}

func (s *Store) ListOrders(filter OrderFilter) ([]Order, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()

    from, to, err := normalizeDateRange(filter.FromDate, filter.ToDate)
    if err != nil {
        return nil, err
    }

    var result []Order
    for _, order := range s.data.Orders {
        if filter.CustomerID != 0 && order.CustomerID != filter.CustomerID {
            continue
        }
        if filter.Status != "" && !strings.EqualFold(order.Status, filter.Status) {
            continue
        }
        if !dateInRange(order.OrderDate, from, to) {
            continue
        }
        result = append(result, order)
    }

    sort.Slice(result, func(i, j int) bool { return result[i].ID < result[j].ID })
    return paginateOrders(result, filter.Offset, filter.MaxResults), nil
}

func (s *Store) GetOrder(id int, orderNo string) (Order, bool) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    if id != 0 {
        for _, order := range s.data.Orders {
            if order.ID == id {
                return order, true
            }
        }
    }
    if orderNo != "" {
        for _, order := range s.data.Orders {
            if strings.EqualFold(order.OrderNo, orderNo) {
                return order, true
            }
        }
    }
    return Order{}, false
}

func (s *Store) CreateOrder(draft OrderDraft) (Order, error) {
    s.mu.Lock()
    defer s.mu.Unlock()

    if err := validateOrderDraft(draft); err != nil {
        return Order{}, err
    }

    customer, err := s.resolveCustomerLocked(draft.CustomerID, draft.CustomerCode)
    if err != nil {
        return Order{}, err
    }

    if draft.Currency == "" {
        draft.Currency = customer.Currency
    }
    draft.Currency = normalizeCurrency(draft.Currency)
    if draft.Currency == "" {
        return Order{}, fmt.Errorf("currency is required")
    }
    if customer.Currency != "" && !strings.EqualFold(customer.Currency, draft.Currency) {
        return Order{}, fmt.Errorf("order currency must match customer currency")
    }

    if err := ensureValidDate(draft.OrderDate, "orderDate"); err != nil {
        return Order{}, err
    }

    status := normalizeStatus(draft.Status, "DRAFT")
    if status != "DRAFT" && status != "APPROVED" {
        return Order{}, fmt.Errorf("invalid order status")
    }

    lines, err := s.buildOrderLinesLocked(draft.Lines)
    if err != nil {
        return Order{}, err
    }

    if status == "APPROVED" {
        if err := s.ensureStockAvailabilityLocked(lines); err != nil {
            return Order{}, err
        }
    }

    orderNo := strings.TrimSpace(draft.OrderNo)
    if orderNo == "" {
        orderNo = fmt.Sprintf("SO-%d", s.data.Meta.NextOrderID)
    }
    if s.orderNoExistsLocked(orderNo, 0) {
        return Order{}, fmt.Errorf("orderNo already exists")
    }

    now := time.Now().UTC().Format(time.RFC3339)
    order := Order{
        ID:         s.data.Meta.NextOrderID,
        OrderNo:    orderNo,
        CustomerID: customer.ID,
        Status:     status,
        OrderDate:  draft.OrderDate,
        Currency:   draft.Currency,
        Notes:      strings.TrimSpace(draft.Notes),
        Lines:      lines,
        CreatedAt:  now,
        UpdatedAt:  now,
    }

    recalcOrderTotals(&order)

    s.data.Orders = append(s.data.Orders, order)
    s.data.Meta.NextOrderID++
    if err := s.saveUnlocked(); err != nil {
        return Order{}, err
    }
    return order, nil
}

func (s *Store) UpdateOrder(id int, orderNo string, draft OrderDraft) (Order, error) {
    s.mu.Lock()
    defer s.mu.Unlock()

    idx := s.findOrderIndexLocked(id, orderNo)
    if idx == -1 {
        return Order{}, fmt.Errorf("order not found")
    }

    order := s.data.Orders[idx]
    if order.Status != "DRAFT" {
        return Order{}, fmt.Errorf("only draft orders can be updated")
    }

    if err := validateOrderDraft(draft); err != nil {
        return Order{}, err
    }

    customer, err := s.resolveCustomerLocked(draft.CustomerID, draft.CustomerCode)
    if err != nil {
        return Order{}, err
    }

    currency := draft.Currency
    if currency == "" {
        currency = customer.Currency
    }
    currency = normalizeCurrency(currency)
    if currency == "" {
        return Order{}, fmt.Errorf("currency is required")
    }
    if customer.Currency != "" && !strings.EqualFold(customer.Currency, currency) {
        return Order{}, fmt.Errorf("order currency must match customer currency")
    }

    if err := ensureValidDate(draft.OrderDate, "orderDate"); err != nil {
        return Order{}, err
    }

    lines, err := s.buildOrderLinesLocked(draft.Lines)
    if err != nil {
        return Order{}, err
    }

    updatedOrderNo := strings.TrimSpace(draft.OrderNo)
    if updatedOrderNo == "" {
        updatedOrderNo = order.OrderNo
    }
    if s.orderNoExistsLocked(updatedOrderNo, order.ID) {
        return Order{}, fmt.Errorf("orderNo already exists")
    }

    order.OrderNo = updatedOrderNo
    order.CustomerID = customer.ID
    order.OrderDate = draft.OrderDate
    order.Currency = currency
    order.Notes = strings.TrimSpace(draft.Notes)
    order.Lines = lines
    order.UpdatedAt = time.Now().UTC().Format(time.RFC3339)

    recalcOrderTotals(&order)

    s.data.Orders[idx] = order
    if err := s.saveUnlocked(); err != nil {
        return Order{}, err
    }
    return order, nil
}

func (s *Store) ApproveOrder(id int, orderNo string) (Order, error) {
    s.mu.Lock()
    defer s.mu.Unlock()

    idx := s.findOrderIndexLocked(id, orderNo)
    if idx == -1 {
        return Order{}, fmt.Errorf("order not found")
    }

    order := s.data.Orders[idx]
    if order.Status != "DRAFT" {
        return Order{}, fmt.Errorf("only draft orders can be approved")
    }

    if err := s.ensureStockAvailabilityLocked(order.Lines); err != nil {
        return Order{}, err
    }

    order.Status = "APPROVED"
    order.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
    s.data.Orders[idx] = order

    if err := s.saveUnlocked(); err != nil {
        return Order{}, err
    }
    return order, nil
}

func (s *Store) CancelOrder(id int, orderNo string) (Order, error) {
    s.mu.Lock()
    defer s.mu.Unlock()

    idx := s.findOrderIndexLocked(id, orderNo)
    if idx == -1 {
        return Order{}, fmt.Errorf("order not found")
    }

    order := s.data.Orders[idx]
    if order.Status == "INVOICED" {
        return Order{}, fmt.Errorf("invoiced orders cannot be cancelled")
    }
    if order.Status == "CANCELLED" {
        return Order{}, fmt.Errorf("order already cancelled")
    }

    order.Status = "CANCELLED"
    order.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
    s.data.Orders[idx] = order

    if err := s.saveUnlocked(); err != nil {
        return Order{}, err
    }
    return order, nil
}

func (s *Store) ListInvoices(filter InvoiceFilter) ([]Invoice, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()

    from, to, err := normalizeDateRange(filter.FromDate, filter.ToDate)
    if err != nil {
        return nil, err
    }

    var result []Invoice
    for _, invoice := range s.data.Invoices {
        if filter.CustomerID != 0 && invoice.CustomerID != filter.CustomerID {
            continue
        }
        if filter.Status != "" && !strings.EqualFold(invoice.Status, filter.Status) {
            continue
        }
        if !dateInRange(invoice.InvoiceDate, from, to) {
            continue
        }
        result = append(result, invoice)
    }

    sort.Slice(result, func(i, j int) bool { return result[i].ID < result[j].ID })
    return paginateInvoices(result, filter.Offset, filter.MaxResults), nil
}

func (s *Store) GetInvoice(id int, invoiceNo string) (Invoice, bool) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    if id != 0 {
        for _, invoice := range s.data.Invoices {
            if invoice.ID == id {
                return invoice, true
            }
        }
    }
    if invoiceNo != "" {
        for _, invoice := range s.data.Invoices {
            if strings.EqualFold(invoice.InvoiceNo, invoiceNo) {
                return invoice, true
            }
        }
    }
    return Invoice{}, false
}

func (s *Store) CreateInvoice(draft InvoiceDraft) (Invoice, error) {
    s.mu.Lock()
    defer s.mu.Unlock()

    if err := validateInvoiceDraft(draft); err != nil {
        return Invoice{}, err
    }

    customer, err := s.resolveCustomerLocked(draft.CustomerID, draft.CustomerCode)
    if err != nil {
        return Invoice{}, err
    }

    if draft.OrderID != 0 || strings.TrimSpace(draft.OrderNo) != "" {
        return Invoice{}, fmt.Errorf("use CreateInvoiceFromOrder for order-linked invoices")
    }

    if err := ensureValidDate(draft.InvoiceDate, "invoiceDate"); err != nil {
        return Invoice{}, err
    }
    if err := ensureValidDate(draft.DueDate, "dueDate"); err != nil {
        return Invoice{}, err
    }
    if err := ensureDueAfterInvoice(draft.InvoiceDate, draft.DueDate); err != nil {
        return Invoice{}, err
    }

    currency := draft.Currency
    if currency == "" {
        currency = customer.Currency
    }
    currency = normalizeCurrency(currency)
    if currency == "" {
        return Invoice{}, fmt.Errorf("currency is required")
    }
    if customer.Currency != "" && !strings.EqualFold(customer.Currency, currency) {
        return Invoice{}, fmt.Errorf("invoice currency must match customer currency")
    }

    lines, err := s.buildInvoiceLinesLocked(draft.Lines)
    if err != nil {
        return Invoice{}, err
    }

    if err := s.ensureStockAvailabilityLocked(invoiceLinesToOrderLines(lines)); err != nil {
        return Invoice{}, err
    }

    invoiceNo := strings.TrimSpace(draft.InvoiceNo)
    if invoiceNo == "" {
        invoiceNo = fmt.Sprintf("INV-%d", s.data.Meta.NextInvoiceID)
    }
    if s.invoiceNoExistsLocked(invoiceNo, 0) {
        return Invoice{}, fmt.Errorf("invoiceNo already exists")
    }

    now := time.Now().UTC().Format(time.RFC3339)
    invoice := Invoice{
        ID:          s.data.Meta.NextInvoiceID,
        InvoiceNo:   invoiceNo,
        OrderID:     draft.OrderID,
        CustomerID:  customer.ID,
        Status:      "ISSUED",
        InvoiceDate: draft.InvoiceDate,
        DueDate:     draft.DueDate,
        Currency:    currency,
        Lines:       lines,
        PaidAmount:  0,
        CreatedAt:   now,
        UpdatedAt:   now,
    }

    recalcInvoiceTotals(&invoice)

    if err := s.applyStockMovementsLocked(invoice.Lines, -1); err != nil {
        return Invoice{}, err
    }

    s.data.Invoices = append(s.data.Invoices, invoice)
    s.data.Meta.NextInvoiceID++
    if err := s.saveUnlocked(); err != nil {
        return Invoice{}, err
    }
    return invoice, nil
}

func (s *Store) CreateInvoiceFromOrder(orderID int, orderNo string, invoiceDate, dueDate string) (Invoice, error) {
    s.mu.Lock()
    defer s.mu.Unlock()

    order, err := s.resolveOrderLocked(orderID, orderNo)
    if err != nil {
        return Invoice{}, err
    }

    if order.Status != "APPROVED" {
        return Invoice{}, fmt.Errorf("only approved orders can be invoiced")
    }

    if err := ensureValidDate(invoiceDate, "invoiceDate"); err != nil {
        return Invoice{}, err
    }
    if dueDate == "" {
        dueDate = invoiceDate
    }
    if err := ensureValidDate(dueDate, "dueDate"); err != nil {
        return Invoice{}, err
    }
    if err := ensureDueAfterInvoice(invoiceDate, dueDate); err != nil {
        return Invoice{}, err
    }

    if err := s.ensureStockAvailabilityLocked(order.Lines); err != nil {
        return Invoice{}, err
    }

    invoiceNo := fmt.Sprintf("INV-%d", s.data.Meta.NextInvoiceID)
    if s.invoiceNoExistsLocked(invoiceNo, 0) {
        return Invoice{}, fmt.Errorf("invoiceNo already exists")
    }

    now := time.Now().UTC().Format(time.RFC3339)
    invoice := Invoice{
        ID:          s.data.Meta.NextInvoiceID,
        InvoiceNo:   invoiceNo,
        OrderID:     order.ID,
        CustomerID:  order.CustomerID,
        Status:      "ISSUED",
        InvoiceDate: invoiceDate,
        DueDate:     dueDate,
        Currency:    order.Currency,
        Lines:       orderLinesToInvoiceLines(order.Lines),
        PaidAmount:  0,
        CreatedAt:   now,
        UpdatedAt:   now,
    }

    recalcInvoiceTotals(&invoice)
    if err := s.applyStockMovementsLocked(invoice.Lines, -1); err != nil {
        return Invoice{}, err
    }

    if err := s.updateOrderStatusLocked(order.ID, "INVOICED"); err != nil {
        return Invoice{}, err
    }

    s.data.Invoices = append(s.data.Invoices, invoice)
    s.data.Meta.NextInvoiceID++
    if err := s.saveUnlocked(); err != nil {
        return Invoice{}, err
    }
    return invoice, nil
}

func (s *Store) CancelInvoice(id int, invoiceNo string) (Invoice, error) {
    s.mu.Lock()
    defer s.mu.Unlock()

    idx := s.findInvoiceIndexLocked(id, invoiceNo)
    if idx == -1 {
        return Invoice{}, fmt.Errorf("invoice not found")
    }

    invoice := s.data.Invoices[idx]
    if invoice.Status != "ISSUED" {
        return Invoice{}, fmt.Errorf("only issued invoices can be cancelled")
    }

    if err := s.applyStockMovementsLocked(invoice.Lines, 1); err != nil {
        return Invoice{}, err
    }

    invoice.Status = "CANCELLED"
    invoice.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
    s.data.Invoices[idx] = invoice

    if invoice.OrderID != 0 {
        _ = s.updateOrderStatusLocked(invoice.OrderID, "APPROVED")
    }

    if err := s.saveUnlocked(); err != nil {
        return Invoice{}, err
    }
    return invoice, nil
}

func (s *Store) ListCashTransactions(filter CashTransactionFilter) ([]CashTransaction, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()

    from, to, err := normalizeDateRange(filter.FromDate, filter.ToDate)
    if err != nil {
        return nil, err
    }

    var result []CashTransaction
    for _, txn := range s.data.CashTransactions {
        if filter.CashAccountID != 0 && txn.CashAccountID != filter.CashAccountID {
            continue
        }
        if filter.Type != "" && !strings.EqualFold(txn.Type, filter.Type) {
            continue
        }
        if !dateInRange(txn.TransactionDate, from, to) {
            continue
        }
        result = append(result, txn)
    }

    sort.Slice(result, func(i, j int) bool { return result[i].ID < result[j].ID })
    return paginateCashTransactions(result, filter.Offset, filter.MaxResults), nil
}

func (s *Store) GetCashTransaction(id int, txnNo string) (CashTransaction, bool) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    if id != 0 {
        for _, txn := range s.data.CashTransactions {
            if txn.ID == id {
                return txn, true
            }
        }
    }
    if txnNo != "" {
        for _, txn := range s.data.CashTransactions {
            if strings.EqualFold(txn.TxnNo, txnNo) {
                return txn, true
            }
        }
    }
    return CashTransaction{}, false
}

func (s *Store) CreateCashTransaction(draft CashTransactionDraft) (CashTransaction, error) {
    s.mu.Lock()
    defer s.mu.Unlock()

    if err := validateCashTransactionDraft(draft); err != nil {
        return CashTransaction{}, err
    }

    cashAccount, err := s.resolveCashAccountLocked(draft.CashAccountID, draft.CashAccountCode)
    if err != nil {
        return CashTransaction{}, err
    }

    customer, err := s.resolveCustomerLocked(draft.CustomerID, draft.CustomerCode)
    if err != nil {
        return CashTransaction{}, err
    }

    invoice, invoiceFound := s.resolveInvoiceLocked(draft.InvoiceID, draft.InvoiceNo)

    if (draft.InvoiceID != 0 || strings.TrimSpace(draft.InvoiceNo) != "") && !invoiceFound {
        return CashTransaction{}, fmt.Errorf("invoice not found")
    }
    if invoiceFound && invoice.CustomerID != customer.ID {
        return CashTransaction{}, fmt.Errorf("invoice customer mismatch")
    }
    if invoiceFound && strings.EqualFold(strings.TrimSpace(draft.Type), "PAYMENT") {
        return CashTransaction{}, fmt.Errorf("payment transactions cannot target invoices")
    }
    if invoiceFound && strings.EqualFold(invoice.Status, "CANCELLED") {
        return CashTransaction{}, fmt.Errorf("invoice is cancelled")
    }

    currency := draft.Currency
    if currency == "" {
        currency = cashAccount.Currency
    }
    currency = normalizeCurrency(currency)
    if currency == "" {
        return CashTransaction{}, fmt.Errorf("currency is required")
    }
    if !strings.EqualFold(currency, cashAccount.Currency) {
        return CashTransaction{}, fmt.Errorf("transaction currency must match cash account")
    }
    if invoiceFound && !strings.EqualFold(currency, invoice.Currency) {
        return CashTransaction{}, fmt.Errorf("transaction currency must match invoice")
    }

    if err := ensureValidDate(draft.TransactionDate, "transactionDate"); err != nil {
        return CashTransaction{}, err
    }

    txnType := strings.ToUpper(strings.TrimSpace(draft.Type))
    if txnType != "COLLECTION" && txnType != "PAYMENT" {
        return CashTransaction{}, fmt.Errorf("invalid transaction type")
    }

    if draft.Amount <= 0 {
        return CashTransaction{}, fmt.Errorf("amount must be greater than zero")
    }

    if txnType == "PAYMENT" && cashAccount.Balance < draft.Amount {
        return CashTransaction{}, fmt.Errorf("insufficient cash balance")
    }

    if invoiceFound && txnType == "COLLECTION" {
        remaining := roundMoney(invoice.GrandTotal - invoice.PaidAmount)
        if draft.Amount > remaining {
            return CashTransaction{}, fmt.Errorf("collection amount exceeds invoice balance")
        }
        invoice.PaidAmount = roundMoney(invoice.PaidAmount + draft.Amount)
        if invoice.PaidAmount >= invoice.GrandTotal {
            invoice.Status = "PAID"
        }
        invoice.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
        s.updateInvoiceLocked(invoice)
    }

    if txnType == "COLLECTION" {
        cashAccount.Balance = roundMoney(cashAccount.Balance + draft.Amount)
    } else {
        cashAccount.Balance = roundMoney(cashAccount.Balance - draft.Amount)
    }
    cashAccount.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
    s.updateCashAccountLocked(cashAccount)

    txnNo := strings.TrimSpace(draft.TxnNo)
    if txnNo == "" {
        txnNo = fmt.Sprintf("CT-%d", s.data.Meta.NextCashTransactionID)
    }
    if s.cashTxnNoExistsLocked(txnNo, 0) {
        return CashTransaction{}, fmt.Errorf("txnNo already exists")
    }

    now := time.Now().UTC().Format(time.RFC3339)
    txn := CashTransaction{
        ID:              s.data.Meta.NextCashTransactionID,
        TxnNo:           txnNo,
        CashAccountID:   cashAccount.ID,
        CustomerID:      customer.ID,
        InvoiceID:       invoice.ID,
        Type:            txnType,
        Amount:          roundMoney(draft.Amount),
        Currency:        currency,
        Method:          strings.ToUpper(strings.TrimSpace(draft.Method)),
        TransactionDate: draft.TransactionDate,
        Description:     strings.TrimSpace(draft.Description),
        Status:          "POSTED",
        ReversalOf:      0,
        CreatedAt:       now,
    }

    s.data.CashTransactions = append(s.data.CashTransactions, txn)
    s.data.Meta.NextCashTransactionID++
    if err := s.saveUnlocked(); err != nil {
        return CashTransaction{}, err
    }
    return txn, nil
}

func (s *Store) ReverseCashTransaction(id int, txnNo string) (CashTransaction, error) {
    s.mu.Lock()
    defer s.mu.Unlock()

    idx := s.findCashTransactionIndexLocked(id, txnNo)
    if idx == -1 {
        return CashTransaction{}, fmt.Errorf("cash transaction not found")
    }

    original := s.data.CashTransactions[idx]
    if original.Status == "REVERSED" {
        return CashTransaction{}, fmt.Errorf("transaction already reversed")
    }

    cashAccount, err := s.resolveCashAccountLocked(original.CashAccountID, "")
    if err != nil {
        return CashTransaction{}, err
    }

    if original.Type == "COLLECTION" {
        if cashAccount.Balance < original.Amount {
            return CashTransaction{}, fmt.Errorf("insufficient cash balance to reverse")
        }
        cashAccount.Balance = roundMoney(cashAccount.Balance - original.Amount)
    } else {
        cashAccount.Balance = roundMoney(cashAccount.Balance + original.Amount)
    }
    cashAccount.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
    s.updateCashAccountLocked(cashAccount)

    if original.InvoiceID != 0 {
        invoice, found := s.resolveInvoiceLocked(original.InvoiceID, "")
        if found {
            invoice.PaidAmount = roundMoney(invoice.PaidAmount - original.Amount)
            if invoice.PaidAmount < invoice.GrandTotal {
                invoice.Status = "ISSUED"
            }
            invoice.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
            s.updateInvoiceLocked(invoice)
        }
    }

    original.Status = "REVERSED"
    s.data.CashTransactions[idx] = original

    now := time.Now().UTC().Format(time.RFC3339)
    reversal := CashTransaction{
        ID:              s.data.Meta.NextCashTransactionID,
        TxnNo:           fmt.Sprintf("CT-%d", s.data.Meta.NextCashTransactionID),
        CashAccountID:   original.CashAccountID,
        CustomerID:      original.CustomerID,
        InvoiceID:       original.InvoiceID,
        Type:            reverseTxnType(original.Type),
        Amount:          original.Amount,
        Currency:        original.Currency,
        Method:          original.Method,
        TransactionDate: time.Now().UTC().Format("2006-01-02"),
        Description:     fmt.Sprintf("Reversal of %s", original.TxnNo),
        Status:          "POSTED",
        ReversalOf:      original.ID,
        CreatedAt:       now,
    }

    s.data.CashTransactions = append(s.data.CashTransactions, reversal)
    s.data.Meta.NextCashTransactionID++
    if err := s.saveUnlocked(); err != nil {
        return CashTransaction{}, err
    }
    return reversal, nil
}

func (s *Store) ListCountries(continent string) []Country {
    s.mu.RLock()
    defer s.mu.RUnlock()
    var result []Country
    for _, country := range s.data.Countries {
        if continent != "" && !strings.EqualFold(country.Continent, continent) {
            continue
        }
        result = append(result, country)
    }
    return result
}

func (s *Store) ListCities(countryCode string) ([]City, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    if strings.TrimSpace(countryCode) == "" {
        return nil, fmt.Errorf("countryCode is required")
    }
    var result []City
    for _, city := range s.data.Cities {
        if strings.EqualFold(city.CountryCode, countryCode) {
            result = append(result, city)
        }
    }
    return result, nil
}

func defaultSeedData() StoreData {
    return StoreData{
        Meta: StoreMeta{
            NextCustomerID:        1,
            NextStockID:           1,
            NextCashAccountID:     1,
            NextOrderID:           1,
            NextInvoiceID:         1,
            NextCashTransactionID: 1,
        },
        Customers:        []Customer{},
        Stocks:           []StockItem{},
        CashAccounts:     []CashAccount{},
        Orders:           []Order{},
        Invoices:         []Invoice{},
        CashTransactions: []CashTransaction{},
        Countries:        []Country{},
        Cities:           []City{},
    }
}

func normalizeStatus(status string, fallback string) string {
    normalized := strings.ToUpper(strings.TrimSpace(status))
    if normalized == "" {
        return fallback
    }
    return normalized
}

func normalizeCurrency(currency string) string {
    normalized := strings.ToUpper(strings.TrimSpace(currency))
    switch normalized {
    case "TRY", "USD", "EUR":
        return normalized
    default:
        return ""
    }
}

func roundMoney(value float64) float64 {
    return math.Round(value*100) / 100
}

func roundQuantity(value float64) float64 {
    return math.Round(value*10000) / 10000
}

func recalcOrderTotals(order *Order) {
    var subtotal, discountTotal, taxTotal float64
    for i := range order.Lines {
        line := &order.Lines[i]
        net := line.Quantity * line.UnitPrice
        discount := net * line.DiscountRate / 100
        taxable := net - discount
        tax := taxable * line.TaxRate / 100
        line.LineTotal = roundMoney(taxable + tax)
        subtotal += net
        discountTotal += discount
        taxTotal += tax
    }
    order.Subtotal = roundMoney(subtotal)
    order.DiscountTotal = roundMoney(discountTotal)
    order.TaxTotal = roundMoney(taxTotal)
    order.GrandTotal = roundMoney(order.Subtotal - order.DiscountTotal + order.TaxTotal)
}

func recalcInvoiceTotals(invoice *Invoice) {
    var subtotal, discountTotal, taxTotal float64
    for i := range invoice.Lines {
        line := &invoice.Lines[i]
        net := line.Quantity * line.UnitPrice
        discount := net * line.DiscountRate / 100
        taxable := net - discount
        tax := taxable * line.TaxRate / 100
        line.LineTotal = roundMoney(taxable + tax)
        subtotal += net
        discountTotal += discount
        taxTotal += tax
    }
    invoice.Subtotal = roundMoney(subtotal)
    invoice.DiscountTotal = roundMoney(discountTotal)
    invoice.TaxTotal = roundMoney(taxTotal)
    invoice.GrandTotal = roundMoney(invoice.Subtotal - invoice.DiscountTotal + invoice.TaxTotal)
}

func invoiceLinesToOrderLines(lines []InvoiceLine) []OrderLine {
    var result []OrderLine
    for _, line := range lines {
        result = append(result, OrderLine{
            LineNo:       line.LineNo,
            StockID:      line.StockID,
            StockCode:    line.StockCode,
            Description:  line.Description,
            Unit:         line.Unit,
            Quantity:     line.Quantity,
            UnitPrice:    line.UnitPrice,
            DiscountRate: line.DiscountRate,
            TaxRate:      line.TaxRate,
            LineTotal:    line.LineTotal,
        })
    }
    return result
}

func orderLinesToInvoiceLines(lines []OrderLine) []InvoiceLine {
    var result []InvoiceLine
    for _, line := range lines {
        result = append(result, InvoiceLine{
            LineNo:       line.LineNo,
            StockID:      line.StockID,
            StockCode:    line.StockCode,
            Description:  line.Description,
            Unit:         line.Unit,
            Quantity:     line.Quantity,
            UnitPrice:    line.UnitPrice,
            DiscountRate: line.DiscountRate,
            TaxRate:      line.TaxRate,
            LineTotal:    line.LineTotal,
        })
    }
    return result
}

func reverseTxnType(txnType string) string {
    if strings.EqualFold(txnType, "COLLECTION") {
        return "PAYMENT"
    }
    return "COLLECTION"
}

func (s *Store) resolveCustomerLocked(id int, code string) (Customer, error) {
    if id != 0 {
        for _, customer := range s.data.Customers {
            if customer.ID == id {
                return customer, nil
            }
        }
    }
    if code != "" {
        for _, customer := range s.data.Customers {
            if strings.EqualFold(customer.Code, code) {
                return customer, nil
            }
        }
    }
    return Customer{}, fmt.Errorf("customer not found")
}

func (s *Store) resolveOrderLocked(id int, orderNo string) (Order, error) {
    if id != 0 {
        for _, order := range s.data.Orders {
            if order.ID == id {
                return order, nil
            }
        }
    }
    if orderNo != "" {
        for _, order := range s.data.Orders {
            if strings.EqualFold(order.OrderNo, orderNo) {
                return order, nil
            }
        }
    }
    return Order{}, fmt.Errorf("order not found")
}

func (s *Store) resolveInvoiceLocked(id int, invoiceNo string) (Invoice, bool) {
    if id != 0 {
        for _, invoice := range s.data.Invoices {
            if invoice.ID == id {
                return invoice, true
            }
        }
    }
    if invoiceNo != "" {
        for _, invoice := range s.data.Invoices {
            if strings.EqualFold(invoice.InvoiceNo, invoiceNo) {
                return invoice, true
            }
        }
    }
    return Invoice{}, false
}

func (s *Store) resolveCashAccountLocked(id int, code string) (CashAccount, error) {
    if id != 0 {
        for _, account := range s.data.CashAccounts {
            if account.ID == id {
                return account, nil
            }
        }
    }
    if code != "" {
        for _, account := range s.data.CashAccounts {
            if strings.EqualFold(account.Code, code) {
                return account, nil
            }
        }
    }
    return CashAccount{}, fmt.Errorf("cash account not found")
}

func (s *Store) buildOrderLinesLocked(lines []OrderLineDraft) ([]OrderLine, error) {
    if len(lines) == 0 {
        return nil, fmt.Errorf("order lines are required")
    }

    var result []OrderLine
    for idx, draft := range lines {
        stock, err := s.resolveStockLocked(draft.StockID, draft.StockCode)
        if err != nil {
            return nil, err
        }
        if draft.Quantity <= 0 {
            return nil, fmt.Errorf("line %d: quantity must be greater than zero", idx+1)
        }
        if draft.UnitPrice <= 0 {
            draft.UnitPrice = stock.Price
        }
        if draft.DiscountRate < 0 || draft.DiscountRate > 100 {
            return nil, fmt.Errorf("line %d: invalid discountRate", idx+1)
        }
        taxRate := draft.TaxRate
        if taxRate == 0 {
            taxRate = stock.VatRate
        }
        if !isAllowedTaxRate(taxRate) {
            return nil, fmt.Errorf("line %d: invalid taxRate", idx+1)
        }
        line := OrderLine{
            LineNo:       idx + 1,
            StockID:      stock.ID,
            StockCode:    stock.Code,
            Description:  stock.Name,
            Unit:         stock.Unit,
            Quantity:     roundQuantity(draft.Quantity),
            UnitPrice:    roundMoney(draft.UnitPrice),
            DiscountRate: roundMoney(draft.DiscountRate),
            TaxRate:      roundMoney(taxRate),
        }
        result = append(result, line)
    }
    return result, nil
}

func (s *Store) buildInvoiceLinesLocked(lines []InvoiceLineDraft) ([]InvoiceLine, error) {
    if len(lines) == 0 {
        return nil, fmt.Errorf("invoice lines are required")
    }

    var result []InvoiceLine
    for idx, draft := range lines {
        stock, err := s.resolveStockLocked(draft.StockID, draft.StockCode)
        if err != nil {
            return nil, err
        }
        if draft.Quantity <= 0 {
            return nil, fmt.Errorf("line %d: quantity must be greater than zero", idx+1)
        }
        if draft.UnitPrice <= 0 {
            draft.UnitPrice = stock.Price
        }
        if draft.DiscountRate < 0 || draft.DiscountRate > 100 {
            return nil, fmt.Errorf("line %d: invalid discountRate", idx+1)
        }
        taxRate := draft.TaxRate
        if taxRate == 0 {
            taxRate = stock.VatRate
        }
        if !isAllowedTaxRate(taxRate) {
            return nil, fmt.Errorf("line %d: invalid taxRate", idx+1)
        }
        line := InvoiceLine{
            LineNo:       idx + 1,
            StockID:      stock.ID,
            StockCode:    stock.Code,
            Description:  stock.Name,
            Unit:         stock.Unit,
            Quantity:     roundQuantity(draft.Quantity),
            UnitPrice:    roundMoney(draft.UnitPrice),
            DiscountRate: roundMoney(draft.DiscountRate),
            TaxRate:      roundMoney(taxRate),
        }
        result = append(result, line)
    }
    return result, nil
}

func (s *Store) ensureStockAvailabilityLocked(lines []OrderLine) error {
    for _, line := range lines {
        idx := s.findStockIndexLocked(line.StockID, "")
        if idx == -1 {
            return fmt.Errorf("stock not found")
        }
        stock := s.data.Stocks[idx]
        if stock.StockOnHand < line.Quantity {
            return fmt.Errorf("insufficient stock for %s", stock.Code)
        }
    }
    return nil
}

func (s *Store) applyStockMovementsLocked(lines []InvoiceLine, direction int) error {
    for _, line := range lines {
        idx := s.findStockIndexLocked(line.StockID, "")
        if idx == -1 {
            return fmt.Errorf("stock not found")
        }
        stock := s.data.Stocks[idx]
        movement := line.Quantity * float64(direction)
        if direction < 0 && stock.StockOnHand < line.Quantity {
            return fmt.Errorf("insufficient stock for %s", stock.Code)
        }
        stock.StockOnHand = roundQuantity(stock.StockOnHand + movement)
        stock.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
        s.data.Stocks[idx] = stock
    }
    return nil
}

func (s *Store) resolveStockLocked(id int, code string) (StockItem, error) {
    if id != 0 {
        for _, stock := range s.data.Stocks {
            if stock.ID == id {
                return stock, nil
            }
        }
    }
    if code != "" {
        for _, stock := range s.data.Stocks {
            if strings.EqualFold(stock.Code, code) {
                return stock, nil
            }
        }
    }
    return StockItem{}, fmt.Errorf("stock not found")
}

func (s *Store) updateOrderStatusLocked(orderID int, status string) error {
    for i := range s.data.Orders {
        if s.data.Orders[i].ID == orderID {
            s.data.Orders[i].Status = status
            s.data.Orders[i].UpdatedAt = time.Now().UTC().Format(time.RFC3339)
            return nil
        }
    }
    return fmt.Errorf("order not found")
}

func (s *Store) updateCashAccountLocked(account CashAccount) {
    for i := range s.data.CashAccounts {
        if s.data.CashAccounts[i].ID == account.ID {
            s.data.CashAccounts[i] = account
            return
        }
    }
}

func (s *Store) updateInvoiceLocked(invoice Invoice) {
    for i := range s.data.Invoices {
        if s.data.Invoices[i].ID == invoice.ID {
            s.data.Invoices[i] = invoice
            return
        }
    }
}

func (s *Store) findCustomerIndexLocked(id int, code string) int {
    for i, customer := range s.data.Customers {
        if id != 0 && customer.ID == id {
            return i
        }
        if code != "" && strings.EqualFold(customer.Code, code) {
            return i
        }
    }
    return -1
}

func (s *Store) findStockIndexLocked(id int, code string) int {
    for i, stock := range s.data.Stocks {
        if id != 0 && stock.ID == id {
            return i
        }
        if code != "" && strings.EqualFold(stock.Code, code) {
            return i
        }
    }
    return -1
}

func (s *Store) findCashAccountIndexLocked(id int, code string) int {
    for i, account := range s.data.CashAccounts {
        if id != 0 && account.ID == id {
            return i
        }
        if code != "" && strings.EqualFold(account.Code, code) {
            return i
        }
    }
    return -1
}

func (s *Store) findOrderIndexLocked(id int, orderNo string) int {
    for i, order := range s.data.Orders {
        if id != 0 && order.ID == id {
            return i
        }
        if orderNo != "" && strings.EqualFold(order.OrderNo, orderNo) {
            return i
        }
    }
    return -1
}

func (s *Store) findInvoiceIndexLocked(id int, invoiceNo string) int {
    for i, invoice := range s.data.Invoices {
        if id != 0 && invoice.ID == id {
            return i
        }
        if invoiceNo != "" && strings.EqualFold(invoice.InvoiceNo, invoiceNo) {
            return i
        }
    }
    return -1
}

func (s *Store) findCashTransactionIndexLocked(id int, txnNo string) int {
    for i, txn := range s.data.CashTransactions {
        if id != 0 && txn.ID == id {
            return i
        }
        if txnNo != "" && strings.EqualFold(txn.TxnNo, txnNo) {
            return i
        }
    }
    return -1
}

func (s *Store) customerCodeExistsLocked(code string, excludeID int) bool {
    for _, customer := range s.data.Customers {
        if strings.EqualFold(customer.Code, code) && customer.ID != excludeID {
            return true
        }
    }
    return false
}

func (s *Store) stockCodeExistsLocked(code string, excludeID int) bool {
    for _, stock := range s.data.Stocks {
        if strings.EqualFold(stock.Code, code) && stock.ID != excludeID {
            return true
        }
    }
    return false
}

func (s *Store) cashAccountCodeExistsLocked(code string, excludeID int) bool {
    for _, account := range s.data.CashAccounts {
        if strings.EqualFold(account.Code, code) && account.ID != excludeID {
            return true
        }
    }
    return false
}

func (s *Store) orderNoExistsLocked(orderNo string, excludeID int) bool {
    for _, order := range s.data.Orders {
        if strings.EqualFold(order.OrderNo, orderNo) && order.ID != excludeID {
            return true
        }
    }
    return false
}

func (s *Store) invoiceNoExistsLocked(invoiceNo string, excludeID int) bool {
    for _, invoice := range s.data.Invoices {
        if strings.EqualFold(invoice.InvoiceNo, invoiceNo) && invoice.ID != excludeID {
            return true
        }
    }
    return false
}

func (s *Store) cashTxnNoExistsLocked(txnNo string, excludeID int) bool {
    for _, txn := range s.data.CashTransactions {
        if strings.EqualFold(txn.TxnNo, txnNo) && txn.ID != excludeID {
            return true
        }
    }
    return false
}

func ensureValidDate(value, field string) error {
    if strings.TrimSpace(value) == "" {
        return fmt.Errorf("%s is required", field)
    }
    if _, err := parseDate(value); err != nil {
        return fmt.Errorf("%s must be YYYY-MM-DD", field)
    }
    return nil
}

func ensureDueAfterInvoice(invoiceDate, dueDate string) error {
    inv, err := parseDate(invoiceDate)
    if err != nil {
        return err
    }
    due, err := parseDate(dueDate)
    if err != nil {
        return err
    }
    if due.Before(inv) {
        return fmt.Errorf("dueDate cannot be earlier than invoiceDate")
    }
    return nil
}

func parseDate(value string) (time.Time, error) {
    if value == "" {
        return time.Time{}, fmt.Errorf("date is required")
    }
    if strings.Contains(value, "T") {
        return time.Parse(time.RFC3339, value)
    }
    return time.Parse("2006-01-02", value)
}

func normalizeDateRange(from, to string) (*time.Time, *time.Time, error) {
    var fromTime *time.Time
    var toTime *time.Time

    if from != "" {
        parsed, err := parseDate(from)
        if err != nil {
            return nil, nil, err
        }
        fromTime = &parsed
    }
    if to != "" {
        parsed, err := parseDate(to)
        if err != nil {
            return nil, nil, err
        }
        toTime = &parsed
    }
    return fromTime, toTime, nil
}

func dateInRange(value string, from, to *time.Time) bool {
    parsed, err := parseDate(value)
    if err != nil {
        return false
    }
    if from != nil && parsed.Before(*from) {
        return false
    }
    if to != nil && parsed.After(*to) {
        return false
    }
    return true
}

func validateCustomerDraft(draft CustomerDraft) error {
    if strings.TrimSpace(draft.Code) == "" {
        return fmt.Errorf("customer code is required")
    }
    if strings.TrimSpace(draft.Name) == "" {
        return fmt.Errorf("customer name is required")
    }
    if normalizeCurrency(draft.Currency) == "" {
        return fmt.Errorf("invalid currency")
    }
    if draft.RiskLimit < 0 {
        return fmt.Errorf("riskLimit cannot be negative")
    }
    if draft.TaxNumber != "" && !isValidTaxNumber(draft.TaxNumber) {
        return fmt.Errorf("invalid taxNumber")
    }
    if draft.Email != "" && !isValidEmail(draft.Email) {
        return fmt.Errorf("invalid email")
    }
    if draft.Status != "" {
        status := normalizeStatus(draft.Status, "")
        if status != "ACTIVE" && status != "INACTIVE" {
            return fmt.Errorf("invalid customer status")
        }
    }
    return nil
}

func validateStockDraft(draft StockDraft) error {
    if strings.TrimSpace(draft.Code) == "" {
        return fmt.Errorf("stock code is required")
    }
    if strings.TrimSpace(draft.Name) == "" {
        return fmt.Errorf("stock name is required")
    }
    if strings.TrimSpace(draft.Unit) == "" {
        return fmt.Errorf("stock unit is required")
    }
    if draft.Price < 0 {
        return fmt.Errorf("price cannot be negative")
    }
    if draft.StockOnHand < 0 {
        return fmt.Errorf("stockOnHand cannot be negative")
    }
    if draft.MinStock < 0 {
        return fmt.Errorf("minStock cannot be negative")
    }
    if !isAllowedTaxRate(draft.VatRate) {
        return fmt.Errorf("invalid vatRate")
    }
    if draft.Status != "" {
        status := normalizeStatus(draft.Status, "")
        if status != "ACTIVE" && status != "INACTIVE" {
            return fmt.Errorf("invalid stock status")
        }
    }
    return nil
}

func validateCashAccountDraft(draft CashAccountDraft) error {
    if strings.TrimSpace(draft.Code) == "" {
        return fmt.Errorf("cash account code is required")
    }
    if strings.TrimSpace(draft.Name) == "" {
        return fmt.Errorf("cash account name is required")
    }
    if normalizeCurrency(draft.Currency) == "" {
        return fmt.Errorf("invalid currency")
    }
    if draft.Balance < 0 {
        return fmt.Errorf("balance cannot be negative")
    }
    if draft.Status != "" {
        status := normalizeStatus(draft.Status, "")
        if status != "OPEN" && status != "CLOSED" {
            return fmt.Errorf("invalid cash account status")
        }
    }
    return nil
}

func validateOrderDraft(draft OrderDraft) error {
    if draft.CustomerID == 0 && strings.TrimSpace(draft.CustomerCode) == "" {
        return fmt.Errorf("customerId or customerCode is required")
    }
    if draft.OrderDate == "" {
        return fmt.Errorf("orderDate is required")
    }
    if len(draft.Lines) == 0 {
        return fmt.Errorf("order lines are required")
    }
    return nil
}

func validateInvoiceDraft(draft InvoiceDraft) error {
    if draft.CustomerID == 0 && strings.TrimSpace(draft.CustomerCode) == "" {
        return fmt.Errorf("customerId or customerCode is required")
    }
    if draft.InvoiceDate == "" {
        return fmt.Errorf("invoiceDate is required")
    }
    if draft.DueDate == "" {
        return fmt.Errorf("dueDate is required")
    }
    if len(draft.Lines) == 0 {
        return fmt.Errorf("invoice lines are required")
    }
    return nil
}

func validateCashTransactionDraft(draft CashTransactionDraft) error {
    if draft.CashAccountID == 0 && strings.TrimSpace(draft.CashAccountCode) == "" {
        return fmt.Errorf("cashAccountId or cashAccountCode is required")
    }
    if draft.CustomerID == 0 && strings.TrimSpace(draft.CustomerCode) == "" {
        return fmt.Errorf("customerId or customerCode is required")
    }
    txnType := strings.ToUpper(strings.TrimSpace(draft.Type))
    if txnType == "" {
        return fmt.Errorf("type is required")
    }
    if txnType != "COLLECTION" && txnType != "PAYMENT" {
        return fmt.Errorf("invalid transaction type")
    }
    if draft.Amount <= 0 {
        return fmt.Errorf("amount must be greater than zero")
    }
    method := strings.ToUpper(strings.TrimSpace(draft.Method))
    if method == "" {
        return fmt.Errorf("method is required")
    }
    switch method {
    case "CASH", "BANK", "POS", "TRANSFER":
    default:
        return fmt.Errorf("invalid payment method")
    }
    if draft.TransactionDate == "" {
        return fmt.Errorf("transactionDate is required")
    }
    return nil
}

func isValidTaxNumber(value string) bool {
    digits := 0
    for _, r := range value {
        if r < '0' || r > '9' {
            return false
        }
        digits++
    }
    return digits == 10 || digits == 11
}

func isValidEmail(value string) bool {
    value = strings.TrimSpace(value)
    if value == "" {
        return true
    }
    parts := strings.Split(value, "@")
    if len(parts) != 2 {
        return false
    }
    if parts[0] == "" || parts[1] == "" {
        return false
    }
    return strings.Contains(parts[1], ".")
}

func isAllowedTaxRate(rate float64) bool {
    switch roundMoney(rate) {
    case 0, 1, 8, 18:
        return true
    default:
        return false
    }
}

func paginateCustomers(items []Customer, offset, max int) []Customer {
    start, end := paginateRange(len(items), offset, max)
    return append([]Customer(nil), items[start:end]...)
}

func paginateStocks(items []StockItem, offset, max int) []StockItem {
    start, end := paginateRange(len(items), offset, max)
    return append([]StockItem(nil), items[start:end]...)
}

func paginateCashAccounts(items []CashAccount, offset, max int) []CashAccount {
    start, end := paginateRange(len(items), offset, max)
    return append([]CashAccount(nil), items[start:end]...)
}

func paginateOrders(items []Order, offset, max int) []Order {
    start, end := paginateRange(len(items), offset, max)
    return append([]Order(nil), items[start:end]...)
}

func paginateInvoices(items []Invoice, offset, max int) []Invoice {
    start, end := paginateRange(len(items), offset, max)
    return append([]Invoice(nil), items[start:end]...)
}

func paginateCashTransactions(items []CashTransaction, offset, max int) []CashTransaction {
    start, end := paginateRange(len(items), offset, max)
    return append([]CashTransaction(nil), items[start:end]...)
}

func paginateRange(total, offset, max int) (int, int) {
    if offset < 0 {
        offset = 0
    }
    if offset > total {
        offset = total
    }
    end := total
    if max > 0 && offset+max < total {
        end = offset + max
    }
    return offset, end
}

func maxCustomerID(customers []Customer) int {
    maxID := 0
    for _, customer := range customers {
        if customer.ID > maxID {
            maxID = customer.ID
        }
    }
    return maxID
}

func maxStockID(stocks []StockItem) int {
    maxID := 0
    for _, stock := range stocks {
        if stock.ID > maxID {
            maxID = stock.ID
        }
    }
    return maxID
}

func maxCashAccountID(accounts []CashAccount) int {
    maxID := 0
    for _, account := range accounts {
        if account.ID > maxID {
            maxID = account.ID
        }
    }
    return maxID
}

func maxOrderID(orders []Order) int {
    maxID := 0
    for _, order := range orders {
        if order.ID > maxID {
            maxID = order.ID
        }
    }
    return maxID
}

func maxInvoiceID(invoices []Invoice) int {
    maxID := 0
    for _, invoice := range invoices {
        if invoice.ID > maxID {
            maxID = invoice.ID
        }
    }
    return maxID
}

func maxCashTransactionID(txns []CashTransaction) int {
    maxID := 0
    for _, txn := range txns {
        if txn.ID > maxID {
            maxID = txn.ID
        }
    }
    return maxID
}

var endpointStores map[string]*Store

func initStores() (map[string]*Store, error) {
    dataDir := os.Getenv("SOAP_MOCK_DATA_DIR")
    if strings.TrimSpace(dataDir) == "" {
        dataDir = "service"
    }

    endpoints := []string{"session", "wsse", "basic", "ntlm", "noauth", "sap"}
    stores := make(map[string]*Store, len(endpoints))
    for _, endpoint := range endpoints {
        path := filepath.Join(dataDir, endpoint+".json")
        store, err := loadStore(path)
        if err != nil {
            return nil, err
        }
        stores[endpoint] = store
    }
    return stores, nil
}

func storeForEndpoint(endpoint string) (*Store, error) {
    if endpointStores == nil {
        return nil, fmt.Errorf("stores not initialized")
    }
    store, ok := endpointStores[endpoint]
    if !ok || store == nil {
        return nil, fmt.Errorf("store not found for endpoint")
    }
    return store, nil
}

