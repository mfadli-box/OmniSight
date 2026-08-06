# LOGIN - Authentication Page

## Ringkasan

- Route: /login
- Tujuan: login user dan menyimpan session OmniSightMemory
- Komponen utama: Card, SearchSelect, InputGroup, Input

## Technical Notes

- Memanggil /proxy/guest/SP00 untuk load HRIS company dan submit login
- Session disimpan ke localStorage dan cookie
- Captcha sederhana divalidasi di client
- Setelah sukses redirect ke /board

## UI Pattern

- Form login terpusat
- Company selector memakai SearchSelect
- Error login ditampilkan inline

## Checklist

- [ ] HRIS company termuat
- [ ] Session tersimpan
- [ ] Redirect ke /board berjalan
