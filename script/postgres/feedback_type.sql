
DROP TABLE  IF EXISTS feedback_type;
create table feedback_type (
    "id" VARCHAR(36) PRIMARY KEY,
	  "name" VARCHAR(100),
    "feedback_type" INT,
		"create_date" BIGINT,
		"update_date" BIGINT,
);

