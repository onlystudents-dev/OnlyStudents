import type {Me} from "../types/api.ts";

export const themes = ["", "dark", "light"] as const;
export type Theme = typeof themes[number];

export function getAutoTheme(): Theme {
    return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light'
}

export function getTheme(me: Me): Theme {
    return me.theme || getAutoTheme()
}

export function applyTheme(arg: Me | Theme) {
    const theme = typeof arg === "string" ? arg : getTheme(arg);
    document.documentElement.setAttribute('data-theme', theme);
}

export async function setTheme(theme: string): Promise<[Response, () => void]> {
    return [await fetch("/api/v1/me/change_theme", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({
            theme
        })
    }), () => applyTheme(theme as Theme)]
}
