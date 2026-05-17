package model

import "time"

// FinanceFeeTemplate - model untuk tbl_finance_fee_templates
type FinanceFeeTemplate struct {
	FeeTemplateID uint64     `gorm:"primaryKey;autoIncrement;column:fee_template_id" json:"fee_template_id"`
	FeeCode       *string    `gorm:"type:varchar(50);column:fee_code" json:"fee_code"`
	FeeName       string     `gorm:"type:varchar(255);not null;column:fee_name" json:"fee_name"`
	DefaultAmount float64    `gorm:"type:decimal(18,2);default:0;column:default_amount" json:"default_amount"`
	Description   *string    `gorm:"type:text;column:description" json:"description"`
	IsActive      int8       `gorm:"default:1;column:is_active" json:"is_active"`
	CreatedAt     *time.Time `gorm:"autoCreateTime;column:created_at" json:"created_at"`
	UpdatedAt     *time.Time `gorm:"autoUpdateTime;column:updated_at" json:"updated_at"`
}

func (FinanceFeeTemplate) TableName() string {
	return "tbl_finance_fee_templates"
}

// FinanceInvoice - model untuk tbl_finance_invoices
type FinanceInvoice struct {
	InvoiceID   uint64     `gorm:"primaryKey;autoIncrement;column:invoice_id" json:"invoice_id"`
	InvoiceNo   string     `gorm:"type:varchar(100);not null;uniqueIndex;column:invoice_no" json:"invoice_no"`
	IssueDate   time.Time  `gorm:"type:date;not null;column:issue_date" json:"issue_date"`
	DueDate     *time.Time `gorm:"type:date;column:due_date" json:"due_date"`
	Description *string    `gorm:"type:text;column:description" json:"description"`
	Status      string     `gorm:"type:enum('draft','published','cancelled');default:'draft';column:status" json:"status"`
	CreatedBy   *uint64    `gorm:"column:created_by" json:"created_by"`
	CreatedAt   *time.Time `gorm:"autoCreateTime;column:created_at" json:"created_at"`
	UpdatedAt   *time.Time `gorm:"autoUpdateTime;column:updated_at" json:"updated_at"`

	// Relasi
	Items       []FinanceInvoiceItem       `gorm:"foreignKey:InvoiceID" json:"items,omitempty"`
	Adjustments []FinanceInvoiceAdjustment `gorm:"foreignKey:InvoiceID" json:"adjustments,omitempty"`
	Allocations []FinancePaymentAllocation `gorm:"foreignKey:InvoiceID" json:"allocations,omitempty"`
}

func (FinanceInvoice) TableName() string {
	return "tbl_finance_invoices"
}

// FinanceInvoiceItem - model untuk tbl_finance_invoice_items
type FinanceInvoiceItem struct {
	InvoiceItemID uint64     `gorm:"primaryKey;autoIncrement;column:invoice_item_id" json:"invoice_item_id"`
	InvoiceID     uint64     `gorm:"not null;column:invoice_id;index:idx_invoice_items_invoice_id" json:"invoice_id"`
	SiswaNIS      string     `gorm:"type:varchar(30);not null;column:siswa_nis;index:idx_invoice_items_siswa_nis" json:"siswa_nis"`
	FeeTemplateID *uint64    `gorm:"column:fee_template_id" json:"fee_template_id"`
	ItemName      string     `gorm:"type:varchar(255);not null;column:item_name" json:"item_name"`
	Description   *string    `gorm:"type:text;column:description" json:"description"`
	Qty           float64    `gorm:"type:decimal(10,2);default:1;column:qty" json:"qty"`
	UnitPrice     float64    `gorm:"type:decimal(18,2);not null;column:unit_price" json:"unit_price"`
	Subtotal      float64    `gorm:"type:decimal(18,2);not null;column:subtotal" json:"subtotal"`
	SortOrder     int        `gorm:"default:0;column:sort_order" json:"sort_order"`
	CreatedAt     *time.Time `gorm:"autoCreateTime;column:created_at" json:"created_at"`

	// Relasi
	Invoice     FinanceInvoice      `gorm:"foreignKey:InvoiceID" json:"-"`
	FeeTemplate *FinanceFeeTemplate `gorm:"foreignKey:FeeTemplateID" json:"fee_template,omitempty"`
}

func (FinanceInvoiceItem) TableName() string {
	return "tbl_finance_invoice_items"
}

// FinanceInvoiceAdjustment - model untuk tbl_finance_invoice_adjustments
type FinanceInvoiceAdjustment struct {
	AdjustmentID   uint64     `gorm:"primaryKey;autoIncrement;column:adjustment_id" json:"adjustment_id"`
	InvoiceID      uint64     `gorm:"not null;column:invoice_id;index:idx_adjustment_invoice_id" json:"invoice_id"`
	SiswaNIS       string     `gorm:"type:varchar(30);not null;column:siswa_nis;index:idx_adjustment_siswa_nis" json:"siswa_nis"`
	AdjustmentType string     `gorm:"type:enum('discount','subsidy','penalty','other');default:'other';column:adjustment_type" json:"adjustment_type"`
	AdjustmentName string     `gorm:"type:varchar(255);not null;column:adjustment_name" json:"adjustment_name"`
	Amount         float64    `gorm:"type:decimal(18,2);not null;column:amount" json:"amount"`
	Notes          *string    `gorm:"type:text;column:notes" json:"notes"`
	CreatedAt      *time.Time `gorm:"autoCreateTime;column:created_at" json:"created_at"`

	Invoice FinanceInvoice `gorm:"foreignKey:InvoiceID" json:"-"`
}

func (FinanceInvoiceAdjustment) TableName() string {
	return "tbl_finance_invoice_adjustments"
}

// FinancePayment - model untuk tbl_finance_payments
type FinancePayment struct {
	PaymentID     uint64     `gorm:"primaryKey;autoIncrement;column:payment_id" json:"payment_id"`
	PaymentNo     string     `gorm:"type:varchar(100);not null;uniqueIndex;column:payment_no" json:"payment_no"`
	SiswaNIS      string     `gorm:"type:varchar(30);not null;column:siswa_nis;index:idx_payment_siswa_nis" json:"siswa_nis"`
	PaymentDate   time.Time  `gorm:"type:date;not null;column:payment_date" json:"payment_date"`
	PaymentMethod string     `gorm:"type:enum('cash','transfer','qris','other');default:'cash';column:payment_method" json:"payment_method"`
	TotalAmount   float64    `gorm:"type:decimal(18,2);not null;column:total_amount" json:"total_amount"`
	Notes         *string    `gorm:"type:text;column:notes" json:"notes"`
	CreatedBy     *uint64    `gorm:"column:created_by" json:"created_by"`
	CreatedAt     *time.Time `gorm:"autoCreateTime;column:created_at" json:"created_at"`

	// Relasi
	Siswa       Siswa                      `gorm:"foreignKey:SiswaNIS;references:SiswaNIS" json:"siswa,omitempty"`
	Allocations []FinancePaymentAllocation `gorm:"foreignKey:PaymentID" json:"allocations,omitempty"`
}

func (FinancePayment) TableName() string {
	return "tbl_finance_payments"
}

// FinancePaymentAllocation - model untuk tbl_finance_payment_allocations
type FinancePaymentAllocation struct {
	AllocationID uint64     `gorm:"primaryKey;autoIncrement;column:allocation_id" json:"allocation_id"`
	PaymentID    uint64     `gorm:"not null;column:payment_id;index:idx_alloc_payment_id" json:"payment_id"`
	InvoiceID    uint64     `gorm:"not null;column:invoice_id;index:idx_alloc_invoice_id" json:"invoice_id"`
	SiswaNIS     string     `gorm:"type:varchar(30);not null;column:siswa_nis;index:idx_alloc_siswa_nis" json:"siswa_nis"`
	Amount       float64    `gorm:"type:decimal(18,2);not null;column:amount" json:"amount"`
	CreatedAt    *time.Time `gorm:"autoCreateTime;column:created_at" json:"created_at"`

	Payment FinancePayment `gorm:"foreignKey:PaymentID" json:"-"`
	Invoice FinanceInvoice `gorm:"foreignKey:InvoiceID" json:"invoice,omitempty"`
}

func (FinancePaymentAllocation) TableName() string {
	return "tbl_finance_payment_allocations"
}
