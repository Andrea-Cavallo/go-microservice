db = db.getSiblingDB('myapp');

db.createCollection('users');

db.users.insertMany([
    { name: "Andrea Cabajo", email: "andrea.cabajo@cabajo.com" },
    { name: "Lucia Uzun", email: "lucia.uz@uzuz.com" }
]);
