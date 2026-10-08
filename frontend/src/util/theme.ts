import type {Me} from "../types/api.ts";

export function getAutoTheme(): string {
    return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light'
}

export function getTheme(me: Me): string {
    return me.preferences.theme || getAutoTheme()
}

export function applyTheme(arg: Me | string) {
    const theme = typeof arg === "string" ? arg ? arg : getAutoTheme() : getTheme(arg); // spaghetti
    document.documentElement.setAttribute('data-theme', theme);
}
