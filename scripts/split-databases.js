// scripts/split-databases.js
//
// Memindahkan data dari satu database lama (default: ecommerce) ke enam database
// bounded context: ec_identity, ec_catalog, ec_merchant, ec_sales,
// ec_experience, ec_email.
//
// Sifat:
//   - IDEMPOTEN — memakai bulkWrite replaceOne + upsert per _id, jadi aman
//     dijalankan berulang. Tidak ada dokumen yang dihapus.
//   - DRY-RUN secara default — hanya mencetak rencana. Set APPLY=1 untuk jalan.
//   - TIDAK menghapus database sumber. Setelah verifikasi hijau, drop manual.
//
// Menyalin `counters` juga penting: ID di-generate lewat counter, jadi kalau
// counter tidak ikut dipindah, insert berikutnya akan menabrak _id yang sudah ada.
//
// Cara menjalankan (mongosh tidak perlu terpasang di host):
//
//   # 1. Lihat rencana (dry-run)
//   docker run --rm -v "$PWD/scripts/split-databases.js:/split.js:ro" \
//     mongo:7 mongosh "mongodb://localhost:27017/?replicaSet=rs0" --file /split.js
//
//   # 2. Jalankan
//   docker run --rm -e APPLY=1 -v "$PWD/scripts/split-databases.js:/split.js:ro" \
//     mongo:7 mongosh "mongodb://localhost:27017/?replicaSet=rs0" --file /split.js
//
// Kalau service berjalan di host, pakai host MongoDB (.env: 127.0.0.1:27017).
// Kalau lewat compose, jalankan dengan network compose dan host `mongodb:27017`.
//
// Konfigurasi lewat env: SRC_DB (default ecommerce), APPLY=1.

const SRC_DB = process.env.SRC_DB || "ecommerce";
const APPLY = process.env.APPLY === "1";

// Pemetaan mengikuti pkg/database/names.go + index plan di pkg/database/migrate.go.
const PLAN = {
  ec_identity: ["users", "roles", "user_roles", "refresh_tokens", "reset_tokens"],
  ec_catalog: ["categories", "products", "product_stock_adjustments"],
  ec_merchant: [
    "merchants", "merchant_documents", "merchant_details",
    "merchant_social_media_links", "merchant_business_information",
    "merchant_policies", "merchant_certifications_and_awards",
  ],
  ec_sales: ["orders", "order_items", "order_stock_reservations", "transactions"],
  ec_experience: ["banners", "sliders", "shipping_addresses", "carts", "reviews", "review_details"],
  ec_email: ["outbox_events", "consumer_inbox"],
};

// Koleksi yang punya counter (lihat index plan / GetNextID).
const COUNTER_COLLECTIONS = [
  "users", "roles", "user_roles",
  "merchants", "merchant_documents", "merchant_details",
  "merchant_social_media_links", "merchant_business_information",
  "merchant_policies", "merchant_certifications_and_awards",
  "products", "orders", "order_items", "order_stock_reservations",
  "refresh_tokens", "reset_tokens",
  "banners", "sliders", "categories",
  "shipping_addresses", "carts", "reviews", "review_details",
  "transactions", "product_stock_adjustments",
];

const src = db.getSiblingDB(SRC_DB);

print("=== split-databases ===");
print("sumber : " + SRC_DB);
print("mode   : " + (APPLY ? "APPLY (menulis)" : "DRY-RUN (tidak menulis; set APPLY=1)"));
print("");

let totalCopied = 0;
let totalSkipped = 0;
let failures = [];

for (const targetName of Object.keys(PLAN)) {
  const target = db.getSiblingDB(targetName);
  print("--- " + targetName);

  for (const coll of PLAN[targetName]) {
    const srcCount = src.getCollection(coll).countDocuments({});
    const dstBefore = target.getCollection(coll).countDocuments({});

    if (!APPLY) {
      print("  " + coll + ": " + srcCount + " dokumen di sumber, " + dstBefore + " sudah ada di target");
      continue;
    }

    let copied = 0;
    let skipped = 0;
    let writeErrors = 0;

    if (srcCount > 0) {
      const docs = src.getCollection(coll).find({}).toArray();
      const ops = docs.map(function (doc) {
        return { replaceOne: { filter: { _id: doc._id }, replacement: doc, upsert: true } };
      });

      // ordered:false -> satu dokumen bermasalah tidak menggagalkan sisanya.
      //
      // PENTING: kalau ada write error, mongosh MELEMPAR MongoBulkWriteError —
      // bukan mengembalikan hasil. Tanpa try/catch, script berhenti di tengah,
      // koleksi berikutnya tidak diproses, dan tidak ada ringkasan sama sekali.
      let res = {};
      try {
        res = target.getCollection(coll).bulkWrite(ops, { ordered: false });
      } catch (e) {
        // hasil parsial menempel di e.result; hitungannya tetap kita pakai.
        res = (e && e.result) ? e.result : {};
        writeErrors = (e && e.writeErrors && e.writeErrors.length) ? e.writeErrors.length : 1;
        failures.push(targetName + "." + coll + ": " + writeErrors +
          " write error — " + ((e && e.message) || String(e)));
      }

      copied = (res.upsertedCount || 0) + (res.modifiedCount || 0);
      skipped = (res.matchedCount || 0) - (res.modifiedCount || 0);
    }

    const dstAfter = target.getCollection(coll).countDocuments({});
    totalCopied += copied;
    totalSkipped += skipped;

    print("  " + coll + ": sumber=" + srcCount + " target=" + dstAfter +
      " (upsert=" + copied + ", sudah sama=" + skipped + ")" +
      (writeErrors > 0 ? "  <-- " + writeErrors + " write error" : ""));

    if (dstAfter < srcCount) {
      failures.push(targetName + "." + coll + ": target (" + dstAfter + ") < sumber (" + srcCount + ")");
    }
  }
}

if (APPLY) {
  print("");
  print("--- counters");
  for (const targetName of Object.keys(PLAN)) {
    const target = db.getSiblingDB(targetName);
    for (const coll of PLAN[targetName]) {
      if (COUNTER_COLLECTIONS.indexOf(coll) === -1) continue;

      const counter = src.getCollection("counters").findOne({ _id: coll });
      if (!counter) {
        print("  " + targetName + "." + coll + ": tidak ada counter di sumber (GetNextID akan mulai dari 1)");
        continue;
      }

      try {
        target.getCollection("counters").updateOne(
          { _id: coll },
          { $max: { seq: counter.seq } },
          { upsert: true }
        );
        print("  " + targetName + "." + coll + ": counter seq=" + counter.seq);
      } catch (e) {
        failures.push(targetName + ".counters[" + coll + "]: gagal tulis counter — " + ((e && e.message) || String(e)));
        print("  " + targetName + "." + coll + ": GAGAL tulis counter — " + ((e && e.message) || String(e)));
      }
    }
  }
}

print("");
print("=== ringkasan ===");
if (!APPLY) {
  print("DRY-RUN selesai. Tidak ada yang ditulis. Jalankan ulang dengan APPLY=1.");
} else {
  print("dokumen di-upsert: " + totalCopied + ", sudah identik: " + totalSkipped);
  if (failures.length > 0) {
    print("PERINGATAN — " + failures.length + " masalah:");
    failures.forEach(function (f) { print("  - " + f); });
  } else {
    print("verifikasi jumlah dokumen: OK untuk semua koleksi");
  }
  print("");
  print("Database sumber `" + SRC_DB + "` TIDAK dihapus. Setelah yakin, drop manual:");
  print("  mongosh \"<uri>\" --eval 'db.getSiblingDB(\"" + SRC_DB + "\").dropDatabase()'");

  // Kode keluar non-nol supaya CI/pipeline berhenti, bukan lanjut dengan data
  // yang belum lengkap.
  if (failures.length > 0) {
    quit(1);
  }
}
