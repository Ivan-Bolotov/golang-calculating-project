import React, {useEffect, useRef, useState} from 'react';
import Item from "./ExpressionItem";
import API from "../classes/API";

const Calculator = (props) => {
    let [arr, setArr] = useState([])
    const getData = async () => {
        const res = await API.GetExpressions()
            if (res !== null) {
                setArr(res)
            }
        }

    useEffect(() => {getData()}, []);
    const ref = useRef(null);
    return (
        <div className="Main">
            <div>
                <input ref={ref} type="text" placeholder="enter the expression to calculate" style={{padding: "10px", margin: "10px", width: "auto", minWidth: 200}}/>
                <button onClick={async () => {
                    const data = await API.PostNewExpression(ref.current.value)
                    ref.current.value = ""
                    if (typeof data === "object") {
                        await getData()
                    } else {
                        alert(data)
                    }
                }} style={{padding: 10}}>Отправить</button>
            </div>
            <div>
                {arr.map((el) => {
                    return <Item key={el.id} data={el}/>
                })}
            </div>
        </div>
    );
};

export default Calculator;