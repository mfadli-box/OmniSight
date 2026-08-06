# Rancangan Submodul: SITE - obx_site Frontend

## Ringkasan

- Scope: frontend web app berbasis Next.js App Router
- Tujuan: UI operator untuk login, board, dan halaman modul backend
- Data source utama: `obx_rest` via `/proxy/...`

## Struktur Dokumen

- Summary teknis: file ini
- Detail page teknis: `obx_docs/blueprint/SITE/*.md`
- Summary user guide: `obx_docs/guide/SITE.md`
- Detail page user guide: `obx_docs/guide/SITE/*.md`

## Inti Platform

- Root layout: `obx_site/src/app/layout.tsx`
- Board shell: `obx_site/src/app/board/layout.tsx`
- Board landing: `obx_site/src/app/board/page.tsx`
- Route guard: `obx_site/proxy.ts`

## Route Domain

- `/login`: login UI
- `/board`: dashboard landing
- `/board/pages/SP01`: profile/company page
- `/board/pages/SP02`: change password
- `/board/pages/SP03`: action history
- `/board/pages/SM01`: user management
- `/board/pages/SM02`: module management
- `/board/pages/SM03`: company management
- `/board/pages/SM04`: signature type management
- `/board/pages/SM05`: session management
- `/board/pages/XX99`: template page (tersedia di source code)

## Pattern Summary

- Login menyimpan session ke cookie + localStorage `OmniSightMemory`
- `clientApi()` menambahkan Authorization dan X-Company-ID
- CRUD page memakai DataTable + DataDialog
- Mobile list memakai DataTableCard bila padat
- Form write action wajib memiliki loading state eksplisit

## Page Index

| Page | Detail Teknis | Detail Guide |
|---|---|---|
| SP01 | SP01.md | SP01.md |
| SP02 | SP02.md | SP02.md |
| SP03 | SP03.md | SP03.md |
| SM01 | SM01.md | SM01.md |
| SM02 | SM02.md | SM02.md |
| SM03 | SM03.md | SM03.md |
| SM04 | SM04.md | SM04.md |
| SM05 | SM05.md | SM05.md |
| XX99 | XX99.md | XX99.md |

## Checklist

- [ ] Proxy guard aktif untuk /board
- [ ] Session tersimpan dan dibaca konsisten
- [ ] CRUD page memakai DataTable/DataDialog
- [ ] SearchSelect/Combobox dipakai sesuai kebutuhan
- [ ] Loading state write action ada
- [ ] Mobile card view aktif bila diperlukan
- [ ] Typecheck/build lulus

## Referensi

- `obx_site/proxy.ts`
- `obx_site/src/app/layout.tsx`
- `obx_site/src/app/login/page.tsx`
- `obx_site/src/app/board/layout.tsx`
- `obx_site/src/app/board/page.tsx`
- `obx_site/src/app/board/pages/*`
- `obx_site/src/uix/*`
