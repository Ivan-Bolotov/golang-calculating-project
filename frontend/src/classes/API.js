import axios from "axios";

export default class API {
    static async PostNewExpression(expression) {
        const data = await axios.post("http://localhost:8080/expression", {"expression": expression})
        return data.data
    }
    static async GetExpressions() {
        const data = await axios.get("http://localhost:8080/expressions")
        return data.data
    }
    static async GetComputingResources() {
        const data = await axios.get("http://127.0.0.1:8080/computing_resources")
        return data.data
    }
    static async SetOperationTime(object) {
        const data = await axios.post("http://127.0.0.1:8080/set_operation_time", object)
        return data.data
    }
    static async GetOperationTime() {
        const data = await axios.get("http://127.0.0.1:8080/operation_time")
        return data.data
    }
}
