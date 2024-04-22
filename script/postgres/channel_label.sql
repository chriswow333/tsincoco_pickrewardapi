
DROP TABLE  IF EXISTS channel_label;
CREATE TABLE channel_label
(
    "label" varchar(36) PRIMARY KEY,
    "name" text,
    "show" INT,
    "order" INT
)



INSERT INTO channel_label ("label", "name", "show", "order") VALUES
('ecommerce', '網路購物', 1, 0),
('food', '美食', 1, 1),
('travel', '旅遊', 1, 2),
('transportation', '交通', 1, 3),
('oversea', '海外消費', 1, 4),
('streaming', '影音/串流', 1, 5),
('payFee', '繳費', 1, 6),
('mall', '百貨/影城', 1, 7),
('insurance', '保險', 1, 8),
('supermarket', '量販/超市', 1, 9),
('delivery', '外送', 1, 10),
('casual', '休閒/運動', 1, 11);