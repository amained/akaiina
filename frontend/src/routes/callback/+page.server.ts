import { redirect } from '@sveltejs/kit';
import { Auth0Provider } from '$lib/auth';
export async function load({ request, cookies }) {
	const url = new URL(request.url);
	const searchParams = new URLSearchParams(url.search);
	const state = searchParams.get('state');
	const code = searchParams.get('code');

	if (!state || !code) throw new Error('Login failed');
	const auth0 = new Auth0Provider(cookies, {
		auth_token: cookies.get('auth-token'),
		state: cookies.get('state')
	});
	const redirectURL = await auth0.validate(state, code);
	return redirect(307, redirectURL);
}
