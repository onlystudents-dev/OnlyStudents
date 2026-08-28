import type {Me} from "./app.tsx";
import Welcome from "./ui/welcome.tsx";

export default function Home({ me }: {me: Me}){
    return (
        <>
            <Welcome me={me} />
        </>
    )
}