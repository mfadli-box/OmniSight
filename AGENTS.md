# OmniSight

**OmniSight** adalah monorepo yang bertujuan untuk membantu dalam mengelola dan mendokumentasikan kegiatan yang dilakukan oleh IT, terutama **DevSecOps**.

## Struktur Folder

| Folder | Deskripsi |
|--------|-----------|
| `obx_base/` | Model dan struktur database **PostgreSQL 18** menggunakan **Prisma ORM 7.8.0** dengan teknik **Multi-file schema layout** |
| `obx_docs/` | Pusat dokumentasi proyek (blueprint, guide, autopilot). Hanya `obx_docs/config/` yang terproteksi dan dilarang diakses AI agent |
| `obx_rest/` | Backend API menggunakan **Go 1.26.4** & **Gin Gonic** dengan arsitektur **Linear Layer Execution**: Template (Struct & Interface) → Repository → Usecase → Handler. Menggunakan **raw SQL** via `database/sql` + `lib/pq` (tanpa ORM di backend) |
| `obx_site/` | Frontend menggunakan **Next.js 16.2.11** (App Router) & **React 19.2.4**, **TypeScript 5** (Strict Mode), serta **Tailwind CSS 4** & **shadcn/ui** (@base-ui) dan library lain sesuai kebutuhan |
| `obx_auto/` | Automation agent (Nginx Log Sync, WAF & Log Rotate) menggunakan **Go 1.26.4** |

## Dokumentasi Aktif Saat Ini

- Indeks blueprint teknis: `obx_docs/blueprint/README.md`
- Rancangan aktif (draft): `obx_docs/blueprint/PLAN/README.md`
- Indeks autopilot: `obx_docs/autopilot/README.md`
- Roadmap rancangan JMS MVP: `obx_docs/autopilot/roadmap_JMS_MVP.md`

## Status Modul Workspace Saat Ini

- Backend terdeteksi: `SP00`, `SP01`, `SP02`, `SP03`, `SM01`, `SM02`, `SM03`, `SM04`, `SM05`, `XX99`
- Frontend terdeteksi: `SP01`, `SP02`, `SP03`, `SM01`, `SM02`, `SM03`, `SM04`, `SM05`, `XX99`

Catatan: gunakan status ini sebagai baseline operasional terkini hingga ada pembaruan dokumentasi modul berikutnya.

# Batasan Akses

## Folder yang DILINDUNGI — DILARANG DIAKSES

Folder berikut **TIDAK BOLEH dibaca, ditulis, dicari (grep/glob), atau diakses dalam bentuk apapun** oleh AI agent:

- `obx_docs/config/`

### Aturan ketat:
- Dilarang membaca isi file di dalam folder tersebut.
- Dilarang mencari file di dalam folder tersebut.
- Dilarang menjalankan perintah bash yang mengakses path tersebut.
- Dilarang menulis atau mengedit file di dalam folder tersebut.

# Aturan Bahasa

## Bahasa untuk Kode Aplikasi (Wajib Inggris)

Seluruh output/user-facing text dalam kode aplikasi WAJIB menggunakan **Bahasa Inggris**, meliputi:

- `alert`, `warning`, `error`, `info`, `message` — semua level notifikasi/pesan
- `data` — properti, field, response API, variabel
- Label UI, tombol, placeholder, tooltip, validasi
- Logging messages (console.log, logger, dll)
- Exception/error messages
- Nama variable, function, class, constant, enum, type
- Komentar dalam kode

### Contoh benar:
```typescript
// ✅ BENAR
throw new Error("User not found");
console.warn("Disk usage exceeds 90%");
toast.success("Data saved successfully");
const errorMessage = "Invalid credentials";
```

### Contoh salah:
```typescript
// ❌ SALAH
throw new Error("Pengguna tidak ditemukan");
console.warn("Penggunaan disk melebihi 90%");
toast.success("Data berhasil disimpan");
const errorMessage = "Kredensial tidak valid";
```

## Bahasa untuk Dokumentasi (Wajib Indonesia)

Seluruh dokumentasi, penjelasan, catatan teknis, komentar arsitektur, dan komunikasi tertulis dalam project ini WAJIB menggunakan **Bahasa Indonesia** yang baik dan benar sesuai PUEBI.

### Contoh benar:
```markdown
<!-- ✅ BENAR -->
<!-- Fungsi ini bertanggung jawab memvalidasi token JWT -->
Fungsi `authenticateUser` memvalidasi token JWT yang dikirim melalui header Authorization.
```

### Contoh salah:
```markdown
<!-- ❌ SALAH -->
<!-- This function is responsible for validating JWT tokens -->
Function `authenticateUser` validates JWT token sent via Authorization header.
```

## Ringkasan

| Konteks | Bahasa | Contoh |
|---------|--------|--------|
| User-facing text (alert, error, dll) | Inggris | `"File not found"` |
| Nama variable/function/kode | Inggris | `getUserData()` |
| Komentar dalam kode | Inggris | `// retry on failure` |
| Dokumentasi/README/AGENTS.md | Indonesia | Penjelasan fitur |
| Komentar arsitektur/desain | Indonesia | Catatan teknis |

# Standar Docker

Setiap komponen aplikasi (`obx_auto`, `obx_rest`, `obx_site`) harus menyertakan `Dockerfile` multi-tahap (*multi-stage build*) untuk memastikan ukuran *image* produksi tetap minimal, optimal, dan aman.

`docker-compose.yml` hanya digunakan untuk **production**. Untuk test dan debug development, gunakan konfigurasi lokal yang ada pada masing-masing folder komponen, bukan dalam bentuk Docker.

# Standar Frontend (obx_site)

Aturan berikut WAJIB diterapkan pada seluruh pengembangan frontend di `obx_site/`:

## DataTable & DataDialog

- Gunakan **DataTable** untuk menampilkan data dalam bentuk tabel dan detail item.
- Gunakan **DataDialog** sebagai wadah untuk action *detail*, *create*, dan *update*. Jangan gunakan page terpisah untuk form CRUD.
- DataTable harus mendukung sorting, searching, dan pagination.

## Form Validation

- Semua form wajib divalidasi menggunakan **zod** schema.
- Tampilkan error validasi menggunakan komponen **FieldError** yang sudah tersedia di `src/uix/`.
- Jangan gunakan alert atau console untuk menampilkan error form.

## Combobox Component

- Component dropdown dengan search filter (combobox) wajib menampilkan **nama** (label), bukan ID internal.
- Nilai yang dikirim ke API adalah ID, tapi yang ditampilkan ke user adalah nama/display value.
- **WAJIB** sertakan prop `itemToStringLabel` pada `<Combobox>` dari `@/uix/combobox` agar input menampilkan label, bukan raw value. Tanpa itu, input akan menampilkan ID/value mentah.

```tsx
// ✅ BENAR — itemToStringLabel wajib
<Combobox
  value={field.value}
  onValueChange={field.onChange}
  itemToStringLabel={(v) => options.find((o) => o.id === v)?.name ?? ""}
>
  <ComboboxInput placeholder="Select..." />
  <ComboboxContent>
    <ComboboxList>
      {options.map((o) => (
        <ComboboxItem key={o.id} value={o.id}>{o.name}</ComboboxItem>
      ))}
    </ComboboxList>
  </ComboboxContent>
</Combobox>

// ❌ SALAH — tanpa itemToStringLabel, input tampilkan raw value
<Combobox value={field.value} onValueChange={field.onChange}>
  <ComboboxInput placeholder="Select..." />
  <ComboboxContent>
    <ComboboxList>
      {options.map((o) => (
        <ComboboxItem key={o.id} value={o.id}>{o.name}</ComboboxItem>
      ))}
    </ComboboxList>
  </ComboboxContent>
</Combobox>
```

**Catatan**:
- `itemToStringLabel` di `Combobox` root berfungsi untuk menampilkan label di input saat value sudah terpilih
- Pattern saat ini (static children dengan `options.map()`) **tidak mendukung filtering otomatis** saat user mengetik
- Jika filtering diperlukan, gunakan `SearchSelect` dari `@/uix/search-select` (lihat section berikutnya)

## SearchSelect Component

- Component dropdown dengan **search filter** (filtering saat user mengetik). Gunakan `SearchSelect` dari `@/uix/search-select` jika diperlukan filtering.
- `SearchSelect` secara otomatis men-filter options berdasarkan input user.
- Nilai yang dikirim ke API adalah `id` (string), tapi yang ditampilkan ke user adalah `name`.

```tsx
import SearchSelect from "@/uix/search-select"
import { Controller } from "react-hook-form"

<Controller
  control={form.control}
  name="company_id"
  render={({ field, fieldState }) => (
    <Field className="gap-1.5" data-invalid={fieldState.invalid}>
      <FieldLabel>Company</FieldLabel>
      <SearchSelect
        items={companies}
        value={field.value}
        onValueChange={field.onChange}
        placeholder="Select company..."
        error={fieldState.invalid}
      />
      {fieldState.invalid && <FieldError errors={[fieldState.error]} />}
    </Field>
  )}
/>
```

**Props**:
| Prop | Type | Deskripsi |
|------|------|-----------|
| `items` | `Array<{id: string, name: string}>` | Array options untuk dropdown |
| `value` | `string \| null \| undefined` | Nilai yang dipilih (ID) |
| `onValueChange` | `(value: string \| null) => void` | Callback saat nilai berubah |
| `placeholder` | `string` | Placeholder saat tidak ada yang dipilih |
| `disabled` | `boolean` | Disable interaksi |
| `className` | `string` | Custom class untuk container |
| `error` | `boolean` | Tampilkan state error (border merah) |

**Fitur**:
- Dropdown tidak auto-open saat focus - hanya buka saat klik tombol chevron
- Filter otomatis saat user mengetik
- Clear button untuk reset pilihan
- Keyboard navigation (arrow keys, enter, escape)
- Touch-friendly dengan `touch-manipulation`
- Mobile responsive dengan viewport height constraints

## Form Field Value

- Field optional (`z.string().optional()`) akan menghasilkan `field.value = undefined` saat form kosong.
- **WAJIB** sertakan `value={field.value ?? ""}` pada `<Input>` dan `<Textarea>` yang spread `{...field}` agar komponen tetap **controlled**. Tanpa ini, Base UI melempar error: *"A component is changing the uncontrolled value state to be controlled"*.
- **WAJIB** sertakan `defaultValues` pada `useForm()` agar semua field memiliki nilai awal (bukan `undefined`). Tanpa ini, field mulai sebagai uncontrolled lalu berubah ke controlled saat `reset()` dipanggil.

```tsx
// ✅ BENAR — defaultValues + value={field.value ?? ""}
const form = useForm<FormValues>({
  resolver: zodResolver(schema),
  defaultValues: {
    code: "",
    name: "",
    description: "",
    is_active: true,
  },
});

<Input {...field} value={field.value ?? ""} />

// ❌ SALAH — tanpa defaultValues, field.value = undefined (uncontrolled)
const form = useForm<FormValues>({ resolver: zodResolver(schema) });

<Input {...field} />
```

## UIX Components

- Selalu gunakan komponen yang sudah ada di `src/uix/` untuk membangun UI. Jangan buat elemen HTML mentah (seperti `<input>`, `<button>`) langsung.
- Jika ada komponen yang belum tersedia, tambahkan ke folder `src/uix/` dengan pola yang konsisten.

## Mobile Compatibility

- Semua halaman harus responsif dan kompatibel untuk **mobile web**.
- Gunakan breakpoint Tailwind (`sm:`, `md:`, `lg:`) untuk menangani layout responsif.
- Sidebar harus bisa disembunyikan/ditutup pada mobile.
- Gunakan hook `useIsMobile()` dari `src/lib/view-mobile.ts` untuk conditional rendering jika diperlukan.

### Pages Mobile View (WAJIB)

Setiap halaman di `src/app/board/pages/{KODE}/` WAJIB menerapkan:

#### 1. useIsMobile() Hook
```tsx
import { useIsMobile } from "@/lib/view-mobile";

// Di dalam komponen:
const isMobile = useIsMobile();
```

#### 2. Conditional View: Table vs Card
Gunakan conditional rendering untuk menampilkan view sesuai screen size:

```tsx
// ✅ BENAR — tampilkan card view pada mobile, table view pada desktop
{isMobile ? (
  <div className="flex flex-col gap-3">
    {data.map((row) => (
      <DataTableCard key={row.id} row={row} columns={columns} onEdit={...} onDelete={...} />
    ))}
  </div>
) : (
  <DataTable ... />  // Desktop: table view
)}
```

#### 3. DataTableCard Component (src/uix/datatable-card.tsx)
Gunakan komponen `DataTableCard` untuk card view pada mobile. Props:
| Prop | Type | Deskripsi |
|------|------|-----------|
| `row` | `T` | Data row |
| `columns` | `Column<T>[]` | Column definitions |
| `onEdit` | `(row: T) => void` | Edit handler |
| `onDelete` | `(row: T) => void` | Delete handler |
| `onDetail` | `(row: T) => void` | Detail handler |

Setiap card menampilkan:
- Kolom pertama sebagai **header card** (nama/title)
- Field-value pairs untuk setiap kolom yang visible
- Action buttons di bagian bawah card
- Border, rounded corners, shadow untuk visual separation

#### 4. Touch-Friendly Action Buttons
- Minimum touch target: **44x44px** (`h-11 w-11`)
- Selalu tambahkan `touch-manipulation` pada tombol aksi
- Jarak antar tombol minimal 8px

```tsx
// ✅ BENAR
<div className="flex items-center gap-2">
  <Button size="sm" className="touch-manipulation h-11 w-11" onClick={() => onDetail(row)}>
    <FileText />
  </Button>
  <Button size="sm" className="touch-manipulation h-11 w-11" onClick={() => onEdit(row)}>
    <Pencil />
  </Button>
</div>
```

#### 5. Responsive Layout Classes
- Container: `p-3 max-w-8xl mx-auto space-y-3`
- Grid untuk form fields: `sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4`
- Stack pada mobile, row pada desktop: `flex flex-col gap-3 sm:flex-row sm:items-center`

### Scrollbar & Overflow

- **Double scrollbar** terjadi jika ada nested scrollable containers. Untuk DataDialog dengan Tabs + DataTable:
  - DataDialog ContentWrapper: `max-h-[85dvh] overflow-y-auto flex-1` (handle scroll utama)
  - TabsContent: `flex flex-col overflow-hidden` dengan wrapper div `flex-1 overflow-auto min-h-0`
  - **JANGAN** tambahkan `max-h-[70vh] overflow-y-auto` di FieldGroup dalam DataDialog - ini menyebabkan double scrollbar
  - Jika ada DataTable dalam tabs, DataTable harus stretch fill available space

### Touch Targets

- Minimum touch target: **44x44px** sesuai Apple HIG dan Google Material Design guidelines
- Tombol aksi di DataTable: `h-9 w-9` atau `h-11 w-11` untuk primary actions
- Dialog buttons: `h-11 px-6` untuk touch-friendly
- Gunakan class `touch-manipulation` untuk disable double-tap zoom pada tombol

### SearchSelect

- Dropdown **tidak auto-open** saat focus - hanya buka saat klik tombol chevron
- Filter berjalan saat user mengetik di input
- Keyboard navigation: Arrow Up/Down, Enter untuk select, Escape untuk close
- On touch devices: gunakan `onTouchStart` untuk highlight, bukan `onMouseEnter`

# Standar Kode Halaman (Module Code)

> **Referensi lengkap**: `README.md`

## Kewajiban

1. **Sebelum pengembangan**: Baca `README.md` untuk memahami kode halaman, privilege level, dan route yang terkait.
2. **Saat pengembangan**: Gunakan kode halaman dari `README.md` sebagai identitas utama. Folder skeleton, route backend, dan path frontend WAJIB menggunakan kode yang terdaftar.
3. **Saat review**: Periksa bahwa kode halaman, privilege level, dan route sesuai dengan `README.md`.
4. **Saat debug**: Cek `README.md` untuk memastikan privilege level dan HTTP method sesuai.

## Konvensi

- **Group/Modul**: 2 huruf kapital — `SP`, `SM`, dll.
- **Halaman Turunan**: Group + 2 digit — `SP01`, `SP02`, `SM01`, dll.
- **Folder Skeleton**: `obx_rest/skeleton/{KODE}/` — satu folder per halaman.
- **Route Backend**: `/rest/pages/{KODE}` — WAJIB sama dengan kode halaman.
- **Route Frontend**: `/board/pages/{KODE}` — WAJIB sama dengan kode halaman.
- **Sidebar**: ID item nav WAJIB menggunakan kode halaman.

## Aturan Privilege

| Tipe Halaman | Middleware | Privilege Level |
|--------------|------------|-----------------|
| Admin Only | `USLock()` | `isAdmin = true` |
| Pages (CRUD) | `USLoad()` | `dat_user_privilege.level >= POST` |
| CRUD (POST) | `USRole("SM01", "POST")` | `dat_user_privilege.level >= POST` |
| View Only | `USRole("SM04", "VIEW")` | `dat_user_privilege.level >= VIEW` |

## Update Terakhir

Format **Oleh** di bagian `## Update Terakhir` pada `README.md`:

| Yang Mengubah | Format |
|---------------|--------|
| AI Agent | `opencode/{model}` — contoh: `opencode/big-pickle` |
| User/Owner | `Owner` |

Contoh:
```markdown
## Update Terakhir

- **Tanggal**: 27 Juli 2026
- **Oleh**: opencode/big-pickle
- **Keterangan**: Rancangan modul SP, SM, IM
```

### Alur Privilege Check

```
Request masuk
  → USLoad()            : autentikasi session token
  → USRole(kode, level) : cek dat_user_privilege berdasarkan module_id
                          - Admin? Langsung lolos
                          - Non-admin? Cari level untuk (user_company_id, module_id)
                          - level == "HIDE" atau tidak ada? → 403 FORBIDDEN
                          - level tidak mencukupi? → 403 FORBIDDEN
  → Handler             : proses request
```

### Level Hierarchy

```
HIDE < VIEW < BOOK < POST
```

- `GET` request memerlukan minimal `VIEW`
- `POST` / `PUT` / `DELETE` request memerlukan `POST`
- `BOOK` untuk operasi reservasi/penandaan (opsional)

---

# Standar Error Handling

## Backend (obx_rest)

### AppError — Satu-Satunya Cara Melempar Error

Seluruh error di backend WAJIB menggunakan `mechanic.AppError` via helper constructor:

| Helper | Kode | HTTP Status | Penggunaan |
|--------|------|-------------|------------|
| `mechanic.ValidationError(msg)` | `VALIDATION_ERROR` | 400 | Input tidak valid, field wajib kosong |
| `mechanic.NotFound(msg)` | `NOT_FOUND` | 404 | Resource tidak ditemukan |
| `mechanic.Conflict(msg)` | `CONFLICT` | 409 | Data duplikat (username sudah ada) |
| `mechanic.Unauthorized(msg)` | `UNAUTHORIZED` | 401 | Token tidak valid/kedaluwarsa |
| `mechanic.Forbidden(msg)` | `FORBIDDEN` | 403 | Akses ditolak (non-admin ke route admin) |
| `mechanic.ExternalServiceError(msg, err)` | `EXTERNAL_SERVICE_ERROR` | 502 | Service eksternal tidak terjangkau |
| `mechanic.InternalError(msg, err)` | `INTERNAL_ERROR` | 500 | Panic recovery, error tak terduga |

### Handler — Selalu `mechanic.Error(c, err)`

```go
// ✅ BENAR
if err := h.usecase.DoSomething(req); err != nil {
    mechanic.Error(c, err)
    return
}

// ❌ SALAH
if err := h.usecase.DoSomething(req); err != nil {
    c.JSON(400, gin.H{"error": err.Error()})
    return
}
```

### Usecase — Wrap DB error dengan `InternalError`

Usecase WAJIB wrap raw error dari repository dengan `mechanic.InternalError()` agar error bisa dilog dengan structured. **Jangan return raw error** dari repo — AI tidak bisa debug tanpa mengetahui error type.

```go
// ✅ BENAR — raw error di-wrap dengan InternalError
func (u *useCase) CreateItem(name string) error {
    if name == "" {
        return mechanic.ValidationError("Name is required")
    }
    if err := u.repo.Insert(name); err != nil {
        return mechanic.InternalError("Failed to create item", err)
    }
    return nil
}

// ❌ SALAH — raw error dari DB, mechanic.Error() tidak bisa identify type
func (u *useCase) CreateItem(name string) error {
    return u.repo.Insert(name)
}
```

### Alur Error Logging

```
Error terjadi di usecase/repo
  → usecase wrap: mechanic.InternalError("msg", err)
  → handler: mechanic.Error(c, err)
  → mechanic.Error():
      1. Deteksi tipe error via errors.As()
      2. LOG structured: method, path, request_id, user_id, error code, error detail
      3. Kirim JSON response ke client (pesan aman untuk user)
```

**Log output**:
```
WRN app error code=VALIDATION_ERROR error="Name is required" status=400 method=POST path=/rest/pages/SM03/area request_id=abc-123 user_id=u1
ERR unhandled error error="pq: null value in column company_id violates not-null constraint" method=POST path=/rest/pages/SM03/area request_id=abc-123 user_id=u1
```

### Repository — Return `error` biasi

Repository tidak perlu `mechanic.Error()`. Cukup return `error` dari query database, atau `mechanic.Conflict(...)` untuk kasus duplikat constraint. Usecase yang wrap.

### Response Format

Semua response error WAJIB dalam format:
```json
{
  "code": "VALIDATION_ERROR",
  "error": "Username is required",
  "request_id": "uuid-dari-request"
}
```

### Logging

Gunakan `backbone.Log` (zerolog) untuk structured logging:
```go
backbone.Log.Info().Str("user_id", id).Msg("User created")
backbone.Log.Warn().Str("token", t).Msg("Token expiring soon")
backbone.Log.Error().Err(err).Msg("Database query failed")
```

Jangan gunakan `log.Println()` atau `fmt.Println()` untuk logging di produksi.

---

## Frontend (obx_site)

### clientApi — Satu-Satunya Cara Fetch API

Seluruh fetch ke backend WAJIB menggunakan `clientApi()` dari `src/lib/client-api.ts`. Jangan gunakan `fetch()` langsung atau `apiFetch()` untuk request yang memerlukan autentikasi.

```typescript
// ✅ BENAR
import { clientApi, ClientApiError } from "@/lib/client-api";

try {
  const data = await clientApi<DataType>("/SM01", { method: "GET" });
} catch (err) {
  if (err instanceof ClientApiError) {
    toast.error(getErrorMessage(err.code));
  }
}

// ❌ SALAH
const res = await fetch("/proxy/pages/SM01");
```

### Toast — Satu-Satunya Cara Menampilkan Error ke User

Gunakan `toast` dari `sonner` untuk semua notifikasi. Jangan gunakan `alert()`, `console.log()` untuk error user-facing, atau state `alertMsg`/`errorMsg` manual.

```typescript
// ✅ BENAR
toast.error(getErrorMessage(err.code));
toast.success("Data saved successfully");

// ❌ SALAH
alert("Data saved successfully");
setAlertMsg("Data saved successfully");
```

### Error Code → User Message

Gunakan `getErrorMessage(code)` dari `src/lib/error-message.ts` untuk mengubah error code ke pesan yang bisa dimengerti user. Jangan hardcode pesan error di halaman.

### ErrorBoundary — Untuk Unhandled Runtime Error

Pasang `<ErrorBoundary>` di `layout.tsx` (root) dan `board/layout.tsx`. Jangan pasang di setiap halaman individual.

### Silent Catch Dilarang

```typescript
// ❌ SALAH — error ditelan, user tidak tahu ada masalah
fetchData().catch(() => {});

// ✅ BENAR — minimal log atau tampilkan toast
clientApi("/SM01").catch(() => {
  setData([]);
});
```

---

## Proxy Routes

### Error Transform

Proxy WAJIB meneruskan response error dari backend dengan format yang konsisten:
- Forward header `X-Request-ID` dari backend ke client
- Jika backend unreachable, return 502 dengan `EXTERNAL_SERVICE_ERROR`
- Jangan ubah format error response dari backend

### Network Error Handling

```typescript
// Proxy harus catch fetch errors (backend unreachable)
try {
  res = await fetch(`${BE_POOL}/rest/pages${apiPath}`, { ... });
} catch (err) {
  const message = err instanceof Error ? err.message : "Backend service unreachable";
  console.error(`[PROXY] ${method} ${apiPath} — backend unreachable:`, message);
  return NextResponse.json(
    { code: "EXTERNAL_SERVICE_ERROR", error: message, request_id: "" },
    { status: 502 },
  );
}
```

### Error Logging di Proxy

Proxy WAJIB log error dari backend untuk debugging:
```typescript
if (res.status >= 400) {
  const errorBody = data ? JSON.stringify(data) : "(empty)";
  console.error(`[PROXY] ${method} ${apiPath} → ${res.status}`, errorBody);
}
```

**Alur debug lengkap**:
```
Frontend submit form
  → clientApi() → proxy route (Next.js)
    → proxy log: [PROXY] POST /SM03/type → 500 {"code":"INTERNAL_ERROR","error":"Failed to create area"}
    → proxy forward response ke client
  → backend log: WRN app error code=INTERNAL_ERROR error="Failed to create area" method=POST path=/rest/pages/SM03/area request_id=abc-123
    → backend log (underlying): ERR error="pq: null value in column company_id violates not-null constraint" method=POST path=/rest/pages/SM03/area
```

---

# Standar Keamanan (Security)

## Backend (obx_rest)

### Input Validation
- Selalu validasi input di usecase layer sebelum proses bisnis
- Gunakan `mechanic.ValidationError()` untuk input yang tidak valid
- Jangan trust data dari request body — sanitasi sebelum insert/query

### SQL Injection Prevention
- Gunakan **parameterized queries** (`$1`, `$2`, ...) — JANGAN pernah string concatenation untuk query
- Repository layer WAJIB menggunakan `sqlx.In()`, `db.QueryRow(query, args...)`, atau `db.Exec(query, args...)`

### Authentication & Session
- Bearer token via `dat_user_session` dengan expiry (24 jam)
- Token di-generate dari `crypto/rand` (32 bytes hex)
- Password di-hash dengan **bcrypt** (cost 10+)
- HRIS password di-encrypt dengan **AES-GCM** via `mechanic.Encrypt()`

### Sensitive Data
- Jangan log password, token, atau key dalam bentuk plain text
- Jangan return password hash ke client
- Gunakan `mechanic.NullableString()` untuk field opsional yang nullable

## Frontend (obx_site)

### XSS Prevention
- React secara otomatis escape JSX content — jangan gunakan `dangerouslySetInnerHTML` kecuali benar-benar perlu
- Jika harus render HTML, sanitasi dengan library yang sesuai

### CSRF & Auth
- Semua request ke backend WAJIB menggunakan `clientApi()` yang otomatis attach Bearer token
- Jangan simpan token di URL query params
- Cookie `OmniSightMemory` hanya untuk preference (bukan auth token)

### Dependency Security
- Jalankan `npm audit` secara berkala — 0 vulnerabilities adalah target
- Gunakan npm overrides untuk memaksa versi aman jika ada vulnerability di dependency tree

---

# Standar Optimasi (Performance)

## Backend (obx_rest)

### Database
- Gunakan **index** pada kolom yang sering di-query (lihat `obx_base` schema)
- Gunakan **pagination** (LIMIT/OFFSET) untuk semua list endpoint
- Gunakan **COUNT(*)** terpisah untuk total — jangan hitung dari result set
- **Connection pooling** via `database/sql` (atur `SetMaxOpenConns`, `SetMaxIdleConns`)

### Caching
- Response yang jarang berubah (module tree, company list) bisa di-cache di-memory
- Session data cached di middleware layer (USLoad)

## Frontend (obx_site)

### Data Fetching
- Gunakan `clientApi()` untuk data fetching di client components
- Gunakan `serverApi()` untuk data fetching di server components (React Server Components)
- Implementasi **loading states** menggunakan Skeleton components

### Bundle Optimization
- Dynamic import untuk komponen yang jarang digunakan (`next/dynamic`)
- Gunakan `React.memo()` untuk komponen yang sering re-render
- Lazy load chart libraries (recharts) hanya saat diperlukan

### Mobile Performance
- Responsive design dengan Tailwind breakpoints (`sm:`, `md:`, `lg:`)
- DataTable: card view pada mobile, table view pada desktop
- Sidebar collapsible pada mobile

---

# Pola & Konvensi Tambahan

## State Management
- **Zustand** untuk global state (preferences, session, companyId)
- **React Hook Form** untuk form state
- Jangan campur aduk state management — pilih yang sesuai konteks

## Error Handling Flow

```
Backend Error
  → mechanic.Error(c, AppError)
  → JSON response { code, error, request_id }
  → Frontend proxy (forward X-Request-ID)
  → clientApi() throws ClientApiError
  → Page: toast.error(getErrorMessage(err.code))
```

## Form Validation Flow

```
Zod Schema
  → React Hook Form (resolver: zodResolver)
  → FieldError component untuk error display
  → Submit ke API → handle error → toast
```

---

# Module Blueprint (Self-Contained)

> **Penting**: Section ini dirancang agar AI agent **cukup membaca AGENTS.md saja** untuk membuat module baru. Blueprint disusun berdasarkan pattern aktual dari modul SP dan SM yang sudah terbukti jalan.

## Workflow Pembuatan Module

```
1. Cek kode halaman di README.md (route & privilege)
2. Buat 4 file backend: template.go → repository.go → usecase.go → handler.go
3. Register DI + routes di backbone/routes.go
4. Buat 1 file frontend: board/pages/{KODE}/page.tsx
5. Build test: go build ./... && npx tsc --noEmit
```

## Referensi Cepat

| Item | Nilai |
|------|-------|
| **Backend path** | `obx_rest/skeleton/{KODE}/` |
| **Frontend path** | `obx_site/src/app/board/pages/{KODE}/` |
| **Route backend** | `/rest/pages/{KODE}` (sesuaikan dengan group) |
| **Route frontend** | `/board/pages/{KODE}` |
| **Package Go** | `package {KODE}` (contoh: `package SM03`) |
| **Write logs** | `USLogs("{KODE}")` wajib di semua POST/PUT/DELETE |

### Route Groups

| Group | Middleware | Penggunaan | Contoh |
|-------|-----------|------------|--------|
| `guest` | Tanpa middleware | Login, public endpoints | SP00 (login) |
| `authu` | `USAuth()` | Autentikasi tanpa company scoping | SP01, SP02, SP03 |
| `admin` | `USAuth()` + `USLock()` | Admin-only CRUD | SM01, SM02, SM03, SM04, SM05 |
| `pages` | `USLoad()` | Multi-tenant dengan company_id | IM01 |

### Context Extraction

Middleware menyuntikkan data ke gin context. Gunakan metode yang sesuai:

```go
// Dari middleware USLoad/USAuth — untuk route yang pakai session
userID, _ := c.Get("userId")        // string, dari session
companyID := c.GetString("companyId") // string, dari X-Company-ID header
isAdmin, _ := c.Get("isAdmin")       // bool, dari session

// Dari route parameter — untuk nested sub-entity
id := c.Param("id")                  // parent ID
subId := c.Param("subId")            // child ID (nama variatif: companyId, areaId, moduleId)
```

---

## File 1: `skeleton/{KODE}/template.go`

```go
package {KODE}

import "obx_rest/mechanic"

// ===== DTOs =====

// Main entity — list response
type {Name}Item struct {
	ID        string `json:"id"`
	Code      string `json:"code"`
	Name      string `json:"name"`
	IsActive  bool   `json:"is_active"`
	CreatedAt string `json:"created_at"`
}

// Request structs
type {Name}CreateRequest struct {
	Code     string `json:"code" binding:"required"`
	Name     string `json:"name" binding:"required"`
	IsActive *bool  `json:"is_active"`
}

type {Name}UpdateRequest struct {
	Code     string `json:"code"`
	Name     string `json:"name"`
	IsActive *bool  `json:"is_active"`
}

// Sub-entity (opsional — jika module punya child table)
type {Name}SubListItem struct {
	ID        string `json:"id"`
	{Name}ID  string `json:"{name_id}"`
	SubName   string `json:"sub_name"`
	IsActive  bool   `json:"is_active"`
	CreatedAt string `json:"created_at"`
}

type {Name}SubCreateRequest struct {
	SubName  string `json:"sub_name" binding:"required"`
	IsActive *bool  `json:"is_active"`
}

type {Name}SubUpdateRequest struct {
	SubName  string `json:"sub_name"`
	IsActive *bool  `json:"is_active"`
}

// Select item (untuk search-select)
type {Name}SelectItem struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ===== Interfaces =====

type Repository interface {
	// Main entity
	Count(search string) (int, error)
	List(search string, page, size int, sortBy, sortOrder string) ([]{Name}Item, error)
	Create(req {Name}CreateRequest) error
	Update(id string, req {Name}UpdateRequest) error
	Delete(id string) error

	// Sub-entity (opsional)
	CountSub(parentID, search string) (int, error)
	ListSub(parentID, search string, page, size int, sortBy, sortOrder string) ([]{Name}SubListItem, error)
	CreateSub(parentID string, req {Name}SubCreateRequest) error
	UpdateSub(id string, req {Name}SubUpdateRequest) error
	DeleteSub(id string) error

	// Select endpoints (opsional)
	ListSelect() ([]{Name}SelectItem, error)
}

type UseCase interface {
	// Main entity
	List(meta mechanic.ActionMeta) ([]{Name}Item, mechanic.GridMeta, error)
	Create(req {Name}CreateRequest) error
	Update(id string, req {Name}UpdateRequest) error
	Delete(id string) error

	// Sub-entity (opsional)
	ListSub(parentID string, meta mechanic.ActionMeta) ([]{Name}SubListItem, mechanic.GridMeta, error)
	CreateSub(parentID string, req {Name}SubCreateRequest) error
	UpdateSub(id string, req {Name}SubUpdateRequest) error
	DeleteSub(id string) error

	// Select endpoints (opsional)
	ListSelect() ([]{Name}SelectItem, error)
}
```

**Ubah**: `{KODE}` → kode halaman, `{Name}` → nama entity, `{name_id}` → FK column name.
Hapus sub-entity methods jika tidak diperlukan. Sesuaikan field struct dengan schema DB.

**Catatan aktual** (lihat SM01, SM03):
- `binding:"required"` hanya untuk field yang benar-benar wajib di create request
- Update request biasanya tanpa `binding:"required"` (field optional)
- `Delete` pada repository cukup return `error` — tidak perlu `RowsAffected()`
- Beberapa module punya DTO tambahan (select items, HRIS items) sesuai kebutuhan

---

## File 2: `skeleton/{KODE}/repository.go`

```go
package {KODE}

import (
	"context"
	"database/sql"
	"fmt"
	"obx_rest/mechanic"
)

type repository struct {
	db *sql.DB
}

func NRepo(db *sql.DB) Repository {
	return &repository{db: db}
}

// ===== Main Entity =====

func (r *repository) Count(search string) (int, error) {
	query := `SELECT COUNT(*) FROM "{table}" WHERE 1=1`
	args := []any{}
	argIdx := 1
	if search != "" {
		query += fmt.Sprintf(" AND (code ILIKE $%d OR name ILIKE $%d)", argIdx, argIdx)
		args = append(args, "%"+search+"%")
	}
	var total int
	err := r.db.QueryRowContext(context.Background(), query, args...).Scan(&total)
	return total, err
}

func (r *repository) List(search string, page, size int, sortBy, sortOrder string) (
	[]{Name}Item, error) {
	query := `
		SELECT id, code, name, is_active,
		       TO_CHAR(created_at, 'YYYY-MM-DD HH24:MI:SS')
		FROM   "{table}" WHERE 1=1
	`
	args := []any{}
	argIdx := 1
	if search != "" {
		query += fmt.Sprintf(" AND (code ILIKE $%d OR name ILIKE $%d)", argIdx, argIdx)
		args = append(args, "%"+search+"%")
		argIdx++
	}
	validSort := map[string]bool{"code": true, "name": true, "created_at": true}
	if !validSort[sortBy] {
		sortBy = "created_at"
	}
	if sortOrder != "asc" {
		sortOrder = "desc"
	}
	query += fmt.Sprintf(" ORDER BY %s %s", sortBy, sortOrder)
	offset := (page - 1) * size
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, size, offset)

	rows, err := r.db.QueryContext(context.Background(), query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []{Name}Item
	for rows.Next() {
		var item {Name}Item
		if err := rows.Scan(
			&item.ID, &item.Code, &item.Name, &item.IsActive, &item.CreatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *repository) Create(req {Name}CreateRequest) error {
	query := `
		INSERT INTO "{table}" (id, code, name, is_active, created_at, updated_at)
		VALUES (gen_random_uuid()::text, $1, $2, $3, NOW(), NOW())
	`
	_, err := r.db.ExecContext(context.Background(), query,
		req.Code, req.Name, mechanic.NullableString("true"),
	)
	return err
}

func (r *repository) Update(id string, req {Name}UpdateRequest) error {
	query := `
		UPDATE "{table}" SET
			code = $2, name = $3, is_active = $4, updated_at = NOW()
		WHERE id = $1
	`
	_, err := r.db.ExecContext(context.Background(), query,
		id, req.Code, req.Name, req.IsActive,
	)
	return err
}

func (r *repository) Delete(id string) error {
	_, err := r.db.ExecContext(context.Background(),
		`DELETE FROM "{table}" WHERE id = $1`, id)
	return err
}

// ===== Sub-entity (opsional) =====

func (r *repository) CountSub(parentID, search string) (int, error) {
	query := `SELECT COUNT(*) FROM "{sub_table}" WHERE {name_id} = $1`
	args := []any{parentID}
	argIdx := 2
	if search != "" {
		query += fmt.Sprintf(" AND (sub_name ILIKE $%d)", argIdx)
		args = append(args, "%"+search+"%")
	}
	var total int
	err := r.db.QueryRowContext(context.Background(), query, args...).Scan(&total)
	return total, err
}

func (r *repository) ListSub(parentID, search string, page, size int, sortBy, sortOrder string) (
	[]{Name}SubListItem, error) {
	query := `
		SELECT s.id, s.{name_id}, s.sub_name, s.is_active,
		       TO_CHAR(s.created_at, 'YYYY-MM-DD HH24:MI:SS')
		FROM "{sub_table}" s
		WHERE s.{name_id} = $1
	`
	args := []any{parentID}
	argIdx := 2
	if search != "" {
		query += fmt.Sprintf(" AND s.sub_name ILIKE $%d", argIdx)
		args = append(args, "%"+search+"%")
		argIdx++
	}
	validSort := map[string]bool{"sub_name": true, "created_at": true}
	if !validSort[sortBy] {
		sortBy = "created_at"
	}
	if sortOrder != "asc" {
		sortOrder = "desc"
	}
	query += fmt.Sprintf(" ORDER BY %s %s", sortBy, sortOrder)
	offset := (page - 1) * size
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, size, offset)

	rows, err := r.db.QueryContext(context.Background(), query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []{Name}SubListItem
	for rows.Next() {
		var item {Name}SubListItem
		if err := rows.Scan(
			&item.ID, &item.{Name}ID, &item.SubName, &item.IsActive, &item.CreatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *repository) CreateSub(parentID string, req {Name}SubCreateRequest) error {
	query := `
		INSERT INTO "{sub_table}" (id, {name_id}, sub_name, is_active, created_at, updated_at)
		VALUES (gen_random_uuid()::text, $1, $2, $3, NOW(), NOW())
	`
	_, err := r.db.ExecContext(context.Background(), query, parentID, req.SubName, mechanic.NullableString("true"))
	return err
}

func (r *repository) UpdateSub(id string, req {Name}SubUpdateRequest) error {
	query := `UPDATE "{sub_table}" SET sub_name = $2, is_active = $3, updated_at = NOW() WHERE id = $1`
	_, err := r.db.ExecContext(context.Background(), query, id, req.SubName, req.IsActive)
	return err
}

func (r *repository) DeleteSub(id string) error {
	_, err := r.db.ExecContext(context.Background(),
		`DELETE FROM "{sub_table}" WHERE id = $1`, id)
	return err
}

// ===== Select (opsional) =====

func (r *repository) ListSelect() ([]{Name}SelectItem, error) {
	query := `SELECT id, name FROM "{table}" WHERE is_active = true ORDER BY name ASC`
	rows, err := r.db.QueryContext(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []{Name}SelectItem
	for rows.Next() {
		var item {Name}SelectItem
		if err := rows.Scan(&item.ID, &item.Name); err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	return list, nil
}
```

**Pola kunci**:
- `gen_random_uuid()::text` untuk ID baru
- `COALESCE()` untuk field NULL di response
- **`rows.Err()` wajib dicek SETELAH loop** — bukan sebelum loop
- `validSort` WAJIB whitelist kolom — hindari SQL injection via sort
- `ORDER BY` harus pakai **SQL alias** jika kolom berasal dari JOIN
- `fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIdx, argIdx+1)` untuk pagination

**Variasi aktual** (lihat SM01, SM03):
- **Join queries**: gunakan `sortExpr` map (bukan `validSort` bool map) untuk map sort key ke SQL expression:
  ```go
  sortExpr := map[string]string{
      "company_name": "company_name",
      "created_at":   "uc.created_at",
  }
  expr, ok := sortExpr[sortBy]
  if !ok {
      expr = "uc.created_at"
  }
  ```
- **Upsert**: gunakan `ON CONFLICT (...) DO UPDATE SET ...` untuk data yang mungkin sudah ada
- **Return ID**: gunakan `RETURNING id` pada INSERT jika perlu ID yang di-generate
- **Search fleksibel**: sesuaikan kolom search (tidak selalu `code ILIKE` — bisa `username ILIKE, fullname ILIKE`)
- **Company scope**: tambahkan parameter `companyID` di awal query untuk multi-tenant

---

## File 3: `skeleton/{KODE}/usecase.go`

```go
package {KODE}

import (
	"obx_rest/mechanic"
	"strings"
)

type useCase struct {
	repo Repository
}

func NCase(r Repository) UseCase {
	return &useCase{repo: r}
}

// ===== Main Entity =====

func (u *useCase) List(meta mechanic.ActionMeta) ([]{Name}Item, mechanic.GridMeta, error) {
	page, size, err := mechanic.CheckMeta(meta.Page, meta.Size)
	if err != nil {
		return nil, mechanic.GridMeta{}, err
	}
	total, err := u.repo.Count(meta.Search)
	if err != nil {
		return nil, mechanic.GridMeta{}, mechanic.InternalError("Failed to count {name}s", err)
	}
	list, err := u.repo.List(meta.Search, page, size, meta.SortBy, meta.SortOrder)
	if err != nil {
		return nil, mechanic.GridMeta{}, mechanic.InternalError("Failed to list {name}s", err)
	}
	return list, mechanic.BuildMeta(page, size, total), nil
}

func (u *useCase) Create(req {Name}CreateRequest) error {
	if strings.TrimSpace(req.Code) == "" {
		return mechanic.ValidationError("Code is required")
	}
	if strings.TrimSpace(req.Name) == "" {
		return mechanic.ValidationError("Name is required")
	}
	if err := u.repo.Create(req); err != nil {
		return mechanic.InternalError("Failed to create {name}", err)
	}
	return nil
}

func (u *useCase) Update(id string, req {Name}UpdateRequest) error {
	if strings.TrimSpace(id) == "" {
		return mechanic.ValidationError("ID is required")
	}
	if err := u.repo.Update(id, req); err != nil {
		return mechanic.InternalError("Failed to update {name}", err)
	}
	return nil
}

func (u *useCase) Delete(id string) error {
	if strings.TrimSpace(id) == "" {
		return mechanic.ValidationError("ID is required")
	}
	if err := u.repo.Delete(id); err != nil {
		return mechanic.InternalError("Failed to delete {name}", err)
	}
	return nil
}

// ===== Sub-entity (opsional) =====

func (u *useCase) ListSub(parentID string, meta mechanic.ActionMeta) ([]{Name}SubListItem, mechanic.GridMeta, error) {
	if strings.TrimSpace(parentID) == "" {
		return nil, mechanic.GridMeta{}, mechanic.ValidationError("Parent ID is required")
	}
	page, size, err := mechanic.CheckMeta(meta.Page, meta.Size)
	if err != nil {
		return nil, mechanic.GridMeta{}, err
	}
	total, err := u.repo.CountSub(parentID, meta.Search)
	if err != nil {
		return nil, mechanic.GridMeta{}, mechanic.InternalError("Failed to count {name} sub-items", err)
	}
	list, err := u.repo.ListSub(parentID, meta.Search, page, size, meta.SortBy, meta.SortOrder)
	if err != nil {
		return nil, mechanic.GridMeta{}, mechanic.InternalError("Failed to list {name} sub-items", err)
	}
	return list, mechanic.BuildMeta(page, size, total), nil
}

func (u *useCase) CreateSub(parentID string, req {Name}SubCreateRequest) error {
	if strings.TrimSpace(parentID) == "" {
		return mechanic.ValidationError("Parent ID is required")
	}
	if strings.TrimSpace(req.SubName) == "" {
		return mechanic.ValidationError("Name is required")
	}
	if err := u.repo.CreateSub(parentID, req); err != nil {
		return mechanic.InternalError("Failed to create {name} sub-item", err)
	}
	return nil
}

func (u *useCase) UpdateSub(id string, req {Name}SubUpdateRequest) error {
	if strings.TrimSpace(id) == "" {
		return mechanic.ValidationError("ID is required")
	}
	if err := u.repo.UpdateSub(id, req); err != nil {
		return mechanic.InternalError("Failed to update {name} sub-item", err)
	}
	return nil
}

func (u *useCase) DeleteSub(id string) error {
	if strings.TrimSpace(id) == "" {
		return mechanic.ValidationError("ID is required")
	}
	if err := u.repo.DeleteSub(id); err != nil {
		return mechanic.InternalError("Failed to delete {name} sub-item", err)
	}
	return nil
}

// ===== Select (opsional) =====

func (u *useCase) ListSelect() ([]{Name}SelectItem, error) {
	list, err := u.repo.ListSelect()
	if err != nil {
		return nil, mechanic.InternalError("Failed to list {name} options", err)
	}
	return list, nil
}
```

**Pola kunci**:
- List: `CheckMeta` → `Count` → `List` → `BuildMeta`
- Create/Update: validasi input (`strings.TrimSpace`) → delegate ke repo
- **Semua error dari repo WAJIB di-wrap** dengan `mechanic.InternalError("pesan", err)` — jangan return raw error
- **Delete** cukup return `error` — tidak perlu cek `RowsAffected()`

---

## File 4: `skeleton/{KODE}/handler.go`

```go
package {KODE}

import (
	"obx_rest/mechanic"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	usecase UseCase
}

func NHand(u UseCase) *Handler {
	return &Handler{usecase: u}
}

// ===== Main Entity =====

func (h *Handler) List(c *gin.Context) {
	var meta mechanic.ActionMeta
	if err := c.ShouldBindQuery(&meta); err != nil {
		mechanic.Error(c, mechanic.ValidationError("Invalid query parameters"))
		return
	}
	list, gridMeta, err := h.usecase.List(meta)
	if err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "{Name}s retrieved",
		"data":    list,
		"meta":    gridMeta,
	})
}

func (h *Handler) Create(c *gin.Context) {
	var req {Name}CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		mechanic.Error(c, mechanic.ValidationError("Invalid request body"))
		return
	}
	if err := h.usecase.Create(req); err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "{Name} created successfully",
	})
}

func (h *Handler) Update(c *gin.Context) {
	id := c.Param("id")
	var req {Name}UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		mechanic.Error(c, mechanic.ValidationError("Invalid request body"))
		return
	}
	if err := h.usecase.Update(id, req); err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "{Name} updated successfully",
	})
}

func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.usecase.Delete(id); err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "{Name} deleted successfully",
	})
}

// ===== Sub-entity (opsional) =====

func (h *Handler) ListSub(c *gin.Context) {
	parentID := c.Param("id")
	var meta mechanic.ActionMeta
	if err := c.ShouldBindQuery(&meta); err != nil {
		mechanic.Error(c, mechanic.ValidationError("Invalid query parameters"))
		return
	}
	list, gridMeta, err := h.usecase.ListSub(parentID, meta)
	if err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Items retrieved",
		"data":    list,
		"meta":    gridMeta,
	})
}

func (h *Handler) CreateSub(c *gin.Context) {
	parentID := c.Param("id")
	var req {Name}SubCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		mechanic.Error(c, mechanic.ValidationError("Invalid request body"))
		return
	}
	if err := h.usecase.CreateSub(parentID, req); err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "Item created successfully",
	})
}

func (h *Handler) UpdateSub(c *gin.Context) {
	id := c.Param("subId")
	var req {Name}SubUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		mechanic.Error(c, mechanic.ValidationError("Invalid request body"))
		return
	}
	if err := h.usecase.UpdateSub(id, req); err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Item updated successfully",
	})
}

func (h *Handler) DeleteSub(c *gin.Context) {
	id := c.Param("subId")
	if err := h.usecase.DeleteSub(id); err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Item deleted successfully",
	})
}

// ===== Select (opsional) =====

func (h *Handler) ListSelect(c *gin.Context) {
	list, err := h.usecase.ListSelect()
	if err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Items retrieved",
		"data":    list,
	})
}
```

**Pola kunci**:
- List: `ShouldBindQuery` → `mechanic.Error(c, err)` untuk semua error
- Create: `ShouldBindJSON` → response 201
- Update: `c.Param("id")` atau `c.Param("subId")` → `ShouldBindJSON`
- Delete: `c.Param("id")` atau `c.Param("subId")` → response 200
- **JANGAN** `c.JSON(400, gin.H{"error": ...})` — selalu `mechanic.Error(c, err)`
- Sub-entity handler ambil `parentID` dari `c.Param("id")`

---

## File 5: Register DI + Routes (`backbone/routes.go`)

```go
import "obx_rest/skeleton/{KODE}"

// DI wiring
{KODE}R := {KODE}.NRepo(PgSQL)
{KODE}U := {KODE}.NCase({KODE}R)
{KODE}H := {KODE}.NHand({KODE}U)

// ===== Main entity routes =====
pages.GET("/{KODE}", {KODE}H.List)
pages.POST("/{KODE}", USLogs("{KODE}"), {KODE}H.Create)
pages.PUT("/{KODE}/:id", USLogs("{KODE}"), {KODE}H.Update)

// ===== Select endpoint (tanpa :id) =====
pages.GET("/{KODE}/select", {KODE}H.ListSelect)

// ===== Sub-entity routes (nested under :id) =====
pages.GET("/{KODE}/:id/sub", {KODE}H.ListSub)
pages.POST("/{KODE}/:id/sub", USLogs("{KODE}"), {KODE}H.CreateSub)
pages.PUT("/{KODE}/:id/sub/:subId", USLogs("{KODE}"), {KODE}H.UpdateSub)
pages.DELETE("/{KODE}/:id/sub/:subId", USLogs("{KODE}"), {KODE}H.DeleteSub)
```

**Rule**:
- `USLogs("{KODE}")` WAJIB untuk semua **write** operations (POST, PUT, DELETE)
- GET tidak perlu `USLogs`
- **Route pattern**: static routes (seperti `/{KODE}/select`) HARUS didaftarkan SEBELUM dynamic routes (seperti `/{KODE}/:id`)
- Group `admin` = `USAuth()` + `USLock()` — untuk halaman admin-only
- Group `authu` = hanya `USAuth()` — untuk halaman tanpa kebutuhan company_id sebagai parameter
- Group `pages` = hanya `USLoad()` — untuk halaman yang bisa diakses non-admin dengan mengirim company_id sebagai parameter

---

## File 6: Frontend `board/pages/{KODE}/page.tsx`

```tsx
"use client";

import { Controller, useForm, useWatch } from "react-hook-form";
import { Suspense, useCallback, useEffect, useMemo, useState } from "react";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { storageKey, parseSession } from "@/lib/utility";
import { clientApi, ClientApiError } from "@/lib/client-api";
import { toast } from "sonner";
import { getErrorMessage } from "@/lib/error-message";
import DataTable, { Column, ActionConfig } from "@/uix/datatable";
import DataDialog from "@/uix/datadialog";
import {
  Field, FieldError, FieldGroup, FieldLabel, FieldSet,
} from "@/uix/field";
import { Input } from "@/uix/input";
import { Checkbox } from "@/uix/checkbox";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/uix/tabs";

// ===== Interfaces =====
interface {Name}Item {
  id: string;
  code: string;
  name: string;
  is_active: boolean;
  created_at: string;
}

interface {Name}SelectItem {
  id: string;
  name: string;
}

// ===== Zod Schemas =====
const createSchema = z.object({
  code: z.string().min(1, "Code is required"),
  name: z.string().min(1, "Name is required"),
  is_active: z.boolean().optional(),
});

const updateSchema = z.object({
  code: z.string().min(1, "Code is required"),
  name: z.string().min(1, "Name is required"),
  is_active: z.boolean().optional(),
});

type CreateForm = z.infer<typeof createSchema>;
type UpdateForm = z.infer<typeof updateSchema>;

export default function Page() {
  const [selectedRow, setSelectedRow] = useState<{Name}Item | null>(null);
  const [refresh, setRefresh] = useState(false);
  const [loading, setLoading] = useState(false);
  const [modal, setModal] = useState<{
    isOpen: boolean;
    mode: "create" | "update" | "detail";
    title: string;
  }>({ isOpen: false, mode: "create", title: "" });

  // Select data (loaded once on mount)
  const [selectItems, setSelectItems] = useState<{Name}SelectItem[]>([]);
  useEffect(() => {
    clientApi<{ data: {Name}SelectItem[] }>("/{KODE}/select")
      .then((data) => setSelectItems(data?.data ?? []))
      .catch(() => setSelectItems([]));
  }, []);

  // ===== Main entity =====
  const loadData = useCallback(async (params: {
    search: string; page: number; size: number;
    sort_by: string; sort_order: "asc" | "desc";
  }) => {
    const session = parseSession(window.localStorage.getItem(storageKey));
    if (!session) return { data: [], meta: { total: 0, page: 1, size: 10 } };
    try {
      const data = await clientApi<{
        data: {Name}Item[];
        meta: { total: number; page: number; size: number };
      }>("/{KODE}", {
        params: {
          search: params.search,
          page: String(params.page),
          size: String(params.size),
          sort_by: params.sort_by,
          sort_order: params.sort_order,
        },
      });
      return {
        data: data?.data ?? [],
        meta: data?.meta ?? { total: 0, page: 1, size: 10 },
      };
    } catch {
      return { data: [], meta: { total: 0, page: 1, size: 10 } };
    }
  }, []);

  const columns: Column<{Name}Item>[] = useMemo(() => [
    { header: "ID", accessor: "id", sortable: false, hidden: true },
    { header: "Code", accessor: "code", sortable: true },
    { header: "Name", accessor: "name", sortable: true },
    {
      header: "Active", accessor: "is_active", sortable: false,
      formatter: (v: boolean) => (
        <span className={v ? "text-emerald-600 font-medium" : "text-red-600 font-medium"}>
          {v ? "Active" : "Inactive"}
        </span>
      ),
    },
  ], []);

  const actions: ActionConfig<{Name}Item> = {
    onSearch: () => setRefresh(!refresh),
    onCreate: openCreate,
    onUpdate: openUpdate,
    onDetail: openDetail,
  };

  const createForm = useForm<CreateForm>({
    resolver: zodResolver(createSchema),
    defaultValues: { code: "", name: "", is_active: true },
  });
  const updateForm = useForm<UpdateForm>({
    resolver: zodResolver(updateSchema),
    defaultValues: { code: "", name: "", is_active: true },
  });

  const openCreate = () => {
    createForm.reset({ code: "", name: "", is_active: true });
    setModal({ isOpen: true, mode: "create", title: "Create {Name}" });
  };

  const openUpdate = (row: {Name}Item) => {
    setSelectedRow(row);
    updateForm.reset({ code: row.code, name: row.name, is_active: row.is_active });
    setModal({ isOpen: true, mode: "update", title: "Update {Name}" });
  };

  const openDetail = (row: {Name}Item) => {
    setSelectedRow(row);
    setModal({ isOpen: true, mode: "detail", title: "{Name} Details" });
  };

  const handleCreate = async (values: CreateForm) => {
    setLoading(true);
    const session = parseSession(window.localStorage.getItem(storageKey));
    if (!session) { setLoading(false); return; }
    try {
      await clientApi("/{KODE}", { method: "POST", body: values });
      toast.success("{Name} created successfully.");
      closeModal();
      setRefresh(!refresh);
    } catch (err) {
      if (err instanceof ClientApiError && err.status === 401) { setLoading(false); return; }
      toast.error(getErrorMessage(err instanceof ClientApiError ? err.code : "UNKNOWN_ERROR"));
    } finally {
      setLoading(false);
    }
  };

  const handleUpdate = async (values: UpdateForm) => {
    setLoading(true);
    if (!selectedRow) { setLoading(false); return; }
    const session = parseSession(window.localStorage.getItem(storageKey));
    if (!session) { setLoading(false); return; }
    try {
      await clientApi(`/{KODE}/${selectedRow.id}`, { method: "PUT", body: values });
      toast.success("{Name} updated successfully.");
      closeModal();
      setRefresh(!refresh);
    } catch (err) {
      if (err instanceof ClientApiError && err.status === 401) { setLoading(false); return; }
      toast.error(getErrorMessage(err instanceof ClientApiError ? err.code : "UNKNOWN_ERROR"));
    } finally {
      setLoading(false);
    }
  };

  const closeModal = () => {
    setModal({ isOpen: false, mode: "create", title: "" });
    setSelectedRow(null);
    createForm.reset();
    updateForm.reset();
  };

  return (
    <Suspense fallback={<p>Loading...</p>}>
      <div className="p-3 max-w-8xl mx-auto space-y-3">
        <DataTable
          fetchData={loadData}
          columns={columns}
          actions={actions}
          hideSearch={false}
          hideSelect={true}
          hidePaging={false}
          hideSort={false}
          hideColumnToggle={true}
          refreshTrigger={refresh}
        />

        {/* Main entity DataDialog — handle create, update, detail */}
        <DataDialog
          isOpen={modal.isOpen}
          mode={modal.mode}
          title={modal.title}
          onClose={closeModal}
          loading={loading}
          onSubmit={
            modal.mode === "create"
              ? createForm.handleSubmit(handleCreate)
              : modal.mode === "update"
                ? updateForm.handleSubmit(handleUpdate)
                : undefined
          }
        >
          {modal.mode === "detail" ? (
            <FieldGroup className="p-3">
              {/* Detail view — Tabs untuk sub-entity */}
              <Tabs defaultValue="information">
                <TabsList variant="line" className="mb-3">
                  <TabsTrigger value="information">Information</TabsTrigger>
                  <TabsTrigger value="sub-entity">Sub Entity</TabsTrigger>
                </TabsList>
                <TabsContent value="information">
                  <FieldGroup className="gap-3">
                    <Field className="gap-1.5 pt-2 px-2">
                      <FieldLabel>Code</FieldLabel>
                      <Input value={selectedRow?.code ?? ""} disabled readOnly />
                    </Field>
                    <Field className="gap-1.5 pt-2 px-2">
                      <FieldLabel>Name</FieldLabel>
                      <Input value={selectedRow?.name ?? ""} disabled readOnly />
                    </Field>
                  </FieldGroup>
                </TabsContent>
                <TabsContent value="sub-entity">
                  {/* Sub-entity DataTable — load via fetchData */}
                </TabsContent>
              </Tabs>
            </FieldGroup>
          ) : (
            <FieldSet>
              <FieldGroup className="gap-3 p-3">
                {/* Form fields — gunakan Controller untuk Combobox/Checkbox */}
                <Controller
                  control={createForm.control}
                  name="code"
                  render={({ field, fieldState }) => (
                    <Field className="gap-1.5 pt-2 px-2" data-invalid={fieldState.invalid}>
                      <FieldLabel>Code *</FieldLabel>
                      <Input {...field} maxLength={100} aria-invalid={fieldState.invalid} />
                      {fieldState.invalid && <FieldError errors={[fieldState.error]} />}
                    </Field>
                  )}
                />
                <Controller
                  control={createForm.control}
                  name="name"
                  render={({ field, fieldState }) => (
                    <Field className="gap-1.5 pt-2 px-2" data-invalid={fieldState.invalid}>
                      <FieldLabel>Name *</FieldLabel>
                      <Input {...field} maxLength={255} aria-invalid={fieldState.invalid} />
                      {fieldState.invalid && <FieldError errors={[fieldState.error]} />}
                    </Field>
                  )}
                />
                <Controller
                  control={createForm.control}
                  name="is_active"
                  render={({ field }) => (
                    <Field className="gap-1.5 pt-2 px-2">
                      <div className="flex items-center gap-2">
                        <Checkbox checked={field.value ?? true} onCheckedChange={field.onChange} />
                        <FieldLabel className="cursor-pointer">Active</FieldLabel>
                      </div>
                    </Field>
                  )}
                />
              </FieldGroup>
            </FieldSet>
          )}
        </DataDialog>
      </div>
    </Suspense>
  );
}
```

**Pola kunci**:
- `clientApi("/{KODE}", { params })` — path tanpa prefix `/rest/pages` (proxy handle)
- `clientApi(\`/{KODE}/${id}\`, { method: "PUT", body })` — untuk update/delete
- `toast.error(getErrorMessage(err.code))` — gunakan `getErrorMessage()` untuk user-friendly message
- `Controller` dari react-hook-form untuk form fields (bukan `register`)
- `useWatch` untuk dependent fields (field A mengubah field B)
- `DataTable` + `DataDialog` sebagai komponen utama
- `Tabs` dalam detail modal untuk sub-entity
- Zod schema terpisah untuk create dan update
- `hidePaging={false}` — semua DataTable WAJIB support pagination
- `useForm` WAJIB sertakan `defaultValues` agar semua field memiliki nilai awal
- `parseSession()` check sebelum API call untuk memastikan session valid
- `closeModal()` panggil di `handleCreate`/`handleUpdate` success (bukan inline `setModal`)
- Optional field gunakan `value={field.value ?? ""}` di Input agar tetap controlled
- **Loading state**: `useState(false)` → `setLoading(true)` diawal handler → `finally { setLoading(false) }` — WAJIB untuk semua write operation
- **DataDialog `loading` prop**:
  - Saat `loading=true`: submit & cancel button disabled, Spinner tampil di tombol submit
  - WAJIB di-pass ke `<DataDialog loading={loading}>` untuk mencegah double-submit
  - Handle early return: `setLoading(false)` SEBELUM `return` agar button tidak stuck disabled

---

## Checklist Lengkap

```
□ 1. Cek kode di README.md (privilege, route, HTTP method)
□ 2. Buat folder: obx_rest/skeleton/{KODE}/
□ 3. Buat template.go (DTOs + interfaces)
□ 4. Buat repository.go (Count, List, Create, Update, Delete + sub-entity)
□ 5. Buat usecase.go (validation + delegate)
□ 6. Buat handler.go (mechanic.Error pattern)
□ 7. Register DI + routes di backbone/routes.go
□ 8. Build test: go build ./...
□ 9. Buat folder: obx_site/src/app/board/pages/{KODE}/
□ 10. Buat page.tsx (DataTable + DataDialog + Tabs + Controller)
□ 11. Typecheck: npx tsc --noEmit
□ 12. Update README.md: status ⬜ → ✅
□ 13. Tambahkan "Backend Pattern Reference" section di skeleton/{KODE}/README.md
```

> **Penting**: Setiap rancangan modul (skeleton/{KODE}/README.md) WAJIB menyertakan section "Backend Pattern Reference" yang memuat:
> 1. DI wiring code snippet (NRepo→NCase→NHand + registrasi routes)
> 2. companyId extraction dari gin context
> 3. InternalError wrapping pattern di usecase
> 4. SQL pattern (validSort, rows.Err(), parameterized queries, COALESCE)
> 5. Response pattern (JSON format untuk List/Create/Update/Delete)
>
> Contoh lengkap lihat di `obx_rest/skeleton/IM/README.md` section 7.

---

## Helper Reference (mechanic package)

| Helper | Lokasi | Fungsi |
|--------|--------|--------|
| `mechanic.ValidationError(msg)` | `mechanic/helper.go` | Return `*AppError` (400) |
| `mechanic.NotFound(msg)` | `mechanic/helper.go` | Return `*AppError` (404) |
| `mechanic.Conflict(msg)` | `mechanic/helper.go` | Return `*AppError` (409) |
| `mechanic.InternalError(msg, err)` | `mechanic/helper.go` | Return `*AppError` (500) |
| `mechanic.Error(c, err)` | `mechanic/helper.go` | Kirim JSON error response |
| `mechanic.NullableString(v)` | `mechanic/typography.go` | Return `nil` jika empty, else trimmed string |
| `mechanic.CheckMeta(page, size)` | `mechanic/typography.go` | Validasi pagination, return `(page, size, error)` |
| `mechanic.BuildMeta(page, size, total)` | `mechanic/typography.go` | Build `GridMeta` dengan `total_pages` |

## DB Patterns

| Pattern | Contoh |
|---------|--------|
| **Primary key** | `id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::text` |
| **Timestamps** | `created_at TIMESTAMPTZ DEFAULT NOW()`, `updated_at TIMESTAMPTZ DEFAULT NOW()` |
| **Nullable string** | `mechanic.NullableString(value)` → Go side, `$N` → SQL side |
| **Soft delete** | `is_active BOOLEAN` atau `status TEXT` (jangan hard delete) |
| **Multi-tenant** | `company_id TEXT` FK ke `dat_company` |
| **Unique per tenant** | `@@unique([company_id, code])` di Prisma |
| **SQL alias** | Gunakan `c.name AS company_name` di SELECT agar bisa `ORDER BY company_name` |

## SQL Anti-Patterns (JANGAN lakukan)

```sql
-- ❌ SALAH: ORDER BY kolom yang tidak ada di SELECT/fROM
SELECT c.name FROM t1 JOIN t2 c ON ... ORDER BY company_name

-- ✅ BENAR: pakai SQL alias
SELECT c.name AS company_name FROM t1 JOIN t2 c ON ... ORDER BY company_name
```

```go
// ❌ SALAH: rows.Err() sebelum loop
rows, _ := r.db.QueryContext(ctx, query, args...)
if rows.Err() != nil { return nil, rows.Err() } // selalu nil di sini
defer rows.Close()
for rows.Next() { ... }

// ✅ BENAR: rows.Err() setelah loop
rows, _ := r.db.QueryContext(ctx, query, args...)
defer rows.Close()
for rows.Next() { ... }
if err := rows.Err(); err != nil { return nil, err }
```

## Frontend Import Reference

```typescript
// Core
import { clientApi, ClientApiError } from "@/lib/client-api";
import { getErrorMessage } from "@/lib/error-message";
import { storageKey, parseSession } from "@/lib/utility";
import { toast } from "sonner";

// Form
import { Controller, useForm, useWatch } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";

// UIX (selalu dari @/uix/)
import DataTable, { Column, ActionConfig } from "@/uix/datatable";
import DataDialog from "@/uix/datadialog";
import { Field, FieldError, FieldGroup, FieldLabel, FieldSet } from "@/uix/field";
import { Input } from "@/uix/input";
import { Textarea } from "@/uix/textarea";
import { Checkbox } from "@/uix/checkbox";
import { Badge } from "@/uix/badge";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/uix/tabs";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/uix/select";
import SearchSelect from "@/uix/search-select";
import { Alert, AlertTitle } from "@/uix/alert";
import { InputGroup, InputGroupAddon, InputGroupInput, InputGroupText } from "@/uix/input-group";
```
