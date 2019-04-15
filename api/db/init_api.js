use asdasd;
db.createCollection("users");
db.users.createIndex({ "login": 1 }, { unique: true });
db.users.createIndex({ "email": 1 }, { unique: true });