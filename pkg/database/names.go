package database

// Nama database fisik per bounded context. Satu konstanta = satu database;
// nama sengaja TIDAK diturunkan dari nama service, supaya mengganti nama
// service tidak diam-diam menunjuk database lain.
//
// Sebelumnya setiap service mendeklarasikan Config.DBCluster (mis. "DB_ORDER"),
// tetapi server.go mengabaikannya dan selalu memakai prefix "MONGO" — jadi
// seluruh service menulis ke satu database `ecommerce` tanpa error apa pun.
// Sekarang service memilih database lewat Config.DBName, yang diisi konstanta
// di bawah dan divalidasi gagal-keras saat start.
//
// Konstanta ini tidak mengandung informasi topologi deployment: kalau suatu saat
// tiap context dipisah ke server Mongo sendiri, yang berubah hanya koneksi
// (pkg/database/mongodb.go), bukan daftar ini.
const (
	IdentityDB   = "ec_identity"
	CatalogDB    = "ec_catalog"
	MerchantDB   = "ec_merchant"
	SalesDB      = "ec_sales"
	ExperienceDB = "ec_experience"
	EmailDB      = "ec_email"
)

// ContextPrefix returns the environment-variable prefix that points at the
// MongoDB server dedicated to the given bounded context.
//
// Enam bounded context sekarang punya server MongoDB-nya sendiri
// ("1 database = 1 MongoDB"). Prefix ini adalah prefix env penuh yang dipakai
// NewMongoClientWithPrefix, jadi nilainya "MONGO_<CONTEXT>" dan key yang dibaca
// adalah MONGO_<CONTEXT>_URI / MONGO_<CONTEXT>_HOST / MONGO_<CONTEXT>_PORT.
//
// Nama yang tidak dikenal mengembalikan "MONGO" supaya pemanggil jatuh ke
// instance tunggal lama (MONGO_URI / MONGO_HOST).
func ContextPrefix(databaseName string) string {
	switch databaseName {
	case IdentityDB:
		return "MONGO_IDENTITY"
	case CatalogDB:
		return "MONGO_CATALOG"
	case MerchantDB:
		return "MONGO_MERCHANT"
	case SalesDB:
		return "MONGO_SALES"
	case ExperienceDB:
		return "MONGO_EXPERIENCE"
	case EmailDB:
		return "MONGO_EMAIL"
	default:
		return "MONGO"
	}
}
