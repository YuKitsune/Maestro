
var maestroUsername = _getEnv('MAESTRO_DATABASE_USERNAME')
var maestroPassword = _getEnv('MAESTRO_DATABASE_PASSWORD')
var maestroDatabase = _getEnv('MAESTRO_DATABASE_DATABASE')

db = db.getSiblingDB(maestroDatabase)
db.createUser({
    user: maestroUsername,
    pwd: maestroPassword,
    roles: [{ role: "readWrite", db: maestroDatabase }]
});
