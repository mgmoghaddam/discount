package gift

import (
	"context"
	"discount/db"
	"discount/internal/serr"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"time"
)

const giftColumns = "id" +
	",code,gift_amount,usage_limit,used_count,expiration_date,start_date_time,created_at,updated_at"

const giftPrefix = "GIFT:%s"
const giftPrefixUpdate = "UPDATED_GIFT:%s"

type Gift struct {
	ID             int64     `db:"id"`
	Code           string    `db:"code"`
	GiftAmount     int64     `db:"gift_amount"`
	UsageLimit     int64     `db:"usage_limit"`
	UsedCount      int64     `db:"used_count"`
	ExpirationDate time.Time `db:"expiration_date"`
	StartDateTime  time.Time `db:"start_date_time"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}

// Create inserts a new gift into the storage.
func (s Storage) Create(g *Gift) error {
	sqlStmt := `
	INSERT INTO gift (code, gift_amount, usage_limit, used_count, expiration_date, start_date_time)
	VALUES ($1, $2, $3, $4, $5, $6) 
	                     RETURNING id, code, created_at, updated_at`
	err := s.db.QueryRow(sqlStmt, g.Code, g.GiftAmount, g.UsageLimit, g.UsedCount, g.ExpirationDate,
		g.StartDateTime).Scan(&g.ID, &g.Code, &g.CreatedAt, &g.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

// CreateBulk inserts a new gift into the storage.
func (s Storage) CreateBulk(gifts []*Gift) error {
	for _, gift := range gifts {
		err := s.Create(gift)
		if err != nil {
			return err
		}
	}
	return nil
}

// UpdateDirectDb updates the gift in the storage. It first calls the SyncRedisWithDB method to synchronize the data with Redis
// and the database. If there is an error during synchronization
func (s Storage) UpdateDirectDb(g *Gift) error {
	err := s.SyncRedisWithDB()
	if err != nil {
		log.Err(err).Msg("SyncRedisWithDB")
	}
	sqlStmt := `
	UPDATE gift SET code = $1, gift_amount = $2, usage_limit = $3, used_count = $4, 
	                   expiration_date = $5, start_date_time = $6, updated_at = now()
	WHERE id = $7 RETURNING updated_at`
	err = s.db.QueryRow(sqlStmt, g.Code, g.GiftAmount, g.UsageLimit, g.UsedCount, g.ExpirationDate,
		g.StartDateTime, g.ID).Scan(&g.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (s Storage) Update(g *Gift) error {
	sqlStmt := `
	UPDATE gift SET code = $1, gift_amount = $2, usage_limit = $3, used_count = $4, 
	                   expiration_date = $5, start_date_time = $6, updated_at = now()
	WHERE id = $7 RETURNING updated_at`
	err := s.db.QueryRow(sqlStmt, g.Code, g.GiftAmount, g.UsageLimit, g.UsedCount, g.ExpirationDate,
		g.StartDateTime, g.ID).Scan(&g.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

// GetByCode retrieves a gift from the storage based on its unique code.
// It first checks if the gift is available in Redis cache using the key pattern giftPrefixUpdate followed by the gift code.
// If it is found, the gift is returned.
// Otherwise, it checks if the gift is available in Redis cache using the key pattern giftPrefix followed by the gift code.
// If it is found, the gift is returned.
// If the gift is not found in Redis cache, it queries the gift from the database based on the code.
func (s Storage) GetByCode(code string) (*Gift, error) {
	keyUpdate := fmt.Sprintf(giftPrefixUpdate, code)
	g, err := s.retrieveGiftFromRedis(keyUpdate)
	if err != nil {
		key := fmt.Sprintf(giftPrefix, code)
		g, err = s.retrieveGiftFromRedis(key)
		if err != nil {
			gift := &Gift{}
			sqlStmt := "SELECT " + giftColumns + " FROM gift WHERE code = $1"
			err := s.db.QueryRow(sqlStmt, code).Scan(&gift.ID, &gift.Code, &gift.GiftAmount, &gift.UsageLimit, &gift.UsedCount,
				&gift.ExpirationDate, &gift.StartDateTime, &gift.CreatedAt, &gift.UpdatedAt)
			if err != nil {
				return nil, serr.ValidationErr("code", "invalid discount code", serr.ErrInvalidGiftCode)
			}
			err = s.updateOrInsertGiftInRedis(key, gift, time.Minute*10)
			if err != nil {
				return nil, err
			}
			return gift, nil
		}
		return g, nil
	}
	return g, nil
}

func (s Storage) GetByID(id int64) (*Gift, error) {
	gift := &Gift{}
	sqlStmt := "SELECT " + giftColumns + " FROM gift WHERE id = $1"
	err := s.db.QueryRow(sqlStmt, id).Scan(&gift.ID, &gift.Code, &gift.GiftAmount, &gift.UsageLimit, &gift.UsedCount,
		&gift.ExpirationDate, &gift.StartDateTime, &gift.CreatedAt, &gift.UpdatedAt)
	if err != nil {
		return nil, serr.ValidationErr("code", "gift", serr.ErrInvalidGiftID)
	}
	return gift, nil
}

// IncreaseUsedCountRedis increases the used count of a gift in Redis cache.
// Consider that the gift save in Redis AOF to prevent data loss.
func (s Storage) IncreaseUsedCountRedis(code string) (*Gift, error) {
	keyUpdate := fmt.Sprintf(giftPrefixUpdate, code)
	gift, err := s.GetByCode(code)
	if err != nil {
		return nil, err
	}

	if gift.UsageLimit > 0 && gift.UsedCount+1 > gift.UsageLimit {
		return nil, serr.ValidationErr("code", "gift usage limit reached", serr.ErrGiftUsageLimitReached)
	}
	gift.UsedCount++
	err = s.updateOrInsertGiftInRedis(keyUpdate, gift, 0)
	if err != nil {
		return nil, err
	}
	return gift, nil
}

func (s Storage) IncreaseUsedCount(code string) error {
	sqlStmt := `
	UPDATE gift SET used_count = used_count + 1, updated_at = now()
	WHERE code = $1 RETURNING updated_at`
	row, err := s.db.Exec(sqlStmt, code)
	if err != nil {
		return err
	}
	if count, err := row.RowsAffected(); err != nil || count == 0 {
		return ErrNoRowToUpdate
	}
	return nil
}

func (s Storage) Delete(id int64) error {
	g, err := s.GetByID(id)
	s.removeGiftFromRedis(g)
	sqlStmt := "DELETE FROM gift WHERE id = $1"
	row, err := s.db.Exec(sqlStmt, id)
	if err != nil {
		return err
	}
	if count, err := row.RowsAffected(); err != nil || count == 0 {
		return ErrNoRowToUpdate
	}
	return nil
}

func (s Storage) DeleteByCode(code string) error {
	s.removeGiftFromRedis(&Gift{Code: code})
	sqlStmt := "DELETE FROM gift WHERE code = $1"
	row, err := s.db.Exec(sqlStmt, code)
	if err != nil {
		return err
	}
	if count, err := row.RowsAffected(); err != nil || count == 0 {
		return ErrNoRowToUpdate
	}
	return nil
}

func (s Storage) DeleteBulkByIDs(ids []int64) error {
	for _, id := range ids {
		err := s.Delete(id)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s Storage) DeleteBulkByCodes(codes []string) error {
	for _, code := range codes {
		s.removeGiftFromRedis(&Gift{Code: code})
		err := s.DeleteByCode(code)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s Storage) getAllByPage(limit, offset int, count bool) ([]*Gift, int, error) {
	var total int
	if count {
		err := s.db.QueryRow("SELECT count(*) FROM gift").Scan(&total)
		if err != nil {
			return nil, 0, serr.DBError("List", "gift", err)
		}
	}
	pagination := fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)
	order := " ORDER BY created_at DESC"
	rows, err := s.db.Query("SELECT " + giftColumns + " FROM gift" + order + pagination)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	gifts := make([]*Gift, 0)
	for rows.Next() {
		g, err := s.scanGift(rows)
		if err != nil {
			return nil, 0, serr.DBError("List", "gift", err)
		}
		gifts = append(gifts, g)
	}
	return gifts, total, nil
}

func (s Storage) SyncRedisWithDB() error {
	keys, err := s.redis.Keys(context.Background(), "UPDATED_GIFT:*").Result()
	if err != nil {
		return err
	}
	for _, key := range keys {
		g, err := s.retrieveGiftFromRedis(key)
		if err != nil {
			return err
		}
		err = s.Update(g)
		if err != nil {
			return err
		}
		s.removeGiftFromRedis(g)
		return nil
	}
	return nil
}

func (s Storage) scanGift(scanner db.Scanner) (*Gift, error) {
	g := &Gift{}
	err := scanner.Scan(&g.ID, &g.Code, &g.GiftAmount, &g.UsageLimit, &g.UsedCount, &g.ExpirationDate,
		&g.StartDateTime, &g.CreatedAt, &g.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return g, nil
}

func (s Storage) updateOrInsertGiftInRedis(key string, g *Gift, exp time.Duration) error {
	err := s.redis.Set(context.Background(), key, g, exp).Err()
	if err != nil {
		return err
	}
	return nil
}

func (s Storage) retrieveGiftFromRedis(key string) (*Gift, error) {
	g, err := s.redis.Get(context.Background(), key).Result()
	if err != nil {
		return nil, err
	}
	gift := &Gift{}
	err = json.Unmarshal([]byte(g), &gift)
	if err != nil {
		return nil, err
	}
	return gift, nil
}

func (s Storage) removeGiftFromRedis(g *Gift) {
	key := fmt.Sprintf(giftPrefix, g.Code)
	keyUpdate := fmt.Sprintf(giftPrefixUpdate, g.Code)
	s.removeWithKey(key)
	s.removeWithKey(keyUpdate)
}

func (s Storage) removeWithKey(k string) {
	s.redis.Del(context.Background(), k)
}

func (g *Gift) MarshalBinary() ([]byte, error) {
	return json.Marshal(g)
}
