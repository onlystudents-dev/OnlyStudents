import "./skeleton.css";
import {cx} from "../cx.ts";

export default function Skeleton({ width, height, color, className }: {width: number, height: number, color: string, className?: string}) {
    return <div className={cx("skeleton", className)} style={{ width: width, height: height, "--skeleton-color": color }} />
}