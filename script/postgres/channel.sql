
DROP TABLE  IF EXISTS channel;
create table channel (
    "id" VARCHAR(36) PRIMARY KEY,
	  "name" VARCHAR(100),
		"create_date" BIGINT,
		"update_date" BIGINT,
		"channel_labels" JSONB, 
		"show_label" varchar(36),
		"order" INT,
		"channel_status" INT
)

