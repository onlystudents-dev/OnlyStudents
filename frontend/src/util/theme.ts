import type {Me} from "../types/api.ts";

export const themes = ["", "dark", "light"] as const;
export type Theme = typeof themes[number];

export function isTheme(value: string): value is Theme {
    return themes.some(theme => theme === value)
}

export function getAutoTheme(): Theme {
    return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light'
}

export function getTheme(me: Me): Theme {
    return me.preferences.theme || getAutoTheme()
}

export function applyTheme(arg: Me | Theme) {
    const theme = typeof arg === "string" ? arg : getTheme(arg);
    document.documentElement.setAttribute('data-theme', theme);
}
