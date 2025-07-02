package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Menus struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	NamaMenu  string             `bson:"nama_menu" json:"nama_menu"`
	Harga     float64            `bson:"harga" json:"harga"`
	Deskripsi string             `bson:"deskripsi" json:"deskripsi"`
	Gambar    string             `bson:"gambar" json:"gambar"`
	Kategori  string             `bson:"kategori" json:"kategori"`
}

type Pesanan struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	NamaPelanggan  string             `bson:"nama_pelanggan" json:"nama_pelanggan"`
	NomorMeja      int                `bson:"nomor_meja" json:"nomor_meja"`
	DaftarMenu     []ItemPesanan      `bson:"daftar_menu" json:"daftar_menu"`
	TotalHarga     float64            `bson:"total_harga" json:"total_harga"`
	StatusPesanan  string             `bson:"status_pesanan" json:"status_pesanan"`
	TanggalPesanan primitive.DateTime `bson:"tanggal_pesanan" json:"tanggal_pesanan"`
	Pembayaran     string             `bson:"pembayaran" json:"pembayaran"`
	CatatanPesanan string             `bson:"catatan_pesanan" json:"catatan_pesanan"`
}

type ItemPesanan struct {
	MenuID      primitive.ObjectID `bson:"menu_id" json:"menu_id"`
	NamaMenu    string             `bson:"nama_menu" json:"nama_menu"`
	Jumlah      int                `bson:"jumlah" json:"jumlah"`
	HargaSatuan float64            `bson:"harga_satuan" json:"harga_satuan"`
	Subtotal    float64            `bson:"subtotal" json:"subtotal"`
}

type Admins struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Username string             `bson:"username,omitempty" json:"username,omitempty"`
	Password string             `bson:"password,omitempty" json:"password,omitempty"`
	Role     string             `bson:"role,omitempty" json:"role,omitempty"`
}

type Tokens struct {
	ID        string    `bson:"_id,omitempty" json:"_id,omitempty"`
	Token     string    `bson:"token,omitempty" json:"token,omitempty"`
	AdminID   string    `bson:"admin_id,omitempty" json:"admin_id,omitempty"`
	CreatedAt time.Time `bson:"created_at,omitempty" json:"created_at,omitempty"`
}

type LoginLogs struct {
	ID        string    `bson:"_id,omitempty" json:"_id,omitempty"`
	Username  string    `bson:"username,omitempty" json:"username,omitempty"`
	Timestamp time.Time `bson:"timestamp,omitempty" json:"timestamp,omitempty"`
	Activity  string    `bson:"activity,omitempry" json:"activity,omitempty"`
}
