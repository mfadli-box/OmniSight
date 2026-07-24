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
    update: {},
    create: {
      username: "admin",
      email: adminEmail,
      password: hashedPassword,
      fullname: "administrator",
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
}

async function main() {
  await seedAdmin();
}

main()
  .catch((e) => {
    console.error(e);
    process.exit(1);
  })
  .finally(async () => {
    await prisma.$disconnect();
  });
