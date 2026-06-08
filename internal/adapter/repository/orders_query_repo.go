package repository

import (
	"context"

	"github.com/yadukrishnan2004/ecommerce-backend/internal/domain"
)

func (o *Order) ToDomain() domain.Order {
	return domain.Order{
		ID:                o.ID,
		UserID:            o.UserID,
		User:              *o.User.ToDomain(),
		AddressID:         o.AddressID,
		Address:           *o.Address.ToDomain(),
		Status:            o.Status,
		Quantity:          uint(len(o.Items)),
		TotalAmount:       o.TotalAmount,
		PaymentMethod:     o.PaymentMethod,
		RazorpayOrderID:   o.RazorpayOrderID,
		RazorpayPaymentID: o.RazorpayPaymentID,
	}
}

func (oi *OrderItem) ToDomain() domain.OrderItem {
	var img string
	if len(oi.Product.Images) > 0 {
		img = oi.Product.Images[0].ImageURL
	}
	return domain.OrderItem{
		OrderId:   oi.OrderId,
		Order:     oi.Order.ToDomain(),
		Image:     img,
		ProductId: oi.ProductId,
		Product:   *oi.Product.ToDomain(),
		Quantity:  oi.Quantity,
		Price:     oi.Price,
	}
}

func FromDomainOrderItem(oi domain.OrderItem) OrderItem {
	return OrderItem{
		OrderId:   oi.OrderId,
		ProductId: oi.ProductId,
		Quantity:  oi.Quantity,
		Price:     oi.Price,
	}
}

func (r *orderRepo) GetOrdersByUserIDAndOrderID(ctx context.Context, userID, OrderID uint) ([]domain.Order, error) {
	var dbOrders []Order
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Address.PostalCode").
		Preload("Items").
		Where("user_id = ? AND id = ?", userID, OrderID).
		Order("created_at desc").
		Find(&dbOrders).Error

	if err != nil {
		return nil, err
	}

	var orders []domain.Order
	for _, o := range dbOrders {
		orders = append(orders, o.ToDomain())
	}
	return orders, nil
}

func (r *orderRepo) GetAllOrders(ctx context.Context, limit, offset int) ([]domain.Order, error) {
	var dbOrders []Order

	dbQuery := r.db.WithContext(ctx).
		Preload("User").
		Preload("Address.PostalCode").
		Preload("Items").
		Order("created_at desc")

	if limit > 0 {
		dbQuery = dbQuery.Limit(limit)
	}
	if offset > 0 {
		dbQuery = dbQuery.Offset(offset)
	}

	err := dbQuery.Find(&dbOrders).Error

	if err != nil {
		return nil, err
	}

	var orders []domain.Order
	for _, o := range dbOrders {
		orders = append(orders, o.ToDomain())
	}

	return orders, nil
}

func (r *orderRepo) GetAllOrdersByUserID(ctx context.Context, userID uint, limit, offset int) ([]domain.Order, error) {
	var dbOrders []Order

	dbQuery := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Preload("User").
		Preload("Address.PostalCode").
		Preload("Items").
		Order("created_at desc")

	if limit > 0 {
		dbQuery = dbQuery.Limit(limit)
	}
	if offset > 0 {
		dbQuery = dbQuery.Offset(offset)
	}

	err := dbQuery.Find(&dbOrders).Error

	if err != nil {
		return nil, err
	}

	var orders []domain.Order
	for _, o := range dbOrders {
		orders = append(orders, o.ToDomain())
	}
	return orders, nil
}

func (r *orderRepo) GetOrdersByOrderID(ctx context.Context, OrderID uint) ([]domain.OrderItem, error) {
	var dbOrderItems []OrderItem
	err := r.db.WithContext(ctx).
		Preload("Product.Images").
		Preload("Order.User").
		Preload("Order.Address.PostalCode").
		Where("order_id = ?", OrderID).
		Order("created_at desc").
		Find(&dbOrderItems).Error

	if err != nil {
		return nil, err
	}

	var orderItems []domain.OrderItem
	for _, item := range dbOrderItems {
		orderItems = append(orderItems, item.ToDomain())
	}

	return orderItems, nil
}

func (r *orderRepo) GetOrderItemsByOrderIDs(ctx context.Context, orderIDs []uint) ([]domain.OrderItem, error) {
	if len(orderIDs) == 0 {
		return []domain.OrderItem{}, nil
	}
	var dbOrderItems []OrderItem
	err := r.db.WithContext(ctx).
		Preload("Product.Images").
		Preload("Order.User").
		Preload("Order.Address.PostalCode").
		Where("order_id IN ?", orderIDs).
		Order("created_at desc").
		Find(&dbOrderItems).Error

	if err != nil {
		return nil, err
	}

	var orderItems []domain.OrderItem
	for _, item := range dbOrderItems {
		orderItems = append(orderItems, item.ToDomain())
	}

	return orderItems, nil
}
