# BOARD - Dashboard Shell

## Ringkasan

- Route: /board
- Tujuan: dashboard landing dan shell navigasi halaman modul
- Komponen utama: Sidebar, Breadcrumb, ErrorBoundary

## Technical Notes

- Board dilindungi proxy.ts
- Sidebar membaca session dan company context
- Layout memakai shared theme preference

## UI Pattern

- Shell dashboard dengan sidebar dan header
- Landing page dapat berisi ringkasan atau widget

## Checklist

- [ ] Redirect guard aktif
- [ ] Sidebar tampil
- [ ] Breadcrumb tampil
