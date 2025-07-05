package utils

import "github.com/madhur1608/aditis-kitchen/customer-service/db"

func EnsureDBConnection() {
    if db.Cluster == nil {
        db.InitCouchbase()
    }
}
