import auth from "./middlewares/auth.js";

export const server = {
    port : 3000,
    live_reload : true,
}

export const build = {
    mode : "spa"
};

export const middlewares = [ auth ];

export const plugins = [];

export const headers = {};