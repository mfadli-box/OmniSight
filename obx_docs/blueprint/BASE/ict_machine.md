# ict_machine - VM & Docker Domain

## Ringkasan

- File schema: `obx_base/prisma/schema/ict_machine.prisma`
- Isi utama: host group, VM host, Docker inventory, update, alert, gitops

## Model Inti

- ict_host_group
- ict_host_group_member
- ict_host_permission
- ict_vm_host
- ict_vm_stat
- ict_docker_compose
- ict_docker_container
- ict_docker_image
- ict_docker_network
- ict_docker_share
- ict_docker_backup
- ict_docker_deploy
- ict_push_device
- ict_alert_rule
- ict_alert_notif
- ict_alert_history
- ict_update_job
- ict_update_history
- ict_update_package
- ict_git_webhook
- ict_git_deploy_mapping
- ict_git_deploy_log

## Catatan Teknis

- Dipakai langsung oleh agent_vm
- Menjadi sumber data VM, Docker, alert, dan deployment

## Checklist

- [ ] Inventory dan job log sinkron
