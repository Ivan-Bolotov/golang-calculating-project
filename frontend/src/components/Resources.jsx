import React, {useEffect, useState} from 'react';
import Item from "./ResourcesItem";
import API from "../classes/API";

const Resources = (props) => {
    let [arr, setArr] = useState([])
    const getData = async () => {
        const res = await API.GetComputingResources()
        if (res !== null) {
            setArr(res)
        }
    }

    useEffect(() => {getData()}, []);
    return (
        <div className="Main">
            <h2>Сервера-вычислители:</h2>
            {arr.map((el) => {
                return <Item key={el.id} data={el}/>
            })}
        </div>
    );
};

export default Resources;