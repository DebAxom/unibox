import { Render } from "nijor/core";
import "nijor/router";
import { Get, setAccessToken } from '@/fetch';

//@Routes()

async function init(){
    const res = await fetch("http://localhost:5000/auth/refresh", { method:"POST", credentials:"include" });
    const data = await res.json();
    setAccessToken(data.access_token);
}

window.refreshUser = async function () {
    try {
        await init();
        const data = await Get("/api/me");
        window.currentUser = data;
    } catch (error) {
        window.currentUser = null;
    }
};

(async()=>{
    try {
        await init();
        const data = await Get("/api/me");
        window.currentUser = data;
    } catch (error) {
        window.currentUser = null;
    }finally{
        let { pathname } = window.location;
        if ((pathname.startsWith("/dashboard") || pathname == "/auth/admin") && window.innerWidth < 768) {
            window.location.pathname = "/auth"
        }
        await Render(document.getElementById('app'));
    }
})();