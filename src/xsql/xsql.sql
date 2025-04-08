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

DROP TABLE IF EXISTS `devices`;
CREATE TABLE `devices` (
                           `id` bigint unsigned NOT NULL AUTO_INCREMENT,
                           `uuid` varchar(255) NOT NULL,
                           `user_id` bigint unsigned NOT NULL,
                           `platform` tinyint unsigned NOT NULL,
                           `info` varchar(255) NOT NULL,
                           `app` json DEFAULT NULL,
                           `language_code` varchar(255) NOT NULL,
                           `status` tinyint unsigned NOT NULL,
                           `synced_message_id` bigint unsigned NOT NULL DEFAULT '0',
                           `firebase_token` varchar(255) NOT NULL DEFAULT '',
                           `last_used_at` timestamp NOT NULL,
                           `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
                           `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                           PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
INSERT INTO `devices` (`id`, `uuid`, `user_id`, `platform`, `info`, `app`, `language_code`, `status`, `synced_message_id`, `firebase_token`, `last_used_at`, `created_at`, `updated_at`) VALUES (1, '0c7f60e8-d7a3-48ae-9139-5128e336736e', 100000010, 1, 'postman', '{\"build\": 1, \"version\": \"v1.0.0\"}', '', 1, 0, '', '1970-01-01 00:00:01', '2025-04-07 07:41:23', '2025-04-07 07:41:23');
INSERT INTO `devices` (`id`, `uuid`, `user_id`, `platform`, `info`, `app`, `language_code`, `status`, `synced_message_id`, `firebase_token`, `last_used_at`, `created_at`, `updated_at`) VALUES (2, '0c7f60e8-d7a3-48ae-9139-5128e336736e', 100000011, 1, 'postman', '{\"build\": 1, \"version\": \"v1.0.0\"}', '', 1, 0, '', '1970-01-01 00:00:01', '2025-04-07 07:41:25', '2025-04-07 07:41:25');
INSERT INTO `devices` (`id`, `uuid`, `user_id`, `platform`, `info`, `app`, `language_code`, `status`, `synced_message_id`, `firebase_token`, `last_used_at`, `created_at`, `updated_at`) VALUES (3, '0c7f60e8-d7a3-48ae-9139-5128e336736e', 100000012, 1, 'postman', '{\"build\": 1, \"version\": \"v1.0.0\"}', '', 1, 0, '', '1970-01-01 00:00:01', '2025-04-07 07:41:26', '2025-04-07 07:41:26');
