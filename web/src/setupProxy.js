const createProxyMiddleware = require("http-proxy-middleware");

module.exports = function (app) {
  app.use(
    "/api",
    createProxyMiddleware({
      target: "http://localhost:8000",
    })
    // createProxyMiddleware({
    //   target: "https://secure.ustracers.com",
    //   pathRewrite: {
    //     "^/api": "",
    //   },
    //   secure: false,
    // })
  );
};
