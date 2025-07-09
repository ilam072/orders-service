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
	Get(key string) (dto.Order, error)
}

type OrderService struct {
	orderRepo OrderRepo
	cache     OrderCache
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
		PaymentDt:    time.Unix(int64(dto.Payment.PaymentDt), 0),
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
