import { apiV1client } from "./axios";

export type LoginRes = {
	ok: boolean;
	email?: string;
	error?: string;
};

export function login(creds: FormData) {
	return apiV1client.post<LoginRes>("/api/v1/auth/login", creds);
}
