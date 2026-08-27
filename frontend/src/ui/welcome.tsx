import type {Me} from "../home.tsx";
import Navbar from "../navbar/navbar.tsx";

export default function Welcome({ me }: {me: Me}) {
    return (
        <>
            <Navbar />
            <div className="w-full full-height flex justify-center items-center">
                <h1 className="text-3xl fredoka">Hello, {me.role}!</h1>
            </div>
        </>
    )
}