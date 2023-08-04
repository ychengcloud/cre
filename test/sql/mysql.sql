CREATE DATABASE IF NOT EXISTS `test`;

DROP TABLE IF EXISTS `numeric`;
CREATE TABLE `numeric` (
  `bigint` bigint NOT NULL AUTO_INCREMENT COMMENT '中文bigint comment',
  `bigint1` bigint(20) NOT NULL,
  `bit` bit NOT NULL,
  `int` int NOT NULL,
  `tinyint` tinyint NOT NULL,
  `smallint` smallint NOT NULL DEFAULT '10',
  `mediumint` mediumint NOT NULL,
  `decimal` decimal NOT NULL,
  `decimal1` decimal(10,2) NOT NULL,
  `numeric` numeric NOT NULL,
  `float` float NOT NULL,
  `float1` float(10,2) NOT NULL,
  `double` double NOT NULL,
  `double1` double(10,2) NOT NULL,
  `real` real NOT NULL,
  `real1` real(10,2) NOT NULL,
  PRIMARY KEY (`bigint`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

DROP TABLE IF EXISTS `string`;
CREATE TABLE `string` (
  `char` char NOT NULL,
  `char1` char(10) NOT NULL,
  `varchar` varchar(255) NOT NULL,
  `binary` binary NOT NULL,
  `varbinary` varbinary(255) NOT NULL,
  `tinyblob` tinyblob NOT NULL,
  `tinytext` tinytext NOT NULL,
  `blob` blob NOT NULL,
  `text` text NOT NULL,
  `mediumblob` mediumblob NOT NULL,
  `mediumtext` mediumtext NOT NULL,
  `longblob` longblob NOT NULL,
  `longtext` longtext NOT NULL,
  `enum` enum('a', 'b', 'c') NOT NULL,
  `set` set('a', 'b', 'c') NOT NULL,
  `json` json NOT NULL
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

DROP TABLE IF EXISTS `time`;
CREATE TABLE `time` (
  `date` date NOT NULL,
  `time` time NOT NULL,
  `timestamp` timestamp NOT NULL,
  `datetime` datetime NOT NULL,
  `year` year NOT NULL
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

DROP TABLE IF EXISTS `spatial`;
CREATE TABLE `spatial` (
  `geometry` geometry NOT NULL,
  `point` point NOT NULL,
  `multipoint` multipoint NOT NULL,
  `linestring` linestring NOT NULL,
  `multilinestring` multilinestring NOT NULL,
  `polygon` polygon NOT NULL
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

DROP TABLE IF EXISTS `fk1`;
CREATE TABLE `fk1` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `fkid` bigint(20) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_fkid` (`fkid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

DROP TABLE IF EXISTS `fk2`;
CREATE TABLE `fk2` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

ALTER TABLE `fk1` ADD FOREIGN KEY (`fkid`) REFERENCES `fk2` (`id`);