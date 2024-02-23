import React from 'react';
import "../styles/App.css"

const ExpressionItem = (props) => {
    return (
        <div className="item" id={props.id}>
            <h2>{props.data.id} expression: {props.data.exp}</h2>
            <strong>State: {props.data.state}, Result: {props.data.res}</strong>
        </div>
    );
};

export default ExpressionItem;