package service

type OrdersRepository interface {
}

type BrokerService interface {
}

type CacheService interface {
}

type Service struct {
	repo   OrdersRepository
	cache  CacheService
	broker BrokerService
}

func (s Service) GetOrder() {

}
