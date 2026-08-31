import "./me.css";
import Navbar from "../../navbar/navbar.tsx";
import type {Me} from "../../app.tsx";
import Button from "./button.tsx";
import {faLock, faRightFromBracket, faUser} from "@fortawesome/free-solid-svg-icons";
import {useRef, useState} from "react";
import {toast} from "react-toastify";
import Loading from "../../util/loading.tsx";

export default function Me({ me }: {me: Me}) {
    const ref = useRef<HTMLDivElement>(null);
    const sref = useRef<HTMLDivElement>(null);

    const [waiting, setWaiting] = useState(false);
    const [hidden, setHidden] = useState(true);

    const [active, setActive] = useState<"user" | "security">("user");

    return (
        <>
            <Navbar me={me} />
            <div className="main">
                {(() => {
                   switch(active) {
                       case "user":
                           return (
                               <>
                                   cica
                               </>
                           )
                       case "security":
                           return (
                               <>
                                   secure cica
                               </>
                           )
                       default:
                           return null
                   }
                })()}
            </div>
            <div className="sidebar">
                <div className="flex flex-row gap-4">
                    <img className="rounded-[50%] h-20" src={me.pfp_url} alt={me.last_name} />
                    <div className="flex flex-col gap-1 justify-center min-w-0">
                        <h1 className="text-2xl truncate">
                            {me.first_name} {me.last_name}
                        </h1>
                        <p>Settings</p>
                    </div>
                </div>
                <br />
                <div className="buttons" onMouseLeave={() => setHidden(true)}>
                    <div ref={ref} className={`sbutton ${hidden && "hid"}`}></div>
                    <div ref={sref} className={`sbutton ssbutton`}></div>
                    <Button onClick={button => {setActive("user"); sreposition(button)}} onMouseEnter={reposition} icon={faUser} text="User settings" />
                    <Button onClick={button => {setActive("security"); sreposition(button)}} onMouseEnter={reposition} icon={faLock} text="Security" />
                    <Button onClick={async (button) => {await logout(); sreposition(button)}} onMouseEnter={reposition} className="text-(--wrong-color)" icon={faRightFromBracket} text="Logout" />
                </div>
            </div>
            {waiting && <Loading />}
        </>
    )

    async function logout() {
        setWaiting(true)
        const response = await fetch("/api/v1/me/logout", {
            method: "POST",
        });

        switch (response.status) {
            case 200:
                location.reload();
                break;
            default:
                toast.error(await response.text() || response.statusText);
        }
        setWaiting(false)
    }

    function reposition(button: HTMLButtonElement) {
        const current = ref.current;
        if (!current) return;
        setHidden(false)

        const buttonRect = button.getBoundingClientRect();
        const containerRect = current.parentElement?.getBoundingClientRect();

        if (!containerRect) return;

        current.style.top = `${buttonRect.top - containerRect.top}px`;
    }

    function sreposition(button: HTMLButtonElement) {
        const current = sref.current;
        if (!current) return;

        const buttonRect = button.getBoundingClientRect();
        const containerRect = current.parentElement?.getBoundingClientRect();

        if (!containerRect) return;

        current.style.top = `${buttonRect.top - containerRect.top}px`;
    }
}