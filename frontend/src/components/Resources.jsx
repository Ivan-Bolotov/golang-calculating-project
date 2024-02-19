import React from 'react';
import Item from "./Item";

const Resources = (props) => {
    return (
        <div>
            {props.array.map((el, i) => {
                return <Item id={i} text={el}/>
            })}
        </div>
    );
};

export default Resources;