import "./welcome.css";
import type {Me} from "../types/api.ts";
import Navbar from "../navbar/navbar.tsx";

export default function Welcome({ me }: {me: Me}) {
    return (
        <>
            <Navbar me={me} />
            <div className="w-full full-height flex justify-center items-center relative">
                <h1 className="hello fredoka">Hello, {me.preferences.nickname || me.last_name}!</h1>
            </div>
        </>
    )
}