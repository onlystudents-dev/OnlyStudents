import { render } from 'preact'
import './index.css'
import {fetchLanguage, getKey} from "./util/language.ts";

const root = document.getElementById('root')
if (!root) throw new Error("#root element not found")

await fetchLanguage(null).then(() =>
    render(
        <>
            <div className="cantar fixed w-5xl h-42 flex flex-col gap-4 justify-between items-center bg-(--hover-color) border-4 border-(--border-color) border-solid rounded-2xl p-4">
                <p className="text-4xl">{getKey("NOT_FOUND")}</p>
                <div className="rounded-2xl bg-(--bg-color) p-4 wrong">
                    GET {location.pathname}
                </div>
            </div>
        </>
    , root)
)