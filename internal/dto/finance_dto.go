package dto

import "time"

// ==================== FEE TEMPLATE DTOs ====================

type FeeTemplateRequest struct {
	FeeCode       *string `json:"fee_code" binding:"omitempty"`
	FeeName       string  `json:"fee_name" binding:"required"`
	DefaultAmount float64 `json:"default_amount" binding:"min=0"`
	Description   *string `json:"description"`
	IsActive      *int8   `json:"is_active"`
}

type FeeTemplateResponse struct {
	FeeTemplateID uint64     `json:"fee_template_id"`
	FeeCode       *string    `json:"fee_code"`
	FeeName       string     `json:"fee_name"`
	DefaultAmount float64    `json:"default_amount"`
	Description   *string    `json:"description"`
	IsActive      int8       `json:"is_active"`
	CreatedAt     *time.Time `json:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at"`
}

// ==================== INVOICE ITEM DTOs ====================

type InvoiceItemRequest struct {
	SiswaNIS      string  `json:"siswa_nis" binding:"required"`
	FeeTemplateID *uint64 `json:"fee_template_id"`
	ItemName      string  `json:"item_name" binding:"required"`
	Description   *string `json:"description"`
	Qty           float64 `json:"qty" binding:"required,min=0.01"`
	UnitPrice     float64 `json:"unit_price" binding:"required,min=0"`
	SortOrder     int     `json:"sort_order"`
}

type InvoiceItemResponse struct {
	InvoiceItemID uint64               `json:"invoice_item_id"`
	SiswaNIS      string               `json:"siswa_nis"`
	ItemName      string               `json:"item_name"`
	Description   *string              `json:"description"`
	Qty           float64              `json:"qty"`
	UnitPrice     float64              `json:"unit_price"`
	Subtotal      float64              `json:"subtotal"`
	SortOrder     int                  `json:"sort_order"`
	FeeTemplate   *FeeTemplateResponse `json:"fee_template,omitempty"`
}

// ==================== INVOICE ADJUSTMENT DTOs ====================

type InvoiceAdjustmentRequest struct {
	SiswaNIS       string  `json:"siswa_nis" binding:"required"`
	AdjustmentType string  `json:"adjustment_type"`
	AdjustmentName string  `json:"adjustment_name" binding:"required"`
	Amount         float64 `json:"amount" binding:"required"`
	Notes          *string `json:"notes"`
}

type InvoiceAdjustmentResponse struct {
	AdjustmentID   uint64     `json:"adjustment_id"`
	SiswaNIS       string     `json:"siswa_nis"`
	AdjustmentType string     `json:"adjustment_type"`
	AdjustmentName string     `json:"adjustment_name"`
	Amount         float64    `json:"amount"`
	Notes          *string    `json:"notes"`
	CreatedAt      *time.Time `json:"created_at"`
}

// ==================== INVOICE REQUEST (untuk Create) ====================

// Di dto/finance.go
type InvoiceData struct {
	IssueDate   string `json:"issue_date" binding:"required"`
	DueDate     string `json:"due_date"` // Hapus binding:"required" kalau ada
	Description string `json:"description"`
}
type CreateInvoiceRequest struct {
	Invoice     InvoiceData                `json:"invoice" binding:"required"`
	Items       []InvoiceItemRequest       `json:"items" binding:"required,min=1"`
	Adjustments []InvoiceAdjustmentRequest `json:"adjustments"`
}

// ==================== INVOICE REQUEST (untuk Update) ====================

type InvoiceRequest struct {
	IssueDate   string  `json:"issue_date" binding:"required"`
	DueDate     *string `json:"due_date"`
	Description *string `json:"description"`
}

// ==================== STUDENT INVOICE DETAIL ====================

type StudentInvoiceDetail struct {
	SiswaNIS        string                      `json:"siswa_nis"`
	SiswaNama       string                      `json:"siswa_nama"`
	Kelas           string                      `json:"kelas"`
	Subtotal        float64                     `json:"subtotal"`
	AdjustmentTotal float64                     `json:"adjustment_total"`
	GrandTotal      float64                     `json:"grand_total"`
	PaidTotal       float64                     `json:"paid_total"`
	RemainingAmount float64                     `json:"remaining_amount"`
	Status          string                      `json:"status"`
	Items           []InvoiceItemResponse       `json:"items"`
	Adjustments     []InvoiceAdjustmentResponse `json:"adjustments"`
	Payments        []PaymentAllocationResponse `json:"payments"`
}

// ==================== INVOICE RESPONSE ====================

type InvoiceResponse struct {
	InvoiceID       uint64                 `json:"invoice_id"`
	InvoiceNo       string                 `json:"invoice_no"`
	IssueDate       time.Time              `json:"issue_date"`
	DueDate         *time.Time             `json:"due_date"`
	Description     *string                `json:"description"`
	Status          string                 `json:"status"`
	CreatedBy       *uint64                `json:"created_by"`
	CreatedAt       *time.Time             `json:"created_at"`
	UpdatedAt       *time.Time             `json:"updated_at"`
	TotalSubtotal   float64                `json:"total_subtotal"`
	TotalAdjustment float64                `json:"total_adjustment"`
	TotalGrandTotal float64                `json:"total_grand_total"`
	TotalPaid       float64                `json:"total_paid"`
	TotalRemaining  float64                `json:"total_remaining"`
	Students        []StudentInvoiceDetail `json:"students"`
}

type InvoiceListItem struct {
	InvoiceID       uint64     `json:"invoice_id"`
	InvoiceNo       string     `json:"invoice_no"`
	IssueDate       time.Time  `json:"issue_date"`
	DueDate         *time.Time `json:"due_date"`
	Description     *string    `json:"description"`
	Status          string     `json:"status"`
	TotalGrandTotal float64    `json:"total_grand_total"`
	TotalPaid       float64    `json:"total_paid"`
	TotalRemaining  float64    `json:"total_remaining"`
	StudentCount    int        `json:"student_count"`
	CreatedAt       *time.Time `json:"created_at"`
}

// ==================== PAYMENT ALLOCATION DTOs ====================

type PaymentAllocationRequest struct {
	InvoiceID uint64  `json:"invoice_id" binding:"required"`
	Amount    float64 `json:"amount" binding:"required,min=0.01"`
}

type PaymentAllocationResponse struct {
	AllocationID uint64     `json:"allocation_id"`
	PaymentID    uint64     `json:"payment_id"`
	InvoiceID    uint64     `json:"invoice_id"`
	InvoiceNo    string     `json:"invoice_no"`
	Amount       float64    `json:"amount"`
	CreatedAt    *time.Time `json:"created_at"`
}

// ==================== PAYMENT DTOs ====================

type PaymentRequest struct {
	PaymentDate   string                     `json:"payment_date" binding:"required"`
	PaymentMethod string                     `json:"payment_method"`
	TotalAmount   float64                    `json:"total_amount" binding:"required,min=0.01"`
	Notes         *string                    `json:"notes"`
	Allocations   []PaymentAllocationRequest `json:"allocations" binding:"required,min=1"`
}

type PaymentResponse struct {
	PaymentID     uint64                      `json:"payment_id"`
	PaymentNo     string                      `json:"payment_no"`
	SiswaNIS      string                      `json:"siswa_nis"`
	SiswaNama     string                      `json:"siswa_nama"`
	PaymentDate   time.Time                   `json:"payment_date"`
	PaymentMethod string                      `json:"payment_method"`
	TotalAmount   float64                     `json:"total_amount"`
	Notes         *string                     `json:"notes"`
	CreatedBy     *uint64                     `json:"created_by"`
	CreatedAt     *time.Time                  `json:"created_at"`
	Allocations   []PaymentAllocationResponse `json:"allocations"`
}

type PaymentListItem struct {
	PaymentID     uint64     `json:"payment_id"`
	PaymentNo     string     `json:"payment_no"`
	SiswaNIS      string     `json:"siswa_nis"`
	SiswaNama     string     `json:"siswa_nama"`
	PaymentDate   time.Time  `json:"payment_date"`
	PaymentMethod string     `json:"payment_method"`
	TotalAmount   float64    `json:"total_amount"`
	Notes         *string    `json:"notes"`
	CreatedAt     *time.Time `json:"created_at"`
}

type PaymentDetailResponse struct {
	PaymentResponse
	Invoices []struct {
		InvoiceID       uint64  `json:"invoice_id"`
		InvoiceNo       string  `json:"invoice_no"`
		AllocatedAmount float64 `json:"allocated_amount"`
		InvoiceTotal    float64 `json:"invoice_total"`
	} `json:"invoices"`
}
