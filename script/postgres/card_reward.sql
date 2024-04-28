

DROP TABLE  IF EXISTS card_reward;
create table card_reward (
    "id" VARCHAR(36) PRIMARY KEY,
	  "card_id" VARCHAR(100),
	  "name"  VARCHAR(100),
		"description" JSON,
		"start_date" BIGINT,
		"end_date" BIGINT,
		"card_reward_type" INT,
		"feedback_type_id" VARCHAR(36),
		"task_labels" JSONB,
		"order" INT,
		"card_reward_status" INT,
		"create_date" BIGINT,
		"update_date" BIGINT
);