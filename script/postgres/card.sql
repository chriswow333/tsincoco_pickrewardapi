

DROP TABLE  IF EXISTS card;
create table card (
    "id" VARCHAR(36) PRIMARY KEY,
	  "name" VARCHAR(100),
	  "descriptions" JSON,
		"link_url" TEXT,
		"bank_id" VARCHAR(36),
		"order" INT,
		"card_status" INT,
		"create_date" BIGINT,
		"update_date" BIGINT
);
