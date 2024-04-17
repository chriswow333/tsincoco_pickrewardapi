
DROP TABLE  IF EXISTS evaluation;
create table evaluation (
    "id" VARCHAR(36) PRIMARY KEY,
	  "name" VARCHAR(100),
    "feedback_type" INT,
		"create_date" BIGINT,
		"update_date" BIGINT
);

