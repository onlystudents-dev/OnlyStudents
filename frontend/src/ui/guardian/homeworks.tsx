import Navbar from "../../navbar/navbar.tsx";
import type {Me} from "../../app.tsx";

export default function GuardianHomeworks({ me }: {me: Me}) {
    return (
        <>
            <Navbar me={me} />
        </>
    )
}