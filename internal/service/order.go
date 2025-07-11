package service

import (
	"context"
	"github.com/google/uuid"
	"time"
	"wb-l0/internal/types/domain"
	"wb-l0/internal/types/dto"
	"wb-l0/pkg/e"
)

type OrderRepo interface {
	CreateOrder(context.Context, domain.Order, domain.Delivery, domain.Payment, []domain.Item) error
	GetOrder(ctx context.Context, ID string) (domain.FullOrder, error)
}

type OrderCache interface {
	Set(key string, order dto.Order)
	Get(key string) (dto.Order, bool)
}

type OrderService struct {
	orderRepo OrderRepo
	cache     OrderCache
}

func NewOrderService(repo OrderRepo, cache OrderCache) *OrderService {
	return &OrderService{
		orderRepo: repo,
		cache:     cache,
	}
}

func (s OrderService) GetOrder(ctx context.Context, orderId string) (dto.Order, error) {
	const op = "OrderService.GetOrder()"

	order, ok := s.cache.Get(orderId)
	if ok {
		return order, nil
	}

	fullOrder, err := s.orderRepo.GetOrder(ctx, orderId)
	if err != nil {
		return dto.Order{}, e.Wrap(op, err)
	}

	return domainToDtoOrder(fullOrder), nil
}

func (s OrderService) CreateOrder(ctx context.Context, order dto.Order) error {
	const op = "OrderService.CreateOrder()"

	domainOrder, delivery, payment, items, err := dtoToDomainOrder(order)
	if err != nil {
		return e.Wrap(op, err)
	}

	if err := s.orderRepo.CreateOrder(ctx, domainOrder, delivery, payment, items); err != nil {
		return e.Wrap(op, err)
	}

	s.cache.Set(order.OrderUID, order)
	return nil
}

func dtoToDomainOrder(dto dto.Order) (domain.Order, domain.Delivery, domain.Payment, []domain.Item, error) {
	orderUID, err := uuid.Parse(dto.OrderUID)
	if err != nil {
		return domain.Order{}, domain.Delivery{}, domain.Payment{}, nil, err
	}

	order := domain.Order{
		ID:                orderUID,
		TrackNumber:       dto.TrackNumber,
		Entry:             dto.Entry,
		Locale:            dto.Locale,
		InternalSignature: dto.InternalSignature,
		CustomerID:        dto.CustomerID,
		DeliveryService:   dto.DeliveryService,
		ShardKey:          dto.Shardkey,
		SmID:              dto.SmID,
		DateCreated:       dto.DateCreated,
		OofShard:          dto.OofShard,
	}

	delivery := domain.Delivery{
		OrderID: order.ID,
		Name:    dto.Delivery.Name,
		Phone:   dto.Delivery.Phone,
		Zip:     dto.Delivery.Zip,
		City:    dto.Delivery.City,
		Address: dto.Delivery.Address,
		Region:  dto.Delivery.Region,
		Email:   dto.Delivery.Email,
	}

	payment := domain.Payment{
		Transaction:  order.ID,
		OrderID:      order.ID,
		RequestID:    dto.Payment.RequestID,
		Currency:     dto.Payment.Currency,
		Provider:     dto.Payment.Provider,
		Amount:       dto.Payment.Amount,
		PaymentDt:    time.Unix(dto.Payment.PaymentDt, 0),
		Bank:         dto.Payment.Bank,
		DeliveryCost: dto.Payment.DeliveryCost,
		GoodsTotal:   dto.Payment.GoodsTotal,
		CustomFee:    dto.Payment.CustomFee,
	}

	var items []domain.Item
	for _, itm := range dto.Items {
		item := domain.Item{
			OrderID:     order.ID,
			TrackNumber: itm.TrackNumber,
			Price:       itm.Price,
			Rid:         itm.Rid,
			Name:        itm.Name,
			Sale:        itm.Sale,
			Size:        itm.Size,
			TotalPrice:  itm.TotalPrice,
			NmID:        int64(itm.NmID),
			Brand:       itm.Brand,
			Status:      itm.Status,
		}
		items = append(items, item)
	}

	return order, delivery, payment, items, nil
}

func domainToDtoOrder(fullOrder domain.FullOrder) dto.Order {

	delivery := dto.Delivery{
		Name:    fullOrder.Delivery.Name,
		Phone:   fullOrder.Delivery.Phone,
		Zip:     fullOrder.Delivery.Zip,
		City:    fullOrder.Delivery.City,
		Address: fullOrder.Delivery.Address,
		Region:  fullOrder.Delivery.Region,
		Email:   fullOrder.Delivery.Email,
	}

	payment := dto.Payment{
		Transaction:  fullOrder.Payment.Transaction.String(),
		RequestID:    fullOrder.Payment.RequestID,
		Currency:     fullOrder.Payment.Currency,
		Provider:     fullOrder.Payment.Provider,
		Amount:       fullOrder.Payment.Amount,
		PaymentDt:    fullOrder.Payment.PaymentDt.Unix(),
		Bank:         fullOrder.Payment.Bank,
		DeliveryCost: fullOrder.Payment.DeliveryCost,
		GoodsTotal:   fullOrder.Payment.GoodsTotal,
		CustomFee:    fullOrder.Payment.CustomFee,
	}

	var items []dto.Item
	for _, itm := range fullOrder.Items {
		item := dto.Item{
			ChrtID:      int(itm.ChrtID),
			TrackNumber: itm.TrackNumber,
			Price:       itm.Price,
			Rid:         itm.Rid,
			Name:        itm.Name,
			Sale:        itm.Sale,
			Size:        itm.Size,
			TotalPrice:  itm.TotalPrice,
			NmID:        int(itm.NmID),
			Brand:       itm.Brand,
			Status:      itm.Status,
		}
		items = append(items, item)
	}

	return dto.Order{
		OrderUID:          fullOrder.Order.ID.String(),
		TrackNumber:       fullOrder.Order.TrackNumber,
		Entry:             fullOrder.Order.Entry,
		Delivery:          delivery,
		Payment:           payment,
		Items:             items,
		Locale:            fullOrder.Order.Locale,
		InternalSignature: fullOrder.Order.InternalSignature,
		CustomerID:        fullOrder.Order.CustomerID,
		DeliveryService:   fullOrder.Order.DeliveryService,
		Shardkey:          fullOrder.Order.ShardKey,
		SmID:              fullOrder.Order.SmID,
		DateCreated:       fullOrder.Order.DateCreated,
		OofShard:          fullOrder.Order.OofShard,
	}

}
