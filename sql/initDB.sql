DROP DATABASE IF EXISTS traQuest;
CREATE DATABASE traQuest;
USE traQuest;

CREATE TABLE IF NOT EXISTS `users` (
  `id` char(36) NOT NULL UNIQUE,
  `name` varchar(20) NOT NULL,
  `score` int(8) NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NUll,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS `quests` (
  `id` char(36) NOT NULL UNIQUE,
  `number` int(8) NOT NULL,
  `title` varchar(40) NOT NULL,
  `description` varchar(100) DEFAULT '',
  `level` int(2) NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NUll,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS `tags` (
  `id` char(36) NOT NULL UNIQUE,
  `name` varchar(20) NOT NULL,
  `created_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS `users_quests` (
  `id` char(36) NOT NULL UNIQUE,
  `user_id` char(36) NOT NULL,
  `quest_id` char(36) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
