# Capability Matrix - Platform Replacement

Matriks ini memetakan kapabilitas target per produk eksternal terhadap domain OmniSight.

## Legenda

- `Ada`: sudah ada di source code aktif
- `Rancang`: baru ada di blueprint / roadmap
- `Belum`: belum ada bukti implementasi atau rancangan detail

## Matriks Kapabilitas

| Domain Kapabilitas | Dockge / Beszel | nginx-ui | Wazuh | LibreNMS | JumpServer | OmniSight Saat Ini |
|---|---|---|---|---|---|---|
| Asset inventory | Ya | Sebagian | Sebagian | Ya | Ya | Rancang |
| Account / credential management | Tidak utama | Tidak | Tidak | Tidak | Ya | Rancang |
| Session audit | Tidak | Tidak | Sebagian | Tidak | Ya | Rancang |
| Web SSH | Tidak | Tidak | Tidak | Tidak | Ya | Rancang |
| Web RDP | Tidak | Tidak | Tidak | Tidak | Ya | Rancang |
| FTP / file browser | Tidak | Tidak | Tidak | Tidak | Ya | Rancang |
| WebAppProxy | Tidak | Tidak | Tidak | Tidak | Ya | Rancang |
| Compose / container deploy | Ya | Tidak | Tidak | Tidak | Tidak | Belum |
| Host/container metrics | Ya | Tidak | Sebagian | Sebagian | Tidak | Belum |
| Nginx vhost management | Tidak | Ya | Tidak | Tidak | Tidak | Belum |
| SSL / certificate management | Tidak | Ya | Tidak | Tidak | Tidak | Belum |
| Network device polling | Tidak | Tidak | Tidak | Ya | Tidak | Belum |
| Security event correlation | Tidak | Tidak | Ya | Tidak | Tidak | Belum |
| SOP / ISO controlled docs | Tidak | Tidak | Tidak | Tidak | Tidak | Rancang-ke-menengah |
| Approval workflow | Tidak | Tidak | Sebagian | Tidak | Ya | Rancang |
| Audit trail dokumen | Tidak | Tidak | Sebagian | Tidak | Sebagian | Rancang-ke-menengah |

## Mapping Modul OmniSight

| Kapabilitas | Target Modul OmniSight |
|---|---|
| SOP / ISO document control | dat_document, dat_signature, dat_request |
| Bastion access | JMS PLAN |
| VM / Docker / Nginx inventory | domain BASE + fase 3 roadmap |
| Monitoring MikroTik | ict_mikrotik |
| Security event | ict_security |

## Kesimpulan Cepat

- Kemenangan tercepat: dokumentasi SOP/ISO dan approval.
- Pengganti tool paling masuk akal tahap awal: bastion berbasis JMS.
- Pengganti domain monitoring dan SIEM masih butuh fase lanjutan yang signifikan.
