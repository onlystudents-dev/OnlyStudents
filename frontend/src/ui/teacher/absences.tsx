import Navbar from "../../navbar/navbar.tsx";
import type {Me} from "../../types/api.ts";

export default function TeacherAbsences({ me }: {me: Me}) {
    return (
        <>
            <Navbar me={me} />
        </>
    )
}
