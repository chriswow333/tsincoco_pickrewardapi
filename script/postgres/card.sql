
DROP TABLE  IF EXISTS bank;
create table bank (
    "id" VARCHAR(36) PRIMARY KEY,
	  "name" VARCHAR(100),
	  "image" TEXT,
		"order" INT
);


DROP TABLE  IF EXISTS card;
create table card (
    "id" VARCHAR(36) PRIMARY KEY,
	  "name" VARCHAR(100),
	  "descriptions" JSON,
		"image" TEXT,
		"create_date" BIGINT,
		"update_date" BIGINT,
		"link_url" TEXT,
		"bank_id" VARCHAR(36),
		"order" INT,
		"card_status" INT,
);
