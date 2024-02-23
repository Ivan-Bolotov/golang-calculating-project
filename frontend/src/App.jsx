import "./styles/App.css"
import icon from "./images/calc.png"
import {useEffect, useRef, useState} from "react";
import Calculator from "./components/Calculator";
import Settings from "./components/Settings";
import Resources from "./components/Resources";

function App() {
    let [current, setCurrent] = useState(1);
    let btns = [useRef(null), useRef(null), useRef(null)]
    useEffect(() => {
        btns.map((el, idx) => {
            if (idx === current - 1) {
                el.current.style.color = "teal"
                el.current.style.fontWeight = "bold"
            } else {
                el.current.style.color = "gray"
                el.current.style.fontWeight = "normal"
            }
        })
    }, [current]);
    return (
        <div>
            <div className="Header">
                <h2 className="label">distributed calculator</h2>
                <img src={icon} width="50" height="50" alt="Icon" style={{margin: 5}}/>
            </div>
            <div className="Panel">
                <button ref={btns[0]} className="btn" onClick={() => {
                    setCurrent(1)
                }}>CALCULATOR
                </button>
                <button ref={btns[1]} className="btn" onClick={() => {
                    setCurrent(2)
                }}>CALCULATION SETTINGS
                </button>
                <button ref={btns[2]} className="btn" onClick={() => {
                    setCurrent(3)
                }}>COMPUTING RESOURCES
                </button>
            </div>
            {current === 1 ? <Calculator/> : current === 2 ? <Settings/> :
                <Resources/>}
        </div>
    );
}

export default App;
