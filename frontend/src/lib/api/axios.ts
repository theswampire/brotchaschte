import axios from "axios";

export const apiV1client = axios.create();

const callbacks = {
	logout: () => {
		console.error("logout callback not registered");
	},
};

export function initApiClient() {
	console.log("register intercepts");
	apiV1client.interceptors.response.use(
		(config) => config,
		(error) => {
			if (axios.isAxiosError(error)) {
				// HTTP Status 401 Unauthorized
				if (error.status === 401) {
					console.log("api intercepted res: unauthenticated");
					callbacks.logout();
				}
			}
			throw error;
		},
	);
}

type RegisterAuthCallbacks = {
	logout: () => void;
};
export function registerAuthCallbacks({ logout }: RegisterAuthCallbacks) {
	callbacks.logout = logout;
}
