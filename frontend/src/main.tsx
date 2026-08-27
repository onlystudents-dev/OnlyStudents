import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'
import "react-toastify/dist/ReactToastify.css";
import Home from './home.tsx'
import {ToastContainer} from "react-toastify";

createRoot(document.getElementById('root')!).render(
    <StrictMode>
        <Home />
        <ToastContainer theme={"dark"} position={"bottom-right"} />
    </StrictMode>,
)
