import "./me.css";
import Navbar from "../../navbar/navbar.tsx";
import type {Me} from "../../app.tsx";
import Button from "./button.tsx";
import {
    faAddressCard,
    faEnvelope, faFloppyDisk,
    faLock, faPaperPlane,
    faRightFromBracket,
    faUser,
    faUserLock
} from "@fortawesome/free-solid-svg-icons";
import {useRef, useState} from "react";
import Loading from "../../util/loading.tsx";
import {getKey} from "../../util/language.ts";
import Config from "./config.tsx";
import {toast} from "react-toastify";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";

export default function Me({ me }: {me: Me}) {
    const ref = useRef<HTMLDivElement>(null)
    const sref = useRef<HTMLDivElement>(null)

    const [waiting, setWaiting] = useState(false)
    const [hidden, setHidden] = useState(true)

    const [active, setActive] = useState<"user" | "security">("user")

    const [email, setEmail] = useState(me.email)
    const [settingEmail, setSettingEmail] = useState(false)
    const [settingEmailA, setSettingEmailA] = useState(false)
    const newEmail = useRef(null)
    const oldEmailCode = useRef(null)
    const newEmailCode = useRef(null)

    const [nickname, setNickname] = useState(me.nickname)
    const [password, setPassword] = useState("")

    return (
        <>
            <Navbar me={me} />
            <div className="main">
                <div className="configs border-(--txt-color) w-full h-fit min-h-80 flex flex-col items-center justify-start p-4 border-4 rounded-2xl gap-2">
                {(() => {
                   switch(active) {
                       case "user":
                           return (
                               <>
                                   <Config text={getKey("EMAIL")} icon={faEnvelope} value={email}>
                                       <h1 className="text-center text-5xl font-bold mb-8 rubik">{getKey("CHANGE_EMAIL")}</h1>
                                       <div className="moving">
                                           {settingEmail ? (<div className="moving-content second">
                                               <input type="text" placeholder={getKey("OLD_EMAIL_CODE")} ref={oldEmailCode} />
                                               <input type="text" placeholder={getKey("NEW_EMAIL_CODE")} ref={newEmailCode} />

                                               <button className="absolute bottom-0 right-0">
                                                   <FontAwesomeIcon icon={faFloppyDisk} /> {getKey("SAVE")}
                                               </button>
                                           </div>) : (<div className={`moving-content ${settingEmailA && "first"}`}>
                                               <input type="email" placeholder={getKey("NEW_EMAIL")} ref={newEmail} />

                                               <button className="absolute bottom-0 right-0" onClick={sendEmail}>
                                                   <FontAwesomeIcon icon={faPaperPlane} /> {getKey("SEND_EMAIL")}
                                               </button>
                                           </div>)}
                                       </div>
                                   </Config>
                                   <Config text={getKey("NICKNAME")} icon={faAddressCard} value={nickname}>
                                       cicamica
                                   </Config>
                               </>
                           )
                       case "security":
                           return (
                               <>
                                   <Config text={getKey("PASSWORD")} icon={faUserLock} value={""}>
                                       jelszocica
                                   </Config>
                               </>
                           )
                       default:
                           return null
                   }
                })()}
                </div>
            </div>
            <div className="sidebar">
                <div className="flex flex-row gap-4">
                    <img className="rounded-[50%] h-20" src={me.pfp_url} alt={me.last_name} />
                    <div className="flex flex-col gap-1 justify-center min-w-0">
                        <h1 className="text-2xl truncate">
                            {me.first_name} {me.last_name}
                        </h1>
                        <p>{getKey("SETTINGS")}</p>
                    </div>
                </div>
                <br />
                <div className="buttons" onMouseLeave={() => setHidden(true)}>
                    <div ref={ref} className={`sbutton ${hidden && "hid"}`}></div>
                    <div ref={sref} className={`sbutton ssbutton`}></div>
                    <Button onClick={button => {setActive("user"); sreposition(button)}} onMouseEnter={reposition} icon={faUser} text={getKey("USER_SETTINGS")} />
                    <Button onClick={button => {setActive("security"); sreposition(button)}} onMouseEnter={reposition} icon={faLock} text={getKey("SECURITY")} />
                    <Button onClick={async (button) => {await logout(); sreposition(button)}} onMouseEnter={reposition} className="text-(--wrong-color)" icon={faRightFromBracket} text={getKey("LOGOUT")} />
                </div>
            </div>
            {waiting && <Loading />}
        </>
    )

    async function sendEmail() {
        setSettingEmailA(true)
        setTimeout(() => setSettingEmail(true), 500)
    }

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