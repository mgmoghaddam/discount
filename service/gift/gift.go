package gift

import (
	"discount/storage/gift"
	"fmt"
	"time"
)

type DTO struct {
	ID             int64     `json:"id"`
	Code           string    `json:"code"`
	GiftAmount     int64     `json:"giftAmount"`
	UsageLimit     int64     `json:"usageLimit"`
	UsedCount      int64     `json:"usedCount"`
	ExpirationDate time.Time `json:"expirationDate"`
	StartDateTime  time.Time `json:"startDateTime"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

type CreateRequest struct {
	CodePrefix     string `json:"codePrefix"`
	Code           string `json:"code"`
	GiftAmount     int64  `json:"giftAmount"`
	UsageLimit     int64  `json:"usageLimit"`
	ExpirationDate string `json:"expirationDate"`
	StartDateTime  string `json:"startDateTime"`
}

type UseGiftRequest struct {
	Code         string
	ResponseChan chan *DTO
	ErrorChan    chan error
}

var useGiftQueue = make(chan UseGiftRequest, 100) // adjust size as needed

func (s *Service) StartProcessingGifts() {
	for request := range useGiftQueue {
		g, err := s.gift.IncreaseUsedCountRedis(request.Code)
		if err != nil {
			request.ErrorChan <- err
		} else {
			request.ResponseChan <- s.FromDBModel(g)
		}
	}
}

func (s *Service) Create(r *CreateRequest) (*DTO, error) {
	var giftRecord *gift.Gift

	giftRecord = s.FromCreateRequest(r)

	err := s.ensureUniqueGiftCode(giftRecord, r.CodePrefix)
	if err != nil {
		return nil, err
	}

	err = s.gift.Create(giftRecord)
	if err != nil {
		return nil, err
	}

	return s.FromDBModel(giftRecord), nil
}

func (s *Service) GetByCode(code string) (*DTO, error) {
	g, err := s.gift.GetByCode(code)
	if err != nil {
		return nil, err
	}
	return s.FromDBModel(g), nil
}

func (s *Service) UpdateByCode(r *DTO) (*DTO, error) {
	giftRecord := s.ToDBModel(r)

	err := s.gift.UpdateDirectDb(giftRecord)
	if err != nil {
		return nil, err
	}

	return s.FromDBModel(giftRecord), nil
}

func (s *Service) UseGift(code string) (*DTO, error) {
	responseChan := make(chan *DTO)
	errorChan := make(chan error)
	useGiftQueue <- UseGiftRequest{Code: code, ResponseChan: responseChan, ErrorChan: errorChan}
	go s.StartProcessingGifts()
	select {
	case response := <-responseChan:
		return response, nil
	case err := <-errorChan:
		return nil, err
	}
}

// SyncGifts syncs updates gifts in redis to the database.
// It is called by the scheduler every 30 seconds.
func (s *Service) syncGift() error {
	err := s.gift.SyncRedisWithDB()
	if err != nil {
		return err
	}
	return nil
}

// ensureUniqueGiftCode ensures that the gift code is unique by generating a unique code and checking against existing codes.
// If the gift already has a code, the function returns nil without generating a new code.
// If the generated code already exists, it generates a new code and continues checking until a unique code is found.
func (s *Service) ensureUniqueGiftCode(giftRecord *gift.Gift, codePrefix string) error {
	if giftRecord.Code != "" {
		return nil
	}
	giftRecord.Code = s.generateCode(codePrefix)
	for {
		g, _ := s.gift.GetByCode(giftRecord.Code)
		if g == nil {
			break
		}
		giftRecord.Code = s.generateCode(codePrefix)
	}
	return nil
}

func (s *Service) generateCode(prefix string) string {
	s.mu.Lock()
	defer s.mu.Unlock()
	timestamp := time.Now().UnixNano() / int64(time.Microsecond)
	return fmt.Sprintf("%s-%06d", prefix, timestamp%1e6)
}
