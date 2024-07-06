import getRequestClient from './getRequestClient';
import type { Cookies } from '@sveltejs/kit';
import type Client from '../client';

export type READONLY_COOKIE = { auth_token: string; state: string };
/**
 * Handles the backend communication for the authentication flow.
 */
export class Auth0Provider {
	cookies: Cookies | null;
	readonly_cookie: READONLY_COOKIE;
	client: Client;
	constructor(cookies: Cookies | null, readonly_cookie: { auth_token: string; state: string }) {
		this.cookies = cookies;
		this.readonly_cookie = readonly_cookie;
		this.client = getRequestClient(readonly_cookie);
	}

	isAuthenticated() {
		// check with the backend if the token is valid
		return !!this.readonly_cookie.auth_token;
	}

	async login(returnTo: string) {
		const response = await this.client.auth.Login();
		if (this.cookies != null) this.cookies.set('state', response.state, { path: '/' });
		return response.auth_code_url;
	}

	async logout() {
		const response = await this.client.auth.Logout();
		if (this.cookies != null) {
			this.cookies.delete('auth-token', { path: '/' });
			this.cookies.delete('state', { path: '/' });
		}

		return response.redirect_url;
	}

	async validate(state: string, authCode: string) {
		if (this.cookies != null && state !== this.cookies.get('state'))
			throw new Error('Invalid state');
		const response = await this.client.auth.Callback({ code: authCode });
		if (this.cookies != null) this.cookies.set('auth-token', response.token, { path: '/' });
		if (typeof window !== 'undefined' && window.sessionStorage) {
			const returnURL = window.sessionStorage.getItem(state) ?? '/';
			window.sessionStorage.removeItem(state);
			return returnURL;
		}
		return '/';
	}
}
