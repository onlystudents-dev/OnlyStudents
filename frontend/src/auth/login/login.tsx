import "./login.css";
import Button from "../../button/button.tsx";
import Overlay from "../../util/overlay.tsx";
import React, {useState} from "react";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {faArrowRight, faQuestion} from "@fortawesome/free-solid-svg-icons";
import PasswordReset from "../pwr/pwr.tsx";

export default function Login({ setLogin }: { setLogin: React.Dispatch<React.SetStateAction<boolean>> }) {
    const [pwr, setPwr] = useState(false);
    const [pwrA, setPwrA] = useState(false);

    const [role, setRole] = useState("guardian");

    const [id, setId] = useState("");
    const [password, setPassword] = useState("");

    return (
        <>
            <div className="box cantar">
                <div className={`content out ${pwrA && "hid"}`}>
                    <h1 className="self-center text-5xl font-bold mb-8 rubik">Login</h1>
                    <div className="loginput fredoka">
                        <select value={role} onChange={(e) => setRole(e.target.value)}>
                            <option value="guardian">Guardian</option>
                            <option value="student">Student</option>
                            <option value="teacher">Teacher</option>
                        </select>
                        <input type="text" placeholder="User ID" value={id} onChange={(e) => setId(e.target.value)} />
                        <input type="password" placeholder="Password" value={password} onChange={(e) => setPassword(e.target.value)} />
                    </div>
                    <div className="logbutton">
                        <Button onClick={sPwr}>
                            <FontAwesomeIcon icon={faQuestion} /> Forgot password
                        </Button>
                        <Button onClick={login}>
                            Login <FontAwesomeIcon icon={faArrowRight} />
                        </Button>
                    </div>
                </div>
                {pwr && <PasswordReset pwrA={pwrA} unsPwr={unsPwr} role={role} setRole={setRole} id={id} setId={setId} />}
            </div>

            <Overlay onClick={() => setLogin(false)} time={400} />
        </>
    )

    function sPwr() {
        setPwr(true)
        setPwrA(true)
    }

    function unsPwr() {
        setPwrA(false)
        setTimeout(() => setPwr(false), 200)
    }

    async function login() {
        // api
    }
}