import React from 'react';
import Item from "./Item";

const Settings = (props) => {
    let operands = ["+", "-", "/", "*"]
    return (
        <div className="Main">
            <h2>Введите значения в секундах: </h2>
            {operands.map((el, i) => {
                return <input type="text" id={i} placeholder={el} style={{margin: "10px", padding: "10px"}}/>
            })}
        </div>
    );
};

export default Settings;