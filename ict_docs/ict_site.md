# Plan: Frontend Website (ict_site)

Dokumen ini berisi rencana pengembangan frontend OmniSight menggunakan **Next.js 16.2.10** (App Router) & **React 19.2.4**, **TypeScript 5** (Strict Mode), serta **Tailwind CSS 4** & **shadcn/ui**.

---

## Daftar Isi

1. [Arsitektur & Tech Stack](#1-arsitektur--tech-stack)
2. [Error Handling Frontend](#2-error-handling-frontend)
3. [Library di `src/lib`](#3-library-di-srclib)
4. [Komponen UIX di `src/uix`](#4-komponen-uix-di-srcuix)
5. [Frontend Pages per Modul](#5-frontend-pages-per-modul)
6. [File yang Direncanakan](#6-file-yang-direncanakan)

---

## 1. Arsitektur & Tech Stack

### Tech Stack

| Package | Version | Status |
|---------|---------|--------|
| Next.js | 16.2.10 (App Router) | ✅ Digunakan |
| React | 19.2.4 | ✅ Digunakan |
| TypeScript | 5.x (Strict Mode) | ✅ Digunakan |
| Tailwind CSS | 4.x | ✅ Digunakan |
| shadcn/ui | @shadcn/react 0.2.0 (Base UI) | ✅ Digunakan |
| Zustand | 5.x | ✅ Digunakan |
| Zod | 4.x | ✅ Digunakan |
| React Hook Form | 7.80.0 | ✅ Digunakan |
| Sonner | 2.x (toast) | ✅ Digunakan |
| TanStack Query | 5.x | ✅ Digunakan |
| TanStack Infinite | latest | ✅ Digunakan |
| class-variance-authority | latest | ✅ Digunakan |
| clsx + tailwind-merge | latest | ✅ Digunakan |

### Backend Connection

| Konfigurasi | Nilai |
|-------------|-------|
| API Base | `/rest/pages` (via proxy `BE_POOL`) |
| Port dev | 36666 (`next dev -p 36666`) |
| Backend | `http://ict_rest:36665` |

### Struktur Folder

```
ict_site/
├── src/
│   ├── app/                    App Router pages & layouts
│   │   ├── layout.tsx          Root layout
│   │   ├── page.tsx            Root redirect
│   │   ├── login/              Login page
│   │   ├── (board)/            Board layout group
│   │   │   ├── board/
│   │   │   │   ├── layout.tsx  Board layout (sidebar + auth)
│   │   │   │   ├── sidebar.tsx Sidebar navigation
│   │   │   │   └── pages/      Page modules (SP01, SM01, etc.)
│   │   └── proxy/              API proxy routes
│   │       ├── pages/[...path]/route.ts   Authenticated proxy
│   │       └── guest/[...path]/route.ts   Guest proxy
│   ├── lib/                    Utility functions & hooks
│   ├── uix/                    Reusable UI components
│   ├── lib/theme.tsx           Theme system (Zustand)
│   └── app/theme/              Theme utilities & providers
├── package.json
├── next.config.ts
└── tsconfig.json
```

---

## 2. Error Handling Frontend

> Status: **SELESAI**

### clientApi — Satu-Satunya Cara Fetch API

Seluruh fetch ke backend WAJIB menggunakan `clientApi()` dari `src/lib/client-api.ts`:

```typescript
import { clientApi, ClientApiError } from "@/lib/client-api";

try {
  const data = await clientApi<DataType>("/SM01", { method: "GET" });
} catch (err) {
  if (err instanceof ClientApiError) {
    toast.error(getErrorMessage(err.code));
  }
}
```

### Toast — Satu-Satunya Cara Menampilkan Error ke User

```typescript
import { toast } from "sonner";

toast.error("Username already exists");
toast.success("Data saved successfully");
```

### Error Code → User Message

Gunakan `getErrorMessage(code)` dari `src/lib/error-message.ts`:

```typescript
const messages: Record<string, string> = {
  VALIDATION_ERROR:       "Please check your input and try again.",
  UNAUTHORIZED:           "Your session has expired. Please log in again.",
  FORBIDDEN:              "You don't have permission to perform this action.",
  NOT_FOUND:              "The requested data was not found.",
  CONFLICT:               "This data already exists. Please use a different value.",
  BAD_REQUEST:            "The request is invalid. Please check your input.",
  TIMEOUT:                "The request timed out. Please try again.",
  RATE_LIMITED:           "Too many requests. Please wait a moment and try again.",
  SERVER_ERROR:           "The server encountered an error. Please try again later.",
  EXTERNAL_SERVICE_ERROR: "An external service is temporarily unavailable.",
  INTERNAL_ERROR:         "An unexpected error occurred. Please try again.",
  NETWORK_ERROR:          "Failed to connect to the server.",
};
```

### ErrorBoundary — Untuk Unhandled Runtime Error

Pasang `<ErrorBoundary>` di `layout.tsx` (root) dan `board/layout.tsx`. Jangan pasang di setiap halaman individual.

### Silent Catch Dilarang

```typescript
// ❌ SALAH — error ditelan
fetchData().catch(() => {});

// ✅ BENAR — minimal log atau tampilkan toast
clientApi("/SM01").catch(() => {
  setData([]);
});
```

### Pola Konsisten Error Display

| Tipe Error | Display | Contoh |
|------------|---------|--------|
| Form validation error | `FieldError` component (inline) | Required field kosong |
| Mutation error (CRUD) | `toast.error(message)` | Username already exists |
| Session expired (401) | Auto-redirect ke `/login` | Token expired |
| Forbidden (403) | `toast.error("You don't have permission")` | Non-admin akses admin |
| Network error | `toast.error("Failed to connect to server")` | Backend down |
| React crash | ErrorBoundary fallback UI | Unexpected rendering error |

---

## 3. Library di `src/lib`

> **17 file** tersedia di `src/lib/`.

### API & Backend

| File | Fungsi |
|------|--------|
| `api.ts` | `ApiError` (class), `apiFetch<T>()` — direct fetch ke backend via `BE_POOL` |
| `backend.ts` | `BE_POOL` — backend pool URL, `WS_CONF` — site config (name, version, copyright, meta) |
| `client-api.ts` | `ClientApiError` (class), `clientApi<T>()` — client-side fetch via proxy `/proxy/pages`, auto-attach Bearer token, session refresh on 401 |
| `server-api.ts` | `ServerApiError` (class), `serverApi<T>()` — server-side fetch langsung ke `BE_POOL` (untuk Server Components/Actions) |
| `error-message.ts` | `getErrorMessage(code, fallback?)` — maps 12 error codes ke human-readable messages |

### Storage & Helpers

| File | Fungsi |
|------|--------|
| `storage.ts` | `storage.get/set/remove/has/clear()` — typed localStorage wrapper; `getPreference<T>()` — validated preference reader |
| `utility.ts` | `cn()` — Tailwind class merger (clsx + twMerge); `storageKey` — cookie/localStorage key; `parseSession()`, `isSessionExpired()`, `forceLogout()`, `getInitials()`, `formatCurrency()`, `formatDateTime()`, `streamToResponse()`, `getProxyHeaders()`, `handleGlobalError()` |

### Hooks

| File | Fungsi |
|------|--------|
| `use-debounce.ts` | `useDebounce<T>(value, delay)` — debounce value; `useDebouncedCallback()` — debounce callback |
| `use-event-listener.ts` | `useEventListener<T, K>(eventName, handler, element?)` — generic DOM event listener |
| `use-isomorphic-layout-effect.ts` | `useIsomorphicLayoutEffect` — SSR-safe `useLayoutEffect` |
| `use-local-storage.ts` | `useLocalStorage<T>(key, initialValue)` — typed localStorage dengan sync antar tab |
| `use-media-query.ts` | `useMediaQuery(query)` — SSR-safe media query match via `useSyncExternalStore` |
| `use-previous.ts` | `usePrevious<T>(value)` — return previous render value |
| `use-update-effect.ts` | `useUpdateEffect(effect, deps)` — skip first render |
| `use-window-size.ts` | `useWindowSize()` — returns `{ width, height }` via `useSyncExternalStore` |
| `view-mobile.ts` | `useIsMobile()` — `matchMedia("(max-width: 768px)")` |
| `view-normal.ts` | `useIsNormal()` — `matchMedia("(min-width: 1024px)")` |

---

## 4. Komponen UIX di `src/uix`

> **84 komponen** tersedia di `src/uix/` (flat directory, tanpa subdirectory).

### Layout & Structure

| Komponen | Fungsi |
|----------|--------|
| `Sidebar` (+ 20+ sub-components) | Sistem navigasi sidebar lengkap dengan mobile responsive, keyboard shortcut (Ctrl+B), cookie persistence |
| `ScrollArea` | Area scroll dengan custom scrollbar |
| `ResizablePanelGroup/Panel/Handle` | Layout panel yang resizable |
| `Separator` | Garis pemisah visual |
| `Skeleton` | Placeholder loading skeleton |

### Core Forms & Inputs

| Komponen | Fungsi |
|----------|--------|
| `Button` | Tombol dengan 6 variant + 8 size |
| `Input` | Text input (controlled/uncontrolled) |
| `InputGroup` (+ sub-components) | Compound input dengan addons, buttons, text |
| `Textarea` | Auto-sizing textarea |
| `InputOTP` | One-time password input |
| `Select` | Dropdown select dengan search, scroll, grouping |
| `NativeSelect` | Native `<select>` wrapper |
| `Combobox` (+ sub-components) | Searchable combobox dengan single/multi-select (chips) |
| `Checkbox` | Toggle checkbox |
| `RadioGroup` | Radio button group |
| `Switch` | Toggle switch |
| `Slider` | Range slider |
| `Label` | Styled `<label>` |
| `Field` (+ sub-components) | Form field layout dengan label, description, error display. **FieldError untuk validasi zod** |
| `Toggle` | Pressable toggle button |
| `ToggleGroup` | Group toggle buttons |
| `ButtonGroup` (+ sub-components) | Grouped buttons dengan separator dan text label |

### Dialogs, Modals & Overlays

| Komponen | Fungsi |
|----------|--------|
| `Dialog` (+ sub-components) | Base modal dialog |
| `DataDialog` | **CRUD dialog wrapper** — create/update/delete/detail dengan auto-colored submit buttons |
| `AlertDialog` | Confirmation dialog untuk destructive actions |
| `Sheet` | Slide-in panel dari edge |
| `Drawer` | Bottom/side drawer (vaul library) |
| `Popover` | Hover/click popover |
| `HoverCard` | Hover-triggered preview card |
| `Tooltip` | Simple tooltip dengan arrow |

### Menus & Navigation

| Komponen | Fungsi |
|----------|--------|
| `DropdownMenu` (+ sub-components) | Full-featured dropdown menu dengan submenus |
| `ContextMenu` | Right-click context menu |
| `Menubar` | Application menu bar |
| `NavigationMenu` | Mega-menu navigation |
| `Breadcrumb` | Breadcrumb navigation trail |
| `Tabs` (+ sub-components) | Tabbed content (horizontal/vertical) |
| `Command` | Command palette / quick search (cmdk) |

### Data Display & Tables

| Komponen | Fungsi |
|----------|--------|
| `Table` (+ sub-components) | Base table primitives |
| `DataTable` | **Main data table** — server-side paginated, sorting, searching, column visibility, row selection, action buttons, page size selector |
| `Pagination` | Pagination controls |

### Cards & Containers

| Komponen | Fungsi |
|----------|--------|
| `Card` (+ sub-components) | General-purpose card container |
| `StatCard` | Dashboard statistics card dengan trend indicator |
| `AspectRatio` | Maintains aspect ratio |

### Data Table Complex

| Komponen | Fungsi |
|----------|--------|
| `Accordion` | Collapsible accordion sections |
| `Collapsible` | Simple collapsible content |
| `Carousel` | Image/content carousel (embla) |
| `KanbanBoard` | Drag-and-drop Kanban board |
| `Heatmap` | Grid heatmap dengan configurable color scales (severity/risk) |

### Media & Uploads

| Komponen | Fungsi |
|----------|--------|
| `Avatar` (+ sub-components) | User avatar dengan image, fallback, status badge, grouped display |
| `Attachment` (+ sub-components) | File attachment card dengan upload states |
| `FileUpload` | Drag-and-drop file upload dengan preview, validation, progress |
| `PhotoUpload` | Image-only upload dengan gallery/camera, caption editing |

### Messaging & Chat

| Komponen | Fungsi |
|----------|--------|
| `MessageGroup/Message` | Chat message layout |
| `BubbleGroup/Bubble` | Chat bubble dengan multiple visual variants |
| `MessageScroller` (+ sub-components) | Auto-scrolling message container |

### Navigation & List

| Komponen | Fungsi |
|----------|--------|
| `Item` (+ sub-components) | Generic list item dengan media, content, actions |

### Feedback & Notification

| Komponen | Fungsi |
|----------|--------|
| `Alert` | Inline alert/notification banner (6 variants) |
| `Badge` | Small label/badge (6 variants) |
| `StatusBadge` | Status indicator dengan colored dot (8 status + 5 severity) |
| `Toaster` | Toast notification provider (sonner) |
| `Spinner` | Loading spinner |
| `Progress` | Progress bar dengan label dan value |

### Empty States & Error

| Komponen | Fungsi |
|----------|--------|
| `Empty` | Empty state placeholder dengan icon, title, description, action |
| `ErrorBoundary` | React error boundary dengan fallback UI |

### Data Visualization

| Komponen | Fungsi |
|----------|--------|
| `ChartContainer/Tooltip/Legend` | Chart wrapper (recharts) dengan themed tooltips/legends |
| `Timeline` | Vertical/horizontal timeline dengan status-colored dots |
| `ApprovalSteps` | Approval workflow visualization |

### Specialized Widgets

| Komponen | Fungsi |
|----------|--------|
| `Calendar` | Date picker calendar (react-day-picker) |
| `EventCalendar` | Monthly calendar dengan event display |
| `LogViewer` | Terminal-style log viewer dengan search, level filtering, download |
| `MarkdownPreview` | Lightweight markdown-to-React renderer |
| `TreeView` | Recursive tree view dengan expand/collapse, selection |
| `Marker` (+ sub-components) | Map marker dengan icon, content, dan visual variants |

### Location Components

| Komponen | Fungsi |
|----------|--------|
| `LocationPicker` | Dropdown lokasi dengan hierarki dan indentasi |
| `LocationBadge` | Badge tipe lokasi (country, city, room, rack, dll) |
| `LocationTree` | Tree view hierarki lokasi dengan code, type badge, contact count |

### Utility Components

| Komponen | Fungsi |
|----------|--------|
| `DirectionProvider` | RTL/LTR direction provider |
| `Kbd` | Keyboard shortcut display |
| `ConfirmDialog` | Confirmation dialog wrapper (default/destructive) |
| `CopyButton` | Copy-to-clipboard button dengan feedback |

### Aturan Komponen

1. **Gunakan komponen yang sudah ada** di `src/uix/` — jangan buat elemen HTML mentah
2. **DataTable** untuk menampilkan data dalam tabel + detail item
3. **DataDialog** sebagai wadah untuk action detail, create, update
4. **Field + FieldError** untuk menampilkan error validasi form
5. **zod** untuk validasi semua form
6. Component select wajib menampilkan **nama** (label), bukan ID
7. Semua komponen menggunakan `data-slot="..."` attribute untuk CSS targeting
8. Semua visual variants dikelola via `class-variance-authority` (cva)

---

## 5. Frontend Pages per Modul

### SP — Session Profile

| Route | Komponen | Deskripsi | Status |
|-------|----------|-----------|--------|
| `/board/pages/SP01` | Profile page | Profile user | ✅ |
| `/board/pages/SP02` | Change password form | Ubah password | ✅ |
| `/board/pages/SP03` | Action history | Riwayat aktivitas | ✅ |

### SM — System Manager (Admin Only)

| Route | Komponen | Deskripsi | Status |
|-------|----------|-----------|--------|
| `/board/pages/SM01` | DataTable + DataDialog | User Management + Company + Privilege + Location Tab | ✅ |
| `/board/pages/SM02` | DataTable + DataDialog | Module Management | ✅ |
| `/board/pages/SM03` | DataTable + DataDialog | Company Management + Module + Location Area Tab | ✅ |
| `/board/pages/SM04` | DataTable + DataDialog | Signature/Approval Type + Steps + Signs | ✅ |
| `/board/pages/SM05` | DataTable + DataDialog | Session Management (list, detail, revoke) | ✅ |

### Modul yang Direncanakan

| Kode | Modul | Keterangan | Status |
|------|-------|------------|--------|
| LM | Location Management | Manajemen lokasi hierarki | ⬜ Direncanakan |
| AM | Asset Management | Manajemen aset IT | ⬜ Direncanakan |
| NK | Network Management | VLAN, subnet, route, NAT | ⬜ Direncanakan |
| MK | Mikrotik Device Management | Device, hotspot, DHCP, firewall, queues, port/VLAN groups, clients, access | ⬜ Direncanakan |
| DN | Domain Management | Domain & DNS records | ⬜ Direncanakan |
| IP | IP Public & Block Management | IP blocks, allocations, NAT mappings | ⬜ Direncanakan |
| VL | Vulnerability Management | Scanner, scan jobs, findings | ⬜ Direncanakan |
| IN | Incident Management | Insiden & response | ⬜ Direncanakan |
| SI | SIEM / Log Monitoring | Log sources, rules, alerts | ⬜ Direncanakan |
| RK | Risk Assessment | Risiko, assessment, treatment | ⬜ Direncanakan |
| CP | Compliance Monitoring | ISO 27001, kontrol, assessment | ⬜ Direncanakan |
| CA | CAPA | Corrective and Preventive Action | ⬜ Direncanakan |
| IA | Internal Audit | Program & instance audit | ⬜ Direncanakan |
| MR | Management Review | Review meeting & keputusan | ⬜ Direncanakan |
| QM | Quality Management System | QMS, KPI, supplier, pelatihan | ⬜ Direncanakan |
| CT | Certificate Management | Sertifikat SSL/TLS | ⬜ Direncanakan |
| DK | Docker Management | Host, container, compose, image, network, volume | ⬜ Direncanakan |
| DM | Document Management | Dokumen, editor, todos, approvals | ⬜ Direncanakan |
| SR | Service Request | Katalog layanan, request, approval | ⬜ Direncanakan |
| BL | Billing & Invoice | Tagihan, vendor, kategori | ⬜ Direncanakan |
| PM | Preventive Maintenance | Ruangan, perangkat, template, kalender | ⬜ Direncanakan |
| DB | Dashboard & Reporting | Overview, widget, report | ⬜ Direncanakan |
| WM | Web Monitoring | UptimeRobot monitors, incidents, SLA | ⬜ Direncanakan |
| WS | Web Security | Nginx logs, WAF, IP whitelist/blacklist | ⬜ Direncanakan |

### Pola Halaman CRUD

Setiap halaman CRUD mengikuti pola:
1. `DataTable` untuk menampilkan data dengan sorting, searching, pagination
2. `DataDialog` untuk action create, update, delete, detail
3. `clientApi()` untuk fetch data dari backend
4. `toast` dari sonner untuk notifikasi
5. `Field + FieldError` untuk form validation dengan zod

---

## 6. File yang Direncanakan

### Frontend — File yang Sudah Ada

| File | Keterangan |
|------|------------|
| `ict_site/src/lib/api.ts` | `ApiError`, `apiFetch` — direct backend fetch |
| `ict_site/src/lib/backend.ts` | `BE_POOL`, `WS_CONF` — backend URL & site config |
| `ict_site/src/lib/client-api.ts` | `ClientApiError`, `clientApi` — client-side fetch via proxy |
| `ict_site/src/lib/server-api.ts` | `ServerApiError`, `serverApi` — server-side fetch |
| `ict_site/src/lib/error-message.ts` | Error code → user message mapping (12 codes) |
| `ict_site/src/lib/storage.ts` | Typed localStorage wrapper + `getPreference` |
| `ict_site/src/lib/utility.ts` | `cn`, `storageKey`, `parseSession`, `formatCurrency`, `formatDateTime`, proxy helpers |
| `ict_site/src/lib/use-debounce.ts` | Debounce hook |
| `ict_site/src/lib/use-event-listener.ts` | Generic DOM event listener hook |
| `ict_site/src/lib/use-isomorphic-layout-effect.ts` | SSR-safe useLayoutEffect |
| `ict_site/src/lib/use-local-storage.ts` | Typed localStorage hook |
| `ict_site/src/lib/use-media-query.ts` | Media query hook |
| `ict_site/src/lib/use-previous.ts` | Previous value hook |
| `ict_site/src/lib/use-update-effect.ts` | Skip-first-render effect hook |
| `ict_site/src/lib/use-window-size.ts` | Window size hook |
| `ict_site/src/lib/view-mobile.ts` | `useIsMobile()` hook |
| `ict_site/src/lib/view-normal.ts` | `useIsNormal()` hook |
| `ict_site/src/uix/button-group.tsx` | Grouped buttons dengan separator |
| `ict_site/src/uix/confirm-dialog.tsx` | Confirmation dialog wrapper |
| `ict_site/src/uix/copy-button.tsx` | Copy-to-clipboard button |
| `ict_site/src/uix/error-boundary.tsx` | ErrorBoundary class component |
| `ict_site/src/uix/file-upload.tsx` | Drag-and-drop file upload |
| `ict_site/src/uix/location-badge.tsx` | Location type badge |
| `ict_site/src/uix/location-picker.tsx` | Hierarchical location selector |
| `ict_site/src/uix/location-tree.tsx` | Location hierarchy tree |
| `ict_site/src/uix/marker.tsx` | Map marker dengan icon dan variants |
| `ict_site/src/app/layout.tsx` | Root layout dengan ErrorBoundary |
| `ict_site/src/app/(board)/board/layout.tsx` | Board layout dengan ErrorBoundary |
| `ict_site/src/app/login/page.tsx` | Login page |
| `ict_site/src/app/proxy/pages/[...path]/route.ts` | Authenticated proxy route |
| `ict_site/src/app/proxy/guest/[...path]/route.ts` | Guest proxy route |

### Frontend — Komponen yang Direncanakan

| Komponen | Modul | Deskripsi |
|----------|-------|-----------|
| `DockerStatusBadge` | DK | Badge warna status container |
| `ContainerStatsChart` | DK | Chart real-time CPU/Memory/Network |
| `ComposeEditor` | DK | Code editor untuk compose.yaml |
| `ContainerTerminal` | DK | Web terminal ke container |
| `ChecklistForm` | PM | Form checklist dinamis dari template |
| `RoomCard` | PM | Card info ruangan |
| `EquipmentBadge` | PM | Badge status perangkat |
| `BillStatusBadge` | BL | Badge status tagihan |
| `InvoiceUploader` | BL | Upload invoice drag & drop |
| `InvoicePreview` | BL | Preview PDF/Image |
