import { useId, useMemo, useState } from "react"
import {getKey, getLanguage} from "../../util/language.ts";
import {formatUnixDate} from "../../util/time.ts";
import type {Grade, Me} from "../../types/api.ts";
import {colors} from "./grades.tsx";
import type {TargetedPointerEvent} from "preact";

const WIDTH = 600
const HEIGHT = 200
const MARGIN = { top: 14, right: 10, bottom: 14, left: 40 }
const CHART_WIDTH = WIDTH - MARGIN.left - MARGIN.right
const CHART_HEIGHT = HEIGHT - MARGIN.top - MARGIN.bottom
const BASELINE = HEIGHT - MARGIN.bottom
const LEVELS = [5, 4, 3, 2, 1] as const
const MAX_VISIBLE_DOTS = 24

type Pt = { x: number; y: number }
type LevelStyle = { color: string; base: string }

const FALLBACK_STYLE: LevelStyle = { color: "#84cc16", base: "rgba(132, 204, 22, 0.22)" }

const clampGrade = (value: number) => Math.max(1, Math.min(5, Math.round(value)))

const xForIndex = (index: number, count: number) => MARGIN.left + (count <= 1 ? CHART_WIDTH / 2 : (index / (count - 1)) * CHART_WIDTH)

const yForAverage = (average: number) => MARGIN.top + ((5 - average) / 4) * CHART_HEIGHT

function at<T>(arr: readonly T[], index: number): T {
    return arr[index] as T
}

function smoothPath(pts: Pt[]): string {
    const n = pts.length
    if (n === 0) return ""
    if (n < 3) return pts.map((p, i) => `${i ? "L" : "M"} ${p.x} ${p.y}`).join(" ")

    const dx: number[] = []
    const slope: number[] = []
    for (let i = 0; i < n - 1; i++) {
        const a = at(pts, i)
        const b = at(pts, i + 1)
        dx.push(b.x - a.x)
        slope.push((b.y - a.y) / (b.x - a.x))
    }

    const tan: number[] = [at(slope, 0)]
    for (let i = 1; i < n - 1; i++) {
        const prev = at(slope, i - 1)
        const next = at(slope, i)
        tan.push(prev * next <= 0 ? 0 : (prev + next) / 2)
    }
    tan.push(at(slope, n - 2))

    for (let i = 0; i < n - 1; i++) {
        const m = at(slope, i)
        if (m === 0) {
            tan[i] = 0
            tan[i + 1] = 0
            continue
        }
        const a = at(tan, i) / m
        const b = at(tan, i + 1) / m
        const s = a * a + b * b
        if (s > 9) {
            const tau = 3 / Math.sqrt(s)
            tan[i] = tau * a * m
            tan[i + 1] = tau * b * m
        }
    }

    const start = at(pts, 0)
    let d = `M ${start.x} ${start.y}`
    for (let i = 0; i < n - 1; i++) {
        const a = at(pts, i)
        const b = at(pts, i + 1)
        const h = at(dx, i) / 3
        d += ` C ${a.x + h} ${a.y + at(tan, i) * h} ${b.x - h} ${b.y - at(tan, i + 1) * h} ${b.x} ${b.y}`
    }
    return d
}

function readThemePalette(): Map<number, LevelStyle> {
    const palette = new Map<number, LevelStyle>()
    const style = typeof window === "undefined" ? null : getComputedStyle(document.documentElement)

    for (const level of LEVELS) {
        const key = (colors as Record<number, string | undefined>)[level]
        palette.set(level, {
            color: (key && style?.getPropertyValue(`--${key}-color`).trim()) || FALLBACK_STYLE.color,
            base: (key && style?.getPropertyValue(`--${key}-base-color`).trim()) || FALLBACK_STYLE.base,
        })
    }
    return palette
}

const levelStyle = (palette: Map<number, LevelStyle>, value: number) => palette.get(clampGrade(value)) ?? FALLBACK_STYLE

export default function SubjectAverageChart({ grades, me }: { grades: Grade[]; me: Me }) {
    const lang = getLanguage(me)?.key || "en-US"
    const gradientId = useId().replaceAll(":", "")
    const [hoverIndex, setHoverIndex] = useState<number | null>(null)

    const { series, startMonth } = useMemo(() => {
        const sorted = grades.filter((g) => g.value >= 1 && g.value <= 5).sort((a, b) => a.date - b.date)
        const only = sorted.at(0)
        const source = sorted.length === 1 && only ? [only, only] : sorted

        const items: { timestamp: number; date: string; average: number }[] = []
        let sum = 0
        for (const g of source) {
            sum += g.value
            items.push({ timestamp: g.date, date: formatUnixDate(g.date, lang), average: Number((sum / (items.length + 1)).toFixed(2)) })
        }

        const first = sorted.at(0)
        const firstDate = first ? new Date(first.date > 1e11 ? first.date : first.date * 1000) : new Date()
        const month = firstDate.toLocaleDateString(lang, { month: "long" })

        return { series: items, startMonth: month.charAt(0).toUpperCase() + month.slice(1) }
    }, [grades, lang])

    const points = series.map((p, i) => ({ ...p, x: xForIndex(i, series.length), y: yForAverage(p.average) }))
    const first = points.at(0)
    const last = points.at(-1)

    if (!first || !last) return null

    const palette = readThemePalette()
    const shown = (hoverIndex === null ? undefined : points.at(hoverIndex)) ?? last
    const shownStyle = levelStyle(palette, shown.average)
    const activeLevel = clampGrade(shown.average)

    const linePath = smoothPath(points)
    const areaPath = `${linePath} L ${last.x} ${BASELINE} L ${first.x} ${BASELINE} Z`
    const showAllDots = points.length <= MAX_VISIBLE_DOTS

    const handlePointerMove = (e: TargetedPointerEvent<SVGSVGElement>) => {
        const rect = e.currentTarget.getBoundingClientRect()
        const x = ((e.clientX - rect.left) / rect.width) * WIDTH
        const index = Math.round(((x - MARGIN.left) / CHART_WIDTH) * (points.length - 1))
        setHoverIndex(Math.max(0, Math.min(points.length - 1, index)))
    }

    return (
        <div className="w-160 max-w-[92vw] rounded-3xl bg-(--hover-color) p-5 border-b-2 border-r-2 border-(--border-color)">
            <div className="mb-2 flex items-baseline justify-between gap-4 px-1 select-none">
                <span className="text-4xl font-bold leading-none tabular-nums" style={{ color: shownStyle.color, fontFamily: "Poppins, sans-serif" }}>
                    {shown.average.toFixed(2)}
                </span>
                <span className="text-sm opacity-70 fredoka">{shown.date}</span>
            </div>

            <svg
                viewBox={`0 0 ${WIDTH} ${HEIGHT}`}
                className="block h-auto w-full touch-pan-y overflow-visible"
                role="img"
                aria-label={`Subject average over time. Current average: ${last.average.toFixed(2)}`}
                onPointerMove={handlePointerMove}
                onPointerDown={handlePointerMove}
                onPointerLeave={() => setHoverIndex(null)}
                onPointerCancel={() => setHoverIndex(null)}
            >
                <defs>
                    <linearGradient id={`${gradientId}-stroke`} gradientUnits="userSpaceOnUse" x1="0" y1={MARGIN.top} x2="0" y2={BASELINE}>
                        {LEVELS.map((level, i) => (
                            <stop key={level} offset={`${i * 25}%`} stopColor={levelStyle(palette, level).color} />
                        ))}
                    </linearGradient>

                    <linearGradient id={`${gradientId}-fill`} gradientUnits="userSpaceOnUse" x1={first.x} y1="0" x2={last.x} y2="0">
                        {points.map((point, index) => (
                            <stop key={`${point.timestamp}-${index}`} offset={`${(index / (points.length - 1)) * 100}%`} stopColor={levelStyle(palette, point.average).color} />
                        ))}
                    </linearGradient>

                    <linearGradient id={`${gradientId}-fade`} gradientUnits="userSpaceOnUse" x1="0" y1={MARGIN.top} x2="0" y2={BASELINE}>
                        <stop offset="0%" stopColor="#fff" stopOpacity={0.4} />
                        <stop offset="100%" stopColor="#fff" stopOpacity={0} />
                    </linearGradient>

                    <mask id={`${gradientId}-mask`} maskUnits="userSpaceOnUse" x="0" y="0" width={WIDTH} height={HEIGHT}>
                        <rect x="0" y={MARGIN.top} width={WIDTH} height={CHART_HEIGHT} fill={`url(#${gradientId}-fade)`} />
                    </mask>
                </defs>

                {LEVELS.map((level) => {
                    const y = yForAverage(level)
                    const isActive = level === activeLevel

                    return (
                        <g key={level}>
                            <line x1={MARGIN.left} x2={WIDTH - MARGIN.right} y1={y} y2={y} stroke="rgba(161,161,170,0.22)" strokeWidth={1.5} strokeDasharray="2 9" strokeLinecap="round" />

                            {isActive && <circle cx={MARGIN.left / 2} cy={y} r={14} fill={shownStyle.base} />}

                            <text
                                x={MARGIN.left / 2}
                                y={y}
                                dy={6}
                                textAnchor="middle"
                                fill={isActive ? shownStyle.color : "rgba(161,161,170,0.55)"}
                                fontSize={18}
                                fontWeight="bold"
                                fontFamily="Poppins, sans-serif"
                            >
                                {level}
                            </text>
                        </g>
                    )
                })}

                <path d={areaPath} fill={`url(#${gradientId}-fill)`} mask={`url(#${gradientId}-mask)`} />
                <path d={linePath} fill="none" stroke={`url(#${gradientId}-stroke)`} strokeWidth={5} strokeLinecap="round" strokeLinejoin="round" />

                {hoverIndex !== null && <line x1={shown.x} x2={shown.x} y1={MARGIN.top} y2={BASELINE} stroke={shownStyle.color} strokeOpacity={0.35} strokeWidth={1.5} />}

                {showAllDots && points.map((point, index) => (
                    <circle key={`${point.timestamp}-${index}`} cx={point.x} cy={point.y} r={3} fill={levelStyle(palette, point.average).color} />
                ))}

                <circle cx={shown.x} cy={shown.y} r={11} fill={shownStyle.base} />
                <circle cx={shown.x} cy={shown.y} r={6} fill={shownStyle.color} stroke="var(--hover-color)" strokeWidth={3} />
            </svg>

            <div className="mt-2 flex select-none items-center justify-between pl-11 pr-2 text-sm opacity-70 fredoka">
                <span>{startMonth}</span>
                <span>{getKey("NOW")}</span>
            </div>
        </div>
    )
}