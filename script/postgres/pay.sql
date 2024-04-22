
DROP TABLE  IF EXISTS pay;
create table pay (
    "id" VARCHAR(36) PRIMARY KEY,
	  "name" VARCHAR(100),
		"order" int,
    "pay_status" INT,
		"create_date" BIGINT,
		"update_date" BIGINT
);

