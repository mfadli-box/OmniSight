# User Guide SITE - obx_site Frontend

## Ringkasan

obx_site adalah antarmuka web utama untuk login, navigasi board, dan pengelolaan data master yang terhubung ke backend obx_rest.

## Prasyarat

- Backend obx_rest berjalan dan dapat diakses melalui proxy frontend
- Session login aktif
- Token dan company context valid
- Browser mendukung cookie dan localStorage

## Menjalankan Frontend

Dari folder obx_site:

```sh
npm install
npm run dev
```

## Alur Masuk

1. Buka halaman `/login`.
2. Pilih company jika diperlukan.
3. Masukkan username, password, dan captcha.
4. Setelah login berhasil, aplikasi mengarahkan ke `/board`.

Jika sesi tidak valid, aplikasi akan mengembalikan user ke `/login`.

## Cara Pakai Board

- Gunakan sidebar untuk berpindah antar halaman.
- Gunakan company selector jika Anda bekerja pada company tertentu.
- Gunakan breadcrumb di header untuk melihat posisi halaman.
- Klik avatar/menu user untuk membuka profil, password, history, dan logout.

## Halaman Dasar

### SP01 - Profile / Company Context

- Dipakai untuk melihat company aktif dan data user.
- Pastikan company yang dipilih sesuai konteks kerja.

### SP02 - Change Password

- Dipakai untuk mengubah password akun.
- Field password lama dan baru harus diisi sesuai validasi.

### SP03 - History

- Dipakai untuk melihat aktivitas akun.
- Cocok untuk audit ringan atau cek riwayat aksi.

### SM01 - User Management

- Dipakai admin untuk membuat/mengubah user.
- Dapat mengatur company, privilege, dan area user.

### SM02 - Module Management

- Dipakai admin untuk mengelola modul aplikasi.

### SM03 - Company Management

- Dipakai admin untuk mengelola company, module assignment, dan area.

### SM04 - Signature Type

- Dipakai admin untuk mengatur tipe approval dan signer.

### SM05 - Session Management

- Dipakai admin untuk melihat dan mencabut session aktif.

### XX99 - Template Page

- Halaman tersedia pada route `/board/pages/XX99`.
- Dipakai sebagai template CRUD berbasis DataTable + DataDialog.
- Jika backend route XX99 belum aktif, aksi data dapat gagal dan perlu aktivasi route backend terlebih dahulu.

Detail guide: `obx_docs/guide/SITE/XX99.md`

## Pola Interaksi UI

- List data biasanya memakai tabel dengan search, sort, dan pagination.
- Create/update/detail biasanya muncul di dialog terpisah.
- Jika halaman menampilkan daftar di layar kecil, aplikasi akan otomatis menampilkan kartu list mobile.
- Beberapa halaman memakai SearchSelect untuk pemilihan company atau data referensi.

## Pesan Error Umum

| Gejala | Arti | Langkah |
|---|---|---|
| Redirect ke login | Session expired / cookie hilang | Login ulang |
| 401 saat simpan data | Token invalid | Pastikan sesi belum expired |
| Company kosong | Konteks company belum dipilih | Pilih company lewat sidebar atau login ulang |
| Data tidak tampil | Filter / koneksi backend bermasalah | Refresh halaman dan cek backend |
| Tombol simpan tidak jalan | Loading/validasi field | Periksa error form di bawah field |

## Praktik Aman Pengguna

- Jangan bagikan cookie sesi ke orang lain.
- Jangan hapus localStorage/cookie saat sedang bekerja kecuali Anda memang logout.
- Pastikan company aktif benar sebelum melakukan perubahan data.
- Gunakan menu logout saat selesai bekerja.

## Troubleshooting

- Jika halaman putih atau error boundary tampil, refresh browser terlebih dahulu.
- Jika sidebar kosong, cek session dan endpoint `/proxy/pages/SP01/module`.
- Jika login gagal, cek company/username/password/captcha.
- Jika perubahan data gagal, cek pesan error di toast dan field form.
- Jika halaman XX99 gagal memuat data, cek registrasi route XX99 di backend.
