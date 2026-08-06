import * as dotenv from "dotenv";
dotenv.config();

import { PrismaClient } from "@prisma/client";
import { PrismaPg } from "@prisma/adapter-pg";
import { Pool } from "pg";
import * as process from "process";
import * as bcrypt from "bcrypt";

const pool = new Pool({ connectionString: process.env.DATABASE_URL });
const adapter = new PrismaPg(pool);
const prisma = new PrismaClient({ adapter });

async function seedAdmin() {
  const adminEmail = process.env.AD_MAIL || "admin@localhost";
  const adminPassword = process.env.AD_PASS || "rahasia";
  const saltRounds = 10;
  const hashedPassword = await bcrypt.hash(adminPassword, saltRounds);
  const admin = await prisma.dat_user.upsert({
    where: {
      company_id_username: {
        company_id: "",
        username: "admin",
      },
    },
    update: {
      email: adminEmail,
      password: hashedPassword,
      fullname: "Administrator",
      is_admin: true,
      is_active: true,
      role: "staff",
      is_hris: false,
      updated_at: new Date(),
    },
    create: {
      username: "admin",
      email: adminEmail,
      password: hashedPassword,
      fullname: "Administrator",
      company_id: "",
      is_admin: true,
      is_active: true,
      phone: null,
      employee_id: null,
      location_id: null,
      department_id: null,
      division_id: null,
      role: "staff",
      job: "",
      is_hris: false,
    },
  });
  console.log("Admin user created:", admin.id);
  return admin.id;
}

type ModuleSeed = {
  code: string;
  name: string;
  path: string;
  is_page: boolean;
};

const moduleGroups: { group: ModuleSeed; pages: ModuleSeed[] }[] = [
  {
    group: { code: "VM", name: "Virtual Machine", path: "/board", is_page: false },
    pages: [
      { code: "VM01", name: "Host Groups", path: "/board/pages/VM01", is_page: true },
      { code: "VM02", name: "Host Permissions", path: "/board/pages/VM02", is_page: true },
      { code: "VM03", name: "Host List", path: "/board/pages/VM03", is_page: true },
      { code: "VM04", name: "Docker Compose", path: "/board/pages/VM04", is_page: true },
      { code: "VM05", name: "Docker Containers", path: "/board/pages/VM05", is_page: true },
      { code: "VM06", name: "Docker Images", path: "/board/pages/VM06", is_page: true },
      { code: "VM07", name: "Docker Networks", path: "/board/pages/VM07", is_page: true },
      { code: "VM08", name: "Deploy History", path: "/board/pages/VM08", is_page: true },
      { code: "VM09", name: "GitOps Webhook", path: "/board/pages/VM09", is_page: true },
    ],
  },
];

async function seedDatModule() {
  let totalPages = 0;
  for (const { group, pages } of moduleGroups) {
    const groupModule = await prisma.dat_module.upsert({
      where: { code: group.code },
      update: { name: group.name, path: group.path, is_page: group.is_page, is_active: true },
      create: { ...group, parent_id: null, is_active: true },
    });
    for (const page of pages) {
      await prisma.dat_module.upsert({
        where: { code: page.code },
        update: { name: page.name, path: page.path, is_page: page.is_page, parent_id: groupModule.id, is_active: true },
        create: { ...page, parent_id: groupModule.id, is_active: true },
      });
      totalPages++;
    }
  }
  console.log("Modules seeded:", moduleGroups.length, "groups,", totalPages, "pages");
}

async function main() {
  await seedAdmin();
  await seedDatModule();
  console.log("\n=== All seed data inserted successfully ===");
}

main()
  .catch((e) => {
    console.error(e);
    process.exit(1);
  })
  .finally(async () => {
    await prisma.$disconnect();
  });
