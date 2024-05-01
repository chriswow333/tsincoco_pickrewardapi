
DROP TABLE  IF EXISTS evaluation;
create table evaluation (
    "id" VARCHAR(36) PRIMARY KEY,
		"owner_id" varchar(36),
		"start_date" integer,
		"end_date" integer,
    "feedback_id" varchar(36),
		"payload" jsonb,
		"create_date" BIGINT,
		"update_date" BIGINT
);


