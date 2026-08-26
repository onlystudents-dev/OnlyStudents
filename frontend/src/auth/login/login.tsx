import "./login.css";
import Button from "../../button/button.tsx";
import Overlay from "../../util/overlay.tsx";
import React from "react";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {faArrowRight, faQuestion} from "@fortawesome/free-solid-svg-icons";

export default function Login({ setLogin }: { setLogin: React.Dispatch<React.SetStateAction<boolean>> }) {
    return (
        <>
            <div className="login cantar">
                <h1 className="self-center text-5xl font-bold mb-8">Login</h1>
                <div className="loginput">
                    <input type="text" placeholder="User ID" />
                    <input type="password" placeholder="Password" />
                </div>
                <div className="logbutton">
                    <Button>
                        <FontAwesomeIcon icon={faQuestion} /> Forgot password
                    </Button>
                    <Button>
                        <FontAwesomeIcon icon={faArrowRight} /> Login
                    </Button>
                </div>
            </div>

            <Overlay onClick={() => setLogin(false)} time={400} />
        </>
    )
}