import getRequestClient from './getRequestClient';

import { browser } from '$app/environment';
type RedirectURL = string;

/**
 * Handles the backend communication for the authentication flow.
 */
export class Auth0Provider {
	constructor(cookies, readonly_cookie) {
		this.cookies = cookies;
		this.readonly_cookie = readonly_cookie;
		this.client = getRequestClient(readonly_cookie);
	}

	isAuthenticated() {
		// check with the backend if the token is valid
		return !!this.readonly_cookie.auth_token;
	}

	async login(returnTo) {
		const response = await this.client.auth.Login();
		this.cookies.set('state', response.state, { path: '/' });
		if (typeof window !== 'undefined' && window.sessionStorage) {
			// Assuming 'browser' refers to a browser environment
			window.sessionStorage.setItem(response.state, returnTo);
		}
		return response.auth_code_url;
	}

	async logout() {
		const response = await this.client.auth.Logout();
		this.cookies.remove('auth-token');
		this.cookies.remove('state');
		return response.redirect_url;
	}

	async validate(state, authCode) {
		if (state !== this.cookies.get('state')) throw new Error('Invalid state');
		const response = await this.client.auth.Callback({ code: authCode });
		this.cookies.set('auth-token', response.token, { path: '/' });
		if (typeof window !== 'undefined' && window.sessionStorage) {
			const returnURL = window.sessionStorage.getItem(state) ?? '/';
			window.sessionStorage.removeItem(state);
			return returnURL;
		}
		return '/';
	}
}
// export const Auth0Provider = (cookies) => {
//   client: getRequestClient(cookies),
//     isAuthenticated: () => !!cookies.get('auth-token'),

//       async login(returnTo: RedirectURL): Promise < RedirectURL > {
//         const response = await this.client.auth.Login();
//         cookies.set('state', response.state);
//         if(browser) sessionStorage.setItem(response.state, returnTo);
//         return response.auth_code_url;
//       },

//         async logout(): Promise < RedirectURL > {
//           const response = await this.client.auth.Logout();

//           cookies.remove('auth-token');
//           cookies.remove('state');

//           return response.redirect_url;
//         },

//           async validate(state: string, authCode: string,): Promise < RedirectURL > {
//             if(state != cookies.get('state')) throw new Error('Invalid state');

//   const response = await this.client.auth.Callback({ code: authCode });
//   cookies.set('auth-token', response.token);
//   if (browser) {
//     const returnURL = sessionStorage.getItem(state) ?? '/';
//     sessionStorage.removeItem(state);
//     return returnURL;
//   }
//   return '';
// }
// };
