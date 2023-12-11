package main

import (
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/bxcodec/faker/v3"
	_ "github.com/lib/pq" // Import driver PostgreSQL
	"gopkg.in/yaml.v2"
)

// Config struct represents the structure of the database configuration
type Config struct {
	Development struct {
		Dialect    string `yaml:"dialect"`
		Datasource string `yaml:"datasource"`
		Dir        string `yaml:"dir"`
		Table      string `yaml:"table"`
	} `yaml:"development"`
}

// Partner struct represents the structure of the "partners" table
type Partner struct {
	ID               int       `faker:"-"` // Ignore during seeding
	PartnerName      string    `faker:"name"`
	ContactPerson    string    `faker:"name"`
	Email            string    `faker:"email"`
	PhoneNumber      string    `faker:"phone_number"`
	Address          string    `faker:"sentence"`
	RegistrationDate time.Time `faker:"-"`
}

type Customer struct {
	ID               int       `faker:"-"`
	CustomerName     string    `faker:"name"`
	Email            string    `faker:"email"`
	PhoneNumber      string    `faker:"phone_number"`
	Address          string    `faker:"sentence"`
	RegistrationDate time.Time `faker:"-"`
	PartnerID        int       `faker:"oneof: 1, 2, 3, 4, 5, 6, 7, 8, 9, 10"`
}

func readConfigFile(filename string) (Config, error) {
	var config Config

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}

func main() {
	var configFile, tableName string
	var total int
	flag.StringVar(&configFile, "configFile", "dbconfig.yml", "Path to the config file")
	flag.StringVar(&tableName, "tableName", "customers", "name of table")
	flag.IntVar(&total, "total", 10, "total of seeders")
	flag.Parse()

	// Baca konfigurasi dari file
	config, err := readConfigFile(configFile)
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("postgres", config.Development.Datasource)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	if tableName == "partners" {
		seedPartnerTable(db, total)
	} else {
		seedCustomerTable(db, total)
	}
}

func seedPartnerTable(db *sql.DB, numRecords int) {
	for i := 0; i < numRecords; i++ {
		// Generate fake data for partner
		partner := Partner{}
		err := faker.FakeData(&partner)
		if err != nil {
			log.Fatal(err)
		}

		// Insert fake data into the "partners" table
		_, err = db.Exec("INSERT INTO partners (partner_name, contact_person, email, phone_number, address, registration_date) VALUES ($1, $2, $3, $4, $5, $6)",
			partner.PartnerName, partner.ContactPerson, partner.Email, partner.PhoneNumber, partner.Address, time.Now())
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Printf("%d records seeded successfully.\n", numRecords)
}

func seedCustomerTable(db *sql.DB, numRecords int) {
	startTime := time.Now()

	for i := 0; i < numRecords; i++ {
		// Generate fake data for customer
		customer := Customer{}
		err := faker.FakeData(&customer)
		if err != nil {
			log.Fatal(err)
		}

		// Insert fake data into the "customers" table
		_, err = db.Exec("INSERT INTO customers (customer_name, email, phone_number, address, registration_date, partner_id) VALUES ($1, $2, $3, $4, $5, $6)",
			customer.CustomerName, customer.Email, customer.PhoneNumber, customer.Address, time.Now(), customer.PartnerID)
		if err != nil {
			log.Fatal(err)
		} else {
			num := i + 1
			if num%25 == 0 {
				// Menghitung persentase
				percentage := calculatePercentage(num, numRecords)

				// Menampilkan hasil
				fmt.Printf("%d dari %d = %.2f%%\n", num, numRecords, percentage)
			}
		}
	}

	fmt.Printf("%d records seeded successfully with duration %f seconds.\n", numRecords, time.Since(startTime).Seconds())
}

func calculatePercentage(value, total int) float64 {
	if total == 0 {
		return 0.0
	}
	return float64(value) / float64(total) * 100
}
