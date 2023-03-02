import axios, {CreateAxiosDefaults} from "axios";

let config: CreateAxiosDefaults = {baseURL: 'development' === process.env.MODE ? "http://127.0.0.1:8888" : ""};
const session = axios.create(config);

export default session
