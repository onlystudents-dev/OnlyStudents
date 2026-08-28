import type {Me} from "../../app.tsx";
import Navbar from "../../navbar/navbar.tsx";

export default function GuardianHomeworks({ me }: {me: Me}) {
    return (
        <>
            <Navbar />
            {me}
        </>
    )
}