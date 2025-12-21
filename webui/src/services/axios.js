import axios from "axios";

const instance = axios.create({
	baseURL: __API_URL__,
	timeout: 1000 * 5
});


export const BASE_URL = "http://localhost:3000";

instance.interceptors.request.use(config => {
	const token = localStorage.getItem("token");
	if (token) {
		config.headers.Authorization = `Bearer ${token}`;
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


export const getConversations = async (userId) => {
	try{
		const response = await instance.get(`users/${userId}/conversations`);
		return response.data;
	} catch (error) {
		console.error("Get conversations error:", error);
		throw error;
	}
}



export const getGroups = async  (userId) => {
	try{
		const response = await instance.get(`users/${userId}/groups`);
		return response.data;
	} catch (error) {
		console.error("Get conversations error:", error);
		throw error;
	}
}


export const getMessages = async (userId,conversationId) => {
	try{
		const response = await instance.get(`users/${userId}/conversations/${conversationId}/messages`);
		return response.data;
	}catch(error){
		console.error("Get conversations error:", error);
		throw error;
	}
}

export const getConversation = async (userId,conversationId) => {
	try{
		const response = await instance.get(`users/${userId}/conversations/${conversationId}`);
		return response.data;
	}catch(error){
		console.error("Get conversations error:", error);
		throw error;
	}
}

















export default instance;
