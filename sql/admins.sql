CREATE TABLE IF NOT EXISTS admins (
  id TEXT PRIMARY KEY,
  username TEXT UNIQUE NOT NULL,
  password TEXT NOT NULL,
  dept TEXT NOT NULL,
  token_version INT NOT NULL DEFAULT 0,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Admins password : joiaiaxom
INSERT INTO admins (id, username, password, dept) VALUES ('x2d-23d-325-32s', 'dean.academic', '$2a$12$RAAbzbfK7RJ/1qCThd9bee9/gxDYfPolOAwYV4fj2Potggg.Qd1tm', 'academic');
INSERT INTO admins (id, username, password, dept) VALUES ('xic-cfg-klt-sde', 'registrar', '$2a$12$RAAbzbfK7RJ/1qCThd9bee9/gxDYfPolOAwYV4fj2Potggg.Qd1tm', 'accounts');
INSERT INTO admins (id, username, password, dept) VALUES ('sku-3er-4rf-wet', 'dean.sw', '$2a$12$RAAbzbfK7RJ/1qCThd9bee9/gxDYfPolOAwYV4fj2Potggg.Qd1tm', 'sw');
INSERT INTO admins (id, username, password, dept) VALUES ('ick-23l-cew-fed', 'warden.abh', '$2a$12$RAAbzbfK7RJ/1qCThd9bee9/gxDYfPolOAwYV4fj2Potggg.Qd1tm', 'hostel-abh');
INSERT INTO admins (id, username, password, dept) VALUES ('swa-ddc-zxc-vzt', 'supervisor.abh', '$2a$12$RAAbzbfK7RJ/1qCThd9bee9/gxDYfPolOAwYV4fj2Potggg.Qd1tm', 'hostel-abh');
INSERT INTO admins (id, username, password, dept) VALUES ('wdx-qwd-xcs-vgd', 'warden.bh8', '$2a$12$RAAbzbfK7RJ/1qCThd9bee9/gxDYfPolOAwYV4fj2Potggg.Qd1tm', 'hostel-bh8');
INSERT INTO admins (id, username, password, dept) VALUES ('sdx-cdf-zaq-23e', 'supervisor.bh8', '$2a$12$RAAbzbfK7RJ/1qCThd9bee9/gxDYfPolOAwYV4fj2Potggg.Qd1tm', 'hostel-bh8');
INSERT INTO admins (id, username, password, dept) VALUES ('wer-cfg-32s-0io', 'warden.gh1', '$2a$12$RAAbzbfK7RJ/1qCThd9bee9/gxDYfPolOAwYV4fj2Potggg.Qd1tm', 'hostel-gh1');
INSERT INTO admins (id, username, password, dept) VALUES ('n0p-edf-de0-sww', 'warden.gh3', '$2a$12$RAAbzbfK7RJ/1qCThd9bee9/gxDYfPolOAwYV4fj2Potggg.Qd1tm', 'hostel-gh3');
