create database if not exists visa_sponsors_db;
use visa_sponsors_db;

drop table if exists organisation_job;
drop table if exists organisation;
drop table if exists job;

CREATE TABLE `organisation` (
  `id` integer PRIMARY KEY AUTO_INCREMENT,
  `name` varchar(255),
  `city` varchar(255),
  `county` varchar(255)
);

CREATE TABLE `job` (
  `id` integer PRIMARY KEY,
  `type` varchar(255),
  `rating` varchar(255),
  `visa_route` varchar(255)
);

CREATE TABLE `organisation_job` (
  `organisation_id` integer,
  `job_id` integer,
  PRIMARY KEY (`organisation_id`, `job_id`)
);

ALTER TABLE `organisation_job` ADD FOREIGN KEY (`organisation_id`) REFERENCES `organisation` (`id`);

ALTER TABLE `organisation_job` ADD FOREIGN KEY (`job_id`) REFERENCES `job` (`id`);

