import React from 'react';
import "../styles/App.css"

const Item = (props) => {
    return (
        <div className="item" id={props.id}>
            <strong>{props.text}</strong>
        </div>
    );
};

export default Item;