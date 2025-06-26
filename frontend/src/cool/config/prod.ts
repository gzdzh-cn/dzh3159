import { proxy } from "./proxy";

export default {
	// 根地址
	host: proxy["/prod"].target,

	// 请求地址
	baseUrl: "http://127.0.0.1:8200"
};
