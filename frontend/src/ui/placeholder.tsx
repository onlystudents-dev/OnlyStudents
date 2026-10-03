import Navbar from "../navbar/navbar.tsx";
import type {Me} from "../types/api.ts";

export default function Placeholder({ me }: {me: Me}) {
    return <Navbar me={me} />
}
