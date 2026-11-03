package services

import (
	"elcomware/app/actions"
	"elcomware/app/http/requests"
	"elcomware/app/models"
	"elcomware/app/queries"
)

type OrderService struct {
	orderQueries   *queries.OrderQueries
	checkoutAction *actions.ExecuteCheckoutAction
}

func NewOrderService() *OrderService {
	return &OrderService{
		orderQueries:   queries.NewOrderQueries(),
		checkoutAction: actions.NewExecuteCheckoutAction(),
	}
}

func (s *OrderService) PlaceOrder(req *requests.CheckoutRequest, user *models.User) (*models.Order, error) {
	return s.checkoutAction.ExecuteFromRequest(req, user)
}

func (s *OrderService) GetAllOrders() ([]models.Order, error) {
	return s.orderQueries.FindAll()
}

func (s *OrderService) GetUserOrders(userID string) ([]models.Order, error) {
	return s.orderQueries.FindByUserID(userID)
}
