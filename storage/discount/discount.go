package discount

import (
	"discount/internal/serr"
	"time"
)

const discountColumns = "id" +
	",code,percent_off,discount_amount,usage_limit,used_count,expiration_date,start_date_time,max_amount" +
	",min_amount,created_at,updated_at"

type Discount struct {
	ID             int64     `db:"id"`
	Code           string    `db:"code"`
	PercentOff     int64     `db:"percent_off"`
	DiscountAmount int64     `db:"discount_amount"`
	UsageLimit     int64     `db:"usage_limit"`
	UsedCount      int64     `db:"used_count"`
	ExpirationDate time.Time `db:"expiration_date"`
	StartDateTime  time.Time `db:"start_date_time"`
	MaxAmount      int64     `db:"max_amount"`
	MinAmount      int64     `db:"min_amount"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}

func (s Storage) Create(d *Discount) error {
	sqlStmt := `
	INSERT INTO discount (code, percent_off, discount_amount, usage_limit, used_count, expiration_date, start_date_time,
	                      max_amount, min_amount)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) 
	                     RETURNING id, code, created_at, updated_at`
	err := s.db.QueryRow(sqlStmt, d.Code, d.PercentOff, d.DiscountAmount, d.UsageLimit, d.UsedCount, d.ExpirationDate,
		d.StartDateTime, d.MaxAmount, d.MinAmount).Scan(&d.ID, &d.Code, &d.CreatedAt, &d.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (s Storage) CreateBulk(discounts []*Discount) error {
	for _, discount := range discounts {
		err := s.Create(discount)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s Storage) Update(d *Discount) error {
	sqlStmt := `
	UPDATE discount SET code = $1, percent_off = $2, discount_amount = $3, usage_limit = $4, used_count = $5, 
	                   expiration_date = $6, start_date_time = $7, max_amount = $8, min_amount = $9, updated_at = now()
	WHERE id = $10 RETURNING updated_at`
	err := s.db.QueryRow(sqlStmt, d.Code, d.PercentOff, d.DiscountAmount, d.UsageLimit, d.UsedCount, d.ExpirationDate,
		d.StartDateTime, d.MaxAmount, d.MinAmount, d.ID).Scan(&d.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (s Storage) GetByCode(code string) (*Discount, error) {
	sqlStmt := "SELECT " + discountColumns + " FROM discount WHERE code = $1"
	d := &Discount{}
	err := s.db.QueryRow(sqlStmt, code).Scan(&d.ID, &d.Code, &d.PercentOff, &d.DiscountAmount, &d.UsageLimit,
		&d.UsedCount, &d.ExpirationDate, &d.StartDateTime, &d.MaxAmount, &d.MinAmount, &d.CreatedAt, &d.UpdatedAt)
	if err != nil {
		return nil, serr.ValidationErr("code", "gift", serr.ErrInvalidDiscountCode)
	}
	return d, nil
}

func (s Storage) GetByID(id int64) (*Discount, error) {
	sqlStmt := "SELECT " + discountColumns + " FROM discount WHERE id = $1"
	d := &Discount{}
	err := s.db.QueryRow(sqlStmt, id).Scan(&d.ID, &d.Code, &d.PercentOff, &d.DiscountAmount, &d.UsageLimit,
		&d.UsedCount, &d.ExpirationDate, &d.StartDateTime, &d.MaxAmount, &d.MinAmount, &d.CreatedAt, &d.UpdatedAt)
	if err != nil {
		return nil, serr.ValidationErr("code", "gift", serr.ErrInvalidDiscountID)
	}
	return d, nil
}

func (s Storage) Delete(id int64) error {
	sqlStmt := "DELETE FROM discount WHERE id = $1"
	row, err := s.db.Exec(sqlStmt, id)
	if err != nil {
		return err
	}
	if count, err := row.RowsAffected(); err != nil || count == 0 {
		return ErrNoRowToUpdate
	}
	return nil
}

func (s Storage) GetAllByPage(limit, offset int, count bool) ([]*Discount, int, error) {
	var total int
	if count {
		err := s.db.QueryRow("SELECT count(*) FROM discount").Scan(&total)
		if err != nil {
			return nil, 0, err
		}
	}
	pagination := " LIMIT $1 OFFSET $2"
	order := " ORDER BY created_at DESC"
	rows, err := s.db.Query("SELECT "+discountColumns+" FROM discount"+order+pagination, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	discounts := make([]*Discount, 0)
	for rows.Next() {
		d := &Discount{}
		err := rows.Scan(&d.ID, &d.Code, &d.PercentOff, &d.DiscountAmount, &d.UsageLimit, &d.UsedCount, &d.ExpirationDate,
			&d.StartDateTime, &d.MaxAmount, &d.MinAmount, &d.CreatedAt, &d.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		discounts = append(discounts, d)
	}
	return discounts, total, nil
}
