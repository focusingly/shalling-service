package util_test

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"io"
	"testing"
	"time"
)

type Address struct {
	Street  string `json:"street"`
	City    string `json:"city"`
	State   string `json:"state"`
	ZipCode string `json:"zip_code"`
}

type ContactInfo struct {
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Website string `json:"website"`
}

type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type Order struct {
	ID        int        `json:"id"`
	Customer  Customer   `json:"customer"`
	Products  []Product  `json:"products"`
	Amount    float64    `json:"amount"`
	PlacedAt  time.Time  `json:"placed_at"`
	ShippedAt *time.Time `json:"shipped_at"`
}

type Customer struct {
	ID        int         `json:"id"`
	FirstName string      `json:"first_name"`
	LastName  string      `json:"last_name"`
	Age       int         `json:"age"`
	Addresses []Address   `json:"addresses"`
	Contact   ContactInfo `json:"contact"`
	Orders    []Order     `json:"orders"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

type Company struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	Location    string     `json:"location"`
	Established time.Time  `json:"established"`
	Employees   []Customer `json:"employees"`
	Revenue     float64    `json:"revenue"`
}

var now = time.Now()
var company = Company{
	ID:          1,
	Name:        "Tech Innovators Inc.",
	Location:    "Silicon Valley, CA",
	Established: now.AddDate(-10, 0, 0), // 成立于10年前
	Employees: []Customer{
		{
			ID:        1,
			FirstName: "Alice",
			LastName:  "Johnson",
			Age:       30,
			Addresses: []Address{
				{"123 Main St", "San Francisco", "CA", "94105"},
			},
			Contact: ContactInfo{
				Email:   "alice.johnson@techinnovators.com",
				Phone:   "+1 555-1234",
				Website: "alicejohnson.com",
			},
			Orders: []Order{
				{
					ID:       101,
					Customer: Customer{FirstName: "Alice"},
					Products: []Product{
						{ID: 1, Name: "Laptop", Price: 1200.99},
						{ID: 2, Name: "Mouse", Price: 25.50},
					},
					Amount:   1226.49,
					PlacedAt: now.AddDate(0, -1, 0),
				},
			},
			CreatedAt: now.AddDate(-5, 0, 0), // 5年前加入
			UpdatedAt: now.AddDate(-1, 0, 0),
		},
		{
			ID:        2,
			FirstName: "Bob",
			LastName:  "Smith",
			Age:       25,
			Addresses: []Address{
				{"456 Elm St", "San Jose", "CA", "95112"},
			},
			Contact: ContactInfo{
				Email:   "bob.smith@techinnovators.com",
				Phone:   "+1 555-5678",
				Website: "bobsmith.com",
			},
			Orders: []Order{
				{
					ID:       102,
					Customer: Customer{FirstName: "Bob"},
					Products: []Product{
						{ID: 3, Name: "Smartphone", Price: 799.99},
						{ID: 4, Name: "Headphones", Price: 199.99},
					},
					Amount:   999.98,
					PlacedAt: now.AddDate(0, -2, 0),
				},
			},
			CreatedAt: now.AddDate(-2, 0, 0),
			UpdatedAt: now.AddDate(-1, 0, 0),
		},
	},
	Revenue: 1000000.50,
}

func BenchmarkGobSerialized(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var bf bytes.Buffer
		encoder := gob.NewEncoder(&bf)
		encoder.Encode(&company)
		io.ReadAll(&bf)
	}

}

func BenchmarkJsonSerialized(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var bf bytes.Buffer
		encoder := json.NewEncoder(&bf)
		encoder.Encode(&company)
		io.ReadAll(&bf)
	}
}
