let _accessToken = "";
let refreshPromise = null;

export const setAccessToken = (token) => {
    _accessToken = token;
};

export const getAccessToken = () => _accessToken;

export function parseCookies(cookieHeader = "") {
    return cookieHeader.split(";").map(c => c.trim()).filter(Boolean).reduce((acc, cookie) => {
        const [key, ...v] = cookie.split("=");
        acc[key] = decodeURIComponent(v.join("="));
        return acc;
    }, {});
}


export async function Fetch(endpoint, options = {}, retry = true) {
    const baseUrl = "http://localhost:5000";
    const url = `${baseUrl}${endpoint}`;

    const headers = {
        ...(options.headers || {}),
    };

    if (!(options.body instanceof FormData)) {
        headers["Content-Type"] = "application/json";
    }

    // Attach access token
    if (_accessToken) {
        headers["Authorization"] = `Bearer ${_accessToken}`;
    }

    let response = await fetch(url, {
        ...options,
        credentials: "include", // important for refresh cookie
        headers,
    });

    if (response.status === 401 && retry) {
        console.warn("Access token expired, attempting refresh...");

        if (!refreshPromise) {
            refreshPromise = fetch(`${baseUrl}/auth/refresh`, {
                method: "POST",
                credentials: "include",
            }).then(res => {
                if (!res.ok) throw new Error("Refresh failed");
                return res.json();
            })
                .then(data => {
                    _accessToken = data.access_token;
                    return _accessToken;
                })
                .catch(err => {
                    _accessToken = "";
                    throw err;
                })
                .finally(() => {
                    refreshPromise = null;
                });
        }

        try {
            const newToken = await refreshPromise;

            return Fetch(endpoint, {
                ...options,
                headers: {
                    ...(options.headers || {}),
                    Authorization: `Bearer ${newToken}`,
                },
            }, false);

        } catch (err) {
            console.error("Session expired. Redirecting to login.");

            _accessToken = "";

            // Optional: preserve redirect
            const redirect = encodeURIComponent(window.location.pathname);
            window.location.href = `/auth/login?redirect=${redirect}`;

            return response;
        }
    }

    return response;
}

export async function Get(endpoint) {
    const res = await Fetch(endpoint, { method: "GET" });
    return res.json();
}

export async function Post(endpoint, body = {}) {
    const res = await Fetch(endpoint, {
        method: "POST",
        body,
    });
    return res.json();
}

export async function Patch(endpoint, body = {}) {
    const res = await Fetch(endpoint, {
        method: "PATCH",
        body,
    });
    return res.json();
}

export async function Delete(endpoint) {
    const res = await Fetch(endpoint, { method: "DELETE" });
    return res.json();
}