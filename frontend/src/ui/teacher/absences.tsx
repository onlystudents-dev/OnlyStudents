import Navbar from "../../navbar/navbar.tsx";
import type {Me} from "../../app.tsx";

export default function TeacherAbsences({ me }: {me: Me}) {
    return (
        <>
            <Navbar me={me} />
        </>
    )
}
