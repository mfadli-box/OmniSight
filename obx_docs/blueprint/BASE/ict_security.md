# ict_security - Security Domain

## Ringkasan

- File schema: `obx_base/prisma/schema/ict_security.prisma`
- Isi utama: firewall, incident, vuln scan, FIM

## Model Inti

- ict_firewall_rule
- ict_firewall_zone
- ict_firewall_zone_rule
- ict_incident
- ict_incident_evidence
- ict_incident_timeline
- ict_incident_ioc
- ict_vuln_scan_schedule
- ict_vuln_scan
- ict_vuln_finding
- ict_fim_path
- ict_fim_snapshot
- ict_fim_alert

## Catatan Teknis

- Dipakai untuk keamanan host dan analitik incident
- Terkait pipeline scan dan alert

## Checklist

- [ ] Security model sinkron
