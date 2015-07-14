package config

import (
    "testing"
)

const dbFilePath = "test_data/db.json"

func TestDbConfig(t *testing.T) {
    InitDatabase(dbFilePath)

    testsDbCon := DbConnectionString
    testsDbName := DbName

    switch {
    case len(testsDbCon) == 0:
        t.Fatal("Database connection string not loaded from config!")
    case len(testsDbName) == 0:
        t.Fatal("Database name could not be loaded from config!")
    }
}
