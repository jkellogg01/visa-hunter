create database if not exists visa_sponsors_db;
use visa_sponsors_db;

drop table if exists companies;
create table companies (
    org_id int auto_increment,
    org_name varchar(255) not null,
    city varchar(255),
    county varchar(255),
    job_type_rating varchar(255),
    visa_route varchar(255),
    primary key (org_id)
);