

DROP TABLE  IF EXISTS task_label;
create table task_label (
    "label" VARCHAR(36) PRIMARY KEY,
	  "name" VARCHAR(100),
	  "show" JSON,
		"order" INT
)
