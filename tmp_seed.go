package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

type spare struct {
	name string
	sku  string
	min  int
}

func main() {
	cfg := mysql.NewConfig()
	cfg.User = "uprogaswmsdev"
	cfg.Passwd = "4nP2gVMFFIAYHv9v"
	cfg.Net = "tcp"
	cfg.Addr = "127.0.0.1:3307"
	cfg.DBName = "dbprogaswmsdev"
	cfg.ParseTime = true
	cfg.Loc = time.Local

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil { log.Fatal(err) }
	defer db.Close()

	spares := []spare{
		{"Morris Regulator", "MORRIS-REGULATOR", 5},
		{"Cigweld Regulator", "CIGWELD-REGULATOR", 5},
		{"Victor Regulator", "VICTOR-REGULATOR", 5},
		{"Harris Regulator", "HARRIS-REGULATOR", 5},
		{"Liquid Nitrogen Container", "LIQUID-NITROGEN-CONTAINER", 2},
		{"Vessel Gas Liquid (VGL)", "VESSEL-GAS-LIQUID-VGL", 2},
		{"Storage Tank", "STORAGE-TANK", 2},
		{"Fire Extinguisher", "FIRE-EXTINGUISHER", 10},
		{"Medical Gas Digital Alarm", "MEDICAL-GAS-DIGITAL-ALARM", 5},
		{"Zone Valve Box", "ZONE-VALVE-BOX", 5},
		{"Bedhead Aluminium Panel", "BEDHEAD-ALUMINIUM-PANEL", 5},
		{"Aluminum Hospital Handrail", "ALUMINUM-HOSPITAL-HANDRAIL", 10},
		{"Hospital Vinyl Floor", "HOSPITAL-VINYL-FLOOR", 50},
	}

	inserted := 0
	skipped := 0
	for _, s := range spares {
		var existing string
		err := db.QueryRow("SELECT id FROM master_item WHERE sku = ?", s.sku).Scan(&existing)
		if err == nil {
			fmt.Printf("SKIP %s (already exists id=%s)\n", s.sku, existing)
			skipped++
			continue
		}
		if err != sql.ErrNoRows {
			log.Fatalf("check sku %s: %v", s.sku, err)
		}

		u, err := uuid.NewV7()
		if err != nil {
			log.Fatalf("uuid: %v", err)
		}
		id := u.String()
		now := time.Now()

		_, err = db.Exec(
			"INSERT INTO master_item (id, created_at, updated_at, deleted_at, created_by, updated_by, deleted_by, name, sku, gas_type, is_serialized, empty_weight_kg, gas_weight_kg, min_stock_alert) VALUES (?, ?, ?, NULL, '', '', NULL, ?, ?, '', 0, NULL, NULL, ?)",
			id, now, now, s.name, s.sku, s.min,
		)
		if err != nil {
			log.Fatalf("insert master_item %s: %v", s.sku, err)
		}

		_, err = db.Exec(
			"INSERT INTO sparepart_stock (id, created_at, updated_at, deleted_at, created_by, updated_by, deleted_by, item_id, quantity) VALUES (?, ?, ?, NULL, '', '', NULL, ?, 0)",
			func() string { u, _ := uuid.NewV7(); return u.String() }(), now, now, id,
		)
		if err != nil {
			log.Fatalf("insert sparepart_stock %s: %v", s.sku, err)
		}

		fmt.Printf("INSERTED %s (%s) min=%d\n", s.sku, s.name, s.min)
		inserted++
	}

	fmt.Printf("\nDone: inserted=%d skipped=%d\n", inserted, skipped)

	// verify
	rows, err := db.Query("SELECT m.name, m.sku, s.quantity FROM master_item m LEFT JOIN sparepart_stock s ON s.item_id = m.id WHERE m.is_serialized = 0 ORDER BY m.name ASC")
	if err != nil { log.Fatal(err) }
	defer rows.Close()
	for rows.Next() {
		var name, sku string
		var qty sql.NullInt64
		rows.Scan(&name, &sku, &qty)
		fmt.Printf("VERIFY %s | %s | qty=%v\n", sku, name, qty.Int64)
	}
}
