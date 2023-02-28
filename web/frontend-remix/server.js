let express = require("express");
let httpProxy = require("http-proxy");
let { createRequestHandler } = require("@remix-run/express");

require('dotenv').config()

async function main() {
    let port = process.env.PORT || "4000";

    let server = express();

    let proxy = httpProxy.createProxyServer();
    server.all("/api/*", (req, res) => {
        // Remove /api prefix
        req.url = req.url.replace("/api/", "/");
        proxy.web(req, res, { target: process.env.API_URL });
    });

    // Then, we need to server the static files on the public folder
    server.use(express.static("public", { immutable: false, maxAge: "1h" }));

    // Everything else goes to remix
    server.all("*", createRequestHandler({ build: require("./build") }));

    let host = "localhost";
    server.listen(Number(port), host, () => {
        console.log(`Ready on http://${host}:${port}`);
    });
}

main().catch((error) => {
    console.error(error);
    process.exit(1);
});