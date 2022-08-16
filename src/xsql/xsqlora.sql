-- ----------------------------
-- Table structure for XSQL
-- ----------------------------
DROP TABLE "TEST"."XSQL";
CREATE TABLE "TEST"."XSQL" (
                             "ID" NUMBER VISIBLE NOT NULL,
                             "FOO" VARCHAR2(255 BYTE) VISIBLE,
                             "BAR" TIMESTAMP(6) VISIBLE
)
    LOGGING
NOCOMPRESS
PCTFREE 10
INITRANS 1
STORAGE (
  INITIAL 1048576 
  NEXT 1048576 
  MINEXTENTS 1
  MAXEXTENTS 2147483645
  BUFFER_POOL DEFAULT
)
PARALLEL 1
NOCACHE
DISABLE ROW MOVEMENT
;

-- ----------------------------
-- Records of XSQL
-- ----------------------------
INSERT INTO "TEST"."XSQL" ("ID", "FOO", "BAR") VALUES ('1', 'v', TO_TIMESTAMP('2022-04-14 23:49:48.000000', 'SYYYY-MM-DD HH24:MI:SS:FF6'));
INSERT INTO "TEST"."XSQL" ("ID", "FOO", "BAR") VALUES ('2', 'v1', TO_TIMESTAMP('2022-04-14 23:50:00.000000', 'SYYYY-MM-DD HH24:MI:SS:FF6'));
COMMIT;
COMMIT;

-- ----------------------------
-- Primary Key structure for table XSQL
-- ----------------------------
ALTER TABLE "TEST"."XSQL" ADD CONSTRAINT "SYS_C00214001" PRIMARY KEY ("ID");

-- ----------------------------
-- Checks structure for table XSQL
-- ----------------------------
ALTER TABLE "TEST"."XSQL" ADD CONSTRAINT "SYS_C00214000" CHECK ("ID" IS NOT NULL) NOT DEFERRABLE INITIALLY IMMEDIATE NORELY VALIDATE;
