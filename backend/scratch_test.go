package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	dbURL := "postgresql://postgres.kjbpidoxpvwrnpurqftk:KlinikCepat123@aws-1-ap-southeast-1.pooler.supabase.com:6543/postgres"
	dbPool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		panic(err)
	}
	defer dbPool.Close()

	for i := 0; i < 5; i++ {
		var role string
		err := dbPool.QueryRow(context.Background(), "SELECT 'klinik_admin'").Scan(&role)
		if err != nil {
			fmt.Printf("Query error on iteration %d: %v\n", i, err)
		} else {
			fmt.Printf("Iteration %d success: %s\n", i, role)
		}
	}
}
