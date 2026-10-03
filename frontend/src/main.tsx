import { render } from 'preact'
import { useState } from 'preact/hooks'
import './index.css'
import "react-toastify/dist/ReactToastify.css"
import App from './app.tsx'

function Root() {
    const [reload, setReload] = useState(false)

    return (
        <App
            key={reload}
            reload={() => setReload(prev => !prev)}
        />
    )
}

const root = document.getElementById('root')
if (!root) throw new Error("#root element not found")

render(<Root />, root)