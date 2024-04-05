

DROP TABLE  IF EXISTS card_reward;
create table card_reward (
    "id" VARCHAR(36) PRIMARY KEY,
	  "card_id" VARCHAR(100),
	  "name"  VARCHAR(100),
		"description" JSON,
		"create_date" BIGINT,
		"update_date" BIGINT,
		"start_date" BIGINT,
		"end_date" BIGINT,
		"currency" INT,
		"reward_type" INT,
		"order" INT,
);