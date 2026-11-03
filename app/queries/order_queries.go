package queries

import (
	"elcomware/app/facades"
	"elcomware/app/models"
)

type OrderQueries struct{}

func NewOrderQueries() *OrderQueries {
	return &OrderQueries{}
}

func (q *OrderQueries) FindAll() ([]models.Order, error) {
	var orders []models.Order
	if err := facades.Orm().Query().OrderByDesc("created_at").Get(&orders); err != nil {
		return nil, err
	}

	for i := range orders {
		var items []models.OrderItem
		_ = facades.Orm().Query().Where("order_id = ?", orders[i].ID).Get(&items)
		orders[i].Items = items
	}

	return orders, nil
}

func (q *OrderQueries) FindByUserID(userID string) ([]models.Order, error) {
	var orders []models.Order
	if err := facades.Orm().Query().Where("user_id = ?", userID).OrderByDesc("created_at").Get(&orders); err != nil {
		return nil, err
	}

	for i := range orders {
		var items []models.OrderItem
		_ = facades.Orm().Query().Where("order_id = ?", orders[i].ID).Get(&items)
		orders[i].Items = items
	}

	return orders, nil
}

func (q *OrderQueries) CountActive() (int64, error) {
	return facades.Orm().Query().Model(&models.Order{}).Count()
}

func (q *OrderQueries) SumRevenue() (float64, error) {
	var total float64
	err := facades.Orm().Query().Model(&models.Order{}).Sum("total_amount", &total)
	return total, err
}
