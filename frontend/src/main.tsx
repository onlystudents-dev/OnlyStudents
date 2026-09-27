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

render(<Root />, document.getElementById('root')!)