# ict_monitoring - Panduan Monitoring Infra Domain

## Ringkasan Fitur

Schema ini dipakai untuk menyimpan metric sample, alert rule, backup job, dan histori aksi infrastruktur.

## Prasyarat Akses

1. Pengguna memiliki akses modul monitoring yang sesuai.
2. Data host, VM, atau device sudah terdaftar.

## Langkah Penggunaan

1. Buat atau perbarui alert rule sesuai threshold operasional.
2. Tinjau metric sample berkala untuk host, VM, dan network.
3. Catat backup job dan action log untuk audit operasional.

## Validasi Hasil

- Metric sample tampil dengan timestamp yang benar.
- Alert rule dapat dipicu sesuai threshold.
- Backup job dan action log tercatat.

## Troubleshooting

- Jika metric tidak muncul, cek alur collector dan company scope.
- Jika alert tidak terpicu, cek nilai threshold dan status rule.
- Jika backup job gagal tercatat, cek status runtime dan jalur log.
