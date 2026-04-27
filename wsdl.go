package main

import (
    "fmt"
    "strings"
)

// PolicyType determines what WS-Policy assertion to include in the WSDL
type PolicyType string

const (
    PolicyNone              PolicyType = ""
    PolicyBasicAuth         PolicyType = "basic"
    PolicyWSSEUsernameToken PolicyType = "wsse"
    PolicyNTLMNegotiate     PolicyType = "ntlm"
)

type wsdlConfig struct {
    Namespace            string
    ServiceName          string
    PortTypeName         string
    BindingName          string
    Address              string
    Endpoint             string
    IncludeSessionOps    bool
    IncludeSessionHeader bool
    IncludePublic        bool
    Policy               PolicyType
}

var erpOperations = []string{
    "GetCustomers",
    "ListCustomers",
    "GetCustomer",
    "CreateCustomer",
    "UpdateCustomer",
    "GetStocks",
    "ListStocks",
    "GetStock",
    "CreateStock",
    "UpdateStock",
    "GetCashAccounts",
    "ListCashAccounts",
    "GetCashAccount",
    "CreateCashAccount",
    "UpdateCashAccount",
    "GetOrders",
    "ListOrders",
    "GetOrder",
    "CreateOrder",
    "UpdateOrder",
    "ApproveOrder",
    "CancelOrder",
    "GetInvoices",
    "ListInvoices",
    "GetInvoice",
    "CreateInvoice",
    "CreateInvoiceFromOrder",
    "CancelInvoice",
    "GetCashTransactions",
    "ListCashTransactions",
    "GetCashTransaction",
    "CreateCashTransaction",
    "ReverseCashTransaction",
}

var publicOperations = []string{
    "GetCountries",
    "GetCities",
}

var sessionOperations = []string{
    "Login",
    "Logout",
}

const sessionSchemaBody = `            <xsd:element name="Login">
                <xsd:complexType>
                    <xsd:sequence>
                        <xsd:element name="username" type="xsd:string"/>
                        <xsd:element name="password" type="xsd:string"/>
                    </xsd:sequence>
                </xsd:complexType>
            </xsd:element>
            <xsd:element name="LoginResponse">
                <xsd:complexType>
                    <xsd:sequence>
                        <xsd:element name="sessionToken" type="xsd:string"/>
                        <xsd:element name="expiresIn" type="xsd:int"/>
                        <xsd:element name="success" type="xsd:boolean"/>
                    </xsd:sequence>
                </xsd:complexType>
            </xsd:element>

            <xsd:element name="SessionHeader">
                <xsd:complexType>
                    <xsd:sequence>
                        <xsd:element name="sessionToken" type="xsd:string"/>
                    </xsd:sequence>
                </xsd:complexType>
            </xsd:element>

            <xsd:element name="Logout">
                <xsd:complexType>
                    <xsd:sequence>
                        <xsd:element name="sessionToken" type="xsd:string" minOccurs="0"/>
                    </xsd:sequence>
                </xsd:complexType>
            </xsd:element>
            <xsd:element name="LogoutResponse">
                <xsd:complexType>
                    <xsd:sequence>
                        <xsd:element name="success" type="xsd:boolean"/>
                    </xsd:sequence>
                </xsd:complexType>
            </xsd:element>
`

const erpSchemaBody = `            <!-- Customers -->
            <xsd:element name="GetCustomers">
                <xsd:complexType>
                    <xsd:sequence>
                        <xsd:element name="code" type="xsd:string" minOccurs="0"/>
                        <xsd:element name="status" type="xsd:string" minOccurs="0"/>
                        <xsd:element name="taxNumber" type="xsd:string" minOccurs="0"/>
                        <xsd:element name="maxResults" type="xsd:int" minOccurs="0"/>
                        <xsd:element name="offset" type="xsd:int" minOccurs="0"/>
                    </xsd:sequence>
                </xsd:complexType>
            </xsd:element>
            <xsd:element name="GetCustomersResponse">
                <xsd:complexType>
                    <xsd:sequence>
                        <xsd:element name="customers" type="tns:CustomerList"/>
                    </xsd:sequence>
                </xsd:complexType>
            </xsd:element>
            <xsd:element name="ListCustomers">
                <xsd:complexType>
                    <xsd:sequence>
                        <xsd:element name="code" type="xsd:string" minOccurs="0"/>
                        <xsd:element name="status" type="xsd:string" minOccurs="0"/>
                        <xsd:element name="taxNumber" type="xsd:string" minOccurs="0"/>
                        <xsd:element name="maxResults" type="xsd:int" minOccurs="0"/>
                        <xsd:element name="offset" type="xsd:int" minOccurs="0"/>
                    </xsd:sequence>
                </xsd:complexType>
            </xsd:element>
            <xsd:element name="ListCustomersResponse">
                <xsd:complexType>
                    <xsd:sequence>
                        <xsd:element name="customers" type="tns:CustomerList"/>
                    </xsd:sequence>
                </xsd:complexType>
            </xsd:element>
            <xsd:element name="GetCustomer">
                <xsd:complexType>
                    <xsd:sequence>
                        <xsd:element name="id" type="xsd:int" minOccurs="0"/>
                        <xsd:element name="code" type="xsd:string" minOccurs="0"/>
                    </xsd:sequence>
                </xsd:complexType>
            </xsd:element>
            <xsd:element name="GetCustomerResponse">
                <xsd:complexType>
                    <xsd:sequence>
                        <xsd:element name="customer" type="tns:Customer"/>
                    </xsd:sequence>
                </xsd:complexType>
            </xsd:element>
            <xsd:element name="CreateCustomer">
                <xsd:complexType>
                    <xsd:sequence>
                        <xsd:element name="customer" type="tns:Customer"/>
                    </xsd:sequence>
                </xsd:complexType>
            </xsd:element>
            <xsd:element name="CreateCustomerResponse">
                <xsd:complexType>
                    <xsd:sequence>
                        <xsd:element name="customer" type="tns:Customer"/>
                    </xsd:sequence>
                </xsd:complexType>
            </xsd:element>
            <xsd:element name="UpdateCustomer">
                <xsd:complexType>
                    <xsd:sequence>
                        <xsd:element name="customer" type="tns:Customer"/>
                    </xsd:sequence>
                </xsd:complexType>
            </xsd:element>
            <xsd:element name="UpdateCustomerResponse">
                <xsd:complexType>
                    <xsd:sequence>
                        <xsd:element name="customer" type="tns:Customer"/>
                    </xsd:sequence>
                </xsd:complexType>
            </xsd:element>

            <xsd:complexType name="CustomerList">
                <xsd:sequence>
                    <xsd:element name="customer" type="tns:Customer" maxOccurs="unbounded"/>
                </xsd:sequence>
            </xsd:complexType>

            <xsd:complexType name="Customer">
                <xsd:sequence>
                    <xsd:element name="id" type="xsd:int"/>
                    <xsd:element name="code" type="xsd:string"/>
                    <xsd:element name="name" type="xsd:string"/>
                    <xsd:element name="taxNumber" type="xsd:string" minOccurs="0"/>
                    <xsd:element name="email" type="xsd:string" minOccurs="0"/>
                    <xsd:element name="phone" type="xsd:string" minOccurs="0"/>
                    <xsd:element name="address" type="xsd:string" minOccurs="0"/>
                    <xsd:element name="currency" type="xsd:string"/>
                    <xsd:element name="riskLimit" type="xsd:decimal" minOccurs="0"/>
                    <xsd:element name="status" type="xsd:string"/>
                    <xsd:element name="createdAt" type="xsd:dateTime" minOccurs="0"/>
                    <xsd:element name="updatedAt" type="xsd:dateTime" minOccurs="0"/>
                </xsd:sequence>
            </xsd:complexType>

            <!-- Stocks -->
            <xsd:element name="GetStocks">
                <xsd:complexType>
                    <xsd:sequence>
                        <xsd:element name="code" type="xsd:string" minOccurs="0"/>
                        <xsd:element name="status" type="xsd:string" minOccurs="0"/>
                        <xsd:element name="nameContains" type="xsd:string" minOccurs="0"/>
                        <xsd:element name="maxResults" type="xsd:int" minOccurs="0"/>
                        <xsd:element name="offset" type="xsd:int" minOccurs="0"/>
                    </xsd:sequence>
                </xsd:complexType>
            </xsd:element>
            <xsd:element name="GetStocksResponse">
                <xsd:complexType>
                    <xsd:sequence>
                        <xsd:element name="stocks" type="tns:StockList"/>
                    </xsd:sequence>
                </xsd:complexType>
            </xsd:element>
            <xsd:element name="ListStocks">
                <xsd:complexType>
                    <xsd:sequence>
                        <xsd:element name="code" type="xsd:string" minOccurs="0"/>
                        <xsd:element name="status" type="xsd:string" minOccurs="0"/>
                        <xsd:element name="nameContains" type="xsd:string" minOccurs="0"/>
                        <xsd:element name="maxResults" type="xsd:int" minOccurs="0"/>
                        <xsd:element name="offset" type="xsd:int" minOccurs="0"/>
                    </xsd:sequence>
                </xsd:complexType>
            </xsd:element>
            <xsd:element name="ListStocksResponse">
                <xsd:complexType>
                    <xsd:sequence>
                        <xsd:element name="stocks" type="tns:StockList"/>
                    </xsd:sequence>
                </xsd:complexType>
            </xsd:element>
            <xsd:element name="GetStock">
                <xsd:complexType>
                    <xsd:sequence>
                        <xsd:element name="id" type="xsd:int" minOccurs="0"/>
                        <xsd:element name="code" type="xsd:string" minOccurs="0"/>
                    </xsd:sequence>
                </xsd:complexType>
            </xsd:element>
            <xsd:element name="GetStockResponse">
                <xsd:complexType>
                    <xsd:sequence>
                        <xsd:element name="stock" type="tns:StockItem"/>
                    </xsd:sequence>
                </xsd:complexType>
            </xsd:element>
            <xsd:element name="CreateStock">
                <xsd:complexType>
                    <xsd:sequence>
                        <xsd:element name="stock" type="tns:StockItem"/>
                    </xsd:sequence>
                </xsd:complexType>
            </xsd:element>
            <xsd:element name="CreateStockResponse">
                <xsd:complexType>
                    <xsd:sequence>
                        <xsd:element name="stock" type="tns:StockItem"/>
                    </xsd:sequence>
                </xsd:complexType>
            </xsd:element>
            <xsd:element name="UpdateStock">
                <xsd:complexType>
                    <xsd:sequence>
                        <xsd:element name="stock" type="tns:StockItem"/>
                    </xsd:sequence>
                </xsd:complexType>
            </xsd:element>
            <xsd:element name="UpdateStockResponse">
                <xsd:complexType>
                    <xsd:sequence>
                        <xsd:element name="stock" type="tns:StockItem"/>
                    </xsd:sequence>
                </xsd:complexType>
            </xsd:element>

            <xsd:complexType name="StockList">
                <xsd:sequence>
                    <xsd:element name="stock" type="tns:StockItem" maxOccurs="unbounded"/>
                </xsd:sequence>
            </xsd:complexType>

            <xsd:complexType name="StockItem">
                <xsd:sequence>
                    <xsd:element name="id" type="xsd:int"/>
                    <xsd:element name="code" type="xsd:string"/>
                    <xsd:element name="name" type="xsd:string"/>
                    <xsd:element name="unit" type="xsd:string"/>
                    <xsd:element name="vatRate" type="xsd:decimal"/>
                    <xsd:element name="price" type="xsd:decimal"/>
                    <xsd:element name="stockOnHand" type="xsd:decimal"/>
                    <xsd:element name="minStock" type="xsd:decimal"/>
                    <xsd:element name="status" type="xsd:string"/>
                    <xsd:element name="createdAt" type="xsd:dateTime" minOccurs="0"/>
                    <xsd:element name="updatedAt" type="xsd:dateTime" minOccurs="0"/>
                </xsd:sequence>
            </xsd:complexType>

            <!-- Cash Accounts -->
            <xsd:element name="GetCashAccounts">
                <xsd:complexType>
                    <xsd:sequence>
                        <xsd:element name="status" type="xsd:string" minOccurs="0"/>
                        <xsd:element name="currency" type="xsd:string" minOccurs="0"/>
                        <xsd:element name="maxResults" type="xsd:int" minOccurs="0"/>
                        <xsd:element name="offset" type="xsd:int" minOccurs="0"/>
                    </xsd:sequence>
                </xsd:complexType>
            </xsd:element>
            <xsd:element name="GetCashAccountsResponse">
                <xsd:complexType>
                    <xsd:sequence>
                        <xsd:element name="cashAccounts" type="tns:CashAccountList"/>
                    </xsd:sequence>
                </xsd:complexType>
            </xsd:element>
            <xsd:element name="ListCashAccounts">
                <xsd:complexType>
                    <xsd:sequence>
                        <xsd:element name="status" type="xsd:string" minOccurs="0"/>
                        <xsd:element name="currency" type="xsd:string" minOccurs="0"/>
                        <xsd:element name="maxResults" type="xsd:int" minOccurs="0"/>
                        <xsd:element name="offset" type="xsd:int" minOccurs="0"/>
                    </xsd:sequence>
                </xsd:complexType>
            </xsd:element>
            <xsd:element name="ListCashAccountsResponse">
                <xsd:complexType>
                    <xsd:sequence>
                        <xsd:element name="cashAccounts" type="tns:CashAccountList"/>
                    </xsd:sequence>
                </xsd:complexType>
            </xsd:element>
            <xsd:element name="GetCashAccount">
                <xsd:complexType>
                    <xsd:sequence>
                        <xsd:element name="id" type="xsd:int" minOccurs="0"/>
                        <xsd:element name="code" type="xsd:string" minOccurs="0"/>
                    </xsd:sequence>
                </xsd:complexType>
            </xsd:element>
            <xsd:element name="GetCashAccountResponse">
                <xsd:complexType>
                    <xsd:sequence>
                        <xsd:element name="cashAccount" type="tns:CashAccount"/>
                    </xsd:sequence>
                </xsd:complexType>
            </xsd:element>
            <xsd:element name="CreateCashAccount">
                <xsd:complexType>
                    <xsd:sequence>
                        <xsd:element name="cashAccount" type="tns:CashAccount"/>
                    </xsd:sequence>
                </xsd:complexType>
            </xsd:element>
            <xsd:element name="CreateCashAccountResponse">
                <xsd:complexType>
                    <xsd:sequence>
                        <xsd:element name="cashAccount" type="tns:CashAccount"/>
                    </xsd:sequence>
                </xsd:complexType>
            </xsd:element>
            <xsd:element name="UpdateCashAccount">
                <xsd:complexType>
                    <xsd:sequence>
                        <xsd:element name="cashAccount" type="tns:CashAccount"/>
                    </xsd:sequence>
                </xsd:complexType>
            </xsd:element>
            <xsd:element name="UpdateCashAccountResponse">
                <xsd:complexType>
                    <xsd:sequence>
                        <xsd:element name="cashAccount" type="tns:CashAccount"/>
                    </xsd:sequence>
                </xsd:complexType>
            </xsd:element>

            <xsd:complexType name="CashAccountList">
                <xsd:sequence>
                    <xsd:element name="cashAccount" type="tns:CashAccount" maxOccurs="unbounded"/>
                </xsd:sequence>
            </xsd:complexType>

            <xsd:complexType name="CashAccount">
                <xsd:sequence>
                    <xsd:element name="id" type="xsd:int"/>
                    <xsd:element name="code" type="xsd:string"/>
                    <xsd:element name="name" type="xsd:string"/>
                    <xsd:element name="currency" type="xsd:string"/>
                    <xsd:element name="balance" type="xsd:decimal"/>
                    <xsd:element name="status" type="xsd:string"/>
                    <xsd:element name="createdAt" type="xsd:dateTime" minOccurs="0"/>
                    <xsd:element name="updatedAt" type="xsd:dateTime" minOccurs="0"/>
                </xsd:sequence>
            </xsd:complexType>

            <!-- Orders -->
            <xsd:element name="GetOrders">
                <xsd:complexType>
                    <xsd:sequence>
                        <xsd:element name="customerId" type="xsd:int" minOccurs="0"/>
                        <xsd:element name="status" type="xsd:string" minOccurs="0"/>
                        <xsd:element name="fromDate" type="xsd:date" minOccurs="0"/>
                        <xsd:element name="toDate" type="xsd:date" minOccurs="0"/>
                        <xsd:element name="maxResults" type="xsd:int" minOccurs="0"/>
                        <xsd:element name="offset" type="xsd:int" minOccurs="0"/>
                    </xsd:sequence>
                </xsd:complexType>
            </xsd:element>
            <xsd:element name="GetOrdersResponse">
                <xsd:complexType>
                    <xsd:sequence>
                        <xsd:element name="orders" type="tns:OrderList"/>
                    </xsd:sequence>
                </xsd:complexType>
            </xsd:element>
            <xsd:element name="ListOrders">
                <xsd:complexType>
                    <xsd:sequence>
                        <xsd:element name="customerId" type="xsd:int" minOccurs="0"/>
                        <xsd:element name="status" type="xsd:string" minOccurs="0"/>
                        <xsd:element name="fromDate" type="xsd:date" minOccurs="0"/>
                        <xsd:element name="toDate" type="xsd:date" minOccurs="0"/>
                        <xsd:element name="maxResults" type="xsd:int" minOccurs="0"/>
                        <xsd:element name="offset" type="xsd:int" minOccurs="0"/>
                    </xsd:sequence>
                </xsd:complexType>
            </xsd:element>
            <xsd:element name="ListOrdersResponse">
                <xsd:complexType>
                    <xsd:sequence>
                        <xsd:element name="orders" type="tns:OrderList"/>
                    </xsd:sequence>
                </xsd:complexType>
            </xsd:element>
            <xsd:element name="GetOrder">
                <xsd:complexType>
                    <xsd:sequence>
                        <xsd:element name="id" type="xsd:int" minOccurs="0"/>
                        <xsd:element name="orderNo" type="xsd:string" minOccurs="0"/>
                    </xsd:sequence>
                </xsd:complexType>
            </xsd:element>
            <xsd:element name="GetOrderResponse">
                <xsd:complexType>
                    <xsd:sequence>
                        <xsd:element name="order" type="tns:Order"/>
                    </xsd:sequence>
                </xsd:complexType>
            </xsd:element>
            <xsd:element name="CreateOrder">
                <xsd:complexType>
                    <xsd:sequence>
                        <xsd:element name="order" type="tns:Order"/>
                    </xsd:sequence>
                </xsd:complexType>
            </xsd:element>
            <xsd:element name="CreateOrderResponse">
                <xsd:complexType>
                    <xsd:sequence>
                        <xsd:element name="order" type="tns:Order"/>
                    </xsd:sequence>
                </xsd:complexType>
            </xsd:element>
            <xsd:element name="UpdateOrder">
                <xsd:complexType>
                    <xsd:sequence>
                        <xsd:element name="order" type="tns:Order"/>
                    </xsd:sequence>
                </xsd:complexType>
            </xsd:element>
            <xsd:element name="UpdateOrderResponse">
                <xsd:complexType>
                    <xsd:sequence>
                        <xsd:element name="order" type="tns:Order"/>
                    </xsd:sequence>
                </xsd:complexType>
            </xsd:element>
            <xsd:element name="ApproveOrder">
                <xsd:complexType>
                    <xsd:sequence>
                        <xsd:element name="id" type="xsd:int" minOccurs="0"/>
                        <xsd:element name="orderNo" type="xsd:string" minOccurs="0"/>
                    </xsd:sequence>
                </xsd:complexType>
            </xsd:element>
            <xsd:element name="ApproveOrderResponse">
                <xsd:complexType>
                    <xsd:sequence>
                        <xsd:element name="order" type="tns:Order"/>
                    </xsd:sequence>
                </xsd:complexType>
            </xsd:element>
            <xsd:element name="CancelOrder">
                <xsd:complexType>
                    <xsd:sequence>
                        <xsd:element name="id" type="xsd:int" minOccurs="0"/>
                        <xsd:element name="orderNo" type="xsd:string" minOccurs="0"/>
                    </xsd:sequence>
                </xsd:complexType>
            </xsd:element>
            <xsd:element name="CancelOrderResponse">
                <xsd:complexType>
                    <xsd:sequence>
                        <xsd:element name="order" type="tns:Order"/>
                    </xsd:sequence>
                </xsd:complexType>
            </xsd:element>

            <xsd:complexType name="OrderLineList">
                <xsd:sequence>
                    <xsd:element name="line" type="tns:OrderLine" maxOccurs="unbounded"/>
                </xsd:sequence>
            </xsd:complexType>

            <xsd:complexType name="OrderLine">
                <xsd:sequence>
                    <xsd:element name="lineNo" type="xsd:int"/>
                    <xsd:element name="stockId" type="xsd:int"/>
                    <xsd:element name="stockCode" type="xsd:string"/>
                    <xsd:element name="description" type="xsd:string"/>
                    <xsd:element name="unit" type="xsd:string"/>
                    <xsd:element name="quantity" type="xsd:decimal"/>
                    <xsd:element name="unitPrice" type="xsd:decimal"/>
                    <xsd:element name="discountRate" type="xsd:decimal"/>
                    <xsd:element name="taxRate" type="xsd:decimal"/>
                    <xsd:element name="lineTotal" type="xsd:decimal"/>
                </xsd:sequence>
            </xsd:complexType>

            <xsd:complexType name="OrderList">
                <xsd:sequence>
                    <xsd:element name="order" type="tns:Order" maxOccurs="unbounded"/>
                </xsd:sequence>
            </xsd:complexType>

            <xsd:complexType name="Order">
                <xsd:sequence>
                    <xsd:element name="id" type="xsd:int"/>
                    <xsd:element name="orderNo" type="xsd:string"/>
                    <xsd:element name="customerId" type="xsd:int"/>
                    <xsd:element name="customerCode" type="xsd:string" minOccurs="0"/>
                    <xsd:element name="status" type="xsd:string"/>
                    <xsd:element name="orderDate" type="xsd:date"/>
                    <xsd:element name="currency" type="xsd:string"/>
                    <xsd:element name="notes" type="xsd:string" minOccurs="0"/>
                    <xsd:element name="lines" type="tns:OrderLineList"/>
                    <xsd:element name="subtotal" type="xsd:decimal"/>
                    <xsd:element name="discountTotal" type="xsd:decimal"/>
                    <xsd:element name="taxTotal" type="xsd:decimal"/>
                    <xsd:element name="grandTotal" type="xsd:decimal"/>
                    <xsd:element name="createdAt" type="xsd:dateTime" minOccurs="0"/>
                    <xsd:element name="updatedAt" type="xsd:dateTime" minOccurs="0"/>
                </xsd:sequence>
            </xsd:complexType>

            <!-- Invoices -->
            <xsd:element name="GetInvoices">
                <xsd:complexType>
                    <xsd:sequence>
                        <xsd:element name="customerId" type="xsd:int" minOccurs="0"/>
                        <xsd:element name="status" type="xsd:string" minOccurs="0"/>
                        <xsd:element name="fromDate" type="xsd:date" minOccurs="0"/>
                        <xsd:element name="toDate" type="xsd:date" minOccurs="0"/>
                        <xsd:element name="maxResults" type="xsd:int" minOccurs="0"/>
                        <xsd:element name="offset" type="xsd:int" minOccurs="0"/>
                    </xsd:sequence>
                </xsd:complexType>
            </xsd:element>
            <xsd:element name="GetInvoicesResponse">
                <xsd:complexType>
                    <xsd:sequence>
                        <xsd:element name="invoices" type="tns:InvoiceList"/>
                    </xsd:sequence>
                </xsd:complexType>
            </xsd:element>
            <xsd:element name="ListInvoices">
                <xsd:complexType>
                    <xsd:sequence>
                        <xsd:element name="customerId" type="xsd:int" minOccurs="0"/>
                        <xsd:element name="status" type="xsd:string" minOccurs="0"/>
                        <xsd:element name="fromDate" type="xsd:date" minOccurs="0"/>
                        <xsd:element name="toDate" type="xsd:date" minOccurs="0"/>
                        <xsd:element name="maxResults" type="xsd:int" minOccurs="0"/>
                        <xsd:element name="offset" type="xsd:int" minOccurs="0"/>
                    </xsd:sequence>
                </xsd:complexType>
            </xsd:element>
            <xsd:element name="ListInvoicesResponse">
                <xsd:complexType>
                    <xsd:sequence>
                        <xsd:element name="invoices" type="tns:InvoiceList"/>
                    </xsd:sequence>
                </xsd:complexType>
            </xsd:element>
            <xsd:element name="GetInvoice">
                <xsd:complexType>
                    <xsd:sequence>
                        <xsd:element name="id" type="xsd:int" minOccurs="0"/>
                        <xsd:element name="invoiceNo" type="xsd:string" minOccurs="0"/>
                    </xsd:sequence>
                </xsd:complexType>
            </xsd:element>
            <xsd:element name="GetInvoiceResponse">
                <xsd:complexType>
                    <xsd:sequence>
                        <xsd:element name="invoice" type="tns:Invoice"/>
                    </xsd:sequence>
                </xsd:complexType>
            </xsd:element>
            <xsd:element name="CreateInvoice">
                <xsd:complexType>
                    <xsd:sequence>
                        <xsd:element name="invoice" type="tns:Invoice"/>
                    </xsd:sequence>
                </xsd:complexType>
            </xsd:element>
            <xsd:element name="CreateInvoiceResponse">
                <xsd:complexType>
                    <xsd:sequence>
                        <xsd:element name="invoice" type="tns:Invoice"/>
                    </xsd:sequence>
                </xsd:complexType>
            </xsd:element>
            <xsd:element name="CreateInvoiceFromOrder">
                <xsd:complexType>
                    <xsd:sequence>
                        <xsd:element name="orderId" type="xsd:int" minOccurs="0"/>
                        <xsd:element name="orderNo" type="xsd:string" minOccurs="0"/>
                        <xsd:element name="invoiceDate" type="xsd:date" minOccurs="0"/>
                        <xsd:element name="dueDate" type="xsd:date" minOccurs="0"/>
                    </xsd:sequence>
                </xsd:complexType>
            </xsd:element>
            <xsd:element name="CreateInvoiceFromOrderResponse">
                <xsd:complexType>
                    <xsd:sequence>
                        <xsd:element name="invoice" type="tns:Invoice"/>
                    </xsd:sequence>
                </xsd:complexType>
            </xsd:element>
            <xsd:element name="CancelInvoice">
                <xsd:complexType>
                    <xsd:sequence>
                        <xsd:element name="id" type="xsd:int" minOccurs="0"/>
                        <xsd:element name="invoiceNo" type="xsd:string" minOccurs="0"/>
                    </xsd:sequence>
                </xsd:complexType>
            </xsd:element>
            <xsd:element name="CancelInvoiceResponse">
                <xsd:complexType>
                    <xsd:sequence>
                        <xsd:element name="invoice" type="tns:Invoice"/>
                    </xsd:sequence>
                </xsd:complexType>
            </xsd:element>

            <xsd:complexType name="InvoiceLineList">
                <xsd:sequence>
                    <xsd:element name="line" type="tns:InvoiceLine" maxOccurs="unbounded"/>
                </xsd:sequence>
            </xsd:complexType>

            <xsd:complexType name="InvoiceLine">
                <xsd:sequence>
                    <xsd:element name="lineNo" type="xsd:int"/>
                    <xsd:element name="stockId" type="xsd:int"/>
                    <xsd:element name="stockCode" type="xsd:string"/>
                    <xsd:element name="description" type="xsd:string"/>
                    <xsd:element name="unit" type="xsd:string"/>
                    <xsd:element name="quantity" type="xsd:decimal"/>
                    <xsd:element name="unitPrice" type="xsd:decimal"/>
                    <xsd:element name="discountRate" type="xsd:decimal"/>
                    <xsd:element name="taxRate" type="xsd:decimal"/>
                    <xsd:element name="lineTotal" type="xsd:decimal"/>
                </xsd:sequence>
            </xsd:complexType>

            <xsd:complexType name="InvoiceList">
                <xsd:sequence>
                    <xsd:element name="invoice" type="tns:Invoice" maxOccurs="unbounded"/>
                </xsd:sequence>
            </xsd:complexType>

            <xsd:complexType name="Invoice">
                <xsd:sequence>
                    <xsd:element name="id" type="xsd:int"/>
                    <xsd:element name="invoiceNo" type="xsd:string"/>
                    <xsd:element name="orderId" type="xsd:int"/>
                    <xsd:element name="customerId" type="xsd:int"/>
                    <xsd:element name="customerCode" type="xsd:string" minOccurs="0"/>
                    <xsd:element name="status" type="xsd:string"/>
                    <xsd:element name="invoiceDate" type="xsd:date"/>
                    <xsd:element name="dueDate" type="xsd:date"/>
                    <xsd:element name="currency" type="xsd:string"/>
                    <xsd:element name="lines" type="tns:InvoiceLineList"/>
                    <xsd:element name="subtotal" type="xsd:decimal"/>
                    <xsd:element name="discountTotal" type="xsd:decimal"/>
                    <xsd:element name="taxTotal" type="xsd:decimal"/>
                    <xsd:element name="grandTotal" type="xsd:decimal"/>
                    <xsd:element name="paidAmount" type="xsd:decimal"/>
                    <xsd:element name="createdAt" type="xsd:dateTime" minOccurs="0"/>
                    <xsd:element name="updatedAt" type="xsd:dateTime" minOccurs="0"/>
                </xsd:sequence>
            </xsd:complexType>

            <!-- Cash Transactions -->
            <xsd:element name="GetCashTransactions">
                <xsd:complexType>
                    <xsd:sequence>
                        <xsd:element name="cashAccountId" type="xsd:int" minOccurs="0"/>
                        <xsd:element name="type" type="xsd:string" minOccurs="0"/>
                        <xsd:element name="fromDate" type="xsd:date" minOccurs="0"/>
                        <xsd:element name="toDate" type="xsd:date" minOccurs="0"/>
                        <xsd:element name="maxResults" type="xsd:int" minOccurs="0"/>
                        <xsd:element name="offset" type="xsd:int" minOccurs="0"/>
                    </xsd:sequence>
                </xsd:complexType>
            </xsd:element>
            <xsd:element name="GetCashTransactionsResponse">
                <xsd:complexType>
                    <xsd:sequence>
                        <xsd:element name="cashTransactions" type="tns:CashTransactionList"/>
                    </xsd:sequence>
                </xsd:complexType>
            </xsd:element>
            <xsd:element name="ListCashTransactions">
                <xsd:complexType>
                    <xsd:sequence>
                        <xsd:element name="cashAccountId" type="xsd:int" minOccurs="0"/>
                        <xsd:element name="type" type="xsd:string" minOccurs="0"/>
                        <xsd:element name="fromDate" type="xsd:date" minOccurs="0"/>
                        <xsd:element name="toDate" type="xsd:date" minOccurs="0"/>
                        <xsd:element name="maxResults" type="xsd:int" minOccurs="0"/>
                        <xsd:element name="offset" type="xsd:int" minOccurs="0"/>
                    </xsd:sequence>
                </xsd:complexType>
            </xsd:element>
            <xsd:element name="ListCashTransactionsResponse">
                <xsd:complexType>
                    <xsd:sequence>
                        <xsd:element name="cashTransactions" type="tns:CashTransactionList"/>
                    </xsd:sequence>
                </xsd:complexType>
            </xsd:element>
            <xsd:element name="GetCashTransaction">
                <xsd:complexType>
                    <xsd:sequence>
                        <xsd:element name="id" type="xsd:int" minOccurs="0"/>
                        <xsd:element name="txnNo" type="xsd:string" minOccurs="0"/>
                    </xsd:sequence>
                </xsd:complexType>
            </xsd:element>
            <xsd:element name="GetCashTransactionResponse">
                <xsd:complexType>
                    <xsd:sequence>
                        <xsd:element name="cashTransaction" type="tns:CashTransaction"/>
                    </xsd:sequence>
                </xsd:complexType>
            </xsd:element>
            <xsd:element name="CreateCashTransaction">
                <xsd:complexType>
                    <xsd:sequence>
                        <xsd:element name="cashTransaction" type="tns:CashTransaction"/>
                    </xsd:sequence>
                </xsd:complexType>
            </xsd:element>
            <xsd:element name="CreateCashTransactionResponse">
                <xsd:complexType>
                    <xsd:sequence>
                        <xsd:element name="cashTransaction" type="tns:CashTransaction"/>
                    </xsd:sequence>
                </xsd:complexType>
            </xsd:element>
            <xsd:element name="ReverseCashTransaction">
                <xsd:complexType>
                    <xsd:sequence>
                        <xsd:element name="id" type="xsd:int" minOccurs="0"/>
                        <xsd:element name="txnNo" type="xsd:string" minOccurs="0"/>
                    </xsd:sequence>
                </xsd:complexType>
            </xsd:element>
            <xsd:element name="ReverseCashTransactionResponse">
                <xsd:complexType>
                    <xsd:sequence>
                        <xsd:element name="cashTransaction" type="tns:CashTransaction"/>
                    </xsd:sequence>
                </xsd:complexType>
            </xsd:element>

            <xsd:complexType name="CashTransactionList">
                <xsd:sequence>
                    <xsd:element name="cashTransaction" type="tns:CashTransaction" maxOccurs="unbounded"/>
                </xsd:sequence>
            </xsd:complexType>

            <xsd:complexType name="CashTransaction">
                <xsd:sequence>
                    <xsd:element name="id" type="xsd:int"/>
                    <xsd:element name="txnNo" type="xsd:string"/>
                    <xsd:element name="cashAccountId" type="xsd:int"/>
                    <xsd:element name="cashAccountCode" type="xsd:string" minOccurs="0"/>
                    <xsd:element name="customerId" type="xsd:int"/>
                    <xsd:element name="customerCode" type="xsd:string" minOccurs="0"/>
                    <xsd:element name="invoiceId" type="xsd:int" minOccurs="0"/>
                    <xsd:element name="invoiceNo" type="xsd:string" minOccurs="0"/>
                    <xsd:element name="type" type="xsd:string"/>
                    <xsd:element name="amount" type="xsd:decimal"/>
                    <xsd:element name="currency" type="xsd:string"/>
                    <xsd:element name="method" type="xsd:string"/>
                    <xsd:element name="transactionDate" type="xsd:date"/>
                    <xsd:element name="description" type="xsd:string" minOccurs="0"/>
                    <xsd:element name="status" type="xsd:string" minOccurs="0"/>
                    <xsd:element name="reversalOf" type="xsd:int" minOccurs="0"/>
                    <xsd:element name="createdAt" type="xsd:dateTime" minOccurs="0"/>
                </xsd:sequence>
            </xsd:complexType>
`

const publicSchemaBody = `            <xsd:element name="GetCountries">
                <xsd:complexType>
                    <xsd:sequence>
                        <xsd:element name="continent" type="xsd:string" minOccurs="0"/>
                    </xsd:sequence>
                </xsd:complexType>
            </xsd:element>
            <xsd:element name="GetCountriesResponse">
                <xsd:complexType>
                    <xsd:sequence>
                        <xsd:element name="countries" type="tns:CountryList"/>
                    </xsd:sequence>
                </xsd:complexType>
            </xsd:element>

            <xsd:element name="GetCities">
                <xsd:complexType>
                    <xsd:sequence>
                        <xsd:element name="countryCode" type="xsd:string"/>
                    </xsd:sequence>
                </xsd:complexType>
            </xsd:element>
            <xsd:element name="GetCitiesResponse">
                <xsd:complexType>
                    <xsd:sequence>
                        <xsd:element name="cities" type="tns:CityList"/>
                    </xsd:sequence>
                </xsd:complexType>
            </xsd:element>

            <xsd:complexType name="CountryList">
                <xsd:sequence>
                    <xsd:element name="country" type="tns:Country" maxOccurs="unbounded"/>
                </xsd:sequence>
            </xsd:complexType>

            <xsd:complexType name="Country">
                <xsd:sequence>
                    <xsd:element name="code" type="xsd:string"/>
                    <xsd:element name="name" type="xsd:string"/>
                    <xsd:element name="continent" type="xsd:string"/>
                </xsd:sequence>
            </xsd:complexType>

            <xsd:complexType name="CityList">
                <xsd:sequence>
                    <xsd:element name="city" type="tns:City" maxOccurs="unbounded"/>
                </xsd:sequence>
            </xsd:complexType>

            <xsd:complexType name="City">
                <xsd:sequence>
                    <xsd:element name="countryCode" type="xsd:string"/>
                    <xsd:element name="name" type="xsd:string"/>
                    <xsd:element name="population" type="xsd:int"/>
                </xsd:sequence>
            </xsd:complexType>
`

func buildWSDL(cfg wsdlConfig) string {
    ops := make([]string, 0, len(erpOperations)+len(publicOperations)+len(sessionOperations))
    if cfg.IncludeSessionOps {
        ops = append(ops, sessionOperations...)
    }
    ops = append(ops, erpOperations...)
    if cfg.IncludePublic {
        ops = append(ops, publicOperations...)
    }

    var sb strings.Builder
    sb.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")

    // Root element with all required namespace declarations
    sb.WriteString(fmt.Sprintf("<definitions xmlns=\"http://schemas.xmlsoap.org/wsdl/\" xmlns:soap=\"http://schemas.xmlsoap.org/wsdl/soap/\" xmlns:soap12=\"http://schemas.xmlsoap.org/wsdl/soap12/\" xmlns:tns=\"%s\" xmlns:xsd=\"http://www.w3.org/2001/XMLSchema\"", cfg.Namespace))
    if cfg.Policy != PolicyNone {
        sb.WriteString(" xmlns:wsp=\"http://schemas.xmlsoap.org/ws/2004/09/policy\"")
        switch cfg.Policy {
        case PolicyWSSEUsernameToken:
            sb.WriteString(" xmlns:sp=\"http://docs.oasis-open.org/ws-sx/ws-securitypolicy/200702\"")
        case PolicyNTLMNegotiate:
            sb.WriteString(" xmlns:http=\"http://schemas.microsoft.com/ws/06/2004/policy/http\"")
        case PolicyBasicAuth:
            sb.WriteString(" xmlns:http=\"http://schemas.microsoft.com/ws/06/2004/policy/http\"")
        }
    }
    sb.WriteString(fmt.Sprintf(" name=\"%s\" targetNamespace=\"%s\">\n", cfg.ServiceName, cfg.Namespace))

    // WS-Policy declaration
    if cfg.Policy != PolicyNone {
        sb.WriteString(buildWSPolicy(cfg.Policy))
        sb.WriteString("\n")
    }

    sb.WriteString("    <types>\n")
    sb.WriteString(fmt.Sprintf("        <xsd:schema targetNamespace=\"%s\" elementFormDefault=\"qualified\">\n", cfg.Namespace))
    if cfg.IncludeSessionOps {
        sb.WriteString(sessionSchemaBody)
    }
    sb.WriteString(erpSchemaBody)
    if cfg.IncludePublic {
        sb.WriteString(publicSchemaBody)
    }
    sb.WriteString("        </xsd:schema>\n")
    sb.WriteString("    </types>\n\n")

    sb.WriteString(buildMessages(ops))
    if cfg.IncludeSessionHeader {
        sb.WriteString("    <message name=\"SessionHeader\"><part name=\"parameters\" element=\"tns:SessionHeader\"/></message>\n")
    }
    sb.WriteString("\n")

    sb.WriteString(buildPortType(cfg.PortTypeName, ops))
    sb.WriteString("\n")

    // SOAP 1.1 binding (with policy reference if applicable)
    sb.WriteString(buildBinding(cfg.BindingName, cfg.PortTypeName, cfg.Endpoint, ops, cfg.IncludeSessionHeader, "soap", cfg.Policy))
    // SOAP 1.2 binding — always included for all endpoints
    sb.WriteString("\n")
    sb.WriteString(buildBinding(cfg.BindingName+"Soap12", cfg.PortTypeName, cfg.Endpoint, ops, cfg.IncludeSessionHeader, "soap12", cfg.Policy))
    sb.WriteString("\n")

    sb.WriteString(buildService(cfg))
    sb.WriteString("</definitions>")
    return sb.String()
}

// buildWSPolicy generates WS-Policy XML block based on the policy type.
// Follows WS-SecurityPolicy 1.2 (OASIS) for WSSE, and Microsoft WS-Policy
// extensions for HTTP-level auth (Basic, NTLM/Negotiate).
func buildWSPolicy(policy PolicyType) string {
    switch policy {
    case PolicyWSSEUsernameToken:
        return `    <wsp:Policy wsu:Id="UsernameTokenPolicy" xmlns:wsu="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-utility-1.0.xsd">
        <wsp:ExactlyOne>
            <wsp:All>
                <sp:TransportBinding>
                    <wsp:Policy>
                        <sp:TransportToken>
                            <wsp:Policy>
                                <sp:HttpsToken RequireClientCertificate="false"/>
                            </wsp:Policy>
                        </sp:TransportToken>
                    </wsp:Policy>
                </sp:TransportBinding>
                <sp:SignedSupportingTokens>
                    <wsp:Policy>
                        <sp:UsernameToken sp:IncludeToken="http://docs.oasis-open.org/ws-sx/ws-securitypolicy/200702/IncludeToken/AlwaysToRecipient">
                            <wsp:Policy>
                                <sp:WssUsernameToken10/>
                            </wsp:Policy>
                        </sp:UsernameToken>
                    </wsp:Policy>
                </sp:SignedSupportingTokens>
            </wsp:All>
        </wsp:ExactlyOne>
    </wsp:Policy>`
    case PolicyBasicAuth:
        return `    <wsp:Policy wsu:Id="BasicAuthPolicy" xmlns:wsu="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-utility-1.0.xsd">
        <wsp:ExactlyOne>
            <wsp:All>
                <http:BasicAuthentication/>
            </wsp:All>
        </wsp:ExactlyOne>
    </wsp:Policy>`
    case PolicyNTLMNegotiate:
        return `    <wsp:Policy wsu:Id="NTLMAuthPolicy" xmlns:wsu="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-utility-1.0.xsd">
        <wsp:ExactlyOne>
            <wsp:All>
                <http:NegotiateAuthentication/>
            </wsp:All>
        </wsp:ExactlyOne>
    </wsp:Policy>`
    default:
        return ""
    }
}

func buildMessages(ops []string) string {
    var sb strings.Builder
    for _, op := range ops {
        sb.WriteString(fmt.Sprintf("    <message name=\"%sInput\"><part name=\"parameters\" element=\"tns:%s\"/></message>\n", op, op))
        sb.WriteString(fmt.Sprintf("    <message name=\"%sOutput\"><part name=\"parameters\" element=\"tns:%sResponse\"/></message>\n", op, op))
    }
    return sb.String()
}

func buildPortType(name string, ops []string) string {
    var sb strings.Builder
    sb.WriteString(fmt.Sprintf("    <portType name=\"%s\">\n", name))
    for _, op := range ops {
        sb.WriteString(fmt.Sprintf("        <operation name=\"%s\">\n", op))
        sb.WriteString(fmt.Sprintf("            <input message=\"tns:%sInput\"/>\n", op))
        sb.WriteString(fmt.Sprintf("            <output message=\"tns:%sOutput\"/>\n", op))
        sb.WriteString("        </operation>\n")
    }
    sb.WriteString("    </portType>\n")
    return sb.String()
}

func buildBinding(name, portTypeName, endpoint string, ops []string, includeHeader bool, soapPrefix string, policy PolicyType) string {
    var sb strings.Builder
    sb.WriteString(fmt.Sprintf("    <binding name=\"%s\" type=\"tns:%s\">\n", name, portTypeName))
    // Reference the WS-Policy if one is defined
    if policy != PolicyNone {
        policyID := policyIDForType(policy)
        sb.WriteString(fmt.Sprintf("        <wsp:PolicyReference URI=\"#%s\"/>\n", policyID))
    }
    sb.WriteString(fmt.Sprintf("        <%s:binding style=\"document\" transport=\"http://schemas.xmlsoap.org/soap/http\"/>\n", soapPrefix))
    for _, op := range ops {
        sb.WriteString(fmt.Sprintf("        <operation name=\"%s\">\n", op))
        sb.WriteString(fmt.Sprintf("            <%s:operation soapAction=\"http://mock.imaxis.com/%s/%s\"/>\n", soapPrefix, endpoint, op))
        sb.WriteString("            <input>\n")
        sb.WriteString(fmt.Sprintf("                <%s:body use=\"literal\"/>\n", soapPrefix))
        if includeHeader && op != "Login" {
            sb.WriteString(fmt.Sprintf("                <%s:header message=\"tns:SessionHeader\" part=\"parameters\" use=\"literal\"/>\n", soapPrefix))
        }
        sb.WriteString("            </input>\n")
        sb.WriteString("            <output>\n")
        sb.WriteString(fmt.Sprintf("                <%s:body use=\"literal\"/>\n", soapPrefix))
        sb.WriteString("            </output>\n")
        sb.WriteString("        </operation>\n")
    }
    sb.WriteString("    </binding>\n")
    return sb.String()
}

func policyIDForType(policy PolicyType) string {
    switch policy {
    case PolicyWSSEUsernameToken:
        return "UsernameTokenPolicy"
    case PolicyBasicAuth:
        return "BasicAuthPolicy"
    case PolicyNTLMNegotiate:
        return "NTLMAuthPolicy"
    default:
        return ""
    }
}

func buildService(cfg wsdlConfig) string {
    var sb strings.Builder
    sb.WriteString(fmt.Sprintf("    <service name=\"%s\">\n", cfg.ServiceName))
    // SOAP 1.1 port
    sb.WriteString(fmt.Sprintf("        <port name=\"%sPort\" binding=\"tns:%s\">\n", cfg.PortTypeName, cfg.BindingName))
    sb.WriteString(fmt.Sprintf("            <soap:address location=\"%s\"/>\n", cfg.Address))
    sb.WriteString("        </port>\n")
    // SOAP 1.2 port — always included
    sb.WriteString(fmt.Sprintf("        <port name=\"%sPortSoap12\" binding=\"tns:%s\">\n", cfg.PortTypeName, cfg.BindingName+"Soap12"))
    sb.WriteString(fmt.Sprintf("            <soap12:address location=\"%s\"/>\n", cfg.Address))
    sb.WriteString("        </port>\n")
    sb.WriteString("    </service>\n")
    return sb.String()
}
