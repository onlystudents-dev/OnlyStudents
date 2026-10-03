import type {Me} from "./types/api.ts";
import Welcome from "./ui/welcome.tsx";

export default function Home({ me }: {me: Me}){
    return (
        <>
            <Welcome me={me} />
        </>
    )
}