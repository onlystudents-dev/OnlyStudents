CREATE TABLE themes (
    name VARCHAR(64) PRIMARY KEY,
    colors JSONB
);

INSERT INTO themes (name, colors)
VALUES (
    'dark',
    '{
     "bg-color": "#0c0910",
     "txt-color": "#eee5e9",
     "hover-color": "#262626",
     "card-color": "#3a3a3a",
     "border-color": "#2b303a",
     "wrong-color": "#ed474a",
     "wrong-base-color": "rgb(207 41 44 / 0.5)",
     "warning-color": "#ff8200",
     "warning-base-color": "rgb(255 130 0 / 0.5)",
     "right-color": "#21fa90"
    }'::jsonb
);