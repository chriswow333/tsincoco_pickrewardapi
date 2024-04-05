DROP TABLE  IF EXISTS bank;
create table bank (
    "id" VARCHAR(36) PRIMARY KEY,
	  "name" VARCHAR(100),
		"order" INT,
    "bank_status" INT,
		"create_date" BIGINT,
		"update_date" BIGINT
);

