CREATE TABLE `user` (
  `id` bigint(20) unsigned NOT NULL,
  `name` varchar(255) NOT NULL,
  `operator_id` bigint(20) unsigned NOT NULL,
  `source` VARCHAR(255) NOT NULL,
  `created` datetime NOT NULL,
  `updated` datetime,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `uniq_name` (`name`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC;

INSERT INTO `user` (`id`, `name`, `operator_id`, `source`, `created`) VALUES 
(1, 'jd', 0, 'inner', '2022-06-23 00:00:00');
