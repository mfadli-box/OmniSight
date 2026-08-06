# dat_user - User Domain

## Ringkasan

- File schema: `obx_base/prisma/schema/dat_user.prisma`
- Isi utama: user, company access, session, privilege, action log

## Model Inti

- dat_user
- dat_user_company
- dat_user_area
- dat_user_privilege
- dat_user_session
- dat_user_action

## Catatan Teknis

- User menjadi basis auth dan audit
- Session dan action log dipakai backend dan frontend

## Checklist

- [ ] Session relation sinkron
- [ ] Audit action tetap tercatat
