import "./styles/App.css"
import icon from "./images/calc.png"
import {useState} from "react";
import Calculator from "./components/Calculator";
import Settings from "./components/Settings";
import Resources from "./components/Resources";

function App() {
    let [current, setCurrent] = useState(1);
    fetch("http://127.0.0.1:8080/computing_resources", {
        method: "GET",
        mode: 'no-cors',
        headers: {"Accept": "application/json"}
    }).then((res) => res.text())
        .then((json) => {
        console.log(json)
    });
    return (
        <div>
            <div className="Header">
                <h2 className="label">distributed calculator</h2>
                <img src={icon} width="50" height="50" alt="Icon" style={{margin: 5}}/>
            </div>
            <div className="Panel">
                <button className="btn" onClick={() => {
                    setCurrent(1)
                }}>CALCULATOR
                </button>
                <button className="btn" onClick={() => {
                    setCurrent(2)
                }}>CALCULATION SETTINGS
                </button>
                <button className="btn" onClick={() => {
                    setCurrent(3)
                }}>COMPUTING RESOURCES
                </button>
            </div>
            {current === 1 ? <Calculator array={[1, 2, 3, 4, "qweqr"]}/> : current === 2 ? <Settings/> :
                <Resources array={["Computer1", "Computer2", "Computer3"]}/>}
        </div>
    );
}

export default App;
