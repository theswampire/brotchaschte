import type { Hooks } from "sv-router";
import { navigate } from "sv-router/generated";
import { isAuthenticated } from "@/lib/auth.svelte";

export default {
	beforeLoad(ctx) {
		// Execute before route loads
		console.log(ctx);
		if (ctx.meta.public) return;

		console.log("check auth");

		if (isAuthenticated()) return;

		console.log("unauthenticated");

		const sp = new URLSearchParams();
		sp.append("callback", ctx.pathname);
		throw navigate("/login", { search: sp.toString() });
	},
	afterLoad() {
		// Execute after route loads
	},
	onPreload() {
		// Execute when the route is preloaded
	},
	onError() {
		// Execute when route components fail to load or beforeLoad hooks throw an error
	},
} satisfies Hooks;
