import "./pwr.css";
import Button from "../../button/button.tsx";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {faArrowLeft, faPaperPlane, faPlane} from "@fortawesome/free-solid-svg-icons";
import React, {useState} from "react";

export default function PasswordReset({ pwrA, unsPwr, role, setRole, id, setId }: { pwrA: boolean, unsPwr: () => void, role: string, setRole: React.Dispatch<React.SetStateAction<string>>, id: string, setId: React.Dispatch<React.SetStateAction<string>> }) {
    const [reset, setReset] = useState(false);
    const [resetA, setResetA] = useState(false);

    return (
        <>
            <div className={`content in outback ${!pwrA && "hid"}`}>
                <h1 className="self-center text-5xl font-bold mb-8 rubik">Reset password</h1>
                {!reset && <div className={`loginput fredoka out ${resetA && "hid"}`}>
                    <select value={role} onChange={(e) => setRole(e.target.value)}>
                        <option value="guardian">Guardian</option>
                        <option value="student">Student</option>
                        <option value="teacher">Teacher</option>
                    </select>
                    <input type="text" placeholder="User ID" value={id} onChange={(e) => setId(e.target.value)} />
                </div>}
                {reset && <div className="loginput fredoka in">
                    <input type="text" placeholder="Code" />
                    <input type="password" placeholder="Password" />
                    <input type="password" placeholder="Confirm password" />
                </div>}

                <div className={`pwrbutton ${reset ? "mt-4" : "mt-20"}`}>
                    <Button onClick={unsPwr}>
                        <FontAwesomeIcon icon={faArrowLeft} /> Back
                    </Button>
                    {!reset && <Button onClick={email} className={`outbottom ${resetA && "hid"}`}>
                        Send email <FontAwesomeIcon icon={faPaperPlane} />
                    </Button>}
                    {reset && <Button onClick={change} className="in">
                        Reset <FontAwesomeIcon icon={faPlane} />
                    </Button>}
                </div>
            </div>
        </>
    )

    async function email() {
        // no api yet

        setResetA(true)
        setTimeout(() => setReset(true), 200)
    }

    async function change() {
        // api
    }
}