DROP TABLE IF EXISTS `xsql`;
CREATE TABLE `xsql` (
                        `id` int unsigned NOT NULL AUTO_INCREMENT,
                        `foo` varchar(255) DEFAULT NULL,
                        `bar` datetime DEFAULT NULL,
                        `bool` int NOT NULL DEFAULT '0',
                        `enum` int NOT NULL DEFAULT '0',
                        `json` json DEFAULT NULL,
                        PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
INSERT INTO `xsql` (`id`, `foo`, `bar`, `bool`, `enum`, `json`) VALUES (1, 'v', '2022-04-12 23:50:00', 1, 1, '{"foo":"bar"}');
INSERT INTO `xsql` (`id`, `foo`, `bar`, `bool`, `enum`, `json`) VALUES (2, 'v1', '2022-04-13 23:50:00', 1, 1, '[1,2]');
INSERT INTO `xsql` (`id`, `foo`, `bar`, `bool`, `enum`, `json`) VALUES (3, 'v2', '2022-04-14 23:50:00', 1, 1, null);