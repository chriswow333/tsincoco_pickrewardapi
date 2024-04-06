
DROP TABLE  IF EXISTS channel;
create table channel (
    "id" VARCHAR(36) PRIMARY KEY,
	  "name" VARCHAR(100),
    "link_url" TEXT,
    "channel_type" INT,
		"create_date" BIGINT,
		"update_date" BIGINT,
		"channel_labels" JSONB, 
		"order" INT,
		"channel_status" INT,
);

