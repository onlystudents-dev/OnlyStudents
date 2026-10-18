CREATE TABLE themes (
    name VARCHAR(64) PRIMARY KEY,
    colors JSONB
);

INSERT INTO themes (name, colors)
VALUES
    (
        'dark',
        '{
          "bg-color": "#1a1b1c",
          "txt-color": "#eee5e9",
          "hover-color": "#262626",
          "card-color": "#3a3a3a",
          "border-color": "#2b303a",
          "wrong-color": "#ed474a",
          "wrong-base-color": "rgb(207 41 44 / 0.5)",
          "warning-color": "#ff8200",
          "warning-base-color": "rgb(255 130 0 / 0.5)",
          "right-color": "#21fa90",
          "right-base-color": "rgb(33 250 144 / 0.5)",
          "edit-color": "#007fff",
          "edit-base-color": "rgb(0 127 255 / 0.5)",
          "yellow-color": "#fec601",
          "yellow-base-color": "rgb(254 198 1 / 0.5)"
        }'::jsonb
    ),
    (
        'oled',
        '{
          "bg-color": "#000000",
          "txt-color": "#eee5e9",
          "hover-color": "#262626",
          "card-color": "#3a3a3a",
          "border-color": "#2b303a",
          "wrong-color": "#ed474a",
          "wrong-base-color": "rgb(207 41 44 / 0.5)",
          "warning-color": "#ff8200",
          "warning-base-color": "rgb(255 130 0 / 0.5)",
          "right-color": "#21fa90",
          "right-base-color": "rgb(33 250 144 / 0.5)",
          "edit-color": "#007fff",
          "edit-base-color": "rgb(0 127 255 / 0.5)",
          "yellow-color": "#fec601",
          "yellow-base-color": "rgb(254 198 1 / 0.5)"
        }'::jsonb
    ),
    (
        'light',
        '{
          "bg-color": "#ffffff",
          "txt-color": "#111111",
          "hover-color": "#e8e8e8",
          "card-color": "#fcfcfc",
          "border-color": "#d0d0d0",
          "wrong-color": "#ed474a",
          "wrong-base-color": "rgb(207 41 44 / 0.5)",
          "warning-color": "#ff8200",
          "warning-base-color": "rgb(255 130 0 / 0.5)",
          "right-color": "#1ab369",
          "right-base-color": "rgb(26 179 105 / 0.5)",
          "edit-color": "#007fff",
          "edit-base-color": "rgb(0 127 255 / 0.5)",
          "yellow-color": "#fec601",
          "yellow-base-color": "rgb(254 198 1 / 0.5)"
        }'::jsonb
    );