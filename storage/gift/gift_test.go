package gift_test

import (
	"database/sql"
	"discount/storage/gift"
	"github.com/redis/go-redis/v9"
	"log"
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	storage := setup()
	//defer teardown(storage) // cleanup your storage here

	g := &gift.Gift{
		Code:           "TEST",
		GiftAmount:     100,
		UsageLimit:     10,
		UsedCount:      0,
		ExpirationDate: time.Now().AddDate(0, 0, 10),
		StartDateTime:  time.Now(),
	}

	err := storage.Create(g)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if g.ID != 0 {
		t.Logf("id set %v", g.ID)
	}
	err = storage.Delete(g.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestCreateBulk(t *testing.T) {
	storage := setup()
	//defer teardown(storage) // cleanup your storage here

	gifts := []*gift.Gift{
		{
			Code:           "TEST1",
			GiftAmount:     100,
			UsageLimit:     10,
			UsedCount:      0,
			ExpirationDate: time.Now().AddDate(0, 0, 10),
			StartDateTime:  time.Now(),
		},
		{
			Code:           "TEST2",
			GiftAmount:     200,
			UsageLimit:     20,
			UsedCount:      0,
			ExpirationDate: time.Now().AddDate(0, 0, 20),
			StartDateTime:  time.Now(),
		},
	}

	err := storage.CreateBulk(gifts)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	for _, g := range gifts {
		err = storage.Delete(g.ID)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	}

	// Verify that all gifts were correctly inserted into the database
}

func TestUpdate(t *testing.T) {
	storage := setup() // setup your storage here
	//defer teardown(storage) // cleanup your storage here

	g := &gift.Gift{
		Code:           "TEST",
		GiftAmount:     100,
		UsageLimit:     10,
		UsedCount:      0,
		ExpirationDate: time.Now().AddDate(0, 0, 10),
		StartDateTime:  time.Now(),
	}

	err := storage.Create(g)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	g.GiftAmount = 200
	err = storage.Update(g)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	err = storage.Delete(g.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestDelete(t *testing.T) {
	storage := setup()
	//defer teardown(storage) // cleanup your storage here

	g := &gift.Gift{
		Code:           "TEST",
		GiftAmount:     100,
		UsageLimit:     10,
		UsedCount:      0,
		ExpirationDate: time.Now().AddDate(0, 0, 10),
		StartDateTime:  time.Now(),
	}

	err := storage.Create(g)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	err = storage.Delete(g.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestGetByCode(t *testing.T) {
	storage := setup()
	//defer teardown(storage) // cleanup your storage here

	g := &gift.Gift{
		Code:           "TEST",
		GiftAmount:     100,
		UsageLimit:     10,
		UsedCount:      0,
		ExpirationDate: time.Now().AddDate(0, 0, 10),
		StartDateTime:  time.Now(),
	}

	err := storage.Create(g)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	got, err := storage.GetByCode(g.Code)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if got.ID != g.ID {
		t.Fatalf("expected %v, got %v", g.ID, got.ID)
	}

	err = storage.Delete(g.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func setup() *gift.Storage {
	// Initialize your database connection
	db, err := sql.Open("postgres", "host=localhost port=5432 user=arv123 password=asd123ASD dbname=test sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	// Initialize your Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// Create a new instance of gift.Storage
	storage := gift.New(db, redisClient)

	return &storage
}
