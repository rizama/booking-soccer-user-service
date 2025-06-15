package seeders

import "gorm.io/gorm"

/*
Mengapa dimulai dari interface:
- Interface mendefinisikan kontrak atau blueprint dari apa yang harus dilakukan
- Memberikan gambaran high-level tentang fungsionalitas sistem
- Menunjukkan public API yang tersedia untuk digunakan
- Dalam hal ini: sistem seeder harus memiliki method Run() untuk menjalankan seeding
*/
type ISeederRegistry interface {
	Run()
}

/*
Mengapa setelah interface:
- Struct adalah implementasi konkret dari interface
- Menunjukkan data/state apa yang diperlukan untuk implementasi
- Dalam hal ini: butuh koneksi database ( *gorm.DB ) untuk menjalankan seeding
- Memberikan gambaran dependencies yang diperlukan
*/
type Registry struct {
	db *gorm.DB
}

// constructor function
/*
Mengapa setelah struct:
- Menunjukkan cara membuat instance dari struct
- Memperlihatkan dependency injection - bagaimana dependencies dimasukkan
- Menunjukkan return type adalah interface (loose coupling)
- Memberikan gambaran lifecycle management
*/
func NewSeederRegistry(db *gorm.DB) ISeederRegistry {
	return &Registry{db: db}
}

/*
Mengapa terakhir:
- Ini adalah implementasi detail dari kontrak interface
- Menunjukkan business logic sebenarnya
- Memperlihatkan orchestration - bagaimana komponen-komponen dikoordinasikan
- Menunjukkan urutan eksekusi dan dependencies antar seeder
*/
func (r *Registry) Run() {
	RunRoleSeeder(r.db)
	RunUserSeeder(r.db)
}
