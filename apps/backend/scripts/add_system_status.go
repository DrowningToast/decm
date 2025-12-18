package main

import (
	"context"
	"decm-database/go/generated"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"apps/backend/common/pgclient"
	"apps/backend/core-api/config"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env if it exists
	_ = godotenv.Load()

	startTimeUnix := flag.Int64("start", 0, "Start time in Unix timestamp (seconds)")
	endTimeUnix := flag.Int64("end", 0, "Planned end time in Unix timestamp (seconds), optional")
	status := flag.Int("status", -1, "Status: 0 for maintenance, 1 for operating")
	isPlanned := flag.Bool("planned", false, "Is this a planned status change")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  go run scripts/add_system_status.go -status [0|1] [options]\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	// Validation
	if *status != 0 && *status != 1 {
		fmt.Fprintln(os.Stderr, "Error: status is required and must be 0 (maintenance) or 1 (operating)")
		flag.Usage()
		os.Exit(1)
	}

	if *startTimeUnix == 0 {
		*startTimeUnix = time.Now().Unix()
		fmt.Printf("Info: start time not provided, using current time: %d\n", *startTimeUnix)
	}

	ctx := context.Background()
	cfg := config.LoadConfig()

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Configuration validation failed: %v", err)
	}

	pgConn, err := pgclient.NewPool(ctx, &cfg.Postgres)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pgConn.Close()

	queries := generated.New(pgConn)

	var plannedEndTime pgtype.Timestamptz
	if *endTimeUnix > 0 {
		plannedEndTime = pgtype.Timestamptz{
			Time:  time.Unix(*endTimeUnix, 0),
			Valid: true,
		}
	} else {
		plannedEndTime = pgtype.Timestamptz{Valid: false}
	}

	result, err := queries.CreateSystemStatusSchedule(ctx, generated.CreateSystemStatusScheduleParams{
		StartTime:      time.Unix(*startTimeUnix, 0),
		PlannedEndTime: plannedEndTime,
		Status:         int32(*status),
		IsPlanned:      *isPlanned,
	})
	if err != nil {
		log.Fatalf("failed to create system status schedule: %v", err)
	}

	fmt.Printf("\nSuccessfully created system status schedule:\n")
	fmt.Printf("-------------------------------------------\n")
	fmt.Printf("ID:               %d\n", result.ID)
	fmt.Printf("Order ID:         %d\n", result.OrderID)
	fmt.Printf("Status:           %d (%s)\n", result.Status, map[int32]string{0: "Maintenance", 1: "Operating"}[result.Status])
	fmt.Printf("Start Time:       %v (%d)\n", result.StartTime, result.StartTime.Unix())
	if result.PlannedEndTime.Valid {
		fmt.Printf("Planned End Time: %v (%d)\n", result.PlannedEndTime.Time, result.PlannedEndTime.Time.Unix())
	} else {
		fmt.Printf("Planned End Time: <not set>\n")
	}
	fmt.Printf("Is Planned:       %v\n", result.IsPlanned)
	fmt.Printf("Created At:       %v\n", result.CreatedAt)
}
