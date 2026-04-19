import { Render } from "nijor/core";
import "nijor/router";
import { Get, setAccessToken } from '@/fetch';

//@Routes()

async function init(){
    const res = await fetch("http://localhost:5000/auth/refresh", { method:"POST", credentials:"include" });
    const data = await res.json();
    setAccessToken(data.access_token);
}

(async()=>{
    try {
        await init();
        const data = await Get("/api/");
        window.currentUser = data;
    } catch (error) {
        if(window.location.pathname.startsWith("/app") || window.location.pathname.startsWith("/dashboard")){
            window.location.href = "/auth";
        }
        window.currentUser = null;
    }finally{
        let { pathname } = window.location;
        if ((pathname.startsWith("/dashboard") || pathname == "/auth/admin") && window.innerWidth < 768) {
            window.location.pathname = "/auth"
        }
        await Render(document.getElementById('app'));
    }
})();