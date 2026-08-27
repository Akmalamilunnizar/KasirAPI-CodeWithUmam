Frontend (Port 5173)
cd e:\Projects\KasirAPI-CodeWithUmam\kasir-fe
npm run dev

Backend (Port 8080)
cd e:\Projects\KasirAPI-CodeWithUmam\kasir-be
go run main.go

http://localhost:5173/.

Database Migration
Jika database belum ada atau ingin me-reset data awal:

1. Masuk ke direktori backend
cd e:\Projects\KasirAPI-CodeWithUmam\kasir-be

2.jalankan Migrate
2.1 Generate Migration (Jika ada perubahan model)
Ini akan membuat file migration baru jika Anda mengubah models.go

powershell
go run cmd/migrate/main.go generate
2.2 Run Migration (Membuat Tabel)
Ini akan menjalankan file 000001_init_schema.up.sql

powershell
go run cmd/migrate/main.go migrate
2.3 Rollback (Menghapus Tabel)
Jika ingin mengembalikan ke kondisi sebelum migration (menghapus tabel):

powershell
go run cmd/migrate/main.go migrate down 1