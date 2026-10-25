// i hate eslint AND typecheck
// @ts-nocheck
/* eslint-disable */
export function selectStyles(width: string = "20rem", height: string = "2.5rem", backgroundColor: string = "var(--card-color)", hoverColor: string = "var(--hover-color)") {
    return {
        control: (base: any) => ({
            ...base,
            width,
            height,
            minHeight: '2.5rem',
            borderRadius: '0.5rem',
            backgroundColor,
            border: 'none',
            boxShadow: 'none',
            cursor: 'text',
            padding: '0 0.25rem',
        }),
        indicatorSeparator: (base: any) => ({
            ...base,
            marginLeft: '0.5rem',
        }),
        valueContainer: (base: any) => ({
            ...base,
            padding: '0',
        }),
        input: (base: any) => ({
            ...base,
            color: 'inherit',
            margin: 0,
            padding: 0,
        }),
        singleValue: (base: any) => ({
            ...base,
            color: 'inherit',
        }),
        menu: (base: any) => ({
            ...base,
            backgroundColor,
            borderRadius: '0.5rem',
            overflow: 'hidden',
            zIndex: 50,
            marginTop: '0.25rem',
        }),
        option: (base: any, state: any) => ({
            ...base,
            backgroundColor: state.isFocused ? hoverColor : 'transparent',
            color: 'inherit',
            cursor: 'pointer',
            padding: '0.5rem 1rem',
        }),
    }
}
// stupid name
export function selectStyles2(height: string = "100%", items: string = "end") {
    return {
        container: (base: any) => ({
            ...base,
            position: 'absolute',
            inset: 0,
        }),
        control: (base: any) => ({
            ...base,
            height,
            minHeight: 0,
            backgroundColor: 'transparent',
            border: 'none',
            boxShadow: 'none',
            borderRadius: '0.5rem',
            cursor: 'pointer',
            padding: '0 0.5rem',
        }),
        valueContainer: (base: any) => ({
            ...base,
            padding: 0,
            justifyItems: items,
        }),
        singleValue: (base: any) => ({
            ...base,
            color: 'inherit',
            margin: 0,
        }),
        dropdownIndicator: (base: any) => ({
            ...base,
            color: 'inherit',
            padding: '0 0 0 0.25rem',
        }),
        menu: (base: any) => ({
            ...base,
            left: 0,
            width: '100%',
            backgroundColor: 'var(--card-color)',
            borderRadius: '0.5rem',
            overflow: 'hidden',
            zIndex: 50,
            marginTop: '0.25rem',
        }),
        option: (base: any, state: any) => ({
            ...base,
            backgroundColor: state.isFocused ? 'var(--hover-color)' : 'transparent',
            color: 'inherit',
            cursor: 'pointer',
            padding: '0.5rem 1rem',
            ':active': { backgroundColor: 'var(--hover-color)' },
        }),
    }
}
/* eslint-enable */

export type Option = {
    value: string;
    label: string;
}
