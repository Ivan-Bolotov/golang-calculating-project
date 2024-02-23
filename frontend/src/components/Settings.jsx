import React, {useEffect, useRef, useState} from 'react';
import API from "../classes/API";

const Settings = (props) => {
    let operands = ["+", "-", "/", "*"]
    let [elements, setElements] = useState([useRef(null), useRef(null), useRef(null), useRef(null)])
    useEffect(() => {
        const func = async () => {
            let data = await API.GetOperationTime()
            elements.forEach((el) => {
                switch (el.current.placeholder) {
                    case "+":
                        el.current.value = Number(data["+"]) / 10 ** 9
                        break
                    case "-":
                        el.current.value = Number(data["-"]) / 10 ** 9
                        break
                    case "/":
                        el.current.value = Number(data["/"]) / 10 ** 9
                        break
                    case "*":
                        el.current.value = Number(data["*"]) / 10 ** 9
                        break
                }
            })
        }
        func()
    }, []);
    return (
        <div className="Main">
            <h2>Введите значения в секундах: </h2>
            {operands.map((el, i) => {
                return (
                    <div key={i} style={{display: "flex", flexDirection: "column", margin: "10px"}}>
                        <label htmlFor={String(i)}><strong>({el})</strong></label>
                        <input id={String(i)} ref={elements[i]} type="text" placeholder={el} style={{flexGrow: 1, marginTop: "10px", padding: "10px"}}/>
                    </div>
                )
            })}
            <button style={{margin: 10, padding: 10}} onClick={async (event) => {
                let object = new Map()
                elements.forEach((el) => {
                    if (!isNaN(Number(el.current.value))) {
                        // object.set(el.current.placeholder, Number(el.current.value))
                        object[el.current.placeholder] = Number(el.current.value)
                    }
                })
                console.log(await API.SetOperationTime(object));
            }}>Submit</button>
        </div>
    );
};

export default Settings;