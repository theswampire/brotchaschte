import { apiV1client } from "./axios";

export async function getKek() {
	await new Promise(resolve => setTimeout(resolve, 3000));
	const { data } = await apiV1client.get("/api/v1/kek");

	return data;
}
