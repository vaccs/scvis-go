-- Add inputfiles table to scvis_local database
-- Same schema as cfiles table but with 'csource' renamed to 'data'

USE scvis_local;

CREATE TABLE `inputfiles` (
  `id` int NOT NULL AUTO_INCREMENT,
  `projectid` int NOT NULL,
  `filename` varchar(255) NOT NULL,
  `data` text NOT NULL,
  `created` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `fk_inputfiles_projectid` (`projectid`),
  CONSTRAINT `fk_inputfiles_projectid` FOREIGN KEY (`projectid`) REFERENCES `projects` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
