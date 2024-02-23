import React from 'react';

const ResourcesItem = (props) => {
    return (
        <div className="item" id={props.id}>
            <h2>Computer {props.data.id}</h2>
            <strong>State: {props.data.state}, Last ping: {props.data.lastping}</strong>
        </div>
    );
};

export default ResourcesItem;