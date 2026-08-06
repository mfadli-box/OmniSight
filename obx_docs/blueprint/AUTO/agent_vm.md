# agent_vm - VM & Infrastructure Agent

## Ringkasan

- Scope: agent Go untuk VM host, Docker, Nginx, SSL, firewall, GitOps, vuln scan, dan update
- Entry: `obx_auto/agent_vm/main.go`
- Backend API: `/robot/vm/*`

## Technical Notes

- Membaca env API_URL, KEY_ROBOT, IDA_HOST, IVL_POOL, IVL_HEARTBEAT, SOC_DOCKER, DIR_NGINX, DIR_LEGO, DIR_TRIVY
- Loop utama: heartbeatLoop, pollLoop, vmPollLoop, inventorySyncLoop, statsCollector
- Job dispatch menuju ComposeExecutor, ContainerExecutor, ImageExecutor, NginxManager, LegoSSL, FirewallManager, UpdateManager, GitOpsManager, TrivyScanner
- agent_vm mengirim heartbeat, inventori, stats, dan hasil job ke backend

## Checklist

- [ ] Env agent terisi
- [ ] Heartbeat aktif
- [ ] Job dispatch aktif
- [ ] Inventory sync aktif
