import "./sidebar.css";
import Skeleton from "../skeleton/skeleton.tsx";
import { cloneElement, isValidElement, type VNode } from "preact";
import type { ComponentChildren } from "preact";
import { useState } from "preact/hooks";
import { getKey } from "../language.ts";
import type { Me } from "../../types/api.ts";
import {Children} from "preact/compat";

type ClickableProps = {
    onClick?: (e: MouseEvent) => void;
    onMouseEnter?: (e: MouseEvent) => void;
};

export default function Sidebar({ children, me }: { children: ComponentChildren; me?: Me }) {
    const [hoverTop, setHoverTop] = useState(0)
    const [selectTop, setSelectTop] = useState(0)

    const [hidden, setHidden] = useState(true)
    const [pfpLoaded, setPfpLoaded] = useState(false)

    function getTop(e: MouseEvent): number {
        const button = e.currentTarget as Element
        const container = button.closest(".buttons")

        if (!container) return 0

        return button.getBoundingClientRect().top - container.getBoundingClientRect().top
    }

    function handleEnter(e: MouseEvent) {
        setHidden(false)
        setHoverTop(getTop(e))
    }

    function handleClick(e: MouseEvent) {
        setSelectTop(getTop(e))
    }

    return (
        <div className="sidebar">
            {me && <>
                <div className="flex flex-row gap-4">
                    <div className="rounded-[50%] h-20 cursor-pointer">
                        <img className={`rounded-[inherit] max-h-full ${pfpLoaded ? "" : "hidden"}`} src={me.preferences.pfp_url} alt="" onLoad={() => setPfpLoaded(true)} />
                        {!pfpLoaded && <Skeleton width={40} height={40} color="var(--bg-color)" className="rounded-[inherit]!"/>}
                    </div>

                    <div className="flex flex-col gap-1 justify-center min-w-0">
                        <h1 className="text-2xl truncate">
                            {me.preferences.nickname || `${me.first_name} ${me.last_name}`}
                        </h1>
                        <p>{getKey("SETTINGS")}</p>
                    </div>
                </div>

                <br />
            </>}

            <div className="buttons" onMouseLeave={() => setHidden(true)}>
                <div
                    className={`sbutton ${hidden ? "hid" : ""}`}
                    style={{ top: `${hoverTop}px` }}
                />

                <div
                    className="sbutton ssbutton"
                    style={{ top: `${selectTop}px` }}
                />

                {Children.toArray(children).map((child) => {
                    if (!isValidElement(child)) return child

                    const element = child as VNode<ClickableProps>

                    return cloneElement(element, {
                        onClick: (e: MouseEvent) => {
                            handleClick(e)
                            element.props.onClick?.(e)
                        },
                        onMouseEnter: (e: MouseEvent) => {
                            handleEnter(e)
                            element.props.onMouseEnter?.(e)
                        },
                    })
                })}
            </div>
        </div>
    );
}