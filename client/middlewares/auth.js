import jwt from "jsonwebtoken";

export default function (req, res, next) {
    const refresh_token = req.cookies.get("refresh_token");
    let role = "";
    let verified = false;

    try {
        const payload = jwt.verify(refresh_token, process.env.JWT_SECRET);
        role = payload.role;
        verified = true;
    } catch (error) {
        verified = false;
    }

    const url = req.url.split("?")[0];

    // Restrict app/
    if (url.startsWith('/app') || url.startsWith('/modules/pages/app-') || url === "/modules/layouts/app.js") {
        if (!verified || role !== "user") {
            res.statusCode = 302;
            res.setHeader("Location", "/auth");
            return res.end();
        }
    }

    // Restrict dashboard/
    if (url.startsWith('/dashboard') || url.startsWith('/modules/pages/dashboard-') || url === "/modules/layouts/dashboard.js") {
        if (!verified || role !== "admin") {
            res.statusCode = 302;
            res.setHeader("Location", "/auth/admin");
            return res.end();
        }
    }

    next();
}