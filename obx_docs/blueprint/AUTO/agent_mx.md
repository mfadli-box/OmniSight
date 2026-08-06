# agent_mx - MikroTik Agent

## Ringkasan

- Scope: agent Go untuk MikroTik stats, firewall, address list, backup, dan syslog
- Entry: `obx_auto/agent_mx/main.go`
- Backend API: `/robot/mx/*`

## Technical Notes

- Membaca env API_URL, KEY_ROBOT, IDA_HOST, IVL_POOL, IVL_HEARTBEAT, IPP_SYSLOG, DIR_SYSLOG
- Loop utama: heartbeatLoop, pollLoop, statsLoop, syslogServer
- RouterOS client digunakan untuk sync status, interface, firewall, address list, backup
- Syslog ingest dikirim ke backend setelah dibuffer

## Checklist

- [ ] Syslog listener aktif
- [ ] Heartbeat aktif
- [ ] Stats terkirim
- [ ] Job sync berhasil
