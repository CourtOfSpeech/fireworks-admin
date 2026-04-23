-- Modify "users" table to match schema
ALTER TABLE "users" RENAME COLUMN "real_name" TO "nickname";
ALTER TABLE "users" ALTER COLUMN "nickname" DROP NOT NULL;
ALTER TABLE "users" ADD COLUMN "avatar" character varying(500) NULL;
