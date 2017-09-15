alter table container_repo change column namebad name varchar(255);
alter table container_image change column repo_idbad repo_id int;
alter table container_image change column protobad proto varchar(255);
