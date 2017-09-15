alter table container_repo change column name namebad varchar(255);
alter table container_image change column repo_id repo_idbad int;
alter table container_image change column proto protobad varchar(255);
