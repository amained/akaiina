import type { PageServerLoad } from './$types';
import type { Cookies } from 'svelte';

export const load: PageServerLoad = async ({ cookies }: { cookies: Cookies }) => {
	return {
		page_server_data: { auth_token: cookies.get('auth-token'), state: cookies.get('state') }
	};
};
