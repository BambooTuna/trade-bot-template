DROP SCHEMA IF EXISTS sample_database;

CREATE SCHEMA sample_database;

USE sample_database;

CREATE TABLE `sample_table` (
  `id` varchar(255) NOT NULL PRIMARY KEY
);

INSERT INTO `sample_table` (`id`) VALUES ("f0c28384-3aa4-3f87-9fba-66a0aa62c504");
