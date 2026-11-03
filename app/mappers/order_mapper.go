package mappers

import (
	"elcomware/app/dto"
	"elcomware/app/models"
)

type OrderMapper struct{}

func NewOrderMapper() *OrderMapper {
	return &OrderMapper{}
}

func (m *OrderMapper) ToItemDTO(item models.OrderItem) dto.OrderItemOutput {
	return dto.OrderItemOutput{
		ID:                 item.ID,
		ProductID:          item.ProductID,
		ProductTitle:       item.ProductTitle,
		ProductType:        item.ProductType,
		Quantity:           item.Quantity,
		UnitPrice:          item.UnitPrice,
		TotalPrice:         item.TotalPrice,
		DigitalLicenseKey:  item.DigitalLicenseKey,
		DigitalDownloadUrl: item.DigitalDownloadUrl,
	}
}

func (m *OrderMapper) ToDTO(order models.Order) dto.OrderOutput {
	items := make([]dto.OrderItemOutput, 0, len(order.Items))
	for _, it := range order.Items {
		items = append(items, m.ToItemDTO(it))
	}

	return dto.OrderOutput{
		ID:            order.ID,
		OrderNumber:   order.OrderNumber,
		UserID:        order.UserID,
		CustomerName:  order.CustomerName,
		CustomerEmail: order.CustomerEmail,
		Items:         items,
		TotalAmount:   order.TotalAmount,
		Currency:      order.Currency,
		Status:        string(order.Status),
		PaymentMethod: order.PaymentMethod,
		HasDigital:    order.HasDigital,
		HasPhysical:   order.HasPhysical,
		CreatedAt:     order.CreatedAt,
	}
}

func (m *OrderMapper) ToDTOList(list []models.Order) []dto.OrderOutput {
	dtos := make([]dto.OrderOutput, 0, len(list))
	for _, o := range list {
		dtos = append(dtos, m.ToDTO(o))
	}
	return dtos
}
