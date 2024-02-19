import React, {useState} from 'react';
import Item from "./Item";
import axios from "axios";

const Resources = (props) => {
    let [array, setArray] = useState([])
    async function GetComputingResources() {
        const res = await axios.get("http://127.0.0.1:8080/computing_resources",)
        return res.data
    }
    GetComputingResources().then((data) => {
        setArray(data)
    })
    return (
        <div className="Main">
            <h2>Сервера-вычислители:</h2>
            {array.map((obj, id) => {
                return <Item text={"Computer " + obj.id} state={obj.state} />
            })}
        </div>
    );
};

export default Resources;