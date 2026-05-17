package finance

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/entertrans/backend-bogor.git/internal/dto"
	"github.com/entertrans/backend-bogor.git/internal/model"
	"gorm.io/gorm"
)

type PaymentController interface {
	GetAllPayments(siswaNIS string, page, limit int) ([]dto.PaymentListItem, int64, error)
	GetPaymentByID(id uint64) (*dto.PaymentDetailResponse, error)
	CreatePayment(siswaNIS string, req *dto.PaymentRequest, createdBy uint64) (*dto.PaymentResponse, error)
	DeletePayment(id uint64) error
}

type paymentController struct {
	db *gorm.DB
}

func NewPaymentController(db *gorm.DB) PaymentController {
	return &paymentController{db: db}
}

func (c *paymentController) generatePaymentNo() (string, error) {
	year := time.Now().Year()
	var lastNumber int = 0

	// Cari payment terakhir berdasarkan payment_no untuk tahun yang sama
	var lastPayment model.FinancePayment
	err := c.db.Where("payment_no LIKE ?", fmt.Sprintf("PAY/%d/%%", year)).
		Order("payment_no DESC").
		First(&lastPayment).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", err
	}

	// Extract number from last payment_no
	if lastPayment.PaymentNo != "" {
		// Format: PAY/2026/002027
		parts := strings.Split(lastPayment.PaymentNo, "/")
		if len(parts) == 3 {
			num, err := strconv.Atoi(parts[2])
			if err == nil {
				lastNumber = num
			}
		}
	}

	lastNumber++

	// Format dengan 6 digit (000001, 000002, dst)
	return fmt.Sprintf("PAY/%d/%06d", year, lastNumber), nil
}

func (c *paymentController) GetAllPayments(siswaNIS string, page, limit int) ([]dto.PaymentListItem, int64, error) {
	var payments []model.FinancePayment
	var total int64

	query := c.db.Model(&model.FinancePayment{})

	if siswaNIS != "" {
		query = query.Where("siswa_nis = ?", siswaNIS)
	}

	query.Count(&total)

	offset := (page - 1) * limit
	if limit > 0 {
		query = query.Offset(offset).Limit(limit)
	}

	err := query.Order("payment_id DESC").Find(&payments).Error
	if err != nil {
		return nil, 0, err
	}

	result := make([]dto.PaymentListItem, len(payments))
	for i, p := range payments {
		siswaNama := ""
		var siswa model.Siswa
		if err := c.db.Where("siswa_nis = ?", p.SiswaNIS).First(&siswa).Error; err == nil && siswa.SiswaNama != nil {
			siswaNama = *siswa.SiswaNama
		}

		result[i] = dto.PaymentListItem{
			PaymentID:     p.PaymentID,
			PaymentNo:     p.PaymentNo,
			SiswaNIS:      p.SiswaNIS,
			SiswaNama:     siswaNama,
			PaymentDate:   p.PaymentDate,
			PaymentMethod: p.PaymentMethod,
			TotalAmount:   p.TotalAmount,
			Notes:         p.Notes,
			CreatedAt:     p.CreatedAt,
		}
	}

	return result, total, nil
}

func (c *paymentController) GetPaymentByID(id uint64) (*dto.PaymentDetailResponse, error) {
	var payment model.FinancePayment
	err := c.db.Preload("Allocations").First(&payment, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("payment not found")
		}
		return nil, err
	}

	siswaNama := ""
	var siswa model.Siswa
	if err := c.db.Where("siswa_nis = ?", payment.SiswaNIS).First(&siswa).Error; err == nil && siswa.SiswaNama != nil {
		siswaNama = *siswa.SiswaNama
	}

	allocations := make([]dto.PaymentAllocationResponse, len(payment.Allocations))
	invoiceDetails := make([]struct {
		InvoiceID       uint64  `json:"invoice_id"`
		InvoiceNo       string  `json:"invoice_no"`
		AllocatedAmount float64 `json:"allocated_amount"`
		InvoiceTotal    float64 `json:"invoice_total"`
	}, len(payment.Allocations))

	for i, alloc := range payment.Allocations {
		var invoice model.FinanceInvoice
		c.db.First(&invoice, alloc.InvoiceID)

		// Calculate invoice total
		var items []model.FinanceInvoiceItem
		c.db.Where("invoice_id = ?", alloc.InvoiceID).Find(&items)
		var subtotal float64
		for _, item := range items {
			subtotal += item.Subtotal
		}

		var adjustments []model.FinanceInvoiceAdjustment
		c.db.Where("invoice_id = ?", alloc.InvoiceID).Find(&adjustments)
		var adjTotal float64
		for _, adj := range adjustments {
			adjTotal += adj.Amount
		}
		invoiceTotal := subtotal + adjTotal

		allocations[i] = dto.PaymentAllocationResponse{
			AllocationID: alloc.AllocationID,
			PaymentID:    alloc.PaymentID,
			InvoiceID:    alloc.InvoiceID,
			InvoiceNo:    invoice.InvoiceNo,
			Amount:       alloc.Amount,
			CreatedAt:    alloc.CreatedAt,
		}

		invoiceDetails[i] = struct {
			InvoiceID       uint64  `json:"invoice_id"`
			InvoiceNo       string  `json:"invoice_no"`
			AllocatedAmount float64 `json:"allocated_amount"`
			InvoiceTotal    float64 `json:"invoice_total"`
		}{
			InvoiceID:       alloc.InvoiceID,
			InvoiceNo:       invoice.InvoiceNo,
			AllocatedAmount: alloc.Amount,
			InvoiceTotal:    invoiceTotal,
		}
	}

	return &dto.PaymentDetailResponse{
		PaymentResponse: dto.PaymentResponse{
			PaymentID:     payment.PaymentID,
			PaymentNo:     payment.PaymentNo,
			SiswaNIS:      payment.SiswaNIS,
			SiswaNama:     siswaNama,
			PaymentDate:   payment.PaymentDate,
			PaymentMethod: payment.PaymentMethod,
			TotalAmount:   payment.TotalAmount,
			Notes:         payment.Notes,
			CreatedBy:     payment.CreatedBy,
			CreatedAt:     payment.CreatedAt,
			Allocations:   allocations,
		},
		Invoices: invoiceDetails,
	}, nil
}

func (c *paymentController) CreatePayment(siswaNIS string, req *dto.PaymentRequest, createdBy uint64) (*dto.PaymentResponse, error) {
	paymentNo, err := c.generatePaymentNo()
	if err != nil {
		return nil, err
	}

	paymentDate, err := time.Parse("2006-01-02", req.PaymentDate)
	if err != nil {
		return nil, errors.New("invalid payment_date format, use YYYY-MM-DD")
	}

	paymentMethod := req.PaymentMethod
	if paymentMethod == "" {
		paymentMethod = "cash"
	}

	payment := model.FinancePayment{
		PaymentNo:     paymentNo,
		SiswaNIS:      siswaNIS,
		PaymentDate:   paymentDate,
		PaymentMethod: paymentMethod,
		TotalAmount:   req.TotalAmount,
		Notes:         req.Notes,
		CreatedBy:     &createdBy,
	}

	err = c.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&payment).Error; err != nil {
			return err
		}

		var totalAllocated float64
		for _, allocReq := range req.Allocations {
			var invoice model.FinanceInvoice
			if err := tx.First(&invoice, allocReq.InvoiceID).Error; err != nil {
				return errors.New("invoice not found")
			}

			var itemCount int64
			tx.Model(&model.FinanceInvoiceItem{}).Where("invoice_id = ? AND siswa_nis = ?", allocReq.InvoiceID, siswaNIS).Count(&itemCount)
			if itemCount == 0 {
				return errors.New("invoice does not belong to this student")
			}

			allocation := model.FinancePaymentAllocation{
				PaymentID: payment.PaymentID,
				InvoiceID: allocReq.InvoiceID,
				SiswaNIS:  siswaNIS,
				Amount:    allocReq.Amount,
			}

			if err := tx.Create(&allocation).Error; err != nil {
				return err
			}

			totalAllocated += allocReq.Amount
		}

		if totalAllocated != req.TotalAmount {
			return errors.New("total allocated amount does not match payment amount")
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	paymentResp, err := c.GetPaymentByID(payment.PaymentID)
	if err != nil {
		return nil, err
	}

	return &paymentResp.PaymentResponse, nil
}

func (c *paymentController) DeletePayment(id uint64) error {
	var payment model.FinancePayment
	err := c.db.First(&payment, id).Error
	if err != nil {
		return err
	}

	return c.db.Delete(&payment).Error
}
