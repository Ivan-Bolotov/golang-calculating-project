import React from 'react';
import "../styles/App.css"

const Item = (props) => {
    return (
        <div className="item">
            <strong>{props.text} state: {props.state}</strong>
        </div>
    );
};

export default Item;