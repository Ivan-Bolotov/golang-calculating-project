import React from 'react';
import Item from "./Item";

const Calculator = (props) => {
    return (
        <div className="Main">
            <input type="text" placeholder={"enter the expression to calculate"} style={{padding: "10px", margin: "10px"}}/>
            <div>
                {props.array.map((el, i) => {
                    return <Item id={i} text={el}/>
                })}
            </div>
        </div>
    );
};

export default Calculator;