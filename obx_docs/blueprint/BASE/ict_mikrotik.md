# ict_mikrotik - MikroTik Domain

## Ringkasan

- File schema: `obx_base/prisma/schema/ict_mikrotik.prisma`
- Isi utama: device, status, interface, firewall, address list, backup, log

## Model Inti

- ict_mikrotik_device
- ict_mikrotik_status
- ict_mikrotik_interface
- ict_mikrotik_firewall_rule
- ict_mikrotik_address_list
- ict_mikrotik_address_entry
- ict_mikrotik_backup
- ict_mikrotik_backup_file
- ict_mikrotik_log

## Catatan Teknis

- Dipakai langsung oleh agent_mx
- Menjadi sumber status dan log perangkat MikroTik

## Checklist

- [ ] Device dan log sinkron
