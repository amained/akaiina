import { redirect } from '@sveltejs/kit';
import { Auth0Provider } from '$lib/auth';
export async function load({ request, cookies }) {
	const url = new URL(request.url);
	const searchParams = new URLSearchParams(url.search);
	const returnToURL = searchParams.get('returnTo') ?? '/';
	const auth0 = new Auth0Provider(cookies, {
		auth_token: cookies.get('auth-token'),
		state: cookies.get('state')
	});
	if (auth0.isAuthenticated()) return redirect(200, returnToURL);

	const returnURL = await auth0.login(returnToURL);
	redirect(307, returnURL);
}
