package repository

import (
	"context"

	"github.com/yadukrishnan2004/ecommerce-backend/internal/domain"
)

func (r *orderRepo) GetTotalSalesByDate(ctx context.Context) ([]domain.SalesData, error) {
	var sales []domain.SalesData
	err := r.db.WithContext(ctx).
		Model(&Order{}).
		Select("TO_CHAR(created_at, 'YYYY-MM-DD') as date, SUM(total_amount) as total").
		Where("status = ?", "Delivered").
		Group("TO_CHAR(created_at, 'YYYY-MM-DD')").
		Order("date asc").
		Limit(30).
		Scan(&sales).Error

	if err != nil {
		return nil, err
	}
	return sales, nil
}

func (r *orderRepo) GetOrderCountsByStatus(ctx context.Context) ([]domain.StatusCount, error) {
	var counts []domain.StatusCount
	err := r.db.WithContext(ctx).
		Model(&Order{}).
		Select("status, count(*) as count").
		Group("status").
		Scan(&counts).Error

	if err != nil {
		return nil, err
	}
	return counts, nil
}

func (r *orderRepo) GetDashboardMetrics(ctx context.Context) (totalRevenue float64, totalOrders int64, averageOrderValue float64, err error) {
	err = r.db.WithContext(ctx).
		Model(&Order{}).
		Select("COALESCE(SUM(total_amount), 0)").
		Where("status NOT IN (?, ?)", "Cancelled", "Pending").
		Scan(&totalRevenue).Error
	if err != nil {
		return
	}

	err = r.db.WithContext(ctx).
		Model(&Order{}).
		Where("status NOT IN (?, ?)", "Cancelled", "Pending").
		Count(&totalOrders).Error
	if err != nil {
		return
	}

	if totalOrders > 0 {
		averageOrderValue = totalRevenue / float64(totalOrders)
	}
	return
}

func (r *orderRepo) GetTopSellingProducts(ctx context.Context, limit int) ([]domain.TopProduct, error) {
	var topProducts []domain.TopProduct
	err := r.db.WithContext(ctx).
		Table("order_items").
		Select("order_items.product_id, products.name, SUM(order_items.quantity) as quantity_sold, SUM(order_items.quantity * order_items.price) as total_revenue").
		Joins("JOIN products ON products.id = order_items.product_id").
		Joins("JOIN orders ON orders.id = order_items.order_id").
		Where("orders.status NOT IN (?, ?)", "Cancelled", "Pending").
		Group("order_items.product_id, products.name").
		Order("quantity_sold DESC").
		Limit(limit).
		Scan(&topProducts).Error

	if err != nil {
		return nil, err
	}
	return topProducts, nil
}
