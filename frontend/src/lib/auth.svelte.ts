import { createContext } from "svelte";

export interface AuthContext {
	is_authenticated: boolean;
	user: any;
}

/*
let authctx = $state<AuthContext>({ is_authenticated: false, user: null });

export function getAuthContext(): AuthContext {
	return authctx;
}
*/

export const [getAuthContext, setAuthContext] = createContext<AuthContext>();

export function isAuthenticated() {
	const authctx = getAuthContext();
	return authctx.is_authenticated;
}
