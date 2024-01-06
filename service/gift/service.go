package gift

import (
	"database/sql"
	"discount/storage/gift"
	"github.com/jasonlvhit/gocron"
	"github.com/rs/zerolog/log"
	"sync"
	"time"
)

type Service struct {
	gift gift.Storage
	mu   sync.Mutex

	inTx bool
}

func New(
	gift gift.Storage,
) *Service {
	s := &Service{
		gift: gift,
	}
	err := gocron.Every(30).Seconds().Do(func() {
		if err := s.syncGift(); err != nil {
			log.Error().Err(err).Msg("failed to sync gifts")
		}
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to start gocron for sync gifts")
	} else {
		gocron.Start()
	}
	return s
}

func (s *Service) withTX(tx *sql.Tx) (*Service, error) {
	service := *s
	g, err := s.gift.WithTX(tx)
	if err != nil {
		return nil, err
	}
	service.gift = g
	service.inTx = true
	return &service, nil
}

func (s *Service) ToDBModel(g *DTO) *gift.Gift {
	return &gift.Gift{
		ID:             g.ID,
		Code:           g.Code,
		GiftAmount:     g.GiftAmount,
		UsageLimit:     g.UsageLimit,
		UsedCount:      g.UsedCount,
		ExpirationDate: g.ExpirationDate,
		StartDateTime:  g.StartDateTime,
	}

}

func (s *Service) FromDBModel(g *gift.Gift) *DTO {
	return &DTO{
		ID:             g.ID,
		Code:           g.Code,
		GiftAmount:     g.GiftAmount,
		UsageLimit:     g.UsageLimit,
		UsedCount:      g.UsedCount,
		ExpirationDate: g.ExpirationDate,
		StartDateTime:  g.StartDateTime,
		CreatedAt:      g.CreatedAt,
		UpdatedAt:      g.UpdatedAt,
	}
}

func (s *Service) FromCreateRequest(r *CreateRequest) *gift.Gift {
	const layout = "2006-01-02"
	exDate, _ := time.Parse(layout, r.ExpirationDate)
	stDate, _ := time.Parse(layout, r.StartDateTime)
	return &gift.Gift{
		Code:           r.Code,
		GiftAmount:     r.GiftAmount,
		UsageLimit:     r.UsageLimit,
		ExpirationDate: exDate,
		StartDateTime:  stDate,
	}
}
