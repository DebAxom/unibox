export default function (protectedRoutes = []) {

    return async function (req, res, next) {
        const refresh_token = req.cookies.get("refresh_token");

        const url = req.url.split("?")[0];
        const normalizedRoute = normalizeRoute(url);
        const modulePath = routeToModule(normalizedRoute);

        const isProtected = matchRoute(normalizedRoute, modulePath, protectedRoutes);

        if (isProtected && !refresh_token) {
            res.statusCode = 302;
            res.setHeader("Location", `/auth/login?redirect=${encodeURIComponent(req.url)}`);
            return res.end();
        }

        next();
    };
}

function normalizeRoute(url) {
    // remove trailing slash except root
    if (url.length > 1 && url.endsWith("/")) {
        url = url.slice(0, -1);
    }
    return url;
}

function routeToModule(route) {
    if (route === "/") return "/modules/pages/index.js";

    const name = route
        .slice(1)              // remove leading '/'
        .replace(/\//g, "-");  // convert subroutes
    return `/modules/pages/${name}.js`;
}

function matchRoute(route, modulePath, rules) {
    return rules.some(rule => {
        // wildcard: /dashboard/*
        if (rule.endsWith("/*")) {
            const base = rule.slice(0, -2);

            const baseModule = routeToModule(base); // exact file
            const baseModulePrefix = baseModule.replace(".js", "-");

            return (
                route === base ||                      // /dashboard
                route.startsWith(base + "/") ||        // /dashboard/*
                modulePath === baseModule ||           // dashboard.js
                modulePath.startsWith(baseModulePrefix) // dashboard-*.js
            );
        }

        // exact match
        return rule === route || rule === modulePath;
    });
}

function parseCookies(cookieHeader = "") {
    return cookieHeader.split(";").map(cookie => cookie.trim()).filter(Boolean).reduce((acc, cookie) => {
        const [key, ...v] = cookie.split("=");
        acc[key] = decodeURIComponent(v.join("="));
        return acc
    }, {});
}