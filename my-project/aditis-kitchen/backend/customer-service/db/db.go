package db

import (
	"log"
	"time"

	"github.com/couchbase/gocb/v2"
	"github.com/madhur1608/aditis-kitchen/customer-service/models"
)

var (
	Cluster *gocb.Cluster
	Bucket  *gocb.Bucket
)

func InitCouchbase() {
	var err error
	Cluster, err = gocb.Connect("couchbase://localhost", gocb.ClusterOptions{
		Username: "Administrator",
		Password: "Aditi17",
	})
	if err != nil {
		log.Fatalf("❌ Couchbase connection failed: %v", err)
	}

	Bucket = Cluster.Bucket("aditiskitchen")
	err = Bucket.WaitUntilReady(5*time.Second, nil)
	if err != nil {
		log.Fatalf("❌ Bucket not ready: %v", err)
	}

	log.Println("✅ Connected to Couchbase bucket: customers")
}

func CreateCustomer(customer models.Customer) error {
	collection := Bucket.DefaultCollection()
	_, err := collection.Insert(customer.ID, customer, nil)
	return err
}

func GetCustomerByMobile(mobile int64) (models.Customer, error) {
	var customer models.Customer
	query := "SELECT c.* FROM `customers` c WHERE c.mobile = $1 LIMIT 1;"
	result, err := Cluster.Query(query, &gocb.QueryOptions{
		PositionalParameters: []interface{}{mobile},
	})
	if err != nil {
		return customer, err
	}

	for result.Next() {
		err := result.Row(&customer)
		return customer, err
	}

	return customer, gocb.ErrDocumentNotFound
}
