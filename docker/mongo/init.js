print("create crud db");
crudDb = db.getSiblingDB("crud");

crudDb.createCollection("authors");
crudDb.authors.createIndex({"id": 1}, {"unique": true, "background": true});

crudDb.createCollection("posts");
crudDb.posts.createIndex({"id": 1}, {"unique": true, "background": true});
crudDb.posts.createIndex({"author_id": 1}, {"background": true});
