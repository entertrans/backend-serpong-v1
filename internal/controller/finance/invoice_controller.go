package finance

import (
	"errors"
	"fmt"
	"time"

	"github.com/entertrans/backend-bogor.git/internal/dto"
	"github.com/entertrans/backend-bogor.git/internal/model"
	"gorm.io/gorm"
)

type InvoiceController interface {
	GetAllInvoices(status string, page, limit int) ([]dto.InvoiceListItem, int64, error)
	GetInvoiceByID(id uint64) (*dto.InvoiceResponse, error)
	CreateInvoice(req *dto.InvoiceRequest, items []dto.InvoiceItemRequest, adjustments []dto.InvoiceAdjustmentRequest, createdBy uint64) (*dto.InvoiceResponse, error)
	UpdateInvoice(id uint64, req *dto.InvoiceRequest) (*dto.InvoiceResponse, error)
	DeleteInvoice(id uint64) error
	PublishInvoice(id uint64) error
	CancelInvoice(id uint64) error
	AddInvoiceItem(invoiceID uint64, req *dto.InvoiceItemRequest) (*dto.InvoiceItemResponse, error)
	UpdateInvoiceItem(itemID uint64, req *dto.InvoiceItemRequest) (*dto.InvoiceItemResponse, error)
	RemoveInvoiceItem(itemID uint64) error
	AddInvoiceAdjustment(invoiceID uint64, req *dto.InvoiceAdjustmentRequest) (*dto.InvoiceAdjustmentResponse, error)
	UpdateInvoiceAdjustment(adjustmentID uint64, req *dto.InvoiceAdjustmentRequest) (*dto.InvoiceAdjustmentResponse, error)
	RemoveInvoiceAdjustment(adjustmentID uint64) error
}

type invoiceController struct {
	db *gorm.DB
}

func NewInvoiceController(db *gorm.DB) InvoiceController {
	return &invoiceController{db: db}
}

func (c *invoiceController) generateInvoiceNo() (string, error) {
	year := time.Now().Year()
	var maxNumber int = 0

	// Query untuk mendapatkan nomor terbesar
	type Result struct {
		MaxNumber int
	}
	var result Result

	// Gunakan query raw untuk mencari max number
	err := c.db.Raw(`
		SELECT COALESCE(MAX(CAST(SUBSTRING_INDEX(invoice_no, '/', -1) AS UNSIGNED)), 0) as max_number
		FROM tbl_finance_invoices 
		WHERE invoice_no LIKE ?
	`, fmt.Sprintf("INV/%d/%%", year)).Scan(&result).Error

	if err != nil {
		return "", fmt.Errorf("failed to get max invoice number: %w", err)
	}

	maxNumber = result.MaxNumber
	maxNumber++

	newInvoiceNo := fmt.Sprintf("INV/%d/%06d", year, maxNumber)

	// Double check: pastikan nomor benar-benar tidak ada
	var count int64
	c.db.Model(&model.FinanceInvoice{}).Where("invoice_no = ?", newInvoiceNo).Count(&count)

	if count > 0 {
		// Jika masih ada, increment lagi
		maxNumber++
		newInvoiceNo = fmt.Sprintf("INV/%d/%06d", year, maxNumber)
	}

	return newInvoiceNo, nil
}

func (c *invoiceController) calculateInvoiceTotals(invoiceID uint64) (subtotal, adjustment, grandTotal float64, err error) {
	// Calculate items subtotal
	var items []model.FinanceInvoiceItem
	err = c.db.Where("invoice_id = ?", invoiceID).Find(&items).Error
	if err != nil {
		return 0, 0, 0, err
	}

	for _, item := range items {
		subtotal += item.Subtotal
	}

	// Calculate adjustments
	var adjustments []model.FinanceInvoiceAdjustment
	err = c.db.Where("invoice_id = ?", invoiceID).Find(&adjustments).Error
	if err != nil {
		return 0, 0, 0, err
	}

	for _, adj := range adjustments {
		adjustment += adj.Amount
	}

	grandTotal = subtotal + adjustment
	return subtotal, adjustment, grandTotal, nil
}

func (c *invoiceController) calculatePaidTotal(invoiceID uint64) (float64, error) {
	var allocations []model.FinancePaymentAllocation
	err := c.db.Where("invoice_id = ?", invoiceID).Find(&allocations).Error
	if err != nil {
		return 0, err
	}

	var total float64
	for _, alloc := range allocations {
		total += alloc.Amount
	}
	return total, nil
}

func (c *invoiceController) getStudentInfo(siswaNIS string) (nama string, kelas string, err error) {
	var siswa model.Siswa
	err = c.db.Where("siswa_nis = ?", siswaNIS).First(&siswa).Error
	if err != nil {
		return "", "", err
	}

	if siswa.SiswaNama != nil {
		nama = *siswa.SiswaNama
	}

	if siswa.SiswaKelasID != nil {
		var kelasModel model.Kelas
		err = c.db.First(&kelasModel, *siswa.SiswaKelasID).Error
		if err == nil {
			kelas = kelasModel.KelasNama
		}
	}

	return nama, kelas, nil
}

func (c *invoiceController) GetAllInvoices(status string, page, limit int) ([]dto.InvoiceListItem, int64, error) {
	var invoices []model.FinanceInvoice
	var total int64

	query := c.db.Model(&model.FinanceInvoice{})

	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)

	offset := (page - 1) * limit
	if limit > 0 {
		query = query.Offset(offset).Limit(limit)
	}

	err := query.Order("invoice_id DESC").Find(&invoices).Error
	if err != nil {
		return nil, 0, err
	}

	result := make([]dto.InvoiceListItem, len(invoices))
	for i, inv := range invoices {
		_, _, grandTotal, _ := c.calculateInvoiceTotals(inv.InvoiceID)
		paidTotal, _ := c.calculatePaidTotal(inv.InvoiceID)

		var studentCount int64
		c.db.Model(&model.FinanceInvoiceItem{}).Where("invoice_id = ?", inv.InvoiceID).Distinct("siswa_nis").Count(&studentCount)

		result[i] = dto.InvoiceListItem{
			InvoiceID:       inv.InvoiceID,
			InvoiceNo:       inv.InvoiceNo,
			IssueDate:       inv.IssueDate,
			DueDate:         inv.DueDate,
			Description:     inv.Description,
			Status:          inv.Status,
			TotalGrandTotal: grandTotal,
			TotalPaid:       paidTotal,
			TotalRemaining:  grandTotal - paidTotal,
			StudentCount:    int(studentCount),
			CreatedAt:       inv.CreatedAt,
		}
	}

	return result, total, nil
}

func (c *invoiceController) GetInvoiceByID(id uint64) (*dto.InvoiceResponse, error) {
	var invoice model.FinanceInvoice
	err := c.db.Preload("Items").Preload("Adjustments").Preload("Allocations").First(&invoice, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invoice not found")
		}
		return nil, err
	}

	// Group items and adjustments by student
	studentMap := make(map[string]*dto.StudentInvoiceDetail)

	// Process items
	for _, item := range invoice.Items {
		if _, exists := studentMap[item.SiswaNIS]; !exists {
			siswaNama, kelas, _ := c.getStudentInfo(item.SiswaNIS)
			studentMap[item.SiswaNIS] = &dto.StudentInvoiceDetail{
				SiswaNIS:    item.SiswaNIS,
				SiswaNama:   siswaNama,
				Kelas:       kelas,
				Items:       []dto.InvoiceItemResponse{},
				Adjustments: []dto.InvoiceAdjustmentResponse{},
				Payments:    []dto.PaymentAllocationResponse{},
			}
		}

		itemResp := dto.InvoiceItemResponse{
			InvoiceItemID: item.InvoiceItemID,
			SiswaNIS:      item.SiswaNIS,
			ItemName:      item.ItemName,
			Description:   item.Description,
			Qty:           item.Qty,
			UnitPrice:     item.UnitPrice,
			Subtotal:      item.Subtotal,
			SortOrder:     item.SortOrder,
		}

		if item.FeeTemplateID != nil {
			var template model.FinanceFeeTemplate
			if err := c.db.First(&template, *item.FeeTemplateID).Error; err == nil {
				itemResp.FeeTemplate = &dto.FeeTemplateResponse{
					FeeTemplateID: template.FeeTemplateID,
					FeeCode:       template.FeeCode,
					FeeName:       template.FeeName,
					DefaultAmount: template.DefaultAmount,
					IsActive:      template.IsActive,
				}
			}
		}

		studentMap[item.SiswaNIS].Items = append(studentMap[item.SiswaNIS].Items, itemResp)
		studentMap[item.SiswaNIS].Subtotal += item.Subtotal
	}

	// Process adjustments
	for _, adj := range invoice.Adjustments {
		if student, exists := studentMap[adj.SiswaNIS]; exists {
			student.Adjustments = append(student.Adjustments, dto.InvoiceAdjustmentResponse{
				AdjustmentID:   adj.AdjustmentID,
				SiswaNIS:       adj.SiswaNIS,
				AdjustmentType: adj.AdjustmentType,
				AdjustmentName: adj.AdjustmentName,
				Amount:         adj.Amount,
				Notes:          adj.Notes,
				CreatedAt:      adj.CreatedAt,
			})
			student.AdjustmentTotal += adj.Amount
		}
	}

	// Process allocations (payments)
	for _, alloc := range invoice.Allocations {
		if student, exists := studentMap[alloc.SiswaNIS]; exists {
			student.Payments = append(student.Payments, dto.PaymentAllocationResponse{
				AllocationID: alloc.AllocationID,
				PaymentID:    alloc.PaymentID,
				InvoiceID:    alloc.InvoiceID,
				Amount:       alloc.Amount,
				CreatedAt:    alloc.CreatedAt,
			})
			student.PaidTotal += alloc.Amount
		}
	}

	// Calculate grand total and status for each student
	students := make([]dto.StudentInvoiceDetail, 0, len(studentMap))
	var totalSubtotal, totalAdjustment, totalGrandTotal, totalPaid float64

	for _, student := range studentMap {
		student.GrandTotal = student.Subtotal + student.AdjustmentTotal
		student.RemainingAmount = student.GrandTotal - student.PaidTotal

		if student.PaidTotal == 0 {
			student.Status = "unpaid"
		} else if student.PaidTotal >= student.GrandTotal {
			student.Status = "paid"
		} else {
			student.Status = "partial"
		}

		totalSubtotal += student.Subtotal
		totalAdjustment += student.AdjustmentTotal
		totalGrandTotal += student.GrandTotal
		totalPaid += student.PaidTotal

		students = append(students, *student)
	}

	return &dto.InvoiceResponse{
		InvoiceID:       invoice.InvoiceID,
		InvoiceNo:       invoice.InvoiceNo,
		IssueDate:       invoice.IssueDate,
		DueDate:         invoice.DueDate,
		Description:     invoice.Description,
		Status:          invoice.Status,
		CreatedBy:       invoice.CreatedBy,
		CreatedAt:       invoice.CreatedAt,
		UpdatedAt:       invoice.UpdatedAt,
		TotalSubtotal:   totalSubtotal,
		TotalAdjustment: totalAdjustment,
		TotalGrandTotal: totalGrandTotal,
		TotalPaid:       totalPaid,
		TotalRemaining:  totalGrandTotal - totalPaid,
		Students:        students,
	}, nil
}

func (c *invoiceController) CreateInvoice(req *dto.InvoiceRequest, items []dto.InvoiceItemRequest, adjustments []dto.InvoiceAdjustmentRequest, createdBy uint64) (*dto.InvoiceResponse, error) {
	// Validate required fields
	if req == nil {
		return nil, errors.New("invoice request is required")
	}

	if len(items) == 0 {
		return nil, errors.New("at least one item is required")
	}

	// Log untuk debug
	fmt.Printf("CreateInvoice called with: IssueDate=%s, DueDate=%v, Description=%v\n",
		req.IssueDate, req.DueDate, req.Description)
	fmt.Printf("Items count: %d\n", len(items))

	invoiceNo, err := c.generateInvoiceNo()
	if err != nil {
		return nil, err
	}

	// Parse issue_date
	issueDate, err := time.Parse("2006-01-02", req.IssueDate)
	if err != nil {
		return nil, fmt.Errorf("invalid issue_date format: %s, error: %w", req.IssueDate, err)
	}

	// Parse due_date (bisa null atau empty string)
	var dueDate *time.Time
	if req.DueDate != nil && *req.DueDate != "" {
		parsed, err := time.Parse("2006-01-02", *req.DueDate)
		if err != nil {
			return nil, fmt.Errorf("invalid due_date format: %s, error: %w", *req.DueDate, err)
		}
		dueDate = &parsed
	}

	invoice := model.FinanceInvoice{
		InvoiceNo:   invoiceNo,
		IssueDate:   issueDate,
		DueDate:     dueDate,
		Description: req.Description,
		Status:      "draft",
		CreatedBy:   &createdBy,
	}

	err = c.db.Transaction(func(tx *gorm.DB) error {
		// Create invoice
		if err := tx.Create(&invoice).Error; err != nil {
			return fmt.Errorf("failed to create invoice: %w", err)
		}

		// Create items
		for i, itemReq := range items {
			// Validate item
			if itemReq.SiswaNIS == "" {
				return fmt.Errorf("item %d: siswa_nis is required", i)
			}
			if itemReq.ItemName == "" {
				return fmt.Errorf("item %d: item_name is required", i)
			}
			if itemReq.Qty <= 0 {
				return fmt.Errorf("item %d: qty must be greater than 0", i)
			}
			if itemReq.UnitPrice < 0 {
				return fmt.Errorf("item %d: unit_price cannot be negative", i)
			}

			// Handle fee_template_id (bisa null)
			var feeTemplateID *uint64
			if itemReq.FeeTemplateID != nil && *itemReq.FeeTemplateID > 0 {
				feeTemplateID = itemReq.FeeTemplateID
			}

			item := model.FinanceInvoiceItem{
				InvoiceID:     invoice.InvoiceID,
				SiswaNIS:      itemReq.SiswaNIS,
				FeeTemplateID: feeTemplateID,
				ItemName:      itemReq.ItemName,
				Description:   itemReq.Description,
				Qty:           itemReq.Qty,
				UnitPrice:     itemReq.UnitPrice,
				Subtotal:      itemReq.Qty * itemReq.UnitPrice,
				SortOrder:     itemReq.SortOrder,
			}
			if err := tx.Create(&item).Error; err != nil {
				return fmt.Errorf("failed to create item %d: %w", i, err)
			}
		}

		// Create adjustments
		for i, adjReq := range adjustments {
			// Skip if amount is 0
			if adjReq.Amount == 0 {
				continue
			}

			if adjReq.SiswaNIS == "" {
				return fmt.Errorf("adjustment %d: siswa_nis is required", i)
			}
			if adjReq.AdjustmentName == "" {
				return fmt.Errorf("adjustment %d: adjustment_name is required", i)
			}

			adjType := adjReq.AdjustmentType
			if adjType == "" {
				adjType = "other"
			}

			adj := model.FinanceInvoiceAdjustment{
				InvoiceID:      invoice.InvoiceID,
				SiswaNIS:       adjReq.SiswaNIS,
				AdjustmentType: adjType,
				AdjustmentName: adjReq.AdjustmentName,
				Amount:         adjReq.Amount,
				Notes:          adjReq.Notes,
			}
			if err := tx.Create(&adj).Error; err != nil {
				return fmt.Errorf("failed to create adjustment %d: %w", i, err)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return c.GetInvoiceByID(invoice.InvoiceID)
}

func (c *invoiceController) UpdateInvoice(id uint64, req *dto.InvoiceRequest) (*dto.InvoiceResponse, error) {
	var invoice model.FinanceInvoice
	err := c.db.First(&invoice, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invoice not found")
		}
		return nil, err
	}

	if invoice.Status != "draft" {
		return nil, errors.New("only draft invoices can be updated")
	}

	issueDate, err := time.Parse("2006-01-02", req.IssueDate)
	if err != nil {
		return nil, errors.New("invalid issue_date format, use YYYY-MM-DD")
	}
	invoice.IssueDate = issueDate

	if req.DueDate != nil && *req.DueDate != "" {
		parsed, err := time.Parse("2006-01-02", *req.DueDate)
		if err != nil {
			return nil, errors.New("invalid due_date format, use YYYY-MM-DD")
		}
		invoice.DueDate = &parsed
	} else {
		invoice.DueDate = nil
	}

	invoice.Description = req.Description

	err = c.db.Save(&invoice).Error
	if err != nil {
		return nil, err
	}

	return c.GetInvoiceByID(id)
}

func (c *invoiceController) DeleteInvoice(id uint64) error {
	var invoice model.FinanceInvoice
	err := c.db.First(&invoice, id).Error
	if err != nil {
		return err
	}

	if invoice.Status != "draft" {
		return errors.New("only draft invoices can be deleted")
	}

	return c.db.Delete(&invoice).Error
}

func (c *invoiceController) PublishInvoice(id uint64) error {
	var invoice model.FinanceInvoice
	err := c.db.First(&invoice, id).Error
	if err != nil {
		return err
	}

	if invoice.Status != "draft" {
		return errors.New("only draft invoices can be published")
	}

	return c.db.Model(&invoice).Update("status", "published").Error
}

func (c *invoiceController) CancelInvoice(id uint64) error {
	var invoice model.FinanceInvoice
	err := c.db.First(&invoice, id).Error
	if err != nil {
		return err
	}

	if invoice.Status != "published" {
		return errors.New("only published invoices can be cancelled")
	}

	return c.db.Model(&invoice).Update("status", "cancelled").Error
}

func (c *invoiceController) AddInvoiceItem(invoiceID uint64, req *dto.InvoiceItemRequest) (*dto.InvoiceItemResponse, error) {
	var invoice model.FinanceInvoice
	err := c.db.First(&invoice, invoiceID).Error
	if err != nil {
		return nil, errors.New("invoice not found")
	}

	if invoice.Status != "draft" {
		return nil, errors.New("items can only be added to draft invoices")
	}

	item := model.FinanceInvoiceItem{
		InvoiceID:     invoiceID,
		SiswaNIS:      req.SiswaNIS,
		FeeTemplateID: req.FeeTemplateID,
		ItemName:      req.ItemName,
		Description:   req.Description,
		Qty:           req.Qty,
		UnitPrice:     req.UnitPrice,
		Subtotal:      req.Qty * req.UnitPrice,
		SortOrder:     req.SortOrder,
	}

	err = c.db.Create(&item).Error
	if err != nil {
		return nil, err
	}

	return &dto.InvoiceItemResponse{
		InvoiceItemID: item.InvoiceItemID,
		SiswaNIS:      item.SiswaNIS,
		ItemName:      item.ItemName,
		Description:   item.Description,
		Qty:           item.Qty,
		UnitPrice:     item.UnitPrice,
		Subtotal:      item.Subtotal,
		SortOrder:     item.SortOrder,
	}, nil
}

func (c *invoiceController) UpdateInvoiceItem(itemID uint64, req *dto.InvoiceItemRequest) (*dto.InvoiceItemResponse, error) {
	var item model.FinanceInvoiceItem
	err := c.db.First(&item, itemID).Error
	if err != nil {
		return nil, errors.New("item not found")
	}

	var invoice model.FinanceInvoice
	err = c.db.First(&invoice, item.InvoiceID).Error
	if err != nil {
		return nil, err
	}

	if invoice.Status != "draft" {
		return nil, errors.New("items can only be updated in draft invoices")
	}

	item.ItemName = req.ItemName
	item.Description = req.Description
	item.Qty = req.Qty
	item.UnitPrice = req.UnitPrice
	item.Subtotal = req.Qty * req.UnitPrice
	item.SortOrder = req.SortOrder

	if req.FeeTemplateID != nil {
		item.FeeTemplateID = req.FeeTemplateID
	}

	err = c.db.Save(&item).Error
	if err != nil {
		return nil, err
	}

	return &dto.InvoiceItemResponse{
		InvoiceItemID: item.InvoiceItemID,
		SiswaNIS:      item.SiswaNIS,
		ItemName:      item.ItemName,
		Description:   item.Description,
		Qty:           item.Qty,
		UnitPrice:     item.UnitPrice,
		Subtotal:      item.Subtotal,
		SortOrder:     item.SortOrder,
	}, nil
}

func (c *invoiceController) RemoveInvoiceItem(itemID uint64) error {
	var item model.FinanceInvoiceItem
	err := c.db.First(&item, itemID).Error
	if err != nil {
		return err
	}

	var invoice model.FinanceInvoice
	err = c.db.First(&invoice, item.InvoiceID).Error
	if err != nil {
		return err
	}

	if invoice.Status != "draft" {
		return errors.New("items can only be removed from draft invoices")
	}

	return c.db.Delete(&item).Error
}

func (c *invoiceController) AddInvoiceAdjustment(invoiceID uint64, req *dto.InvoiceAdjustmentRequest) (*dto.InvoiceAdjustmentResponse, error) {
	var invoice model.FinanceInvoice
	err := c.db.First(&invoice, invoiceID).Error
	if err != nil {
		return nil, errors.New("invoice not found")
	}

	if invoice.Status != "draft" {
		return nil, errors.New("adjustments can only be added to draft invoices")
	}

	adjustment := model.FinanceInvoiceAdjustment{
		InvoiceID:      invoiceID,
		SiswaNIS:       req.SiswaNIS,
		AdjustmentType: req.AdjustmentType,
		AdjustmentName: req.AdjustmentName,
		Amount:         req.Amount,
		Notes:          req.Notes,
	}

	if adjustment.AdjustmentType == "" {
		adjustment.AdjustmentType = "other"
	}

	err = c.db.Create(&adjustment).Error
	if err != nil {
		return nil, err
	}

	return &dto.InvoiceAdjustmentResponse{
		AdjustmentID:   adjustment.AdjustmentID,
		SiswaNIS:       adjustment.SiswaNIS,
		AdjustmentType: adjustment.AdjustmentType,
		AdjustmentName: adjustment.AdjustmentName,
		Amount:         adjustment.Amount,
		Notes:          adjustment.Notes,
		CreatedAt:      adjustment.CreatedAt,
	}, nil
}

func (c *invoiceController) UpdateInvoiceAdjustment(adjustmentID uint64, req *dto.InvoiceAdjustmentRequest) (*dto.InvoiceAdjustmentResponse, error) {
	var adjustment model.FinanceInvoiceAdjustment
	err := c.db.First(&adjustment, adjustmentID).Error
	if err != nil {
		return nil, errors.New("adjustment not found")
	}

	var invoice model.FinanceInvoice
	err = c.db.First(&invoice, adjustment.InvoiceID).Error
	if err != nil {
		return nil, err
	}

	if invoice.Status != "draft" {
		return nil, errors.New("adjustments can only be updated in draft invoices")
	}

	adjustment.AdjustmentType = req.AdjustmentType
	adjustment.AdjustmentName = req.AdjustmentName
	adjustment.Amount = req.Amount
	adjustment.Notes = req.Notes

	if adjustment.AdjustmentType == "" {
		adjustment.AdjustmentType = "other"
	}

	err = c.db.Save(&adjustment).Error
	if err != nil {
		return nil, err
	}

	return &dto.InvoiceAdjustmentResponse{
		AdjustmentID:   adjustment.AdjustmentID,
		SiswaNIS:       adjustment.SiswaNIS,
		AdjustmentType: adjustment.AdjustmentType,
		AdjustmentName: adjustment.AdjustmentName,
		Amount:         adjustment.Amount,
		Notes:          adjustment.Notes,
		CreatedAt:      adjustment.CreatedAt,
	}, nil
}

func (c *invoiceController) RemoveInvoiceAdjustment(adjustmentID uint64) error {
	var adjustment model.FinanceInvoiceAdjustment
	err := c.db.First(&adjustment, adjustmentID).Error
	if err != nil {
		return err
	}

	var invoice model.FinanceInvoice
	err = c.db.First(&invoice, adjustment.InvoiceID).Error
	if err != nil {
		return err
	}

	if invoice.Status != "draft" {
		return errors.New("adjustments can only be removed from draft invoices")
	}

	return c.db.Delete(&adjustment).Error
}
