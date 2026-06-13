import { navigate } from "sv-router/generated";
import { createContext } from "svelte";

export interface AuthContext {
	is_authenticated: boolean;
	user: any;
}

export const [getAuthContext, setAuthContext] = createContext<AuthContext>();

export const auth = {
	isAuthenticated: () => {
		console.error("auth not initialized");
		return false;
	},
	setAuthenticated: (_user: AuthContext["user"]) => {
		console.log("auth not initialized");
	},
	setUnauthenticated: () => {
		console.log("auth not initialized");
	},
	getUser: () => {
		console.log("auth not initialized");
		return null as AuthContext["user"];
	},
};

export function initAuthStuff() {
	const authctx = getAuthContext();

	const isAuthenticated = () => {
		return authctx.is_authenticated;
	};
	const setAuthenticated = (user: AuthContext["user"]) => {
		authctx.is_authenticated = true;
		authctx.user = user;
	};
	const setUnauthenticated = () => {
		authctx.is_authenticated = false;
	};
	const getUser = () => {
		return authctx.user;
	};

	auth.isAuthenticated = isAuthenticated;
	auth.setAuthenticated = setAuthenticated;
	auth.setUnauthenticated = setUnauthenticated;
	auth.getUser = getUser;
}

export function navigateToLoginWithCallback(current: string) {
	const sp = new URLSearchParams();
	sp.append("callback", current);
	return navigate("/login", { search: sp.toString() });
}
