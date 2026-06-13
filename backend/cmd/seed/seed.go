package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

// Passwords used (plaintext -> hashed with bcrypt before insert):
//   virgil@purgatorio.com   -> "Password123!"
//   beatrice@purgatorio.com -> "Password123!"
//   dante@purgatorio.com    -> "Password123!"

func hash(pw string) string {
	b, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("bcrypt error: %v", err)
	}
	return string(b)
}

func main() {
	_ = godotenv.Load()

	dsn := os.Getenv("DB_URL")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5432/purgatorio?sslmode=disable"
	}

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("unable to connect: %v", err)
	}
	defer pool.Close()

	ctx := context.Background()

	tx, err := pool.Begin(ctx)
	if err != nil {
		log.Fatalf("begin tx: %v", err)
	}
	defer tx.Rollback(ctx)

	// ---------------------------------------------------------------
	// 1. Buildings
	// ---------------------------------------------------------------
	type buildingDef struct {
		id       string
		name     string
		size     int
		price    int
		currency string
		category string
	}

	buildings := []buildingDef{
		{"bastion", "Bastion", 1, 50, "penitence", "defense"},
		{"angel-spire", "Angel Spire", 2, 500, "penitence", "defense"},
		{"lament-basin", "Lament Basin", 2, 300, "penitence", "resource"},
		{"sanctum", "Sanctum", 3, 10000, "penitence", "other"},
	}

	for _, b := range buildings {
		_, err = tx.Exec(ctx, `
			INSERT INTO buildings (id, name, size, price, currency, category)
			VALUES ($1, $2, $3, $4, $5, $6)
			ON CONFLICT (id) DO NOTHING
		`, b.id, b.name, b.size, b.price, b.currency, b.category)
		if err != nil {
			log.Fatalf("insert building %s: %v", b.id, err)
		}
	}

	// ---------------------------------------------------------------
	// 2. Building limits (per terrace level)
	// ---------------------------------------------------------------
	type limitDef struct {
		buildingID string
		terrace    int
		maxAllowed int
	}

	limits := []limitDef{
		// Bastion: lots of walls available, scaling with terrace
		{"bastion", 1, 10},
		{"bastion", 2, 20},
		{"bastion", 3, 35},

		// Angel Spire: scarce defense towers
		{"angel-spire", 1, 1},
		{"angel-spire", 2, 2},
		{"angel-spire", 3, 3},

		// Lament Basin: resource generators
		{"lament-basin", 1, 1},
		{"lament-basin", 2, 2},
		{"lament-basin", 3, 3},

		// Sanctum: terrace progression — one per terrace level
		{"sanctum", 1, 1},
		{"sanctum", 2, 1},
		{"sanctum", 3, 1},
	}

	for _, l := range limits {
		_, err = tx.Exec(ctx, `
			INSERT INTO building_limits (building_id, terrace_level, max_allowed)
			VALUES ($1, $2, $3)
			ON CONFLICT (building_id, terrace_level) DO NOTHING
		`, l.buildingID, l.terrace, l.maxAllowed)
		if err != nil {
			log.Fatalf("insert building_limit %s/%d: %v", l.buildingID, l.terrace, err)
		}
	}

	// ---------------------------------------------------------------
	// 3. Building levels
	// ---------------------------------------------------------------
	type levelDef struct {
		buildingID    string
		level         int
		hp            *int
		dps           *int
		productionRate *int
		storageCap    *int
		attackRange   *float64
		upgradeCost   int
		upgradeTime   int
	}

	intp := func(v int) *int { return &v }
	floatp := func(v float64) *float64 { return &v }

	levels := []levelDef{
		// Bastion (wall): only HP scales, no damage/production/storage/range
		{"bastion", 1, intp(300), nil, nil, nil, floatp(0), 50, 60},
		{"bastion", 2, intp(500), nil, nil, nil, floatp(0), 150, 300},
		{"bastion", 3, intp(800), nil, nil, nil, floatp(0), 400, 900},
		{"bastion", 4, intp(1200), nil, nil, nil, floatp(0), 900, 1800},

		// Angel Spire: ranged defense attack_range 5.0
		{"angel-spire", 1, intp(450), intp(12), nil, nil, floatp(5.0), 500, 600},
		{"angel-spire", 2, intp(600), intp(18), nil, nil, floatp(5.0), 1200, 1800},
		{"angel-spire", 3, intp(800), intp(26), nil, nil, floatp(5.0), 2500, 3600},
		{"angel-spire", 4, intp(1050), intp(36), nil, nil, floatp(5.0), 5000, 7200},
		{"angel-spire", 5, intp(1400), intp(50), nil, nil, floatp(5.0), 9000, 14400},

		// Lament Basin: resource collector, production_rate in penitence/sec (int), no attack range
		{"lament-basin", 1, intp(400), nil, intp(2), intp(500), floatp(0), 300, 600},
		{"lament-basin", 2, intp(550), nil, intp(4), intp(1000), floatp(0), 800, 1800},
		{"lament-basin", 3, intp(700), nil, intp(6), intp(1800), floatp(0), 2000, 3600},
		{"lament-basin", 4, intp(900), nil, intp(8), intp(3000), floatp(0), 4500, 7200},
		{"lament-basin", 5, intp(1150), nil, intp(12), intp(5000), floatp(0), 9000, 14400},

		// Sanctum: terrace progression — minimal HP, no combat/production stats
		{"sanctum", 1, intp(500), nil, nil, nil, floatp(0), 5000, 3600},
		{"sanctum", 2, intp(600), nil, nil, nil, floatp(0), 12000, 7200},
		{"sanctum", 3, intp(700), nil, nil, nil, floatp(0), 25000, 14400},
		{"sanctum", 4, intp(800), nil, nil, nil, floatp(0), 50000, 28800},
		{"sanctum", 5, intp(900), nil, nil, nil, floatp(0), 100000, 57600},
	}

	for _, lv := range levels {
		_, err = tx.Exec(ctx, `
			INSERT INTO building_levels
				(building_id, level, hp, damage_per_second, production_rate, storage_capacity, attack_range, upgrade_cost, upgrade_time)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			ON CONFLICT (building_id, level) DO NOTHING
		`, lv.buildingID, lv.level, lv.hp, lv.dps, lv.productionRate, lv.storageCap, lv.attackRange, lv.upgradeCost, lv.upgradeTime)
		if err != nil {
			log.Fatalf("insert building_level %s/%d: %v", lv.buildingID, lv.level, err)
		}
	}

	// ---------------------------------------------------------------
	// 4. Troops
	// ---------------------------------------------------------------
	type troopDef struct {
		id              string
		name            string
		trainingCost    int
		space           int
		hp              int
		dps             int
		speed           float64
		attackRange     float64
		preferredTarget string
	}

	troops := []troopDef{
		{"stone-bearer", "Stone Bearer", 80, 6, 300, 12, 1.0, 1.0, "defense"},
		{"ashwalker", "Ashwalker", 60, 3, 80, 25, 2.0, 2.0, "defense"},
		{"hoarder", "Hoarder", 50, 4, 120, 18, 1.5, 1.5, "resource"},
		{"ravager", "Ravager", 90, 7, 200, 40, 1.2, 1.0, "defense"},
		{"coveter", "Coveter", 30, 2, 60, 15, 2.5, 3.0, "any"},
	}

	for _, t := range troops {
		_, err = tx.Exec(ctx, `
			INSERT INTO troops (id, name, training_cost, space, hp, dps, speed, attack_range, preferred_target)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			ON CONFLICT (id) DO NOTHING
		`, t.id, t.name, t.trainingCost, t.space, t.hp, t.dps, t.speed, t.attackRange, t.preferredTarget)
		if err != nil {
			log.Fatalf("insert troop %s: %v", t.id, err)
		}
	}

	// ---------------------------------------------------------------
	// 5. Users (auth + users + user_stats + user_economy + user_combat
	//    + game_state + user_army)
	// ---------------------------------------------------------------
	type userDef struct {
		email        string
		username     string
		password     string
		xp           int
		terraceLevel int
		penitence    int
		grace        int
		maxPenitence int
		sinMeter     int
		troops       map[string]int
		maxCapacity  int
	}

	users := []userDef{
		{
			email:        "virgil@purgatorio.com",
			username:     "virgil",
			password:     "Password123!",
			xp:           1200,
			terraceLevel: 3,
			penitence:    2500,
			grace:        120,
			maxPenitence: 8000,
			sinMeter:     10,
			troops: map[string]int{
				"stone-bearer": 4,
				"ashwalker":    6,
				"coveter":      8,
			},
			maxCapacity: 100,
		},
		{
			email:        "beatrice@purgatorio.com",
			username:     "beatrice",
			password:     "Password123!",
			xp:           300,
			terraceLevel: 2,
			penitence:    900,
			grace:        40,
			maxPenitence: 6000,
			sinMeter:     0,
			troops: map[string]int{
				"hoarder":  3,
				"ashwalker": 2,
			},
			maxCapacity: 50,
		},
		{
			email:        "dante@purgatorio.com",
			username:     "dante",
			password:     "Password123!",
			xp:           50,
			terraceLevel: 1,
			penitence:    500,
			grace:        50,
			maxPenitence: 5000,
			sinMeter:     35,
			troops: map[string]int{
				"stone-bearer": 1,
			},
			maxCapacity: 20,
		},
	}

	type insertedUser struct {
		id           string
		terraceLevel int
	}
	var insertedUsers []insertedUser

	for _, u := range users {
		var authID string
		err = tx.QueryRow(ctx, `
			INSERT INTO auth (email, password_hash)
			VALUES ($1, $2)
			RETURNING id
		`, u.email, hash(u.password)).Scan(&authID)
		if err != nil {
			log.Fatalf("insert auth %s: %v", u.email, err)
		}

		_, err = tx.Exec(ctx, `
			INSERT INTO users (id, username, xp, terrace_level)
			VALUES ($1, $2, $3, $4)
		`, authID, u.username, u.xp, u.terraceLevel)
		if err != nil {
			log.Fatalf("insert user %s: %v", u.username, err)
		}

		_, err = tx.Exec(ctx, `
			INSERT INTO user_stats (user_id)
			VALUES ($1)
		`, authID)
		if err != nil {
			log.Fatalf("insert user_stats %s: %v", u.username, err)
		}

		_, err = tx.Exec(ctx, `
			INSERT INTO user_economy (user_id, penitence, grace, max_penitence)
			VALUES ($1, $2, $3, $4)
		`, authID, u.penitence, u.grace, u.maxPenitence)
		if err != nil {
			log.Fatalf("insert user_economy %s: %v", u.username, err)
		}

		_, err = tx.Exec(ctx, `
			INSERT INTO user_combat (user_id, sin_meter)
			VALUES ($1, $2)
		`, authID, u.sinMeter)
		if err != nil {
			log.Fatalf("insert user_combat %s: %v", u.username, err)
		}

		_, err = tx.Exec(ctx, `
			INSERT INTO game_state (user_id, scene_id, fallback)
			VALUES ($1, $2, $3)
		`, authID, "base", "base")
		if err != nil {
			log.Fatalf("insert game_state %s: %v", u.username, err)
		}

		// user_army: build troops JSONB and used_capacity
		usedCapacity := 0
		var troopsJSON strings.Builder; troopsJSON.WriteString("{")
		first := true
		for _, t := range troops {
			count, ok := u.troops[t.id]
			if !ok {
				count = 0
			}
			usedCapacity += count * t.space
			if !first {
				troopsJSON.WriteString(",")
			}
			fmt.Fprintf(&troopsJSON, `"%s":%d`, t.id, count)
			first = false
		}
		troopsJSON.WriteString("}")

		_, err = tx.Exec(ctx, `
			INSERT INTO user_army (user_id, troops, used_capacity, max_capacity)
			VALUES ($1, $2::jsonb, $3, $4)
		`, authID, troopsJSON.String(), usedCapacity, u.maxCapacity)
		if err != nil {
			log.Fatalf("insert user_army %s: %v", u.username, err)
		}

		insertedUsers = append(insertedUsers, insertedUser{id: authID, terraceLevel: u.terraceLevel})

		// ---------------------------------------------------------------
		// user_buildings + base_layouts
		// ---------------------------------------------------------------
		// Each user gets:
		//   1 lament-basin (2x2), 1 angel-spire (2x2), and a ring of bastions (1x1)
		// scaled roughly by terrace level.
		numBastions := 6 + (u.terraceLevel-1)*4 // 6, 10, 14...

		// user_buildings counts
		_, err = tx.Exec(ctx, `
			INSERT INTO user_buildings (user_id, building_id, quantity)
			VALUES ($1, 'lament-basin', 1), ($1, 'angel-spire', 1), ($1, 'bastion', $2), ($1, 'sanctum', 1)
		`, authID, numBastions)
		if err != nil {
			log.Fatalf("insert user_buildings %s: %v", u.username, err)
		}

		// base_layouts on a 30x30 grid
		// Place lament-basin near top-left (resource), angel-spire near center (defense core)
		// Sanctum (3x3) in the top-right corner
		// Bastions form a perimeter ring around the angel-spire.

		basinX, basinY := 4, 4
		spireX, spireY := 14, 14
		sanctumX, sanctumY := 24, 5

		_, err = tx.Exec(ctx, `
			INSERT INTO base_layouts (user_id, building_id, x, y, metadata)
			VALUES ($1, 'lament-basin', $2, $3, '{}'::jsonb)
		`, authID, basinX, basinY)
		if err != nil {
			log.Fatalf("insert base_layout lament-basin %s: %v", u.username, err)
		}

		_, err = tx.Exec(ctx, `
			INSERT INTO base_layouts (user_id, building_id, x, y, metadata)
			VALUES ($1, 'angel-spire', $2, $3, '{}'::jsonb)
		`, authID, spireX, spireY)
		if err != nil {
			log.Fatalf("insert base_layout angel-spire %s: %v", u.username, err)
		}

		_, err = tx.Exec(ctx, `
			INSERT INTO base_layouts (user_id, building_id, x, y, metadata)
			VALUES ($1, 'sanctum', $2, $3, '{}'::jsonb)
		`, authID, sanctumX, sanctumY)
		if err != nil {
			log.Fatalf("insert base_layout sanctum %s: %v", u.username, err)
		}

		// Bastion ring around the angel-spire (which occupies spireX..spireX+1, spireY..spireY+1)
		// Build a perimeter at offset -1 and +2 from the spire's footprint, taking only `numBastions` cells
		ringOffsets := [][2]int{
			{-1, -1}, {0, -1}, {1, -1}, {2, -1},
			{2, 0}, {2, 1}, {2, 2},
			{1, 2}, {0, 2}, {-1, 2},
			{-1, 1}, {-1, 0},
		}

		placed := 0
		for _, off := range ringOffsets {
			if placed >= numBastions {
				break
			}
			bx := spireX + off[0]
			by := spireY + off[1]
			// keep within 30x30 grid bounds (0-29)
			if bx < 0 || by < 0 || bx > 29 || by > 29 {
				continue
			}
			_, err = tx.Exec(ctx, `
				INSERT INTO base_layouts (user_id, building_id, x, y, metadata)
				VALUES ($1, 'bastion', $2, $3, '{}'::jsonb)
			`, authID, bx, by)
			if err != nil {
				log.Fatalf("insert base_layout bastion %s (%d,%d): %v", u.username, bx, by, err)
			}
			placed++
		}

		if placed < numBastions {
			log.Printf("warning: only placed %d/%d bastions for %s (ring exhausted)", placed, numBastions, u.username)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		log.Fatalf("commit tx: %v", err)
	}

	fmt.Println("Seed completed successfully.")
	fmt.Println()
	fmt.Println("Seeded users (email / password):")
	for _, u := range users {
		fmt.Printf("  %-28s / %s\n", u.email, u.password)
	}
}
