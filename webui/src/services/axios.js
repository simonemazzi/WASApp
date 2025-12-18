import axios from "axios";

const instance = axios.create({
	baseURL: __API_URL__,
	timeout: 1000 * 5
});

export const setAuthHeader = token => {
	instance.defaults.headers.common.Authorization = `Bearer ${token}`;
};

instance.interceptors.request.use(config => {
	if (!config.headers.Date) {
		config.headers.Date = new Date().toUTCString();
	}
	return config;
});

export const doLogin = async (username) => {
	try {
		const response = await instance.post("/session", {
			name: username
		});

		return {
			userId: response.data.userId,
			token: response.data.token,
			time: response.data.time
		};
	} catch (error) {
		console.error("Login error:", error);
		throw error;
	}
};


export const getConversations = async (userId,token) => {
	setAuthHeader(token);
	try{
		const response = await instance.get(`users/${userId}/conversations`);
		return response.data;
	} catch (error) {
		console.error("Get conversations error:", error);
		throw error;
	}
}



export const getGroups = async (userId,token) => {
	setAuthHeader(token);
	try{
		const response = await instance.get(`users/${userId}/groups`);
		return response.data;
	} catch (error) {
		console.error("Get conversations error:", error);
		throw error;
	}
}



















export default instance;
