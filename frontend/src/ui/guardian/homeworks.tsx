import Navbar from "../../navbar/navbar.tsx";
import type {Me} from "../../types/api.ts";

export default function GuardianHomeworks({ me }: {me: Me}) {
    return (
        <>
            <Navbar me={me} />
        </>
    )
}
