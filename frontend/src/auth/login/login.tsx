import "./login.css";
import Button from "../../util/button/button.tsx";
import {useState} from "react";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {faArrowRight, faQuestion} from "@fortawesome/free-solid-svg-icons";
import PasswordReset from "../pwr/pwr.tsx";
import {toast} from "react-toastify";
import Loading from "../../util/loading.tsx";

export default function Login() {
    const [pwr, setPwr] = useState(false);
    const [pwrA, setPwrA] = useState(false);

    const [red, setRed] = useState(false);
    const [wrong, setWrong] = useState(false);

    const [role, setRole] = useState("guardian");

    const [id, setId] = useState("");
    const [password, setPassword] = useState("");

    const [waiting, setWaiting] = useState(false);

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
                        <input className={`${(red || wrong) && "wrong"}`} type="text" placeholder="User ID" value={id} onChange={(e) => {setId(e.target.value); checkUserID(e.target.value); setWrong(false)}} />
                        <input className={`${wrong && "wrong"}`} type="password" placeholder="Password" value={password} onChange={(e) => {setPassword(e.target.value); setWrong(false)}} />
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
                {pwr && <PasswordReset pwrA={pwrA} unsPwr={unsPwr} role={role} setRole={setRole} id={id} setId={setId} red={red} checkUserID={checkUserID} />}
                {waiting && <Loading />}
            </div>
        </>
    )

    function checkUserID(id: string) {
        if (id && Number.isNaN(Number(id))) {
            setRed(true);
        } else {
            setRed(false);
        }
    }

    function sPwr() {
        setPwr(true)
        setPwrA(true)
    }

    function unsPwr() {
        setPwrA(false)
        setTimeout(() => setPwr(false), 200)
    }

    async function login() {
        if (red) return
        if (wrong) return
        if (!id) return
        if (!password) return
        setWaiting(true)
        try {
            const response = await fetch("/api/login", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({
                    user: Number(id),
                    password: password,
                    role: role,
                })
            })
            switch (response.status) {
                case 200:
                    location.reload()
                    break
                case 401:
                    toast.error("The username doesn't pair with the password!")
                    setWrong(true)
                    break
                default:
                    toast.error(response.statusText)
            }
        } finally {
            setWaiting(false)
        }
    }
}