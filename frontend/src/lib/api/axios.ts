import axios from "axios";

export const apiV1client = axios.create();

const callbacks = {
	logout: () => {
		console.error("logout callback not registered");
	},
};

export function initApiClient() {
	apiV1client.interceptors.response.use(
		(config) => config,
		(error) => {
			// HTTP Status 401 Unauthorized
			if (error.code === 401) {
				console.log("api intercepted: unauthenticated");
				callbacks.logout();
			}
		},
	);
}

type RegisterAuthCallbacks = {
	logout: () => void;
};
export function registerAuthCallbacks({ logout }: RegisterAuthCallbacks) {
	callbacks.logout = logout;
}
