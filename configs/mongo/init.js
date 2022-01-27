
var maestroUsername = _getEnv('MONGO_MAESTRO_USERNAME')
var maestroPassword = _getEnv('MONGO_MAESTRO_PASSWORD')
var maestroDatabase = _getEnv('MONGO_MAESTRO_DATABASE')

db = db.getSiblingDB(maestroDatabase)
db.createUser({
    user: maestroUsername,
    pwd: maestroPassword,
    roles: [{ role: "readWrite", db: maestroDatabase }]
});
