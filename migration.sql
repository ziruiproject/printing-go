CREATE TABLE IF NOT EXISTS settings
(
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at      DATETIME,
    updated_at      DATETIME,
    deleted_at      DATETIME,
    device_id       TEXT,
    choosen_printer TEXT,
    websocket_url   TEXT,
    websocket_key   TEXT,
    channel_name    TEXT,
    store_name      TEXT,
    address         TEXT
);

CREATE INDEX IF NOT EXISTS idx_settings_deleted_at ON settings (deleted_at);