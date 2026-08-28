import type {Me} from "../app.tsx";
import Navbar from "../navbar/navbar.tsx";

export default function Welcome({ me }: {me: Me}) {
    return (
        <>
            <Navbar />
            <div className="w-full full-height flex justify-center items-center">
                <h1 className="text-3xl fredoka">Hello, {me.last_name}!</h1>
            </div>
        </>
    )
}