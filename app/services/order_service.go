package services

import (
	"elcomware/app/actions"
	"elcomware/app/dto"
	"elcomware/app/mappers"
	"elcomware/app/queries"
)

type OrderService struct {
	orderQueries   *queries.OrderQueries
	orderMapper    *mappers.OrderMapper
	checkoutAction *actions.ExecuteCheckoutAction
}

func NewOrderService() *OrderService {
	return &OrderService{
		orderQueries:   queries.NewOrderQueries(),
		orderMapper:    mappers.NewOrderMapper(),
		checkoutAction: actions.NewExecuteCheckoutAction(),
	}
}

func (s *OrderService) Checkout(input dto.CheckoutInput) (dto.OrderOutput, error) {
	return s.checkoutAction.Execute(input)
}

func (s *OrderService) GetAllOrders() ([]dto.OrderOutput, error) {
	orders, err := s.orderQueries.FindAll()
	if err != nil {
		return nil, err
	}
	return s.orderMapper.ToDTOList(orders), nil
}

func (s *OrderService) GetUserOrders(userID string) ([]dto.OrderOutput, error) {
	orders, err := s.orderQueries.FindByUserID(userID)
	if err != nil {
		return nil, err
	}
	return s.orderMapper.ToDTOList(orders), nil
}
