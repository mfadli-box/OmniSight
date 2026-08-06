# ict_website - Website Infra Domain

## Ringkasan

- File schema: `obx_base/prisma/schema/ict_website.prisma`
- Isi utama: Nginx, SSL, uptime, whitelist/blacklist

## Model Inti

- ict_nginx_log
- ict_nginx_app
- ict_nginx_atc
- ict_nginx_atc_sum
- ict_nginx_sla
- ict_ip_whitelist
- ict_ip_blacklist
- ict_nginx_site
- ict_ssl_certificate
- ict_uptimerobot_log
- ict_uptimerobot_sla

## Catatan Teknis

- Dipakai oleh obx_auto agent_ws dan obx_rest/backend docs
- Menyimpan log, SLA, sertifikat, dan pengaturan site

## Checklist

- [ ] Log dan SLA sinkron
- [ ] Site dan SSL terdokumentasi
