package tests

import (
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/MamangRust/microservice-ecommerce-pkg/database"
)

// expectedDBName memetakan direktori service ke konstanta database bounded
// context yang WAJIB dipakai di cmd/main.go-nya. Tabel ini adalah kontrak:
// kalau ada main.go yang salah memetakan service ke database lain, test ini
// gagal — persis seperti pointofsale/payment-gateway (compile-time tidak bisa
// menangkap salah pilih konstanta).
var expectedDBName = map[string]string{
	"auth":              "IdentityDB",
	"user":              "IdentityDB",
	"role":              "IdentityDB",
	"category":          "CatalogDB",
	"product":           "CatalogDB",
	"merchant":          "MerchantDB",
	"merchant_award":    "MerchantDB",
	"merchant_business": "MerchantDB",
	"merchant_detail":   "MerchantDB",
	"merchant_policy":   "MerchantDB",
	"order":             "SalesDB",
	"order_item":        "SalesDB",
	"transaction":       "SalesDB",
	"cart":              "ExperienceDB",
	"shipping_address":  "ExperienceDB",
	"banner":            "ExperienceDB",
	"slider":            "ExperienceDB",
	"review":            "ExperienceDB",
	"review_detail":     "ExperienceDB",
}

// allowedDBConstants adalah himpunan konstanta database yang sah di
// pkg/database/names.go.
var allowedDBConstants = map[string]bool{
	"IdentityDB":   true,
	"CatalogDB":    true,
	"MerchantDB":   true,
	"SalesDB":      true,
	"ExperienceDB": true,
	"EmailDB":      true,
}

var (
	reDBName      = regexp.MustCompile(`DBName:\s*database\.(\w+)`)
	reServiceName = regexp.MustCompile(`ServiceName:\s*"([^"]+)"`)
	reDBCluster   = regexp.MustCompile(`DBCluster:`)
)

func findMainFiles(t *testing.T) []string {
	t.Helper()
	serviceDir := filepath.Join("..", "service")
	entries, err := os.ReadDir(serviceDir)
	if err != nil {
		t.Fatalf("gagal baca direktori service: %v", err)
	}
	var mains []string
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		// Hanya service bounded-context yang punya field DBName: database.<Context>DB
		// yang masuk kontrak ini. Service infrastruktur (apigateway, email, migrate,
		// seeder, stats_*) tidak menyetel DBName tunggal sehingga dikecualikan.
		if _, ok := expectedDBName[e.Name()]; !ok {
			continue
		}
		candidate := filepath.Join(serviceDir, e.Name(), "cmd", "main.go")
		if info, err := os.Stat(candidate); err == nil && !info.IsDir() {
			mains = append(mains, candidate)
		}
	}
	if len(mains) == 0 {
		t.Fatal("tidak ada service/*/cmd/main.go yang ditemukan")
	}
	return mains
}

func TestServiceDBMapping(t *testing.T) {
	mains := findMainFiles(t)
	for _, main := range mains {
		main := main
		t.Run(filepath.Base(filepath.Dir(filepath.Dir(main))), func(t *testing.T) {
			src, err := os.ReadFile(main)
			if err != nil {
				t.Fatalf("gagal baca %s: %v", main, err)
			}
			content := string(src)

			// DBCluster tidak lagi ada — harus diganti DBName.
			if reDBCluster.MatchString(content) {
				t.Errorf("%s masih memakai DBCluster; harus DBName", main)
			}

			m := reDBName.FindStringSubmatch(content)
			if m == nil {
				t.Fatalf("%s tidak menyetel DBName: database.<Context>DB", main)
			}
			got := m[1]
			if !allowedDBConstants[got] {
				t.Errorf("%s memakai konstanta database tidak dikenal: %s", main, got)
			}

			svc := filepath.Base(filepath.Dir(filepath.Dir(main)))
			want, ok := expectedDBName[svc]
			if !ok {
				t.Fatalf("service %q tidak ada di tabel expectedDBName", svc)
			}
			if got != want {
				t.Errorf("service %q memetakan ke %s, padahal seharusnya %s", svc, got, want)
			}

			if sm := reServiceName.FindStringSubmatch(content); sm == nil {
				t.Errorf("%s tidak menyetel ServiceName", main)
			}
		})
	}
}

func TestDatabaseConstantValues(t *testing.T) {
	// Konstanta harus persis bernilai nama database fisik yang dipakai script
	// migrasi (scripts/split-databases.js PLAN). Kalau salah satu diganti, test
	// ini gagal dan migrasi + indeks tidak akan cocok.
	cases := []struct {
		name string
		got  string
		want string
	}{
		{"IdentityDB", database.IdentityDB, "ec_identity"},
		{"CatalogDB", database.CatalogDB, "ec_catalog"},
		{"MerchantDB", database.MerchantDB, "ec_merchant"},
		{"SalesDB", database.SalesDB, "ec_sales"},
		{"ExperienceDB", database.ExperienceDB, "ec_experience"},
		{"EmailDB", database.EmailDB, "ec_email"},
	}
	for _, c := range cases {
		if c.got != c.want {
			t.Errorf("database.%s = %q, want %q", c.name, c.got, c.want)
		}
	}
}
