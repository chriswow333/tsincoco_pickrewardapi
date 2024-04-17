
DROP TABLE  IF EXISTS pay;
create table pay (
    "id" VARCHAR(36) PRIMARY KEY,
	  "name" VARCHAR(100),
    "feedback_type" INT,
		"create_date" BIGINT,
		"update_date" BIGINT
);

