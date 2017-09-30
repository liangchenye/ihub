

INSERT INTO container_repo(id, name, star, download_num, description) VALUES 
  (1, "first", 1, 1, "this is the first repo with one level name"), 
  (2, "second/second", 2, 2, "this is the second repo with two levels name"), 
  (3, "third/third", 3, 3, "this is the third repo with two levels name"), 
  (4, "fourth/fourth/fourth", 4, 4, "this is the fourth repo with three levels name"); 

INSERT INTO container_image(id, tag, size, repo_id, proto, proto_version) VALUES 
 (1, "v0.1", 0, 2, "docker", "v2"), 
 (2, "v0.2", 0, 2, "docker", "v2"),
 (3, "v0.1", 1024, 2, "oci", "v1"), 
 (4, "v0.2", 1024, 2, "oci", "v1"),
 (5, "v0.3", 1024, 3, "oci", "v1"), 
 (6, "v0.4", 1024, 4, "oci", "v1");

INSERT INTO pkg_repo(id, name, star, description) VALUES
 (1, "isula", 1, "this is isula repo"),
 (2, "euleros", 3, "this is euleros repo");

INSERT INTO pkg(id, name, size, repo_id, downloads, type) VALUES
 (1, "ostree", 100, 1, 100, "rpm"),
 (2, "runc", 10, 1, 0, "rpm"),
 (3, "ostree",  101, 2, 0, "rpm");
